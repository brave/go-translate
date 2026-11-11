package language

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the language code mapping.
// en should be used as it, but zh-Hans should be mapped to zh-CN
func TestToGoogleLanguageList(t *testing.T) {
	msList := []byte(`[{"code_alpha_1": "en", "codeName": "English", "rtl": false}, {"code_alpha_1": "zh-Hans", "codeName": "Chinese (Simplified)", "rtl": false}]`)
	googleList := []byte(`{"sl":{"en":"English","zh-CN":"Chinese (Simplified)"},"tl":{"en":"English","zh-CN":"Chinese (Simplified)"}}`)
	list, err := ToGoogleLanguageList(msList)
	assert.Equal(t, nil, err)

	var expected GoogleLanguageList
	err = json.Unmarshal(googleList, &expected)
	assert.Equal(t, nil, err)

	assert.Equal(t, expected, *list)
}

// Modern BCP 47 codes sent by newer Chromium versions should map to the legacy
// Google codes, and from there to the Lingvanex codes.
func TestNormalizeGoogleLanguageCode(t *testing.T) {
	cases := map[string]string{
		"he":    "iw",
		"jv":    "jw",
		"fil":   "tl",
		"iw":    "iw",
		"fr":    "fr",
		"zh-CN": "zh-CN",
		"auto":  "auto",
	}
	for input, expected := range cases {
		assert.Equal(t, expected, NormalizeGoogleLanguageCode(input), input)
	}

	lnx, err := ToLnxLanguageCode(NormalizeGoogleLanguageCode("he"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "he", lnx)
}
