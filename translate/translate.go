package translate

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/brave/go-translate/language"
)

type TranslationServiceRequestBody struct {
	From          string   `json:"source,omitempty"`
	To            string   `json:"target"`
	Data          []string `json:"q"`
	TranslateMode string   `json:"translateMode"`
}

type TranslationServiceResponseBody struct {
	SourceText     []string `json:"sourceText"`
	TranslatedText []string `json:"translatedText"`
}

func BuildTargetRequestBody(reqFromSource *http.Request) ([]byte, error) {
	// Parse source language code
	slParameterValues := reqFromSource.URL.Query()["sl"]
	if len(slParameterValues) != 1 {
		return nil, errors.New("There should be exactly one sl parameter value")
	}
	sourceLanguageCode := slParameterValues[0]

	if !language.IsSupportedLanguage(sourceLanguageCode) && !language.ShouldAutoDetectSourceLanguage(sourceLanguageCode) {
		return nil, errors.New("Source language code not supported")
	}

	// Parse target language code
	tlParameterValues := reqFromSource.URL.Query()["tl"]
	if len(tlParameterValues) != 1 {
		return nil, errors.New("There should be exactly one tl parameter")
	}
	targetLanguageCode := tlParameterValues[0]

	if !language.IsSupportedLanguage(targetLanguageCode) {
		return nil, errors.New("Target language code not supported")
	}

	// Parse request body
	err := reqFromSource.ParseForm()
	if err != nil {
		return nil, err
	}
	qParameterValues := reqFromSource.PostForm["q"]

	var requestBody TranslationServiceRequestBody

	// Only set `source` if source language is specified. Lingvanex attempts
	// auto source language detection if `source` is missing
	if !language.ShouldAutoDetectSourceLanguage(sourceLanguageCode) {
		requestBody.From = sourceLanguageCode
	}
	requestBody.To = targetLanguageCode
	requestBody.TranslateMode = "html"
	requestBody.Data = qParameterValues

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	return body, err
}

func ToResourceResponseBody(body []byte) ([]byte, error) {
	var resBody TranslationServiceResponseBody
	err := json.Unmarshal(body, &resBody)
	if err != nil {
		return nil, err
	}

	return json.Marshal(resBody.TranslatedText)
}
