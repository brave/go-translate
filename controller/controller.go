// Package controller provides handlers and routing for the translation service API
package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/brave-intl/bat-go/libs/logging"
	"github.com/brave-intl/bat-go/libs/middleware"
	"github.com/brave/go-translate/language"
	"github.com/brave/go-translate/translate"
	"github.com/go-chi/chi/v5"
)

var (
	// LnxEndpoint stores the configuration for Lingvanex translation service endpoints
	LnxEndpoint   *LnxEndpointConfiguration
	// LnxAPIKey is the API key for accessing Lingvanex translation services
	LnxAPIKey     = os.Getenv("LNX_API_KEY")
	languagePath  = "/get-languages"
	translatePath = "/translate"
)

// LnxEndpointConfiguration describes a configuration of lingvanex endpoints, their supported
// languages and weights.
type LnxEndpointConfiguration struct {
	// A list of endpoint URLs.
	Endpoints []string
	// A list of default endpoint weights.
	DefaultWeights []float64
	// A GoogleLanguageList containing source language descriptions and target language descriptions.
	LanguagePairList language.GoogleLanguageList
	// A nested map of endpoint weights for a language pair.
	// The first key represents the source language, the second key represents the target language, the
	// third key represents the endpoint URL and the value the corresponding weight for that endpoint.
	LanguagePairWeights map[string]map[string]map[string]float64
}

// NewLnxEndpointConfiguration returns a new endpoint configuration based on a list of endpoints, weights and list of supported languages
func NewLnxEndpointConfiguration(endpoints []string, weights []float64, languageLists []language.GoogleLanguageList) (*LnxEndpointConfiguration, error) {
	if len(endpoints) != len(weights) || len(weights) != len(languageLists) {
		return nil, fmt.Errorf("number of endpoints must match number of weights and number of language lists")
	}

	conf := LnxEndpointConfiguration{
		Endpoints:           endpoints,
		DefaultWeights:      weights,
		LanguagePairList:    language.GoogleLanguageList{Sl: make(map[string]string), Tl: make(map[string]string)},
		LanguagePairWeights: make(map[string]map[string]map[string]float64),
	}

	for i, endpoint := range endpoints {
		// get the list of supported languages for the current endpoint
		list := languageLists[i]

		// iterate through the source languages the current endpoint supports
		for sl, sldesc := range list.Sl {
			// add the source language description to the merged language pair list
			conf.LanguagePairList.Sl[sl] = sldesc

			// check if the source language weight map already exists in the language pair weights
			if _, ok := conf.LanguagePairWeights[sl]; !ok {
				// if not, create a new weight map for the source language
				conf.LanguagePairWeights[sl] = make(map[string]map[string]float64)
			}

			for tl, tldesc := range list.Tl {
				// add the target language description to the merged language pair list
				conf.LanguagePairList.Tl[tl] = tldesc

				// check if the weight map for the source / target language pair already exists
				if _, ok := conf.LanguagePairWeights[sl][tl]; !ok {
					// if not, create a new weight map for it
					conf.LanguagePairWeights[sl][tl] = make(map[string]float64)
				}
				// set the default weight for the current endpoint for the source-target language pair
				conf.LanguagePairWeights[sl][tl][endpoint] = conf.DefaultWeights[i]
			}
		}
	}
	return &conf, nil
}

// GetEndpoint returns the endpoint which should be used based on the weights and languages supported.
func (c *LnxEndpointConfiguration) GetEndpoint(from, to string) string {
	// initialize total weight and incrementals.
	total := 0.0
	incrementals := []float64{}

	// retrieve the nested map of language pair weights.
	weights := c.LanguagePairWeights[from][to]

	// iterate through the Endpoints array, accumulating the total weight and storing the intermediate sums in incrementals.
	for _, endpoint := range c.Endpoints {
		total += weights[endpoint]
		incrementals = append(incrementals, total)
	}

	// generate a random number between 0 and total.
	r := rand.Float64() * total

	// find the endpoint with the smallest incremental weight greater than r.
	for i, incremental := range incrementals {
		if r < incremental {
			return c.Endpoints[i]
		}
	}
	// otherwise default to the first endpoint
	return c.Endpoints[0]
}

// TranslateRouter add routers for translate requests and translate script
// requests.
func TranslateRouter(ctx context.Context) (chi.Router, error) {
	r := chi.NewRouter()

	var weights []float64
	endpoints := strings.Split(os.Getenv("LNX_HOST"), ",")
	for _, weight := range strings.Split(os.Getenv("LNX_WEIGHTS"), ",") {
		if len(weight) > 0 {
			weight, err := strconv.ParseFloat(weight, 64)
			if err != nil {
				return r, fmt.Errorf("must pass at least one endpoint via LNX_HOST and one weight via LNX_WEIGHTS: %v", err)
			}
			weights = append(weights, weight)
		}
	}
	if len(endpoints) == 1 && len(weights) == 0 {
		weights = append(weights, 1)
	}

	var lists []language.GoogleLanguageList
	for _, endpoint := range endpoints {
		list, err := getLanguageList(ctx, endpoint)
		if err != nil {
			panic(err)
		}
		lists = append(lists, *list)
	}

	var err error
	LnxEndpoint, err = NewLnxEndpointConfiguration(endpoints, weights, lists)
	if err != nil {
		return r, fmt.Errorf("failed to setup endpoint configuration: %v", err)
	}

	r.Post("/translate_a/t", middleware.InstrumentHandler("Translate", http.HandlerFunc(Translate)).ServeHTTP)
	r.Get("/translate_a/l", middleware.InstrumentHandler("GetLanguageList", http.HandlerFunc(GetLanguageList)).ServeHTTP)

	r.Get("/static/v1/element.js", middleware.InstrumentHandler("ServeStaticFile", http.HandlerFunc(ServeStaticFile)).ServeHTTP)
	r.Get("/static/v1/js/element/main.js", middleware.InstrumentHandler("ServeStaticFile", http.HandlerFunc(ServeStaticFile)).ServeHTTP)
	r.Get("/static/v1/css/translateelement.css", middleware.InstrumentHandler("ServeStaticFile", http.HandlerFunc(ServeStaticFile)).ServeHTTP)

	return r, nil
}

// ServeStaticFile serves static files from the assets directory for the translation script
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

func getLanguageList(ctx context.Context, endpoint string) (*language.GoogleLanguageList, error) {
	logger := logging.FromContext(ctx)

	// Send a get language list request to Lnx
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+languagePath, nil)
	req.Header.Add("Authorization", "Bearer "+LnxAPIKey)

	if err != nil {
		return nil, fmt.Errorf("error creating Lnx request: %v", err)
	}

	client := getHTTPClient()
	lnxResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Lnx server: %v", err)
	}
	defer func() {
		err := lnxResp.Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Error closing response body stream")
		}
	}()

	// Convert to google format language list and write it back
	lnxBody, err := io.ReadAll(lnxResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Lnx response body: %v", err)
	}
	list, err := language.ToGoogleLanguageList(lnxBody)
	if err != nil {
		return nil, fmt.Errorf("error converting to google language list: %v", err)
	}

	return list, nil
}

// GetLanguageList send a request to Lingvanex server and convert the response
// into google format and reply back to the client.
func GetLanguageList(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())

	body, err := json.Marshal(LnxEndpoint.LanguagePairList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		logger.Error().Err(err).Msg("Error writing response body for translate requests")
	}
}

// handleBadRequestError writes a 400 Bad Request error response
func handleBadRequestError(w http.ResponseWriter, message string, err error) {
	http.Error(w, fmt.Sprintf("%s: %v", message, err), http.StatusBadRequest)
}

// handleInternalServerError writes a 500 Internal Server Error response
func handleInternalServerError(w http.ResponseWriter, message string, err error) {
	http.Error(w, fmt.Sprintf("%s: %v", message, err), http.StatusInternalServerError)
}

// handleNonOKResponse handles responses with non-OK status codes
func handleNonOKResponse(w http.ResponseWriter, lnxResp *http.Response) {
	w.WriteHeader(lnxResp.StatusCode)
	_, err := w.Write([]byte("LNX-ERROR:\n"))
	if err != nil {
		handleInternalServerError(w, "Error writing error message", err)
		return
	}
	_, err = io.Copy(w, lnxResp.Body)
	if err != nil {
		handleInternalServerError(w, "Error copying LnxEndpoint response body", err)
	}
}

// Translate converts a Google format translate request into a Lingvanex format
// one which will be send to the Lingvanex server, and write a Google format
// response back to the client.
func Translate(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())

	w.Header().Set("Access-Control-Allow-Origin", "*") // same as Google response

	to, from, err := translate.GetLanguageParams(r)
	if err != nil {
		handleBadRequestError(w, "error converting to LnxEndpoint request", err)
		return
	}

	endpoint := LnxEndpoint.GetEndpoint(from, to)
	req, isAuto, err := translate.ToLingvanexRequest(r, endpoint+translatePath)
	if err != nil {
		handleBadRequestError(w, "error converting to LnxEndpoint request", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+LnxAPIKey)

	// Send translate request to Lnx server
	client := getHTTPClient()
	lnxResp, err := client.Do(req)
	if err != nil {
		handleInternalServerError(w, "error sending request to LnxEndpoint", err)
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

	// Handle non-OK responses
	if lnxResp.StatusCode != http.StatusOK {
		handleNonOKResponse(w, lnxResp)
		return
	}

	// Set google format response body
	lnxBody, err := io.ReadAll(lnxResp.Body)
	if err != nil {
		handleInternalServerError(w, "Error reading LnxEndpoint response body", err)
		return
	}
	body, err := translate.ToGoogleResponseBody(lnxBody, isAuto)
	if err != nil {
		handleInternalServerError(w, "Error converting to google response body", err)
		return
	}
	w.WriteHeader(lnxResp.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		logger.Error().Err(err).Msg("Error writing response body for translate requests")
	}
}
