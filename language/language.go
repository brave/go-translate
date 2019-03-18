package language

import (
	"encoding/json"
)

// Language represents the format of language in Microsoft language list
// Note that nativeName and dir are not used in this package since Google
// doesn't have them.
type Language struct {
	Name string `json:"name"`
}

// MicrosoftLanguageList represents the JSON format for Microsoft language list
// Example:
//	{
//		"translation": {
//			"af": {
//				"name": "Afrikaans",
//				"nativeName": "Afrikaans",
//				"dir": "ltr"
//			},
//			"en": {
//				"name": "English",
//				"nativeName": "English",
//				"dir": "ltr"
//			}
//		}
//	}
type MicrosoftLanguageList struct {
	Translation map[string]Language
}

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

// ToGoogleLanguageList unmarshal a MS language list and marshal a corresponding
// google language list and return it.
func ToGoogleLanguageList(body []byte) ([]byte, error) {
	// Unmarshal the MS language list from the response
	var msLangList MicrosoftLanguageList
	err := json.Unmarshal(body, &msLangList)
	if err != nil {
		return nil, err
	}

	// Marshal the language list using google's format
	var googleLangList GoogleLanguageList
	googleLangList.Sl = make(map[string]string, len(msLangList.Translation))
	googleLangList.Tl = make(map[string]string, len(msLangList.Translation))
	for k, lang := range msLangList.Translation {
		googleLangList.Sl[k] = lang.Name
		googleLangList.Tl[k] = lang.Name
	}
	return json.Marshal(googleLangList)
}
