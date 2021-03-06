package controller

import (
	"fmt"
	"io"
	"io/ioutil"
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
	// Send a get language list request to MS
	req, err := http.NewRequest("GET", MSTranslateServer+languageEndpoint, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating MS request: %v", err), http.StatusInternalServerError)
		return
	}

	client := getHTTPClient()
	msResp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending request to MS server: %v", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		err := msResp.Body.Close()
		if err != nil {
			log.Errorf("Error closing response body stream: %v", err)
		}
	}()

	// Set response header
	w.Header().Set("Content-Type", msResp.Header["Content-Type"][0])
	w.WriteHeader(msResp.StatusCode)

	// Copy resonse body if status is not OK
	if msResp.StatusCode != http.StatusOK {
		_, err = io.Copy(w, msResp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error copying MS response body: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Convert to google format language list and write it back
	msBody, err := ioutil.ReadAll(msResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading MS response body: %v", err), http.StatusInternalServerError)
		return
	}
	body, err := language.ToGoogleLanguageList(msBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to google language list: %v", err), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Errorf("Error writing response body for translate requests: %v", err)
	}
}

// Translate converts a Google format translate request into a Microsoft format
// one which will be send to the Microsoft server, and write a Google format
// response back to the client.
func Translate(w http.ResponseWriter, r *http.Request) {
	// Convert google format request to MS format
	req, isAuto, err := translate.ToMicrosoftRequest(r, MSTranslateServer)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to MS request: %v", err), http.StatusBadRequest)
		return
	}

	// Send translate request to MS server
	client := getHTTPClient()
	msResp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending request to MS server: %v", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		err := msResp.Body.Close()
		if err != nil {
			log.Errorf("Error closing response body stream: %v", err)
		}
	}()

	// Set Header
	w.Header().Set("Content-Type", msResp.Header["Content-Type"][0])
	w.Header().Set("Access-Control-Allow-Origin", "*") // same as Google response

	// Copy resonse body if status is not OK
	if msResp.StatusCode != http.StatusOK {
		w.WriteHeader(msResp.StatusCode)
		_, err = io.Copy(w, msResp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error copying MS response body: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Set google format response body
	msBody, err := ioutil.ReadAll(msResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading MS response body: %v", err), http.StatusInternalServerError)
		return
	}
	body, err := translate.ToGoogleResponseBody(msBody, isAuto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to google response body: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(msResp.StatusCode)
	_, err = w.Write(body)
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
