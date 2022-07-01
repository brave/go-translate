package translate

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTargetRequestBody(t *testing.T) {
	// Explicit source language
	sourceRequestUrl := "https://translate.brave.com/translate_a/t?anno=3&client=te_lib&format=html&v=1.0&key=qjVKcxtUybh8WpKNoQ7EbgbkJTMu7omjDHKk%3DVrPApb8PwJyPE9eqchxedTsMEWg&logld=vTE_20220615&sl=de&tl=en&tc=7&sr=1&tk=82273.411120&mode=1"
	sourceRequestBody := []byte(`q=Hallo&q=Welt`)
	sourceRequest, _ := http.NewRequest("POST", sourceRequestUrl, bytes.NewBuffer(sourceRequestBody))
	sourceRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	body, err := BuildTargetRequestBody(sourceRequest)

	expectedBody := []byte(`{"source":"de","target":"en","q":["Hallo","Welt"],"translateMode":"html"}`)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, body)

	// Auto detect source language
	sourceRequestUrl = "https://translate.brave.com/translate_a/t?anno=3&client=te_lib&format=html&v=1.0&key=qjVKcxtUybh8WpKNoQ7EbgbkJTMu7omjDHKk%3DVrPApb8PwJyPE9eqchxedTsMEWg&logld=vTE_20220615&sl=auto&tl=en&tc=7&sr=1&tk=82273.411120&mode=1"
	sourceRequestBody = []byte(`q=Hallo&q=Welt`)
	sourceRequest, _ = http.NewRequest("POST", sourceRequestUrl, bytes.NewBuffer(sourceRequestBody))
	sourceRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	body, err = BuildTargetRequestBody(sourceRequest)
	expectedBody = []byte(`{"target":"en","q":["Hallo","Welt"],"translateMode":"html"}`)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, body)
}

// func TestToGoogleResponseBody(t *testing.T) {
// 	// auto-detect source language
// 	msBody := []byte(`[{"detectedLanguage":{"language":"de","score":1.0},"translations":[{"text":"Hallo","to":"en"}]},{"detectedLanguage":{"language":"de","score":1.0},"translations":[{"text":"Welt","to":"en"}]}]`)
// 	googleBody := []byte(`[["Hallo","de"],["Welt","de"]]`)
// 	body, err := ToGoogleResponseBody(msBody)
// 	assert.Nil(t, err)
// 	assert.Equal(t, googleBody, body)

// 	// known source language
// 	msBody = []byte(`[{"translations":[{"text":"Hallo","to":"en"}]},{"translations":[{"text":"Welt","to":"en"}]}]`)
// 	googleBody = []byte(`["Hallo","Welt"]`)
// 	body, err = ToGoogleResponseBody(msBody)
// 	assert.Nil(t, err)
// 	assert.Equal(t, googleBody, body)
// }
