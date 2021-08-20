package controller

import (
	"encoding/json"
	"fmt"
	"github.com/grokify/html-strip-tags-go"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/brave/go-translate/language"
	"github.com/brave/go-translate/translate"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// MSTranslateServer specifies the remote MS translate server used by
// brave-core, and it can be set to a mock server during testing.
var MSTranslateServer = "https://api.cognitive.microsofttranslator.com"

// GoogleTranslateServerProxy specifies the proxy server for requesting
// resource from google translate server, and it can be set to a mock server
// during testing.
var GoogleTranslateServerProxy = "https://translate.brave.com"

const (
	// GoogleTranslateServer specifies the remote google translate server.
	GoogleTranslateServer = "https://translate.googleapis.com"

	// GStaticServerProxy specifies the proxy server for requesting resource
	// from google gstatic server.
	GStaticServerProxy = "https://translate-static.brave.com"

	languageEndpoint = "/languages?api-version=3.0&scope=translation"
)

// TranslateRouter add routers for translate requests and translate script
// requests.
func TranslateRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/translate", Translate)
	r.Get("/language", GetLanguageList)

	r.Get("/translate_a/element.js", GetTranslateScript)
	r.Get("/element/*/js/element/element_main.js", GetTranslateScript)
	r.Get("/translate_static/js/element/main.js", GetTranslateScript)

	r.Get("/translate_static/css/translateelement.css", GetGoogleTranslateResource)
	r.Get("/images/branding/product/1x/translate_24dp.png", GetGStaticResource)
	r.Get("/images/branding/product/2x/translate_24dp.png", GetGStaticResource)

	return r
}

func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}

// GetLanguageList send a request to Microsoft server and convert the response
// into google format and reply back to the client.
func GetLanguageList(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*") // same as Google response
	w.WriteHeader(http.StatusOK)

	body, err := language.GetBergamotLanguageList()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting language list: %v", err), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Errorf("Error writing response body for translate requests: %v", err)
	}
}

func Translate(w http.ResponseWriter, r *http.Request) {
	slVals := r.URL.Query()["sl"]
	if len(slVals) != 1 {
		http.Error(w, fmt.Sprintf("Error parsing the request"), http.StatusBadRequest)
		return;
	}
	tlVals := r.URL.Query()["tl"]
	if len(tlVals) != 1 {
		http.Error(w, fmt.Sprintf("Error parsing the request"), http.StatusBadRequest)
		return;
	}
	from := slVals[0]
	to := tlVals[0]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing the request: %v", err), http.StatusBadRequest)
		return
	}

	allTexts := r.PostForm["q"]
	textCount := len(allTexts)

	originalTextsByLanguages := map[string][]string {}
	originalTextIdsByLanguages := map[string][]int {}
	for id, text := range allTexts {
		text = strip.StripTags(text)
		textLanguage := from
		if textLanguage == "auto" {
			detectedLanguage, cldErr := translate.DetectLanguage(text)
			if cldErr != nil {
				textLanguage = "unknown"
			} else {
				textLanguage = detectedLanguage
			}
		}
		textsForLanguage, ok := originalTextsByLanguages[textLanguage]
		if ok {
			originalTextsByLanguages[textLanguage] = append(textsForLanguage, text)
		} else {
			originalTextsByLanguages[textLanguage] = []string{text}
		}
		idsForLanguage, ok := originalTextIdsByLanguages[textLanguage]
		if ok {
			originalTextIdsByLanguages[textLanguage] = append(idsForLanguage, id)
		} else {
			originalTextIdsByLanguages[textLanguage] = []int{id}
		}
	}

	var translatedTexts = make([]string, textCount)
	for language, texts := range originalTextsByLanguages {
		idList := originalTextIdsByLanguages[language]
		processedTexts := []string{}
		if language != "unknown" {
			translateOutput, err := translate.TranslateTexts(texts, language, to)
			if err == nil {
				processedTexts = translateOutput
			}
		}
		if len(processedTexts) == 0 {
			processedTexts = texts
		}
		for ind, processedText := range processedTexts {
			if ind >= len(idList) {
				break
			}
			translatedTexts[idList[ind]] = processedText
		}
	}

	body := make([]string, textCount)
	for i, translatedText := range translatedTexts {
		body[i] = translatedText
	}

	mbody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error forming response body: %v", err), http.StatusBadRequest)
		return
	}

	// Set Header
	w.Header().Set("Content-Type", r.Header["Content-Type"][0])
	w.Header().Set("Access-Control-Allow-Origin", "*") // same as Google response
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(mbody)
	if err != nil {
		log.Errorf("Error writing response body for translate requests: %v", err)
	}
}

// GetTranslateScript use a reverse proxy to forward handle translate script
// requests to brave's proxy server. We're not replying a HTTP redirect
// directly because the Originin the response header will be cleared to null
// instead of having the value "https://translate.googleapis.com" when the
// client follows the redirect to a cross origin, so we would violate the CORS
// policy when the client try to access other resources from google server.
func GetTranslateScript(w http.ResponseWriter, r *http.Request) {
	target, _ := url.Parse(GoogleTranslateServerProxy)
	// Use a custom director so req.Host will be changed.
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

// GetGoogleTranslateResource redirect the resource requests from google
// translate server to brave's proxy.
func GetGoogleTranslateResource(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, GoogleTranslateServerProxy+r.URL.Path, http.StatusTemporaryRedirect)
}

// GetGStaticResource redirect the requests from gstatic to brave's proxy.
func GetGStaticResource(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, GStaticServerProxy+r.URL.Path, http.StatusTemporaryRedirect)
}
