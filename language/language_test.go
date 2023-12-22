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
