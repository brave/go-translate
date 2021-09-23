package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/brave/go-translate/language"
	"github.com/brave/go-translate/translate"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

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

// GetLanguageList generates language list in google format and replies back to the client.
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

func EncodeTags(text string) (string, []string) {
	balance := 0
	token := ""
	tokens := []string {}
	cleanText := ""
	elements := []rune {}
	for _, char := range text {
		elements = append(elements, char)
	}
	i := 0
	last_close := -2
	last_open := -2
	for i < len(elements) {
		if elements[i] == '<' {
			if balance == 0 {
				last_open = i
			}
			balance += 1
		}
		if balance == 0 {
			if (i + 2 < len(elements)) && (elements[i] == '(') && (elements[i + 1] == '1') && (elements[i + 2] == ')') {
				tokens = append(tokens, " (1) ")
				cleanText += " (1) "
				i += 3
				continue
			} else {
				cleanText += string(elements[i])
			}
		} else {
			token += string(elements[i])
		}
		if elements[i] == '>' {
			balance -= 1
			if balance == 0 {
				if last_close + 1 == last_open {
					tokens[len(tokens) - 1] += token
				} else {
					cleanText += " (1) "
					tokens = append(tokens, token)
				}
				token = ""
				last_close = i
			}
		}
		i += 1
	}
	return cleanText, tokens
}

func DecodeTags(text string, tokens []string) (string) {
	tokenIndex := 0
	decodedText := ""
	elements := []rune {}
	for _, char := range text {
		elements = append(elements, char)
	}
	i := 0
	for i < len(elements) {
		if (i + 2 < len(elements)) && (elements[i] == '(') && (elements[i + 1] == '1') && (elements[i + 2] == ')') {
			if tokenIndex == len(tokens) {
				decodedText += ""
			} else {
				decodedText += tokens[tokenIndex]
				tokenIndex += 1
			}
			i += 3
		} else {
			decodedText += string(elements[i])
			i += 1
		}
	}
	return decodedText
}

// Translate translates input texts with Bergamot and sends back to the client.
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
		var encodedTexts []string
		var tagTokens [][]string
		encodedTexts = make([]string, len(texts))
		tagTokens = make([][]string, len(texts))
		for ind, text := range texts {
			encodedTexts[ind], tagTokens[ind] = EncodeTags(text)
		}
		idList := originalTextIdsByLanguages[language]
		processedTexts := []string{}
		if language != "unknown" {
			translateOutput, err := translate.TranslateTexts(encodedTexts, language, to)
			if err == nil {
				processedTexts = translateOutput
			}
		}
		if len(processedTexts) == 0 {
			processedTexts = encodedTexts
		}
		for ind, processedText := range processedTexts {
			if ind >= len(idList) {
				break
			}
			decProcessedText := DecodeTags(processedText, tagTokens[ind])
			translatedTexts[idList[ind]] = decProcessedText
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
	http.Redirect(w, r, GoogleTranslateServerProxy + r.URL.Path, http.StatusTemporaryRedirect)
}

// GetGStaticResource redirect the requests from gstatic to brave's proxy.
func GetGStaticResource(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, GStaticServerProxy + r.URL.Path, http.StatusTemporaryRedirect)
}
