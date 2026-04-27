// Package utils — country↔currency mapping shared by tenant-service
// (which derives a tenant's primary currency from its address country)
// and billing-service (which validates the derived value at cache-write
// time). Lives in common so the same lookup runs everywhere; previously
// each service had its own copy and they drifted.
package utils

import "strings"

// countryCurrencyByAlpha2 maps ISO 3166-1 alpha-2 country codes to
// ISO 4217 currency codes. Keep this list in sync with the active rows
// in the billing-service `currencies` table — adding a country here
// without registering its currency in money.minorUnitsByCurrency means
// arithmetic falls back to a 100 minor-unit factor, which is wrong for
// JPY (1) and KWD/BHD/OMR/TND (1000).
var countryCurrencyByAlpha2 = map[string]string{
	"IN": "INR", // India
	"US": "USD", // United States
	"GB": "GBP", // United Kingdom
	"AE": "AED", // United Arab Emirates
	"SA": "SAR", // Saudi Arabia
	"BD": "BDT", // Bangladesh
	"PK": "PKR", // Pakistan
	"LK": "LKR", // Sri Lanka
	"NP": "NPR", // Nepal
	"MY": "MYR", // Malaysia
	"SG": "SGD", // Singapore
	"AU": "AUD", // Australia
	"CA": "CAD", // Canada
	"JP": "JPY", // Japan
	"DE": "EUR", // Germany
	"FR": "EUR", // France
	"IT": "EUR", // Italy
	"ES": "EUR", // Spain
	"NL": "EUR", // Netherlands
	"BR": "BRL", // Brazil
	"MX": "MXN", // Mexico
	"ZA": "ZAR", // South Africa
	"KE": "KES", // Kenya
	"NG": "NGN", // Nigeria
	"EG": "EGP", // Egypt
	"TH": "THB", // Thailand
	"ID": "IDR", // Indonesia
	"PH": "PHP", // Philippines
	"VN": "VND", // Vietnam
	"KR": "KRW", // South Korea
	"CN": "CNY", // China
	"HK": "HKD", // Hong Kong
	"TW": "TWD", // Taiwan
	"NZ": "NZD", // New Zealand
	"SE": "SEK", // Sweden
	"NO": "NOK", // Norway
	"DK": "DKK", // Denmark
	"CH": "CHF", // Switzerland
	"KW": "KWD", // Kuwait
	"BH": "BHD", // Bahrain
	"OM": "OMR", // Oman
	"QA": "QAR", // Qatar
	"TN": "TND", // Tunisia
}

// countryCurrencyByName maps lowercase common country names to ISO
// 4217 codes. The address forms in admin-dashboard store country as
// a free-text label ("India", "United States"), not an ISO code, so a
// name-based path is required — relying on the alpha-2 path alone
// silently fell back to USD for every Indian tenant.
//
// Keys are lowercased at lookup time. Add aliases (e.g. "uk" alongside
// "united kingdom") here rather than in callers.
var countryCurrencyByName = map[string]string{
	"india":                "INR",
	"united states":        "USD",
	"united states of america": "USD",
	"usa":                  "USD",
	"united kingdom":       "GBP",
	"uk":                   "GBP",
	"great britain":        "GBP",
	"united arab emirates": "AED",
	"uae":                  "AED",
	"saudi arabia":         "SAR",
	"bangladesh":           "BDT",
	"pakistan":             "PKR",
	"sri lanka":            "LKR",
	"nepal":                "NPR",
	"malaysia":             "MYR",
	"singapore":            "SGD",
	"australia":            "AUD",
	"canada":               "CAD",
	"japan":                "JPY",
	"germany":              "EUR",
	"france":               "EUR",
	"italy":                "EUR",
	"spain":                "EUR",
	"netherlands":          "EUR",
	"brazil":               "BRL",
	"mexico":               "MXN",
	"south africa":         "ZAR",
	"kenya":                "KES",
	"nigeria":              "NGN",
	"egypt":                "EGP",
	"thailand":             "THB",
	"indonesia":            "IDR",
	"philippines":          "PHP",
	"vietnam":              "VND",
	"south korea":          "KRW",
	"korea":                "KRW",
	"china":                "CNY",
	"hong kong":            "HKD",
	"taiwan":               "TWD",
	"new zealand":          "NZD",
	"sweden":               "SEK",
	"norway":               "NOK",
	"denmark":              "DKK",
	"switzerland":          "CHF",
	"kuwait":               "KWD",
	"bahrain":              "BHD",
	"oman":                 "OMR",
	"qatar":                "QAR",
	"tunisia":              "TND",
}

// CurrencyFromCountry resolves a country to its ISO-4217 currency.
// Accepts either an alpha-2 ISO code ("IN") or a common country name
// ("India", "United States", "uk"); matching is case-insensitive and
// trims surrounding whitespace.
//
// Falls back to "INR" — not "USD" — because the dataset is
// India-dominant and an unrecognised country today almost certainly
// means a free-text typo on an Indian address. Audit the resulting
// TenantCache.Currency in production logs for the warning emitted by
// the resolver.
func CurrencyFromCountry(country string) string {
	c := strings.TrimSpace(country)
	if c == "" {
		return "INR"
	}
	if len(c) == 2 {
		if v, ok := countryCurrencyByAlpha2[strings.ToUpper(c)]; ok {
			return v
		}
	}
	if v, ok := countryCurrencyByName[strings.ToLower(c)]; ok {
		return v
	}
	return "INR"
}

// IsKnownCountry reports whether the country (alpha-2 or name) maps
// to a currency. Use at the schema/API boundary to flag operator
// typos before they propagate into TenantCache.
func IsKnownCountry(country string) bool {
	c := strings.TrimSpace(country)
	if c == "" {
		return false
	}
	if len(c) == 2 {
		if _, ok := countryCurrencyByAlpha2[strings.ToUpper(c)]; ok {
			return true
		}
	}
	_, ok := countryCurrencyByName[strings.ToLower(c)]
	return ok
}
