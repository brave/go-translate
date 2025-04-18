// Package translate provides functionality for translating text between languages
// and converting between Google Translate and Lingvanex API formats.
package translate

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/brave/go-translate/language"
)

var (
	charsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "translate_processed_chars_total",
		Help: "The total number of characters processed for translation",
	})
	reqsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "translate_processed_requests_total",
		Help: "The total number of requests processed for translation by language",
	},
		[]string{"to_lang", "from_lang"},
	)
)

// RequestBody represents JSON format of Lingvanex requests.
type RequestBody struct {
	From          string   `json:"source,omitempty"`
	To            string   `json:"target"`
	Data          []string `json:"q"`
	TranslateMode string   `json:"translateMode"`
}

// LingvanexResponseBody represents JSON format of Lingvanex response bodies.
// Translations's size is limited to 1 since multiple translations is not
// compatible with Google.
// Format with auto-detect source language:
//
//	[
//		{
//			"detectedLanguage": {"language": "de", "score": 1.0},
//			"translations": [{"text": "Hallo", "to": "en"}]
//		},
//		{
//			"detectedLanguage": {"language": "de", "score": 1.0},
//			"translations": [{"text": "Welt", "to": "en"}]
//		}
//	]
//
// Format without auto-detect source language:
//
//	[
//		{
//			"translations": [{"text": "Hallo", "to": "en"}]
//		},
//		{
//			"translations": [{"text": "Welt", "to": "en"}]
//		}
//	]
//
// score and to are not saved in this struct because we don't need them to
// convert to a google format response.
type LingvanexResponseBody struct {
	SourceText     []string `json:"sourceText"`
	TranslatedText []string `json:"translatedText"`
}

// GetLanguageParams extracts source and target language parameters from the request
func GetLanguageParams(r *http.Request) (string, string, error) {
	// Parse google format query parameters
	slVals := r.URL.Query()["sl"]
	if len(slVals) != 1 {
		return "", "", errors.New("invalid query parameter format: There should be one sl parameter")
	}
	tlVals := r.URL.Query()["tl"]
	if len(tlVals) != 1 {
		return "", "", errors.New("invalid query parameter format: There should be one tl parameter")
	}
	return slVals[0], tlVals[0], nil
}

// ToLingvanexRequest parses the input Google format translate request and
// return a corresponding Lingvanex format request.
func ToLingvanexRequest(r *http.Request, serverURL string) (*http.Request, bool, error) {
	lnxURL := serverURL

	from, to, err := GetLanguageParams(r)
	if err != nil {
		return nil, false, err
	}

	reqsProcessed.With(prometheus.Labels{
		"from_lang": from,
		"to_lang":   to,
	}).Inc()

	// Set Lnx format query parameters
	u, err := url.Parse(lnxURL)
	if err != nil {
		return nil, false, err
	}

	// Convert Google format request body into Lnx format request body
	err = r.ParseForm()
	if err != nil {
		return nil, false, err
	}
	qVals := r.PostForm["q"]

	for _, q := range qVals {
		charsProcessed.Add(float64(len(q)))
	}

	lnxTo, err := language.ToLnxLanguageCode(to)
	if err != nil {
		return nil, false, errors.New("No matching lnxTo language code:" + err.Error())
	}

	var reqBody RequestBody
	if from != "auto" {
		lnxFrom, err := language.ToLnxLanguageCode(from)
		if err != nil {
			return nil, false, errors.New("No matching lnxFrom language code:" + err.Error())
		}
		reqBody.From = lnxFrom
	}
	reqBody.To = lnxTo
	reqBody.TranslateMode = "html"
	reqBody.Data = qVals

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, false, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, false, err
	}

	// Set request headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	return req, from == "auto", nil
}

// ToGoogleResponseBody parses the input Lingvanex response and return the JSON
// response body in Google format.
func ToGoogleResponseBody(body []byte, _ bool) ([]byte, error) {
	// Parse Lnx response body
	var lnxResp LingvanexResponseBody
	err := json.Unmarshal(body, &lnxResp)
	if err != nil {
		return nil, err
	}

	return json.Marshal(lnxResp.TranslatedText)
}
