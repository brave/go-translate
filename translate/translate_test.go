package translate

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func init() {
	err := os.Setenv("MS_TRANSLATE_API_KEY", "DUMMY_KEY")
	if err != nil {
		panic("Test init failed when setting API Key env.")
	}
}

func TestToMicrosoftRequest(t *testing.T) {
	// auto-detect source language
	u := "https://example.com/translate/t?anno=3&client=te_lib&format=html&v=1.0&key=DUMMY_KEY&logld=vTE_20181015_01&sl=auto&tl=en&sp=nmt&tc=1&sr=1&tk=568026.932804&mode=1"
	msBody := []byte(`[{"Text":"guten Abend"},{"Text":"Hallo Welt"}]`)
	googleBody := []byte(`q=guten%20Abend&q=Hallo%20Welt`)

	googleReq, err := http.NewRequest("POST", u, bytes.NewBuffer(googleBody))
	googleReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.Nil(t, err)

	r, isAuto, err := ToMicrosoftRequest(googleReq, "https://api.cognitive.microsofttranslator.com")
	assert.Nil(t, err)
	assert.Equal(t, true, isAuto)
	body, err := ioutil.ReadAll(r.Body)
	assert.Nil(t, err)
	assert.Equal(t, msBody, body)
	assert.Equal(t, "https", r.URL.Scheme)
	assert.Equal(t, "api.cognitive.microsofttranslator.com", r.URL.Host)
	assert.Equal(t, "api-version=3.0&textType=html&to=en", r.URL.RawQuery)
	assert.Equal(t, "application/json", r.Header["Content-Type"][0])
	assert.Equal(t, "DUMMY_KEY", r.Header["Ocp-Apim-Subscription-Key"][0])

	// known source language
	u = "https://example.com/translate/t?anno=3&client=te_lib&format=html&v=1.0&key=DUMMY_KEY&logld=vTE_20181015_01&sl=de&tl=en&sp=nmt&tc=1&sr=1&tk=568026.932804&mode=1"
	googleReq, err = http.NewRequest("POST", u, bytes.NewBuffer(googleBody))
	googleReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.Nil(t, err)
	assert.Nil(t, err)
	r, isAuto, err = ToMicrosoftRequest(googleReq, "https://api.cognitive.microsofttranslator.com")
	assert.Nil(t, err)
	assert.Equal(t, false, isAuto)
	body, err = ioutil.ReadAll(r.Body)
	assert.Nil(t, err)
	assert.Equal(t, msBody, body)
	assert.Equal(t, "https", r.URL.Scheme)
	assert.Equal(t, "api.cognitive.microsofttranslator.com", r.URL.Host)
	assert.Equal(t, "api-version=3.0&from=de&textType=html&to=en", r.URL.RawQuery)
	assert.Equal(t, "application/json", r.Header["Content-Type"][0])
	assert.Equal(t, "DUMMY_KEY", r.Header["Ocp-Apim-Subscription-Key"][0])
}

func TestToGoogleResponseBody(t *testing.T) {
	// auto-detect source language
	isAuto := true
	msBody := []byte(`[{"detectedLanguage":{"language":"de","score":1.0},"translations":[{"text":"Hallo","to":"en"}]},{"detectedLanguage":{"language":"de","score":1.0},"translations":[{"text":"Welt","to":"en"}]}]`)
	googleBody := []byte(`[["Hallo","de"],["Welt","de"]]`)
	body, err := ToGoogleResponseBody(msBody, isAuto)
	assert.Nil(t, err)
	assert.Equal(t, googleBody, body)

	// known source language
	isAuto = false
	msBody = []byte(`[{"translations":[{"text":"Hallo","to":"en"}]},{"translations":[{"text":"Welt","to":"en"}]}]`)
	googleBody = []byte(`["Hallo","Welt"]`)
	body, err = ToGoogleResponseBody(msBody, isAuto)
	assert.Nil(t, err)
	assert.Equal(t, googleBody, body)
}
