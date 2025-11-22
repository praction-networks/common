package constants

// LanguageCode represents a language code (ISO 639-1/2/3)
type LanguageCode string

// LanguageInfo contains language code, name, and native name
type LanguageInfo struct {
	Code       LanguageCode `json:"code"`
	Name       string       `json:"name"`
	NativeName string       `json:"nativeName"`
}

// All supported language codes (WhatsApp supported languages)
const (
	LanguageCodeAfrikaans          LanguageCode = "af"
	LanguageCodeAlbanian           LanguageCode = "sq"
	LanguageCodeArabic             LanguageCode = "ar"
	LanguageCodeAzerbaijani        LanguageCode = "az"
	LanguageCodeBengali            LanguageCode = "bn"
	LanguageCodeBulgarian          LanguageCode = "bg"
	LanguageCodeCatalan            LanguageCode = "ca"
	LanguageCodeChineseSimplified  LanguageCode = "zh_CN"
	LanguageCodeChineseHongKong    LanguageCode = "zh_HK"
	LanguageCodeChineseTraditional LanguageCode = "zh_TW"
	LanguageCodeCroatian           LanguageCode = "hr"
	LanguageCodeCzech              LanguageCode = "cs"
	LanguageCodeDanish             LanguageCode = "da"
	LanguageCodeDutch              LanguageCode = "nl"
	LanguageCodeEnglish            LanguageCode = "en"
	LanguageCodeEnglishUK          LanguageCode = "en_GB"
	LanguageCodeEnglishUS          LanguageCode = "en_US"
	LanguageCodeEstonian           LanguageCode = "et"
	LanguageCodeFilipino           LanguageCode = "fil"
	LanguageCodeFinnish            LanguageCode = "fi"
	LanguageCodeFrench             LanguageCode = "fr"
	LanguageCodeFrenchCanada       LanguageCode = "fr_CA"
	LanguageCodeGeorgian           LanguageCode = "ka"
	LanguageCodeGerman             LanguageCode = "de"
	LanguageCodeGreek              LanguageCode = "el"
	LanguageCodeGujarati           LanguageCode = "gu"
	LanguageCodeHausa              LanguageCode = "ha"
	LanguageCodeHebrew             LanguageCode = "he"
	LanguageCodeHindi              LanguageCode = "hi"
	LanguageCodeHungarian          LanguageCode = "hu"
	LanguageCodeIndonesian         LanguageCode = "id"
	LanguageCodeIrish              LanguageCode = "ga"
	LanguageCodeItalian            LanguageCode = "it"
	LanguageCodeJapanese           LanguageCode = "ja"
	LanguageCodeKannada            LanguageCode = "kn"
	LanguageCodeKazakh             LanguageCode = "kk"
	LanguageCodeKinyarwanda        LanguageCode = "rw_RW"
	LanguageCodeKorean             LanguageCode = "ko"
	LanguageCodeKyrgyz             LanguageCode = "ky_KG"
	LanguageCodeLao                LanguageCode = "lo"
	LanguageCodeLatvian            LanguageCode = "lv"
	LanguageCodeLithuanian         LanguageCode = "lt"
	LanguageCodeMacedonian         LanguageCode = "mk"
	LanguageCodeMalay              LanguageCode = "ms"
	LanguageCodeMalayalam          LanguageCode = "ml"
	LanguageCodeMarathi            LanguageCode = "mr"
	LanguageCodeNorwegian          LanguageCode = "nb"
	LanguageCodePashto             LanguageCode = "ps"
	LanguageCodePersian            LanguageCode = "fa"
	LanguageCodePolish             LanguageCode = "pl"
	LanguageCodePortugueseBrazil   LanguageCode = "pt_BR"
	LanguageCodePortuguesePortugal LanguageCode = "pt_PT"
	LanguageCodePunjabi            LanguageCode = "pa"
	LanguageCodeRomanian           LanguageCode = "ro"
	LanguageCodeRussian            LanguageCode = "ru"
	LanguageCodeSerbian            LanguageCode = "sr"
	LanguageCodeSlovak             LanguageCode = "sk"
	LanguageCodeSlovenian          LanguageCode = "sl"
	LanguageCodeSpanish            LanguageCode = "es"
	LanguageCodeSpanishArgentina   LanguageCode = "es_AR"
	LanguageCodeSpanishSpain       LanguageCode = "es_ES"
	LanguageCodeSpanishMexico      LanguageCode = "es_MX"
	LanguageCodeSwahili            LanguageCode = "sw"
	LanguageCodeSwedish            LanguageCode = "sv"
	LanguageCodeTamil              LanguageCode = "ta"
	LanguageCodeTelugu             LanguageCode = "te"
	LanguageCodeThai               LanguageCode = "th"
	LanguageCodeTurkish            LanguageCode = "tr"
	LanguageCodeUkrainian          LanguageCode = "uk"
	LanguageCodeUrdu               LanguageCode = "ur"
	LanguageCodeUzbek              LanguageCode = "uz"
	LanguageCodeVietnamese         LanguageCode = "vi"
	LanguageCodeZulu               LanguageCode = "zu"
)

// String returns the string representation of the language code
func (c LanguageCode) String() string {
	return string(c)
}

// IsValid checks if the language code is valid
func (c LanguageCode) IsValid() bool {
	_, exists := GetAllLanguageCodes()[c]
	return exists
}

// GetInfo returns the LanguageInfo for a language code
func (c LanguageCode) GetInfo() *LanguageInfo {
	allInfo := GetAllLanguageInfo()
	if info, exists := allInfo[c]; exists {
		return &info
	}
	return nil
}

// GetName returns the English name for a language code
func (c LanguageCode) GetName() string {
	info := c.GetInfo()
	if info != nil {
		return info.Name
	}
	return ""
}

// GetNativeName returns the native name for a language code
func (c LanguageCode) GetNativeName() string {
	info := c.GetInfo()
	if info != nil {
		return info.NativeName
	}
	return ""
}

// GetAllLanguageCodes returns a map of all valid language codes
func GetAllLanguageCodes() map[LanguageCode]bool {
	return map[LanguageCode]bool{
		LanguageCodeAfrikaans:          true,
		LanguageCodeAlbanian:           true,
		LanguageCodeArabic:             true,
		LanguageCodeAzerbaijani:        true,
		LanguageCodeBengali:            true,
		LanguageCodeBulgarian:          true,
		LanguageCodeCatalan:            true,
		LanguageCodeChineseSimplified:  true,
		LanguageCodeChineseHongKong:    true,
		LanguageCodeChineseTraditional: true,
		LanguageCodeCroatian:           true,
		LanguageCodeCzech:              true,
		LanguageCodeDanish:             true,
		LanguageCodeDutch:              true,
		LanguageCodeEnglish:            true,
		LanguageCodeEnglishUK:          true,
		LanguageCodeEnglishUS:          true,
		LanguageCodeEstonian:           true,
		LanguageCodeFilipino:           true,
		LanguageCodeFinnish:            true,
		LanguageCodeFrench:             true,
		LanguageCodeFrenchCanada:       true,
		LanguageCodeGeorgian:           true,
		LanguageCodeGerman:             true,
		LanguageCodeGreek:              true,
		LanguageCodeGujarati:           true,
		LanguageCodeHausa:              true,
		LanguageCodeHebrew:             true,
		LanguageCodeHindi:              true,
		LanguageCodeHungarian:          true,
		LanguageCodeIndonesian:         true,
		LanguageCodeIrish:              true,
		LanguageCodeItalian:            true,
		LanguageCodeJapanese:           true,
		LanguageCodeKannada:            true,
		LanguageCodeKazakh:             true,
		LanguageCodeKinyarwanda:        true,
		LanguageCodeKorean:             true,
		LanguageCodeKyrgyz:             true,
		LanguageCodeLao:                true,
		LanguageCodeLatvian:            true,
		LanguageCodeLithuanian:         true,
		LanguageCodeMacedonian:         true,
		LanguageCodeMalay:              true,
		LanguageCodeMalayalam:          true,
		LanguageCodeMarathi:            true,
		LanguageCodeNorwegian:          true,
		LanguageCodePashto:             true,
		LanguageCodePersian:            true,
		LanguageCodePolish:             true,
		LanguageCodePortugueseBrazil:   true,
		LanguageCodePortuguesePortugal: true,
		LanguageCodePunjabi:            true,
		LanguageCodeRomanian:           true,
		LanguageCodeRussian:            true,
		LanguageCodeSerbian:            true,
		LanguageCodeSlovak:             true,
		LanguageCodeSlovenian:          true,
		LanguageCodeSpanish:            true,
		LanguageCodeSpanishArgentina:   true,
		LanguageCodeSpanishSpain:       true,
		LanguageCodeSpanishMexico:      true,
		LanguageCodeSwahili:            true,
		LanguageCodeSwedish:            true,
		LanguageCodeTamil:              true,
		LanguageCodeTelugu:             true,
		LanguageCodeThai:               true,
		LanguageCodeTurkish:            true,
		LanguageCodeUkrainian:          true,
		LanguageCodeUrdu:               true,
		LanguageCodeUzbek:              true,
		LanguageCodeVietnamese:         true,
		LanguageCodeZulu:               true,
	}
}

// GetAllLanguageInfo returns a map of all language codes to their LanguageInfo
func GetAllLanguageInfo() map[LanguageCode]LanguageInfo {
	return map[LanguageCode]LanguageInfo{
		LanguageCodeAfrikaans:          {Code: LanguageCodeAfrikaans, Name: "Afrikaans", NativeName: "Afrikaans"},
		LanguageCodeAlbanian:           {Code: LanguageCodeAlbanian, Name: "Albanian", NativeName: "Shqip"},
		LanguageCodeArabic:             {Code: LanguageCodeArabic, Name: "Arabic", NativeName: "العربية"},
		LanguageCodeAzerbaijani:        {Code: LanguageCodeAzerbaijani, Name: "Azerbaijani", NativeName: "Azərbaycan"},
		LanguageCodeBengali:            {Code: LanguageCodeBengali, Name: "Bengali", NativeName: "বাংলা"},
		LanguageCodeBulgarian:          {Code: LanguageCodeBulgarian, Name: "Bulgarian", NativeName: "Български"},
		LanguageCodeCatalan:            {Code: LanguageCodeCatalan, Name: "Catalan", NativeName: "Català"},
		LanguageCodeChineseSimplified:  {Code: LanguageCodeChineseSimplified, Name: "Chinese (Simplified)", NativeName: "简体中文"},
		LanguageCodeChineseHongKong:    {Code: LanguageCodeChineseHongKong, Name: "Chinese (Hong Kong)", NativeName: "繁體中文（香港）"},
		LanguageCodeChineseTraditional: {Code: LanguageCodeChineseTraditional, Name: "Chinese (Traditional)", NativeName: "繁體中文"},
		LanguageCodeCroatian:           {Code: LanguageCodeCroatian, Name: "Croatian", NativeName: "Hrvatski"},
		LanguageCodeCzech:              {Code: LanguageCodeCzech, Name: "Czech", NativeName: "Čeština"},
		LanguageCodeDanish:             {Code: LanguageCodeDanish, Name: "Danish", NativeName: "Dansk"},
		LanguageCodeDutch:              {Code: LanguageCodeDutch, Name: "Dutch", NativeName: "Nederlands"},
		LanguageCodeEnglish:            {Code: LanguageCodeEnglish, Name: "English", NativeName: "English"},
		LanguageCodeEnglishUK:          {Code: LanguageCodeEnglishUK, Name: "English (United Kingdom)", NativeName: "English (UK)"},
		LanguageCodeEnglishUS:          {Code: LanguageCodeEnglishUS, Name: "English (United States)", NativeName: "English (US)"},
		LanguageCodeEstonian:           {Code: LanguageCodeEstonian, Name: "Estonian", NativeName: "Eesti"},
		LanguageCodeFilipino:           {Code: LanguageCodeFilipino, Name: "Filipino", NativeName: "Filipino"},
		LanguageCodeFinnish:            {Code: LanguageCodeFinnish, Name: "Finnish", NativeName: "Suomi"},
		LanguageCodeFrench:             {Code: LanguageCodeFrench, Name: "French", NativeName: "Français"},
		LanguageCodeFrenchCanada:       {Code: LanguageCodeFrenchCanada, Name: "French (Canada)", NativeName: "Français (Canada)"},
		LanguageCodeGeorgian:           {Code: LanguageCodeGeorgian, Name: "Georgian", NativeName: "ქართული"},
		LanguageCodeGerman:             {Code: LanguageCodeGerman, Name: "German", NativeName: "Deutsch"},
		LanguageCodeGreek:              {Code: LanguageCodeGreek, Name: "Greek", NativeName: "Ελληνικά"},
		LanguageCodeGujarati:           {Code: LanguageCodeGujarati, Name: "Gujarati", NativeName: "ગુજરાતી"},
		LanguageCodeHausa:              {Code: LanguageCodeHausa, Name: "Hausa", NativeName: "Hausa"},
		LanguageCodeHebrew:             {Code: LanguageCodeHebrew, Name: "Hebrew", NativeName: "עברית"},
		LanguageCodeHindi:              {Code: LanguageCodeHindi, Name: "Hindi", NativeName: "हिन्दी"},
		LanguageCodeHungarian:          {Code: LanguageCodeHungarian, Name: "Hungarian", NativeName: "Magyar"},
		LanguageCodeIndonesian:         {Code: LanguageCodeIndonesian, Name: "Indonesian", NativeName: "Bahasa Indonesia"},
		LanguageCodeIrish:              {Code: LanguageCodeIrish, Name: "Irish", NativeName: "Gaeilge"},
		LanguageCodeItalian:            {Code: LanguageCodeItalian, Name: "Italian", NativeName: "Italiano"},
		LanguageCodeJapanese:           {Code: LanguageCodeJapanese, Name: "Japanese", NativeName: "日本語"},
		LanguageCodeKannada:            {Code: LanguageCodeKannada, Name: "Kannada", NativeName: "ಕನ್ನಡ"},
		LanguageCodeKazakh:             {Code: LanguageCodeKazakh, Name: "Kazakh", NativeName: "Қазақ"},
		LanguageCodeKinyarwanda:        {Code: LanguageCodeKinyarwanda, Name: "Kinyarwanda", NativeName: "Ikinyarwanda"},
		LanguageCodeKorean:             {Code: LanguageCodeKorean, Name: "Korean", NativeName: "한국어"},
		LanguageCodeKyrgyz:             {Code: LanguageCodeKyrgyz, Name: "Kyrgyz", NativeName: "Кыргызча"},
		LanguageCodeLao:                {Code: LanguageCodeLao, Name: "Lao", NativeName: "ລາວ"},
		LanguageCodeLatvian:            {Code: LanguageCodeLatvian, Name: "Latvian", NativeName: "Latviešu"},
		LanguageCodeLithuanian:         {Code: LanguageCodeLithuanian, Name: "Lithuanian", NativeName: "Lietuvių"},
		LanguageCodeMacedonian:         {Code: LanguageCodeMacedonian, Name: "Macedonian", NativeName: "Македонски"},
		LanguageCodeMalay:              {Code: LanguageCodeMalay, Name: "Malay", NativeName: "Bahasa Melayu"},
		LanguageCodeMalayalam:          {Code: LanguageCodeMalayalam, Name: "Malayalam", NativeName: "മലയാളം"},
		LanguageCodeMarathi:            {Code: LanguageCodeMarathi, Name: "Marathi", NativeName: "मराठी"},
		LanguageCodeNorwegian:          {Code: LanguageCodeNorwegian, Name: "Norwegian", NativeName: "Norsk"},
		LanguageCodePashto:             {Code: LanguageCodePashto, Name: "Pashto", NativeName: "پښتو"},
		LanguageCodePersian:            {Code: LanguageCodePersian, Name: "Persian", NativeName: "فارسی"},
		LanguageCodePolish:             {Code: LanguageCodePolish, Name: "Polish", NativeName: "Polski"},
		LanguageCodePortugueseBrazil:   {Code: LanguageCodePortugueseBrazil, Name: "Portuguese (Brazil)", NativeName: "Português (Brasil)"},
		LanguageCodePortuguesePortugal: {Code: LanguageCodePortuguesePortugal, Name: "Portuguese (Portugal)", NativeName: "Português (Portugal)"},
		LanguageCodePunjabi:            {Code: LanguageCodePunjabi, Name: "Punjabi", NativeName: "ਪੰਜਾਬੀ"},
		LanguageCodeRomanian:           {Code: LanguageCodeRomanian, Name: "Romanian", NativeName: "Română"},
		LanguageCodeRussian:            {Code: LanguageCodeRussian, Name: "Russian", NativeName: "Русский"},
		LanguageCodeSerbian:            {Code: LanguageCodeSerbian, Name: "Serbian", NativeName: "Српски"},
		LanguageCodeSlovak:             {Code: LanguageCodeSlovak, Name: "Slovak", NativeName: "Slovenčina"},
		LanguageCodeSlovenian:          {Code: LanguageCodeSlovenian, Name: "Slovenian", NativeName: "Slovenščina"},
		LanguageCodeSpanish:            {Code: LanguageCodeSpanish, Name: "Spanish", NativeName: "Español"},
		LanguageCodeSpanishArgentina:   {Code: LanguageCodeSpanishArgentina, Name: "Spanish (Argentina)", NativeName: "Español (Argentina)"},
		LanguageCodeSpanishSpain:       {Code: LanguageCodeSpanishSpain, Name: "Spanish (Spain)", NativeName: "Español (España)"},
		LanguageCodeSpanishMexico:      {Code: LanguageCodeSpanishMexico, Name: "Spanish (Mexico)", NativeName: "Español (México)"},
		LanguageCodeSwahili:            {Code: LanguageCodeSwahili, Name: "Swahili", NativeName: "Kiswahili"},
		LanguageCodeSwedish:            {Code: LanguageCodeSwedish, Name: "Swedish", NativeName: "Svenska"},
		LanguageCodeTamil:              {Code: LanguageCodeTamil, Name: "Tamil", NativeName: "தமிழ்"},
		LanguageCodeTelugu:             {Code: LanguageCodeTelugu, Name: "Telugu", NativeName: "తెలుగు"},
		LanguageCodeThai:               {Code: LanguageCodeThai, Name: "Thai", NativeName: "ไทย"},
		LanguageCodeTurkish:            {Code: LanguageCodeTurkish, Name: "Turkish", NativeName: "Türkçe"},
		LanguageCodeUkrainian:          {Code: LanguageCodeUkrainian, Name: "Ukrainian", NativeName: "Українська"},
		LanguageCodeUrdu:               {Code: LanguageCodeUrdu, Name: "Urdu", NativeName: "اردو"},
		LanguageCodeUzbek:              {Code: LanguageCodeUzbek, Name: "Uzbek", NativeName: "Oʻzbek"},
		LanguageCodeVietnamese:         {Code: LanguageCodeVietnamese, Name: "Vietnamese", NativeName: "Tiếng Việt"},
		LanguageCodeZulu:               {Code: LanguageCodeZulu, Name: "Zulu", NativeName: "isiZulu"},
	}
}

// GetAllLanguageCodeStrings returns a slice of all language code strings
func GetAllLanguageCodeStrings() []string {
	codes := GetAllLanguageCodes()
	result := make([]string, 0, len(codes))
	for code := range codes {
		result = append(result, code.String())
	}
	return result
}

// GetAllLanguageInfoList returns a slice of all LanguageInfo
func GetAllLanguageInfoList() []LanguageInfo {
	allInfo := GetAllLanguageInfo()
	result := make([]LanguageInfo, 0, len(allInfo))
	for _, info := range allInfo {
		result = append(result, info)
	}
	return result
}
