package language

import (
	"encoding/json"
)

type Language struct {
	IsoCode string `json:"code_alpha_1"`
	Name    string `json:"codeName"`
}

type SourceLanguageList struct {
	Sl map[string]string `json:"sl"`
	Tl map[string]string `json:"tl"`
}

func ToSourceResponseBody(body []byte) ([]byte, error) {
	var targetLangList []Language
	err := json.Unmarshal(body, &targetLangList)
	if err != nil {
		return nil, err
	}

	var sourceLangList SourceLanguageList
	sourceLangList.Sl = make(map[string]string, len(targetLangList))
	sourceLangList.Tl = make(map[string]string, len(targetLangList))
	for _, lang := range targetLangList {
		sourceLangList.Sl[lang.IsoCode] = lang.Name
		sourceLangList.Tl[lang.IsoCode] = lang.Name
	}
	return json.Marshal(sourceLangList)
}

func IsSupportedLanguage(languageCode string) bool {
	supportedLanguageCodes := []string{"en", "fr", "pl", "zh-Hans", "vi", "pt", "nl", "de", "ja", "es", "ro", "tr", "it", "hi", "ru"}

	for _, supportedLanguageCode := range supportedLanguageCodes {
		if supportedLanguageCode == languageCode {
			return true
		}
	}

	return false
}

func ShouldAutoDetectSourceLanguage(languageCode string) bool {
	if languageCode == "auto" {
		return true
	}

	return false
}
