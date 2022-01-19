package translate

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	err := os.Setenv("MS_TRANSLATE_API_KEY", "DUMMY_KEY")
	if err != nil {
		panic("Test init failed when setting API Key env.")
	}
}

func TestToMicrosoftRequest(t *testing.T) {
	// auto-detect source language
	// u := "https://example.com/translate?anno=3&client=te_lib&format=html&v=1.0&key=dummytoken&logld=vTE_20210503_00&sl=en&tl=de&tc=2&sr=1&tk=683550.820145&mode=1"
	// msBody := []byte(`["guten Abend","Hallo Welt"]`)
	// googleBody := []byte(`q=guten%20Abend&q=Hallo%20Welt`)

	// googleReq, err := http.NewRequest("POST", u, bytes.NewBuffer(googleBody))
	// googleReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// assert.Nil(t, err)

	// r, isAuto, err := ToMicrosoftRequest(googleReq, "https://api-b2b.backenster.com/b1/api/v3")
	// assert.Nil(t, err)
	// assert.Equal(t, true, isAuto)
	// body, err := ioutil.ReadAll(r.Body)
	// assert.Nil(t, err)
	// assert.Equal(t, msBody, body)
	// assert.Equal(t, "https", r.URL.Scheme)
	// assert.Equal(t, "https://api-b2b.backenster.com/b1/api/v3/", r.URL.Host)
	// assert.Equal(t, "translate", r.URL.RawQuery)
	// assert.Equal(t, "application/json", r.Header["Content-Type"][0])
	// assert.Equal(t, "DUMMY_KEY", r.Header["Ocp-Apim-Subscription-Key"][0])

	// known source language
	u := "https://example.com/translate?anno=3&client=te_lib&format=html&v=1.0&key=dummytoken&logld=vTE_20210503_00&sl=de&tl=en&tc=2&sr=1&tk=683550.820145&mode=1"
	msBody := []byte(`{"from":"de_DE","to":"en_US","text":["guten Abend","Hallo Welt"],"platform":"api","translateMode":"text"}`)
	googleBody := []byte(`q=guten%20Abend&q=Hallo%20Welt`)
	googleReq, err := http.NewRequest("POST", u, bytes.NewBuffer(googleBody))
	googleReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.Nil(t, err)
	assert.Nil(t, err)
	r, isAuto, err := ToMicrosoftRequest(googleReq, "https://api-b2b.backenster.com/b1/api/v3")
	assert.Nil(t, err)
	assert.Equal(t, false, isAuto)
	body, err := ioutil.ReadAll(r.Body)
	assert.Nil(t, err)
	assert.Equal(t, msBody, body)
	assert.Equal(t, "https", r.URL.Scheme)
	assert.Equal(t, "api-b2b.backenster.com", r.URL.Host)
	assert.Equal(t, "/b1/api/v3/translate", r.URL.Path)
	assert.Equal(t, "", r.URL.RawQuery)
	assert.Equal(t, "application/json", r.Header["Content-Type"][0])
}

func TestToGoogleResponseBody(t *testing.T) {
	// auto-detect source language
	// isAuto := true
	// msBody := []byte(`[{"detectedLanguage":{"language":"de","score":1.0},"translations":[{"text":"Hallo","to":"en"}]},{"detectedLanguage":{"language":"de","score":1.0},"translations":[{"text":"Welt","to":"en"}]}]`)
	// googleBody := []byte(`[["Hallo","de"],["Welt","de"]]`)
	// body, err := ToGoogleResponseBody(msBody, isAuto)
	// assert.Nil(t, err)
	// assert.Equal(t, googleBody, body)

	// known source language
	isAuto := false
	lnxBody := []byte(`{"err":null,"result":["Hallo","Welt"]}`)
	googleBody := []byte(`["Hallo","Welt"]`)
	body, err := ToGoogleResponseBody(lnxBody, isAuto)
	assert.Nil(t, err)
	assert.Equal(t, googleBody, body)
}
