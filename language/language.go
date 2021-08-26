package language

import (
	"encoding/json"
)

// GoogleLanguageList represents the JSON format for Google language list
// Example:
//	{
//		"sl":{"auto": "Detect language", "af": "Afrikaans"},
//		"tl":{"af": "Afrikaans", "sq": "Albanian"}
//		"al":{}
//	}
type GoogleLanguageList struct {
	Sl map[string]string `json:"sl"`
	Tl map[string]string `json:"tl"`
	Al map[string]string `json:"al"`
}

func GetBergamotLanguageList() ([]byte, error) {
	sourceLanguages := map[string]string {
		"auto": "Detect language",
		"en": "English",
		"es": "Spanish",
		"et": "Estonian",
		"it": "Italian",
		"pt": "Portuguese",
		"ru": "Russian",
	}
	targetLanguages := map[string]string {
		"de": "German",
		"en": "English",
		"es": "Spanish",
		"et": "Estonian",
		"ru": "Russian",
	}
	var googleLangList GoogleLanguageList
	googleLangList.Sl = sourceLanguages
	googleLangList.Tl = targetLanguages
	googleLangList.Al = map[string]string {}
	return json.Marshal(googleLangList)
}
