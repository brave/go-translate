package language

import (
	"encoding/json"
	"errors"
)

// Language represents the format of language in Lingvanex language list
// Note that nativeName and dir are not used in this package since Google
// doesn't have them.
type Language struct {
	IsoCode string `json:"code_alpha_1"`
	Name    string `json:"codeName"`
}

// GoogleLanguageList represents the JSON format for Google language list
// Example:
//	{
//		"sl":{"auto":"Detect language","af":"Afrikaans"},
//		"tl":{"af": "Afrikaans", "sq": "Albanian"}
//		"al":{}
//	}
// Note that al is not used in this package since Lingvanex don't have it.
type GoogleLanguageList struct {
	Sl map[string]string `json:"sl"`
	Tl map[string]string `json:"tl"`
}

// ToGoogleLanguageList unmarshal a Lnx language list and marshal a corresponding
// google language list and return it.
func ToGoogleLanguageList(body []byte) ([]byte, error) {
	var lnxLangList []Language
	err := json.Unmarshal(body, &lnxLangList)
	if err != nil {
		return nil, err
	}

	// Marshal the language list using google's format
	var googleLangList GoogleLanguageList
	googleLangList.Sl = make(map[string]string, len(lnxLangList))
	googleLangList.Tl = make(map[string]string, len(lnxLangList))
	for _, lang := range lnxLangList {
		googleLangList.Sl[lang.IsoCode] = lang.Name
		googleLangList.Tl[lang.IsoCode] = lang.Name
	}
	return json.Marshal(googleLangList)
}

func ToLnxLanguageCode(gLangCode string) (string, error) {
	// Also deal with duplicate mappings for ar, en, es, fr, pt
	supportedLangList := []byte(`[{"code_alpha_1":"en","codeName":"English","rtl":false},{"code_alpha_1":"fr","codeName":"French","rtl":false},{"code_alpha_1":"pl","codeName":"Polish","rtl":false},{"code_alpha_1":"zh-Hans","codeName":"Chinese (Simplified)","rtl":false},{"code_alpha_1":"vi","codeName":"Vietnamese","rtl":false},{"code_alpha_1":"pt","codeName":"Portuguese","rtl":false},{"code_alpha_1":"nl","codeName":"Dutch","rtl":false},{"code_alpha_1":"de","codeName":"German","rtl":false},{"code_alpha_1":"ja","codeName":"Japanese","rtl":false},{"code_alpha_1":"es","codeName":"Spanish","rtl":false},{"code_alpha_1":"ro","codeName":"Romanian","rtl":false},{"code_alpha_1":"tr","codeName":"Turkish","rtl":false},{"code_alpha_1":"it","codeName":"Italian","rtl":false},{"code_alpha_1":"hi","codeName":"Hindi","rtl":false},{"code_alpha_1":"ru","codeName":"Russian","rtl":false}]`)
	supportedLangList := []byte(`[{"code_alpha_1":"auto","codeName":"Auto","rtl":false},{"code_alpha_1":"en","codeName":"English","rtl":false},{"code_alpha_1":"fr","codeName":"French","rtl":false},{"code_alpha_1":"pl","codeName":"Polish","rtl":false},{"code_alpha_1":"zh-Hans","codeName":"Chinese (Simplified)","rtl":false},{"code_alpha_1":"vi","codeName":"Vietnamese","rtl":false},{"code_alpha_1":"pt","codeName":"Portuguese","rtl":false},{"code_alpha_1":"nl","codeName":"Dutch","rtl":false},{"code_alpha_1":"de","codeName":"German","rtl":false},{"code_alpha_1":"ja","codeName":"Japanese","rtl":false},{"code_alpha_1":"es","codeName":"Spanish","rtl":false},{"code_alpha_1":"ro","codeName":"Romanian","rtl":false},{"code_alpha_1":"tr","codeName":"Turkish","rtl":false},{"code_alpha_1":"it","codeName":"Italian","rtl":false},{"code_alpha_1":"hi","codeName":"Hindi","rtl":false},{"code_alpha_1":"ru","codeName":"Russian","rtl":false}]`)
	langs := make([]Language, 0)
	err := json.Unmarshal(supportedLangList, &langs)
	if err != nil {
		return "", err
	}

	for _, v := range langs {
		if v.IsoCode == gLangCode {
			return v.IsoCode, nil
		}
	}

	return "", errors.New("No matching Lnx Language code")
}
