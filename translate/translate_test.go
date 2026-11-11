package translate

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLanguageParamsNormalizesModernCodes(t *testing.T) {
	r := httptest.NewRequest("POST", "/translate_a/t?sl=he&tl=en", nil)
	from, to, err := GetLanguageParams(r)
	assert.Equal(t, nil, err)
	assert.Equal(t, "iw", from)
	assert.Equal(t, "en", to)

	r = httptest.NewRequest("POST", "/translate_a/t?sl=en&tl=fil", nil)
	from, to, err = GetLanguageParams(r)
	assert.Equal(t, nil, err)
	assert.Equal(t, "en", from)
	assert.Equal(t, "tl", to)
}
