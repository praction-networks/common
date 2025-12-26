package helpers

import (
	"fmt"
	"time"
)

type TimeProvider interface {
	Now() time.Time
	NowInTimezone(timezone string) (time.Time, error)
	ConvertToTimezone(t time.Time, timezone string) (time.Time, error)
}

type RealTimeProvider struct{}

func (r *RealTimeProvider) Now() time.Time {
	// Always return UTC to match MongoDB (dates stored as UTC milliseconds since epoch).
	return time.Now().UTC()
}

// NowInTimezone returns the current time in the specified timezone
// Timezone examples: "Asia/Kolkata", "America/New_York", "Europe/Madrid", "Asia/Bangkok", "Asia/Colombo"
func (r *RealTimeProvider) NowInTimezone(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	return time.Now().In(loc), nil
}

// ConvertToTimezone converts a given time to the specified timezone
// Timezone examples: "Asia/Kolkata", "America/New_York", "Europe/Madrid", "Asia/Bangkok", "Asia/Colombo"
func (r *RealTimeProvider) ConvertToTimezone(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	return t.In(loc), nil
}

// ProvideTimeProvider returns a new TimeProvider instance
func ProvideTimeProvider() TimeProvider {
	return &RealTimeProvider{}
}

// Common timezone constants for convenience
// Use IANA timezone database names for reliability
const (
	// UTC/GMT
	TimezoneUTC = "UTC"
	TimezoneGMT = "GMT"

	// Asia - South Asia
	TimezoneIndia      = "Asia/Kolkata"     // India (IST, UTC+5:30)
	TimezonePakistan   = "Asia/Karachi"     // Pakistan (PKT, UTC+5:00)
	TimezoneSriLanka   = "Asia/Colombo"     // Sri Lanka (IST, UTC+5:30)
	TimezoneBangladesh = "Asia/Dhaka"       // Bangladesh (BST, UTC+6:00)
	TimezoneNepal      = "Asia/Kathmandu"   // Nepal (NPT, UTC+5:45)
	TimezoneMaldives   = "Indian/Maldives"  // Maldives (MVT, UTC+5:00)

	// Asia - Southeast Asia
	TimezoneThailand   = "Asia/Bangkok"    // Thailand (ICT, UTC+7:00)
	TimezoneVietnam    = "Asia/Ho_Chi_Minh" // Vietnam (ICT, UTC+7:00)
	TimezoneSingapore  = "Asia/Singapore"  // Singapore (SGT, UTC+8:00)
	TimezoneMalaysia   = "Asia/Kuala_Lumpur" // Malaysia (MYT, UTC+8:00)
	TimezoneIndonesia  = "Asia/Jakarta"    // Indonesia - West (WIB, UTC+7:00)
	TimezonePhilippines = "Asia/Manila"    // Philippines (PST, UTC+8:00)
	TimezoneMyanmar    = "Asia/Yangon"     // Myanmar (MMT, UTC+6:30)

	// Asia - East Asia
	TimezoneChina      = "Asia/Shanghai"   // China (CST, UTC+8:00)
	TimezoneHongKong   = "Asia/Hong_Kong"  // Hong Kong (HKT, UTC+8:00)
	TimezoneTaiwan     = "Asia/Taipei"     // Taiwan (CST, UTC+8:00)
	TimezoneJapan      = "Asia/Tokyo"      // Japan (JST, UTC+9:00)
	TimezoneSouthKorea = "Asia/Seoul"      // South Korea (KST, UTC+9:00)

	// Asia - Middle East
	TimezoneDubai      = "Asia/Dubai"      // UAE (GST, UTC+4:00)
	TimezoneSaudiArabia = "Asia/Riyadh"    // Saudi Arabia (AST, UTC+3:00)
	TimezoneQatar      = "Asia/Qatar"      // Qatar (AST, UTC+3:00)
	TimezoneIsrael     = "Asia/Jerusalem"  // Israel (IST, UTC+2:00/+3:00)
	TimezoneTurkey     = "Europe/Istanbul" // Turkey (TRT, UTC+3:00)

	// Europe - Western Europe
	TimezoneUK         = "Europe/London"   // UK (GMT/BST, UTC+0:00/+1:00)
	TimezoneIreland    = "Europe/Dublin"   // Ireland (GMT/IST, UTC+0:00/+1:00)
	TimezonePortugal   = "Europe/Lisbon"   // Portugal (WET/WEST, UTC+0:00/+1:00)

	// Europe - Central Europe
	TimezoneSpain      = "Europe/Madrid"   // Spain (CET/CEST, UTC+1:00/+2:00)
	TimezoneFrance     = "Europe/Paris"    // France (CET/CEST, UTC+1:00/+2:00)
	TimezoneGermany    = "Europe/Berlin"   // Germany (CET/CEST, UTC+1:00/+2:00)
	TimezoneItaly      = "Europe/Rome"     // Italy (CET/CEST, UTC+1:00/+2:00)
	TimezoneNetherlands = "Europe/Amsterdam" // Netherlands (CET/CEST, UTC+1:00/+2:00)
	TimezoneBelgium    = "Europe/Brussels" // Belgium (CET/CEST, UTC+1:00/+2:00)
	TimezoneSwitzerland = "Europe/Zurich"  // Switzerland (CET/CEST, UTC+1:00/+2:00)
	TimezoneAustria    = "Europe/Vienna"   // Austria (CET/CEST, UTC+1:00/+2:00)
	TimezonePoland     = "Europe/Warsaw"   // Poland (CET/CEST, UTC+1:00/+2:00)
	TimezoneCzech      = "Europe/Prague"   // Czech Republic (CET/CEST, UTC+1:00/+2:00)

	// Europe - Eastern Europe
	TimezoneGreece     = "Europe/Athens"   // Greece (EET/EEST, UTC+2:00/+3:00)
	TimezoneRomania    = "Europe/Bucharest" // Romania (EET/EEST, UTC+2:00/+3:00)
	TimezoneFinland    = "Europe/Helsinki" // Finland (EET/EEST, UTC+2:00/+3:00)
	TimezoneUkraine    = "Europe/Kyiv"     // Ukraine (EET/EEST, UTC+2:00/+3:00)
	TimezoneRussiaMoscow = "Europe/Moscow" // Russia - Moscow (MSK, UTC+3:00)

	// North America
	TimezoneUSAEast    = "America/New_York"      // USA - Eastern (EST/EDT, UTC-5:00/-4:00)
	TimezoneUSACentral = "America/Chicago"       // USA - Central (CST/CDT, UTC-6:00/-5:00)
	TimezoneUSAMountain = "America/Denver"       // USA - Mountain (MST/MDT, UTC-7:00/-6:00)
	TimezoneUSAWest    = "America/Los_Angeles"   // USA - Pacific (PST/PDT, UTC-8:00/-7:00)
	TimezoneCanadaEast = "America/Toronto"       // Canada - Eastern
	TimezoneCanadaWest = "America/Vancouver"     // Canada - Pacific
	TimezoneMexico     = "America/Mexico_City"   // Mexico (CST/CDT, UTC-6:00/-5:00)

	// South America
	TimezoneBrazil     = "America/Sao_Paulo"     // Brazil (BRT/BRST, UTC-3:00/-2:00)
	TimezoneArgentina  = "America/Argentina/Buenos_Aires" // Argentina (ART, UTC-3:00)
	TimezoneChile      = "America/Santiago"      // Chile (CLT/CLST, UTC-4:00/-3:00)
	TimezoneColombia   = "America/Bogota"        // Colombia (COT, UTC-5:00)
	TimezonePeru       = "America/Lima"          // Peru (PET, UTC-5:00)

	// Oceania
	TimezoneAustraliaSydney = "Australia/Sydney"    // Australia - NSW (AEDT/AEST, UTC+11:00/+10:00)
	TimezoneAustraliaMelbourne = "Australia/Melbourne" // Australia - VIC
	TimezoneAustraliaPerth = "Australia/Perth"      // Australia - WA (AWST, UTC+8:00)
	TimezoneNewZealand = "Pacific/Auckland"         // New Zealand (NZDT/NZST, UTC+13:00/+12:00)
	TimezoneFiji       = "Pacific/Fiji"             // Fiji (FJT/FJST, UTC+12:00/+13:00)

	// Africa
	TimezoneSouthAfrica = "Africa/Johannesburg" // South Africa (SAST, UTC+2:00)
	TimezoneEgypt      = "Africa/Cairo"         // Egypt (EET/EEST, UTC+2:00/+3:00)
	TimezoneNigeria    = "Africa/Lagos"         // Nigeria (WAT, UTC+1:00)
	TimezoneKenya      = "Africa/Nairobi"       // Kenya (EAT, UTC+3:00)
	TimezoneMorocco    = "Africa/Casablanca"    // Morocco (WET/WEST, UTC+0:00/+1:00)
)

// GetAllSupportedTimezones returns a map of timezone names to their IANA identifiers
func GetAllSupportedTimezones() map[string]string {
	return map[string]string{
		"UTC/GMT": TimezoneUTC,
		
		// Asia - South Asia
		"India": TimezoneIndia,
		"Pakistan": TimezonePakistan,
		"Sri Lanka": TimezoneSriLanka,
		"Bangladesh": TimezoneBangladesh,
		"Nepal": TimezoneNepal,
		"Maldives": TimezoneMaldives,
		
		// Asia - Southeast Asia
		"Thailand": TimezoneThailand,
		"Vietnam": TimezoneVietnam,
		"Singapore": TimezoneSingapore,
		"Malaysia": TimezoneMalaysia,
		"Indonesia": TimezoneIndonesia,
		"Philippines": TimezonePhilippines,
		"Myanmar": TimezoneMyanmar,
		
		// Asia - East Asia
		"China": TimezoneChina,
		"Hong Kong": TimezoneHongKong,
		"Taiwan": TimezoneTaiwan,
		"Japan": TimezoneJapan,
		"South Korea": TimezoneSouthKorea,
		
		// Asia - Middle East
		"UAE/Dubai": TimezoneDubai,
		"Saudi Arabia": TimezoneSaudiArabia,
		"Qatar": TimezoneQatar,
		"Israel": TimezoneIsrael,
		"Turkey": TimezoneTurkey,
		
		// Europe
		"United Kingdom": TimezoneUK,
		"Ireland": TimezoneIreland,
		"Spain": TimezoneSpain,
		"France": TimezoneFrance,
		"Germany": TimezoneGermany,
		"Italy": TimezoneItaly,
		"Netherlands": TimezoneNetherlands,
		"Belgium": TimezoneBelgium,
		"Switzerland": TimezoneSwitzerland,
		"Austria": TimezoneAustria,
		"Poland": TimezonePoland,
		"Czech Republic": TimezoneCzech,
		"Greece": TimezoneGreece,
		"Romania": TimezoneRomania,
		"Finland": TimezoneFinland,
		"Ukraine": TimezoneUkraine,
		"Russia (Moscow)": TimezoneRussiaMoscow,
		
		// North America
		"USA - Eastern": TimezoneUSAEast,
		"USA - Central": TimezoneUSACentral,
		"USA - Mountain": TimezoneUSAMountain,
		"USA - Pacific": TimezoneUSAWest,
		"Canada - Eastern": TimezoneCanadaEast,
		"Canada - Pacific": TimezoneCanadaWest,
		"Mexico": TimezoneMexico,
		
		// South America
		"Brazil": TimezoneBrazil,
		"Argentina": TimezoneArgentina,
		"Chile": TimezoneChile,
		"Colombia": TimezoneColombia,
		"Peru": TimezonePeru,
		
		// Oceania
		"Australia - Sydney": TimezoneAustraliaSydney,
		"Australia - Melbourne": TimezoneAustraliaMelbourne,
		"Australia - Perth": TimezoneAustraliaPerth,
		"New Zealand": TimezoneNewZealand,
		"Fiji": TimezoneFiji,
		
		// Africa
		"South Africa": TimezoneSouthAfrica,
		"Egypt": TimezoneEgypt,
		"Nigeria": TimezoneNigeria,
		"Kenya": TimezoneKenya,
		"Morocco": TimezoneMorocco,
	}
}

// ValidateTimezone checks if a timezone string is valid
func ValidateTimezone(timezone string) bool {
	_, err := time.LoadLocation(timezone)
	return err == nil
}

// GetTimezoneOffset returns the offset in hours for a given timezone at the current time
func GetTimezoneOffset(timezone string) (float64, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	_, offset := time.Now().In(loc).Zone()
	return float64(offset) / 3600.0, nil
}

// FormatTimeForDisplay formats a UTC time for display in a specific timezone
// Returns formatted string like "2024-12-25 15:30:45 IST"
func FormatTimeForDisplay(utcTime time.Time, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	localTime := utcTime.In(loc)
	zoneName, _ := localTime.Zone()
	
	return fmt.Sprintf("%s %s", localTime.Format("2006-01-02 15:04:05"), zoneName), nil
}

// Frontend Helper Functions - For timezone selection and conversion

// GetTimezoneListForFrontend returns a list of timezones formatted for frontend dropdowns
// Returns slice of {value, label, offset} for easy frontend consumption
type TimezoneOption struct {
	Value  string `json:"value"`  // IANA timezone name (e.g., "Asia/Kolkata")
	Label  string `json:"label"`  // Display name (e.g., "India (IST, UTC+5:30)")
	Offset string `json:"offset"` // Current offset (e.g., "+5:30")
}

func GetTimezoneListForFrontend() ([]TimezoneOption, error) {
	timezones := GetAllSupportedTimezones()
	options := make([]TimezoneOption, 0, len(timezones))
	
	for label, tz := range timezones {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			continue // Skip invalid timezones
		}
		
		now := time.Now().In(loc)
		_, offset := now.Zone()
		
		// Format offset as "+5:30" or "-8:00"
		hours := offset / 3600
		minutes := (offset % 3600) / 60
		var offsetStr string
		if offset >= 0 {
			offsetStr = fmt.Sprintf("+%d:%02d", hours, minutes)
		} else {
			offsetStr = fmt.Sprintf("%d:%02d", hours, minutes)
		}
		
		zoneName, _ := now.Zone()
		fullLabel := fmt.Sprintf("%s (%s, UTC%s)", label, zoneName, offsetStr)
		
		options = append(options, TimezoneOption{
			Value:  tz,
			Label:  fullLabel,
			Offset: offsetStr,
		})
	}
	
	return options, nil
}

// ParseFrontendDateTime parses a datetime string from frontend and converts to UTC
// Frontend sends: "2024-12-25T10:30:00" (local time) + timezone: "Asia/Kolkata"
// Returns: UTC time for storage
func ParseFrontendDateTime(dateTimeStr string, timezone string) (time.Time, error) {
	// Parse the datetime string (assumes format: "2006-01-02T15:04:05" or "2006-01-02 15:04:05")
	var t time.Time
	var err error
	
	// Try ISO format first
	if t, err = time.Parse("2006-01-02T15:04:05", dateTimeStr); err != nil {
		// Try space-separated format
		if t, err = time.Parse("2006-01-02 15:04:05", dateTimeStr); err != nil {
			// Try with seconds optional
			if t, err = time.Parse("2006-01-02T15:04", dateTimeStr); err != nil {
				return time.Time{}, fmt.Errorf("invalid datetime format '%s': %w", dateTimeStr, err)
			}
		}
	}
	
	// Load timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	// Create time in the specified timezone
	localTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, loc)
	
	// Convert to UTC for storage
	return localTime.UTC(), nil
}

// ParseFrontendDate parses a date string from frontend and converts to UTC
// Frontend sends: "2024-12-25" (date only) + timezone: "Asia/Kolkata"
// Returns: UTC time at start of day (00:00:00) in that timezone
func ParseFrontendDate(dateStr string, timezone string) (time.Time, error) {
	// Parse date string (format: "2006-01-02")
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format '%s': %w", dateStr, err)
	}
	
	// Load timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	// Create start of day in the specified timezone
	localTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	
	// Convert to UTC for storage
	return localTime.UTC(), nil
}

// FormatUTCForFrontend formats a UTC time for frontend display in a specific timezone
// Returns ISO 8601 string in the specified timezone (e.g., "2024-12-25T15:30:45+05:30")
func FormatUTCForFrontend(utcTime time.Time, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	localTime := utcTime.In(loc)
	return localTime.Format(time.RFC3339), nil
}

// FormatUTCDateForFrontend formats a UTC time as date only for frontend display
// Returns date string in the specified timezone (e.g., "2024-12-25")
func FormatUTCDateForFrontend(utcTime time.Time, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	localTime := utcTime.In(loc)
	return localTime.Format("2006-01-02"), nil
}

// FormatUTCDateTimeForFrontend formats a UTC time for frontend display (without timezone suffix)
// Returns datetime string in the specified timezone (e.g., "2024-12-25T15:30:45")
func FormatUTCDateTimeForFrontend(utcTime time.Time, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	localTime := utcTime.In(loc)
	return localTime.Format("2006-01-02T15:04:05"), nil
}

// GetStartOfDayUTC returns the start of day (00:00:00) in the specified timezone, converted to UTC
// Useful for date range queries: "Get all records from Dec 25 in India timezone"
func GetStartOfDayUTC(date time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	// Create start of day in the specified timezone
	localStartOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	
	// Convert to UTC
	return localStartOfDay.UTC(), nil
}

// GetEndOfDayUTC returns the end of day (23:59:59) in the specified timezone, converted to UTC
// Useful for date range queries
func GetEndOfDayUTC(date time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}
	
	// Create end of day in the specified timezone
	localEndOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, loc)
	
	// Convert to UTC
	return localEndOfDay.UTC(), nil
}

// GetDateRangeUTC returns start and end of day in UTC for a given date and timezone
// Returns (startUTC, endUTC, error)
func GetDateRangeUTC(date time.Time, timezone string) (time.Time, time.Time, error) {
	start, err := GetStartOfDayUTC(date, timezone)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	
	end, err := GetEndOfDayUTC(date, timezone)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	
	return start, end, nil
}
