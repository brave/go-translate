package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToGoogleLanguageList(t *testing.T) {
	lnxList := []byte(`{"err":null,"result":[{"full_code":"af_ZA","code_alpha_1":"af","englishName":"Afrikaans","codeName":"Afrikaans","flagPath":"static/flags/afrikaans","testWordForSyntezis":"Hallo","rtl":"false","modes":[{"name":"Translation","value":true},{"name":"Translation document","value":true},{"name":"Image recognition","value":true},{"name":"Image object recognition","value":true},{"name":"Translate web site","value":true}]},{"full_code":"en_US","code_alpha_1":"en","englishName":"English (USA)","codeName":"English","flagPath":"static/flags/english_us","testWordForSyntezis":"test","rtl":"false","modes":[{"name":"Speech recognition","value":true},{"name":"Translation document","value":true},{"name":"Image recognition","value":true},{"name":"Image object recognition","value":true},{"name":"Translation","value":true},{"name":"Translate web site","value":true}]}]}`)
	googleList := []byte(`{"sl":{"af":"Afrikaans","en":"English"},"tl":{"af":"Afrikaans","en":"English"}}`)
	list, err := ToGoogleLanguageList(lnxList)
	assert.Equal(t, nil, err)
	assert.Equal(t, googleList, list)
}

func TestToLnxLanguageCode(t *testing.T) {
	gLangCode := "ar"
	lnxLangCode, err := ToLnxLanguageCode(gLangCode)
	expectedLangCode := "ar_AE"
	assert.Nil(t, err)
	assert.Equal(t, expectedLangCode, lnxLangCode)

	gLangCode = "en"
	lnxLangCode, err = ToLnxLanguageCode(gLangCode)
	expectedLangCode = "en_US"
	assert.Nil(t, err)
	assert.Equal(t, expectedLangCode, lnxLangCode)

	gLangCode = "es"
	lnxLangCode, err = ToLnxLanguageCode(gLangCode)
	expectedLangCode = "es_ES"
	assert.Nil(t, err)
	assert.Equal(t, expectedLangCode, lnxLangCode)

	gLangCode = "fr"
	lnxLangCode, err = ToLnxLanguageCode(gLangCode)
	expectedLangCode = "fr_FR"
	assert.Nil(t, err)
	assert.Equal(t, expectedLangCode, lnxLangCode)

	gLangCode = "pt"
	lnxLangCode, err = ToLnxLanguageCode(gLangCode)
	expectedLangCode = "pt_PT"
	assert.Nil(t, err)
	assert.Equal(t, expectedLangCode, lnxLangCode)

	gLangCode = "de"
	lnxLangCode, err = ToLnxLanguageCode(gLangCode)
	expectedLangCode = "de_DE"
	assert.Nil(t, err)
	assert.Equal(t, expectedLangCode, lnxLangCode)

	gLangCode = "asdf"
	lnxLangCode, err = ToLnxLanguageCode(gLangCode)
	expectedLangCode = ""
	assert.Error(t, err)
	assert.Equal(t, expectedLangCode, lnxLangCode)
}
