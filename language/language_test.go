package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSourceResponseBody(t *testing.T) {
	targetResBody := []byte(`{"translation":{"af":{"name":"Afrikaans","nativeName":"Afrikaans","dir":"ltr"},"en":{"name":"English","nativeName":"English","dir":"ltr"}}}`)
	convertedResBody, err := ToSourceResponseBody(targetResBody)
	expectedResBody := []byte(`{"sl":{"af":"Afrikaans","en":"English"},"tl":{"af":"Afrikaans","en":"English"}}`)

	assert.Nil(t, err)
	assert.Equal(t, expectedResBody, convertedResBody)
}

func TestIsSupportedLanguage(t *testing.T) {
	assert.Equal(t, true, IsSupportedLanguage("de"))
	assert.Equal(t, false, IsSupportedLanguage("ab"))
}

func TestShouldAutoDetectSourceLanguage(t *testing.T) {
	assert.Equal(t, true, ShouldAutoDetectSourceLanguage("auto"))
	assert.Equal(t, false, ShouldAutoDetectSourceLanguage("en"))
}
