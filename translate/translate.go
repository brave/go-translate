package translate

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// RequestBody represents JSON format of Microsoft requests.
type RequestBody struct {
	Text string `json:"Text"`
}

// MicrosoftResponseBody represents JSON format of Microsoft response bodies.
// Translations's size is limited to 1 since multiple translations is not
// compatible with Google.
// Format with auto-detect source language:
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
// Format without auto-detect source language:
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
type MicrosoftResponseBody []struct {
	DetectedLang struct {
		Language string `json:"language"`
	} `json:"detectedLanguage,omitempty"`
	Translations [1]struct {
		Text string `json:"text"`
	} `json:"translations"`
}

const (
	translateEndpoint = "/translate?api-version=3.0"
)

// ToMicrosoftRequest parses the input Google format translate request and
// return a corresponding Microsoft format request.
func ToMicrosoftRequest(r *http.Request, serverURL string) (*http.Request, bool, error) {
	msURL := serverURL + translateEndpoint
	// Parse google format query parameters
	slVals := r.URL.Query()["sl"]
	if len(slVals) != 1 {
		return nil, false, errors.New("invalid query parameter format: There should be one sl parameter")
	}
	tlVals := r.URL.Query()["tl"]
	if len(tlVals) != 1 {
		return nil, false, errors.New("invalid query parameter format: There should be one tl parameter")
	}
	from := slVals[0]
	to := tlVals[0]

	// Set MS format query parameters
	u, err := url.Parse(msURL)
	if err != nil {
		return nil, false, err
	}
	q := u.Query()
	if from != "auto" {
		q.Add("from", from)
	}
	q.Add("to", to)
	q.Add("textType", "html")
	u.RawQuery = q.Encode()

	// Convert Google format request body into MS format request body
	err = r.ParseForm()
	if err != nil {
		return nil, false, err
	}
	qVals := r.PostForm["q"]

	// Set the request body
	reqBody := make([]RequestBody, len(qVals))
	for i, q := range qVals {
		reqBody[i] = RequestBody{q}
	}

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
	req.Header.Add("Ocp-Apim-Subscription-Key", os.Getenv("MS_TRANSLATE_API_KEY"))
	return req, from == "auto", nil
}

// ToGoogleResponseBody parses the input Microsoft response and return the JSON
// response body in Google format.
func ToGoogleResponseBody(body []byte, isAuto bool) ([]byte, error) {
	// Parse MS response body
	var msResp MicrosoftResponseBody
	err := json.Unmarshal(body, &msResp)
	if err != nil {
		return nil, err
	}

	// Source language is specified, google result format: ["aa", "bb", ...]
	if !isAuto {
		body := make([]string, len(msResp))
		for i, responseBody := range msResp {
			body[i] = responseBody.Translations[0].Text
		}
		return json.Marshal(body)
	}

	// Source language is auto detected,
	// google result format: [["aa", "from_len_a"], ["bb", "from_len_b"], ...]
	bodyAuto := make([][2]string, len(msResp))
	for i, responseBody := range msResp {
		bodyAuto[i][0] = responseBody.Translations[0].Text
		bodyAuto[i][1] = responseBody.DetectedLang.Language
	}
	return json.Marshal(bodyAuto)
}
