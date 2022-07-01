package controller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/brave/go-translate/language"
	"github.com/brave/go-translate/translate"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

var TranslationServiceHost = os.Getenv("LNX_HOST")
var TranslationServiceApiKey = os.Getenv("LNX_API_KEY")

func Router() chi.Router {
	r := chi.NewRouter()

	// Relay
	r.Get("/translate_a/l", RelayLanguageListRequest)
	r.Post("/translate_a/t", RelayTranslateRequest)

	// Static
	r.Get("/static/v1/element.js", ServeStaticFile)
	r.Get("/static/v1/js/element/main.js", ServeStaticFile)
	r.Get("/static/v1/css/translateelement.css", ServeStaticFile)

	return r
}

func RelayLanguageListRequest(w http.ResponseWriter, _ *http.Request) {
	// Relay request to target
	res, err := RequestLanguageListFromTarget()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error requesting language list from target: %v", err), http.StatusInternalServerError)
	}

	if res.StatusCode != http.StatusOK { // TODO: needed? Copy resonse body if status is not OK)
		_, err = io.Copy(w, res.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error copying response body: %v", err), http.StatusInternalServerError)
		}
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert response
	convertedResBody, err := language.ToSourceResponseBody(resBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond to source
	w.Header().Set("Content-Type", res.Header["Content-Type"][0])
	w.WriteHeader(res.StatusCode)

	_, err = w.Write(convertedResBody)
	if err != nil {
		log.Errorf("Error writing response body: %v", err)
	}
}

func RequestLanguageListFromTarget() (*http.Response, error) {
	req, err := http.NewRequest("GET", TranslationServiceHost+"/get-languages", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+TranslationServiceApiKey)

	client := getHTTPClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Errorf("Error closing response body stream: %v", err)
		}
	}()

	return res, err
}

func RelayTranslateRequest(w http.ResponseWriter, r *http.Request) {
	// Build request body for target
	body, err := translate.BuildTargetRequestBody(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error building target request body: %v", err), http.StatusBadRequest)
	}

	// Relay request to target
	res, err := RequestTranslationFromTarget(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error requesting translation from target: %v", err), http.StatusInternalServerError)
	}

	if res.StatusCode != http.StatusOK { // TODO: needed? Copy resonse body if status is not OK)
		w.WriteHeader(res.StatusCode)
		_, err = io.Copy(w, res.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error copying target response body: %v", err), http.StatusInternalServerError)
		}
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading target response body: %v", err), http.StatusInternalServerError)
		return
	}

	convertedBody, err := translate.ToResourceResponseBody(resBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to response body: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond to source
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", res.Header["Content-Type"][0])
	w.WriteHeader(res.StatusCode)

	_, err = w.Write(convertedBody)
	if err != nil {
		log.Errorf("Error writing body: %v", err)
	}
}

func RequestTranslationFromTarget(body []byte) (*http.Response, error) {
	targetUrl := TranslationServiceHost + "/translate"
	u, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	req.Header.Add("Authorization", "Bearer "+TranslationServiceApiKey)

	client := getHTTPClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Errorf("Error closing target response body stream: %v", err)
		}
	}()

	return res, err
}

func ServeStaticFile(w http.ResponseWriter, r *http.Request) {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	fs := http.FileServer(filesDir)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Security-Policy", "require-trusted-types-for 'script'")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Cache-Control", "max-age=86400")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")

	fs.ServeHTTP(w, r)
}

func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}
