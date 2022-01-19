package language

import (
	"encoding/json"
	"errors"
)

// Language represents the format of language in Microsoft language list
// Note that nativeName and dir are not used in this package since Google
// doesn't have them.
type Language struct {
	IsoCode string `json:"code_alpha_1"`
	LnxCode string `json:"lnx_code"`
	Name    string `json:"codeName"`
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
	Languages []Language `json:"result"`
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

	// TODO(Moritz Haller): Deal with duplicates for ar, en, es, fr, pt
	googleLangList.Sl = make(map[string]string, len(msLangList.Languages))
	googleLangList.Tl = make(map[string]string, len(msLangList.Languages))
	for _, lang := range msLangList.Languages {
		googleLangList.Sl[lang.IsoCode] = lang.Name
		googleLangList.Tl[lang.IsoCode] = lang.Name
	}
	return json.Marshal(googleLangList)
}

func ToLnxLanguageCode(gLangCode string) (string, error) {
	supportedLangList := []byte(`[{"code_alpha_1":"af","lnx_code":"af_ZA"},{"code_alpha_1":"am","lnx_code":"am_ET"},{"code_alpha_1":"ar","lnx_code":"ar_AE"},{"code_alpha_1":"az","lnx_code":"az_AZ"},{"code_alpha_1":"be","lnx_code":"be_BY"},{"code_alpha_1":"bg","lnx_code":"bg_BG"},{"code_alpha_1":"bn","lnx_code":"bn_BD"},{"code_alpha_1":"bs","lnx_code":"bs_BA"},{"code_alpha_1":"ca","lnx_code":"ca_ES"},{"code_alpha_1":"co","lnx_code":"co_FR"},{"code_alpha_1":"cs","lnx_code":"cs_CZ"},{"code_alpha_1":"cy","lnx_code":"cy_GB"},{"code_alpha_1":"da","lnx_code":"da_DK"},{"code_alpha_1":"de","lnx_code":"de_DE"},{"code_alpha_1":"el","lnx_code":"el_GR"},{"code_alpha_1":"en","lnx_code":"en_US"},{"code_alpha_1":"es","lnx_code":"es_ES"},{"code_alpha_1":"et","lnx_code":"et_EE"},{"code_alpha_1":"eu","lnx_code":"eu_ES"},{"code_alpha_1":"fa","lnx_code":"fa_IR"},{"code_alpha_1":"fi","lnx_code":"fi_FI"},{"code_alpha_1":"fr","lnx_code":"fr_FR"},{"code_alpha_1":"fy","lnx_code":"fy_NL"},{"code_alpha_1":"ga","lnx_code":"ga_IE"},{"code_alpha_1":"gd","lnx_code":"gd_GB"},{"code_alpha_1":"gl","lnx_code":"gl_ES"},{"code_alpha_1":"gu","lnx_code":"gu_IN"},{"code_alpha_1":"ha","lnx_code":"ha_NE"},{"code_alpha_1":"he","lnx_code":"he_IL"},{"code_alpha_1":"hi","lnx_code":"hi_IN"},{"code_alpha_1":"hr","lnx_code":"hr_HR"},{"code_alpha_1":"ht","lnx_code":"ht_HT"},{"code_alpha_1":"hu","lnx_code":"hu_HU"},{"code_alpha_1":"hy","lnx_code":"hy_AM"},{"code_alpha_1":"id","lnx_code":"id_ID"},{"code_alpha_1":"ig","lnx_code":"ig_NG"},{"code_alpha_1":"is","lnx_code":"is_IS"},{"code_alpha_1":"it","lnx_code":"it_IT"},{"code_alpha_1":"ja","lnx_code":"ja_JP"},{"code_alpha_1":"jv","lnx_code":"jv_ID"},{"code_alpha_1":"ka","lnx_code":"ka_GE"},{"code_alpha_1":"kk","lnx_code":"kk_KZ"},{"code_alpha_1":"km","lnx_code":"km_KH"},{"code_alpha_1":"kn","lnx_code":"kn_IN"},{"code_alpha_1":"ko","lnx_code":"ko_KR"},{"code_alpha_1":"ku","lnx_code":"ku_IR"},{"code_alpha_1":"ky","lnx_code":"ky_KG"},{"code_alpha_1":"lb","lnx_code":"lb_LU"},{"code_alpha_1":"lo","lnx_code":"lo_LA"},{"code_alpha_1":"lt","lnx_code":"lt_LT"},{"code_alpha_1":"lv","lnx_code":"lv_LV"},{"code_alpha_1":"mg","lnx_code":"mg_MG"},{"code_alpha_1":"mi","lnx_code":"mi_NZ"},{"code_alpha_1":"mk","lnx_code":"mk_MK"},{"code_alpha_1":"ml","lnx_code":"ml_IN"},{"code_alpha_1":"mn","lnx_code":"mn_MN"},{"code_alpha_1":"mr","lnx_code":"mr_IN"},{"code_alpha_1":"ms","lnx_code":"ms_MY"},{"code_alpha_1":"mt","lnx_code":"mt_MT"},{"code_alpha_1":"my","lnx_code":"my_MM"},{"code_alpha_1":"ne","lnx_code":"ne_NP"},{"code_alpha_1":"nl","lnx_code":"nl_NL"},{"code_alpha_1":"no","lnx_code":"no_NO"},{"code_alpha_1":"ny","lnx_code":"ny_MW"},{"code_alpha_1":"or","lnx_code":"or_OR"},{"code_alpha_1":"pa","lnx_code":"pa_PK"},{"code_alpha_1":"pl","lnx_code":"pl_PL"},{"code_alpha_1":"ps","lnx_code":"ps_AF"},{"code_alpha_1":"pt","lnx_code":"pt_PT"},{"code_alpha_1":"ro","lnx_code":"ro_RO"},{"code_alpha_1":"ru","lnx_code":"ru_RU"},{"code_alpha_1":"rw","lnx_code":"rw_RW"},{"code_alpha_1":"sd","lnx_code":"sd_PK"},{"code_alpha_1":"si","lnx_code":"si_LK"},{"code_alpha_1":"sk","lnx_code":"sk_SK"},{"code_alpha_1":"sl","lnx_code":"sl_SI"},{"code_alpha_1":"sm","lnx_code":"sm_WS"},{"code_alpha_1":"sn","lnx_code":"sn_ZW"},{"code_alpha_1":"so","lnx_code":"so_SO"},{"code_alpha_1":"sq","lnx_code":"sq_AL"},{"code_alpha_1":"st","lnx_code":"st_LS"},{"code_alpha_1":"su","lnx_code":"su_ID"},{"code_alpha_1":"sv","lnx_code":"sv_SE"},{"code_alpha_1":"sw","lnx_code":"sw_TZ"},{"code_alpha_1":"ta","lnx_code":"ta_IN"},{"code_alpha_1":"te","lnx_code":"te_IN"},{"code_alpha_1":"tg","lnx_code":"tg_TJ"},{"code_alpha_1":"th","lnx_code":"th_TH"},{"code_alpha_1":"tk","lnx_code":"tk_TK"},{"code_alpha_1":"tl","lnx_code":"tl_PH"},{"code_alpha_1":"tr","lnx_code":"tr_TR"},{"code_alpha_1":"tt","lnx_code":"tt_TT"},{"code_alpha_1":"ug","lnx_code":"ug_UG"},{"code_alpha_1":"uk","lnx_code":"uk_UA"},{"code_alpha_1":"ur","lnx_code":"ur_PK"},{"code_alpha_1":"uz","lnx_code":"uz_UZ"},{"code_alpha_1":"vi","lnx_code":"vi_VN"},{"code_alpha_1":"xh","lnx_code":"xh_ZA"},{"code_alpha_1":"yi","lnx_code":"yi_IL"},{"code_alpha_1":"yo","lnx_code":"yo_NG"},{"code_alpha_1":"zu","lnx_code":"zu_ZA"}]`)
	langs := make([]Language, 0)
	err := json.Unmarshal(supportedLangList, &langs)
	if err != nil {
		return "", err
	}

	for _, v := range langs {
		if v.IsoCode == gLangCode {
			return v.LnxCode, nil
		}
	}

	return "", errors.New("No matching Lnx Language code")
}
