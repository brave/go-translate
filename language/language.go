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
	supportedLangList := []byte(`[{"code_alpha_1":"auto","codeName":"Auto","rtl":false},{"code_alpha_1":"af","codeName":"Afrikaans","rtl":false},{"code_alpha_1":"am","codeName":"Amharic","rtl":false},{"code_alpha_1":"ar","codeName":"Arabic","rtl":false},{"code_alpha_1":"az","codeName":"Azerbaijani","rtl":false},{"code_alpha_1":"be","codeName":"Belarusian","rtl":false},{"code_alpha_1":"bg","codeName":"Bulgarian","rtl":false},{"code_alpha_1":"bn","codeName":"Bengali","rtl":false},{"code_alpha_1":"bs","codeName":"Bosnian","rtl":false},{"code_alpha_1":"ca","codeName":"Catalan","rtl":false},{"code_alpha_1":"ceb","codeName":"Cebuano","rtl":false},{"code_alpha_1":"co","codeName":"Corsican","rtl":false},{"code_alpha_1":"cs","codeName":"Czech","rtl":false},{"code_alpha_1":"cy","codeName":"Welsh","rtl":false},{"code_alpha_1":"da","codeName":"Danish","rtl":false},{"code_alpha_1":"de","codeName":"German","rtl":false},{"code_alpha_1":"el","codeName":"Greek","rtl":false},{"code_alpha_1":"en","codeName":"English","rtl":false},{"code_alpha_1":"eo","codeName":"Esperanto","rtl":false},{"code_alpha_1":"es","codeName":"Spanish","rtl":false},{"code_alpha_1":"et","codeName":"Estonian","rtl":false},{"code_alpha_1":"eu","codeName":"Basque","rtl":false},{"code_alpha_1":"fa","codeName":"Persian","rtl":false},{"code_alpha_1":"fi","codeName":"Finnish","rtl":false},{"code_alpha_1":"fr","codeName":"French","rtl":false},{"code_alpha_1":"fy","codeName":"Frisian","rtl":false},{"code_alpha_1":"ga","codeName":"Irish","rtl":false},{"code_alpha_1":"gd","codeName":"Scots Gaelic","rtl":false},{"code_alpha_1":"gl","codeName":"Galician","rtl":false},{"code_alpha_1":"gu","codeName":"Gujarati","rtl":false},{"code_alpha_1":"ha","codeName":"Hausa","rtl":false},{"code_alpha_1":"haw","codeName":"Hawaiian","rtl":false},{"code_alpha_1":"he","codeName":"Hebrew","rtl":false},{"code_alpha_1":"hi","codeName":"Hindi","rtl":false},{"code_alpha_1":"hmn","codeName":"Hmong","rtl":false},{"code_alpha_1":"hr","codeName":"Croatian","rtl":false},{"code_alpha_1":"ht","codeName":"Haitian Creole","rtl":false},{"code_alpha_1":"hu","codeName":"Hungarian","rtl":false},{"code_alpha_1":"hy","codeName":"Armenian","rtl":false},{"code_alpha_1":"id","codeName":"Indonesian","rtl":false},{"code_alpha_1":"ig","codeName":"Igbo","rtl":false},{"code_alpha_1":"is","codeName":"Icelandic","rtl":false},{"code_alpha_1":"it","codeName":"Italian","rtl":false},{"code_alpha_1":"ja","codeName":"Japanese","rtl":false},{"code_alpha_1":"jv","codeName":"Javanese","rtl":false},{"code_alpha_1":"ka","codeName":"Georgian","rtl":false},{"code_alpha_1":"kk","codeName":"Kazakh","rtl":false},{"code_alpha_1":"km","codeName":"Khmer","rtl":false},{"code_alpha_1":"kn","codeName":"Kannada","rtl":false},{"code_alpha_1":"ko","codeName":"Korean","rtl":false},{"code_alpha_1":"ku","codeName":"Kurdish (Kurmanji)","rtl":false},{"code_alpha_1":"ky","codeName":"Kyrgyz","rtl":false},{"code_alpha_1":"la","codeName":"Latin","rtl":false},{"code_alpha_1":"lb","codeName":"Luxembourgish","rtl":false},{"code_alpha_1":"lo","codeName":"Lao","rtl":false},{"code_alpha_1":"lt","codeName":"Lithuanian","rtl":false},{"code_alpha_1":"lv","codeName":"Latvian","rtl":false},{"code_alpha_1":"mg","codeName":"Malagasy","rtl":false},{"code_alpha_1":"mi","codeName":"Maori","rtl":false},{"code_alpha_1":"mk","codeName":"Macedonian","rtl":false},{"code_alpha_1":"ml","codeName":"Malayalam","rtl":false},{"code_alpha_1":"mn","codeName":"Mongolian","rtl":false},{"code_alpha_1":"mr","codeName":"Marathi","rtl":false},{"code_alpha_1":"ms","codeName":"Malay","rtl":false},{"code_alpha_1":"mt","codeName":"Maltese","rtl":false},{"code_alpha_1":"my","codeName":"Myanmar (Burmese)","rtl":false},{"code_alpha_1":"ne","codeName":"Nepali","rtl":false},{"code_alpha_1":"nl","codeName":"Dutch","rtl":false},{"code_alpha_1":"no","codeName":"Norwegian","rtl":false},{"code_alpha_1":"ny","codeName":"Chichewa","rtl":false},{"code_alpha_1":"or","codeName":"Odia","rtl":false},{"code_alpha_1":"pa","codeName":"Punjabi","rtl":false},{"code_alpha_1":"pl","codeName":"Polish","rtl":false},{"code_alpha_1":"ps","codeName":"Pashto","rtl":false},{"code_alpha_1":"pt","codeName":"Portuguese","rtl":false},{"code_alpha_1":"ro","codeName":"Romanian","rtl":false},{"code_alpha_1":"ru","codeName":"Russian","rtl":false},{"code_alpha_1":"rw","codeName":"Kinyarwanda","rtl":false},{"code_alpha_1":"sd","codeName":"Sindhi","rtl":false},{"code_alpha_1":"si","codeName":"Sinhala","rtl":false},{"code_alpha_1":"sk","codeName":"Slovak","rtl":false},{"code_alpha_1":"sl","codeName":"Slovenian","rtl":false},{"code_alpha_1":"sm","codeName":"Samoan","rtl":false},{"code_alpha_1":"sn","codeName":"Shona","rtl":false},{"code_alpha_1":"so","codeName":"Somali","rtl":false},{"code_alpha_1":"sq","codeName":"Albanian","rtl":false},{"code_alpha_1":"sr-Cyrl","codeName":"Serbian Cyrilic","rtl":false},{"code_alpha_1":"sr-Cyrl","codeName":"Serbian Cyrilic","rtl":false},{"code_alpha_1":"st","codeName":"Sesotho","rtl":false},{"code_alpha_1":"su","codeName":"Sundanese","rtl":false},{"code_alpha_1":"sv","codeName":"Swedish","rtl":false},{"code_alpha_1":"sw","codeName":"Swahili","rtl":false},{"code_alpha_1":"ta","codeName":"Tamil","rtl":false},{"code_alpha_1":"te","codeName":"Telugu","rtl":false},{"code_alpha_1":"tg","codeName":"Tajik","rtl":false},{"code_alpha_1":"th","codeName":"Thai","rtl":false},{"code_alpha_1":"tk","codeName":"Turkmen","rtl":false},{"code_alpha_1":"tl","codeName":"Tagalog","rtl":false},{"code_alpha_1":"tr","codeName":"Turkish","rtl":false},{"code_alpha_1":"tt","codeName":"Tatar","rtl":false},{"code_alpha_1":"ug","codeName":"Uyghur","rtl":false},{"code_alpha_1":"uk","codeName":"Ukrainian","rtl":false},{"code_alpha_1":"ur","codeName":"Urdu","rtl":false},{"code_alpha_1":"uz","codeName":"Uzbek","rtl":false},{"code_alpha_1":"vi","codeName":"Vietnamese","rtl":false},{"code_alpha_1":"xh","codeName":"Xhosa","rtl":false},{"code_alpha_1":"yi","codeName":"Yiddish","rtl":false},{"code_alpha_1":"yo","codeName":"Yoruba","rtl":false},{"code_alpha_1":"zh-CN","codeName":"Chinese (Simplified)","rtl":false},{"code_alpha_1":"zh-Hans","codeName":"Chinese (Simplified)","rtl":false},{"code_alpha_1":"zh-Hant","codeName":"Chinese (Traditional)","rtl":false},{"code_alpha_1":"zh-TW","codeName":"Chinese (Traditional)","rtl":false},{"code_alpha_1":"zu","codeName":"Zulu","rtl":false}]`)
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
