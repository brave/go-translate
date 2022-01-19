package translate

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/brave/go-translate/language"
)

// RequestBody represents JSON format of Microsoft requests.
type RequestBody struct {
	From          string   `json:"from"`
	To            string   `json:"to"`
	Data          []string `json:"text"`
	Platform      string   `json:"platform"`
	TranslateMode string   `json:"translateMode"`
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

type LnxResponseBody struct {
	Error  string   `json:"err"`
	Result []string `json:"result"`
}

const (
	translateEndpoint = "/translate"
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

	// Convert Google format request body into MS format request body
	err = r.ParseForm()
	if err != nil {
		return nil, false, err
	}
	qVals := r.PostForm["q"]

	lnx_from, err := language.ToLnxLanguageCode(from)
	if err != nil {
		return nil, false, errors.New("No matching lnx_from language code")
	}

	lnx_to, err := language.ToLnxLanguageCode(to)
	if err != nil {
		return nil, false, errors.New("No matching lnx_to language code")
	}

	var reqBody RequestBody
	reqBody.From = lnx_from
	reqBody.To = lnx_to
	reqBody.Platform = "api"
	reqBody.TranslateMode = "text" // TODO(Moritz Haller): blocked by lingvanex, change to "html"
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

// ToGoogleResponseBody parses the input Microsoft response and return the JSON
// response body in Google format.
func ToGoogleResponseBody(body []byte, isAuto bool) ([]byte, error) {
	// Parse MS response body
	var msResp LnxResponseBody
	err := json.Unmarshal(body, &msResp)
	if err != nil {
		return nil, err
	}

	return json.Marshal(msResp.Result)

	// if !isAuto {
	// 	body := make([]string, len(msResp))
	// 	for i, responseBody := range msResp {
	// 		body[i] = responseBody.Translations[0].Text
	// 	}
	// 	return json.Marshal(body)
	// }

	// bodyAuto := make([][2]string, len(msResp))
	// for i, responseBody := range msResp {
	// 	bodyAuto[i][0] = responseBody.Translations[0].Text
	// 	bodyAuto[i][1] = responseBody.DetectedLang.Language
	// }
	// return json.Marshal(bodyAuto)
}
