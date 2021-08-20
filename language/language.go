package language

import (
	"encoding/json"
)

// GoogleLanguageList represents the JSON format for Google language list
// Example:
//	{
//		"sl":{"auto":"Detect language","af":"Afrikaans"},
//		"tl":{"af": "Afrikaans", "sq": "Albanian"}
//		"al":{}
//	}
// Note that al is not used in this package since Microsoft don't have it.
type GoogleLanguageList struct {
	Sl map[string]string `json:"sl"`
	Tl map[string]string `json:"tl"`
}

func GetBergamotLanguageList() ([]byte, error) {
	sourceLanguages := map[string]string {
		"auto": "Detect automatically",
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
	return json.Marshal(googleLangList)
}
