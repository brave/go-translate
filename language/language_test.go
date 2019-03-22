package language

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToGoogleLanguageList(t *testing.T) {
	msList := []byte(`{"translation":{"af":{"name":"Afrikaans","nativeName":"Afrikaans","dir":"ltr"},"en":{"name":"English","nativeName":"English","dir":"ltr"}}}`)
	googleList := []byte(`{"sl":{"af":"Afrikaans","en":"English"},"tl":{"af":"Afrikaans","en":"English"}}`)
	list, err := ToGoogleLanguageList(msList)
	assert.Equal(t, nil, err)
	assert.Equal(t, googleList, list)
}
