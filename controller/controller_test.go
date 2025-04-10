package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/brave/go-translate/language"
)

func TestNewLnxEndpointConfiguration(t *testing.T) {
	lists := []language.GoogleLanguageList{
		language.GoogleLanguageList{
			Sl: map[string]string{"en":"English", "es":"Spanish", "it": "Italian"},
			Tl: map[string]string{"en":"English", "es":"Spanish", "it": "Italian"},
		},
		language.GoogleLanguageList{
			Sl: map[string]string{"en":"English", "es":"Spanish", "de": "Deutsch"},
			Tl: map[string]string{"en":"English", "es":"Spanish", "de": "Deutsch"},
		},
	}
	endpoints := []string{"endpoint1.com", "endpoint2.com"}
	weights := []float64{0.5, 0.5}

	conf, err := NewLnxEndpointConfiguration(endpoints, weights, lists)
	assert.NoError(t, err)

	assert.Equal(t, endpoints, conf.Endpoints)
	assert.Equal(t, weights, conf.DefaultWeights)

	assert.NotNil(t, conf.LanguagePairList)
	assert.NotEmpty(t, conf.LanguagePairWeights)

	assert.Equal(t, conf.LanguagePairWeights["en"]["it"], map[string]float64{"endpoint1.com": 0.5})
	assert.Equal(t, conf.LanguagePairWeights["it"]["en"], map[string]float64{"endpoint1.com": 0.5})
	assert.Equal(t, conf.LanguagePairWeights["en"]["de"], map[string]float64{"endpoint2.com": 0.5})
	assert.Equal(t, conf.LanguagePairWeights["de"]["en"], map[string]float64{"endpoint2.com": 0.5})
	assert.Equal(t, conf.LanguagePairWeights["en"]["es"], map[string]float64{"endpoint1.com": 0.5, "endpoint2.com": 0.5})
	assert.Equal(t, conf.LanguagePairWeights["es"]["en"], map[string]float64{"endpoint1.com": 0.5, "endpoint2.com": 0.5})
}

func TestLnxEndpointConfiguration_GetEndpoint(t *testing.T) {
	lists := []language.GoogleLanguageList{
		language.GoogleLanguageList{
			Sl: map[string]string{"en":"English", "es":"Spanish", "it": "Italian"},
			Tl: map[string]string{"en":"English", "es":"Spanish", "it": "Italian"},
		},
		language.GoogleLanguageList{
			Sl: map[string]string{"en":"English", "es":"Spanish", "de": "Deutsch"},
			Tl: map[string]string{"en":"English", "es":"Spanish", "de": "Deutsch"},
		},
	}
	endpoints := []string{"endpoint1.com", "endpoint2.com"}
	weights := []float64{0.5, 0.5}

	conf, err := NewLnxEndpointConfiguration(endpoints, weights, lists)
	assert.NoError(t, err)

	t.Run("random selection", func(t *testing.T) {
		from := "en"
		to := "es"
		countOne := 0
		countTwo := 0

		for i := 0; i<2000; i++ {
			got := conf.GetEndpoint(from, to)
			switch got {
			case "endpoint1.com":
				countOne++
			case "endpoint2.com":
				countTwo++
			}
		}
		assert.Less(t, 900, countTwo)
		assert.Greater(t, 1100, countTwo)
		assert.Less(t, 900, countOne)
		assert.Greater(t, 1100, countOne)
	})

	t.Run("first endpoint", func(t *testing.T) {
		from := "en"
		to := "it"
		expected := "endpoint1.com"

		for i := 0; i<100; i++ {
			got := conf.GetEndpoint(from, to)
			assert.Equal(t, expected, got)
		}
	})

	t.Run("second endpoint", func(t *testing.T) {
		from := "en"
		to := "de"
		expected := "endpoint2.com"

		for i := 0; i<100; i++ {
			got := conf.GetEndpoint(from, to)
			assert.Equal(t, expected, got)
		}
	})
}
