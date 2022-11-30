package language

import (
	"encoding/json"
	"errors"
)

var ChromiumLanguageList = []string{
	"af",
	"am",
	"ar",
	"az",
	"be",
	"bg",
	"bn",
	"bs",
	"ca",
	"ceb",
	"co",
	"cs",
	"cy",
	"da",
	"de",
	"el",
	"en",
	"eo",
	"es",
	"et",
	"eu",
	"fa",
	"fi",
	"fr",
	"fy",
	"ga",
	"gd",
	"gl",
	"gu",
	"ha",
	"haw",
	"hi",
	"hmn",
	"hr",
	"ht",
	"hu",
	"hy",
	"id",
	"ig",
	"is",
	"it",
	"iw",
	"ja",
	"jw",
	"ka",
	"kk",
	"km",
	"kn",
	"ko",
	"ku",
	"ky",
	"la",
	"lb",
	"lo",
	"lt",
	"lv",
	"mg",
	"mi",
	"mk",
	"ml",
	"mn",
	"mr",
	"ms",
	"mt",
	"my",
	"ne",
	"nl",
	"no",
	"ny",
	"or",
	"pa",
	"pl",
	"ps",
	"pt",
	"ro",
	"ru",
	"rw",
	"sd",
	"si",
	"sk",
	"sl",
	"sm",
	"sn",
	"so",
	"sq",
	"sr",
	"st",
	"su",
	"sv",
	"sw",
	"ta",
	"te",
	"tg",
	"th",
	"tk",
	"tl",
	"tr",
	"tt",
	"ug",
	"uk",
	"ur",
	"uz",
	"vi",
	"xh",
	"yi",
	"yo",
	"zh-CN",
	"zh-TW",
	"zu",
}

var GoogleToLnxLangExclusions = map[string]string{
	"zh-CN": "zh-Hans",
	"zh-TW": "zh-Hant",
	"sr":    "sr-Cyrl",
	"iw":    "he",
	"jw":    "jv",
}

func MakeLnxToGoogleLangMapping() map[string]string {
	result := make(map[string]string)
	for _, glang := range ChromiumLanguageList {
		if lnxLang, ok := GoogleToLnxLangExclusions[glang]; ok {
			result[lnxLang] = glang
		} else {
			result[glang] = glang
		}
	}
	return result
}

func MakeGoogleToLnxLangMapping() map[string]string {
	result := make(map[string]string)
	for _, glang := range ChromiumLanguageList {
		if lnxLang, ok := GoogleToLnxLangExclusions[glang]; ok {
			result[glang] = lnxLang
		} else {
			result[glang] = glang
		}
	}
	return result
}

var googleToLnxLangMapping map[string]string
var lnxToGoogleLangMapping map[string]string

func ToGoogleLanguageCode(lnxLangCode string) (string, error) {
	if val, ok := lnxToGoogleLangMapping[lnxLangCode]; ok {
		return val, nil
	}
	return "", errors.New("Invalid language code " + lnxLangCode)
}

func ToLnxLanguageCode(gLangCode string) (string, error) {
	if len(googleToLnxLangMapping) == 0 {
		return "", errors.New("uninitialized map")
	}

	if val, ok := googleToLnxLangMapping[gLangCode]; ok {
		return val, nil
	}
	return "", errors.New("Invalid language code " + gLangCode)
}

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
		val, err := ToGoogleLanguageCode(lang.IsoCode)
		if err != nil {
			continue
		}
		googleLangList.Sl[val] = lang.Name
		googleLangList.Tl[val] = lang.Name
	}
	return json.Marshal(googleLangList)
}

func init() {
	googleToLnxLangMapping = MakeGoogleToLnxLangMapping()
	lnxToGoogleLangMapping = MakeLnxToGoogleLangMapping()
}
