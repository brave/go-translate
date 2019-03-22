package server

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brave/go-translate/controller"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

var handler http.Handler
var mockMSServerHandler http.Handler
var mockGoogleServerHandler http.Handler

func init() {
	err := os.Setenv("MS_TRANSLATE_API_KEY", "DUMMY_KEY")
	if err != nil {
		panic("Test init failed when setting API Key env.")
	}
}

// HandleTranslate is a HTTP handler function used by the mock MS translate
// server using in tests which responds to translate requests with pre-defined
// responses.
func HandleTranslate(w http.ResponseWriter, r *http.Request) {
	slVals := r.URL.Query()["from"]
	isAuto := len(slVals) == 0
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var body []byte
	if isAuto {
		body = []byte(`[{"detectedLanguage":{"language":"de","score":1.0},"translations":[{"text":"Good Evening","to":"en"}]},{"detectedLanguage":{"language":"de","score":1.0},"translations":[{"text":"Hello World","to":"en"}]}]`)
	} else {
		body = []byte(`[{"translations":[{"text":"Good Evening","to":"en"}]},{"translations":[{"text":"Hello World","to":"en"}]}]`)
	}

	_, err := w.Write(body)
	if err != nil {
		fmt.Printf("write response error: %v", err)
	}
}

// HandleLanguage is a HTTP handler function used by the mock MS translate
// server using in tests which responds to get language list requests with
// pre-defined responses.
func HandleLanguages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body := []byte(`{"translation":{"af":{"name":"Afrikaans","nativeName":"Afrikaans","dir":"ltr"},"en":{"name":"English","nativeName":"English","dir":"ltr"}}}`)
	_, err := w.Write(body)
	if err != nil {
		fmt.Printf("write response error: %v", err)
	}
}

// HandleGetTranslateScript is a HTTP handler function used by the mock proxy
// of google translate server in tests which returns OK for translate script
// requests.
func HandleGetTranslateScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", controller.GoogleTranslateServer)
	w.WriteHeader(http.StatusOK)
}

// Mock MS translate server for testing
func setupMockMSServerRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/translate", HandleTranslate)
	r.Get("/languages", HandleLanguages)
	return r
}

func setupMockGoogleServerRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/translate_a/element.js", HandleGetTranslateScript)
	r.Get("/element/*/js/element/element_main.js", HandleGetTranslateScript)
	r.Get("/translate_static/js/element/main.js", HandleGetTranslateScript)
	return r
}

func init() {
	testCtx, logger := setupLogger(context.Background())
	serverCtx, mux := setupRouter(testCtx, logger)
	handler = chi.ServerBaseContext(serverCtx, mux)
	mockMSServerHandler = chi.ServerBaseContext(serverCtx, setupMockMSServerRouter())
	mockGoogleServerHandler = chi.ServerBaseContext(serverCtx, setupMockGoogleServerRouter())
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, err := http.Get(server.URL)
	assert.Nil(t, err)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	expected := "."
	actual, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	if expected != string(actual) {
		t.Errorf("Expected the message '%s'\n", expected)
	}
}

func translate(t *testing.T, server *httptest.Server, url string, expectedBody []byte) {
	googleBody := []byte(`q=guten%20Abend&q=Hallo%20Welt`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(googleBody))
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, body)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestTranslate(t *testing.T) {
	server := httptest.NewServer(handler)
	defer server.Close()
	mockServer := httptest.NewServer(mockMSServerHandler)
	defer mockServer.Close()
	controller.MSTranslateServer = mockServer.URL

	// auto-detect source language
	url := fmt.Sprintf("%v/translate?anno=3&client=te_lib&format=html&v=1.0&key=DUMMY_KEY&logld=vTE_20181015_01&tl=en&sp=nmt&tc=1&sr=1&tk=568026.932804&mode=1&sl=", server.URL)
	expectedBody := []byte(`[["Good Evening","de"],["Hello World","de"]]`)
	translate(t, server, url+"auto", expectedBody)

	// known source language
	expectedBody = []byte(`["Good Evening","Hello World"]`)
	translate(t, server, url+"de", expectedBody)
}

func getLanguage(t *testing.T, server *httptest.Server, url string, expectedBody []byte) {
	resp, err := http.Get(url)
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, body)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestGetLanguage(t *testing.T) {
	server := httptest.NewServer(handler)
	defer server.Close()
	mockMSServer := httptest.NewServer(mockMSServerHandler)
	defer mockMSServer.Close()
	controller.MSTranslateServer = mockMSServer.URL

	url := fmt.Sprintf("%v/language?client=chrome&hl=en", server.URL)
	expectedBody := []byte(`{"sl":{"af":"Afrikaans","en":"English"},"tl":{"af":"Afrikaans","en":"English"}}`)
	getLanguage(t, server, url, expectedBody)
}

func getTranslateScript(t *testing.T, serverURL string, path string, expectedStatus int) {
	url := fmt.Sprintf("%v%v", serverURL, path)
	resp, err := http.Get(url)
	assert.Nil(t, err)
	assert.Equal(t, expectedStatus, resp.StatusCode, http.StatusOK)
	assert.Equal(t, controller.GoogleTranslateServer, resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestGetTranslateScript(t *testing.T) {
	server := httptest.NewServer(handler)
	defer server.Close()
	mockGoogleServer := httptest.NewServer(mockGoogleServerHandler)
	defer mockGoogleServer.Close()
	controller.GoogleTranslateServerProxy = mockGoogleServer.URL

	path := "/translate_a/element.js"
	getTranslateScript(t, server.URL, path, http.StatusOK)
	path = "/translate_static/js/element/main.js"
	getTranslateScript(t, server.URL, path, http.StatusOK)
	path = "/element/TE_20181015_01/e/js/element/element_main.js"
	getTranslateScript(t, server.URL, path, http.StatusOK)
}

func testRedirect(t *testing.T, path string, redirectHost string) {
	server := httptest.NewServer(handler)
	defer server.Close()
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := fmt.Sprintf("%v%v", server.URL, path)
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)

	redirectURL := fmt.Sprintf("%v%v", redirectHost, path)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, redirectURL, resp.Header.Get("Location"))
}

func TestGetGoogleTranslateResource(t *testing.T) {
	path := "/translate_static/css/translateelement.css"
	testRedirect(t, path, controller.GoogleTranslateServerProxy)
}

func TestGetGStaticResource(t *testing.T) {
	path := "/images/branding/product/1x/translate_24dp.png"
	testRedirect(t, path, controller.GStaticServerProxy)
	path = "/images/branding/product/2x/translate_24dp.png"
	testRedirect(t, path, controller.GStaticServerProxy)
}
