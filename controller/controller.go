package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/brave-intl/bat-go/libs/logging"
	"github.com/brave-intl/bat-go/libs/middleware"
	"github.com/brave/go-translate/language"
	"github.com/brave/go-translate/translate"
	"github.com/go-chi/chi"
)

// LnxEndpoint specifies the remote Lnx translate server used by
// brave-core, and it can be set to a mock server during testing.
var LnxEndpoint = os.Getenv("LNX_HOST")
var LnxAPIKey = os.Getenv("LNX_API_KEY")
var languagePath = "/get-languages"
var translatePath = "/translate"

// TranslateRouter add routers for translate requests and translate script
// requests.
func TranslateRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/translate_a/t", middleware.InstrumentHandler("Translate", http.HandlerFunc(Translate)).ServeHTTP)
	r.Get("/translate_a/l", middleware.InstrumentHandler("GetLanguageList", http.HandlerFunc(GetLanguageList)).ServeHTTP)

	r.Get("/static/v1/element.js", middleware.InstrumentHandler("ServeStaticFile", http.HandlerFunc(ServeStaticFile)).ServeHTTP)
	r.Get("/static/v1/js/element/main.js", middleware.InstrumentHandler("ServeStaticFile", http.HandlerFunc(ServeStaticFile)).ServeHTTP)
	r.Get("/static/v1/css/translateelement.css", middleware.InstrumentHandler("ServeStaticFile", http.HandlerFunc(ServeStaticFile)).ServeHTTP)

	return r
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
		Timeout: time.Second * 60,
	}
}

// GetLanguageList send a request to Lingvanex server and convert the response
// into google format and reply back to the client.
func GetLanguageList(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())

	// Send a get language list request to Lnx
	req, err := http.NewRequest("GET", LnxEndpoint+languagePath, nil)
	req.Header.Add("Authorization", "Bearer "+LnxAPIKey)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Lnx request: %v", err), http.StatusInternalServerError)
		return
	}

	client := getHTTPClient()
	lnxResp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending request to Lnx server: %v", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		err := lnxResp.Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Error closing response body stream")
		}
	}()

	// Set response header
	w.Header().Set("Content-Type", lnxResp.Header["Content-Type"][0])
	w.WriteHeader(lnxResp.StatusCode)

	// Copy resonse body if status is not OK
	if lnxResp.StatusCode != http.StatusOK {
		_, err = io.Copy(w, lnxResp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error copying Lnx response body: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Convert to google format language list and write it back
	lnxBody, err := ioutil.ReadAll(lnxResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading Lnx response body: %v", err), http.StatusInternalServerError)
		return
	}
	body, err := language.ToGoogleLanguageList(lnxBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to google language list: %v", err), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		logger.Error().Err(err).Msg("Error writing response body for translate requests")
	}
}

// Translate converts a Google format translate request into a Lingvanex format
// one which will be send to the Lingvanex server, and write a Google format
// response back to the client.
func Translate(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())

	w.Header().Set("Access-Control-Allow-Origin", "*") // same as Google response

	req, isAuto, err := translate.ToLingvanexRequest(r, LnxEndpoint+translatePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to LnxEndpoint request: %v", err), http.StatusBadRequest)
		return
	}

	req.Header.Add("Authorization", "Bearer "+LnxAPIKey)

	// Send translate request to Lnx server
	client := getHTTPClient()
	lnxResp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending request to LnxEndpoint: %v", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		err := lnxResp.Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Error closing response body stream")
		}
	}()

	// Set Header
	w.Header().Set("Content-Type", lnxResp.Header["Content-Type"][0])
	w.Header().Set("Access-Control-Allow-Origin", "*") // same as Google response

	// Copy resonse body if status is not OK
	if lnxResp.StatusCode != http.StatusOK {
		w.WriteHeader(lnxResp.StatusCode)
		_, err = io.Copy(w, lnxResp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error copying LnxEndpoint response body: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Set google format response body
	lnxBody, err := ioutil.ReadAll(lnxResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading LnxEndpoint response body: %v", err), http.StatusInternalServerError)
		return
	}
	body, err := translate.ToGoogleResponseBody(lnxBody, isAuto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to google response body: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(lnxResp.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		logger.Error().Err(err).Msg("Error writing response body for translate requests")
	}
}
