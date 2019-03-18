package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/brave/go-translate/translate"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// MSTranslateServer specifies the remote MS translate server used by
// brave-core, and it can be set to other hosts during testing.
var MSTranslateServer = "https://api.cognitive.microsofttranslator.com"

// TranslateRouter add routers for translate requests and translate script
// requests.
func TranslateRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/translate", Translate)
	return r
}

func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
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
