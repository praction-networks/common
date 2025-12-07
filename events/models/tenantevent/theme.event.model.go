package tenantevent

import "time"

// ThemeInsertEventModel defines the model for theme created events
// Unified subject: theme.created
// Consumers filter by portalType in event payload
// Structure matches TenantThemeConfig for easy BSON/JSON conversion
type ThemeInsertEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantId" json:"tenantId"`

	// Portal-specific theme configs (only one will be set based on portalType)
	HotspotTheme     *HotspotPortalThemeEvent `bson:"hotspotTheme,omitempty" json:"hotspotTheme,omitempty"`
	TenantAdminTheme *TenantAdminThemeEvent   `bson:"tenantAdminTheme,omitempty" json:"tenantAdminTheme,omitempty"`
	LandingPageTheme *LandingPageThemeEvent   `bson:"landingPageTheme,omitempty" json:"landingPageTheme,omitempty"`

	// Metadata
	Version   int       `bson:"version" json:"version"` // Version for matching/validation
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	CreatedBy string    `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
}

// ThemeUpdateEventModel defines the model for theme updated events
// Unified subject: theme.updated
// Consumers filter by portalType in event payload and match version
type ThemeUpdateEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantId" json:"tenantId"`

	// Portal-specific theme configs (only one will be set based on portalType)
	HotspotTheme     *HotspotPortalThemeEvent `bson:"hotspotTheme,omitempty" json:"hotspotTheme,omitempty"`
	TenantAdminTheme *TenantAdminThemeEvent   `bson:"tenantAdminTheme,omitempty" json:"tenantAdminTheme,omitempty"`
	LandingPageTheme *LandingPageThemeEvent   `bson:"landingPageTheme,omitempty" json:"landingPageTheme,omitempty"`

	// Metadata
	Version   int       `bson:"version" json:"version"` // Version for matching/validation
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy string    `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

// ThemeDeleteEventModel defines the model for theme deleted events
// Unified subject: theme.deleted
// Consumers filter by portalType in event payload
type ThemeDeleteEventModel struct {
	ID         string `bson:"_id" json:"id"`
	TenantID   string `bson:"tenantId" json:"tenantId"`
	PortalType string `bson:"portalType" json:"portalType"` // "hotspot" | "tenant_admin" | "landing"
	Version    int    `bson:"version" json:"version"`       // Version for matching/validation
}

// ThemeSetDefaultEventModel defines the model for theme set default events
// Unified subject: theme.set_default
// Consumers filter by portalType in event payload
type ThemeSetDefaultEventModel struct {
	ID                string `bson:"_id" json:"id"`
	TenantID          string `bson:"tenantId" json:"tenantId"`
	PortalType        string `bson:"portalType" json:"portalType"` // "hotspot" | "tenant_admin" | "landing"
	Version           int    `bson:"version" json:"version"`       // Version for matching/validation
	PreviousDefaultID string `bson:"previousDefaultId,omitempty" json:"previousDefaultId,omitempty"`
}

// HotspotPortalThemeEvent - Matches HotspotPortalTheme structure
type HotspotPortalThemeEvent struct {
	Branding   BrandingConfigEvent  `bson:"branding" json:"branding"`
	Content    ContentConfigEvent   `bson:"content" json:"content"`
	Behavior   BehaviorConfigEvent  `bson:"behavior" json:"behavior"`
	Advanced   AdvancedConfigEvent  `bson:"advanced,omitempty" json:"advanced,omitempty"`
	Components ComponentStylesEvent `bson:"components,omitempty" json:"components,omitempty"`
}

// TenantAdminThemeEvent - Matches TenantAdminTheme structure
type TenantAdminThemeEvent struct {
	Branding   BrandingConfigEvent   `bson:"branding" json:"branding"`
	Typography TypographyConfigEvent `bson:"typography" json:"typography"`
	UIDensity  string                `bson:"uiDensity" json:"uiDensity"`
}

// LandingPageThemeEvent - Matches LandingPageTheme structure
type LandingPageThemeEvent struct {
	Branding   BrandingConfigEvent  `bson:"branding" json:"branding"`
	Content    ContentConfigEvent   `bson:"content" json:"content"`
	Advanced   AdvancedConfigEvent  `bson:"advanced,omitempty" json:"advanced,omitempty"`
	Components ComponentStylesEvent `bson:"components,omitempty" json:"components,omitempty"`
}

// BrandingConfigEvent - Matches BrandingConfig structure
type BrandingConfigEvent struct {
	LogoURL         string `bson:"logoUrl,omitempty" json:"logoUrl,omitempty"`
	FaviconURL      string `bson:"faviconUrl,omitempty" json:"faviconUrl,omitempty"`
	PrimaryColor    string `bson:"primaryColor" json:"primaryColor"`
	SecondaryColor  string `bson:"secondaryColor" json:"secondaryColor"`
	BackgroundColor string `bson:"backgroundColor" json:"backgroundColor"`
	BackgroundImage string `bson:"backgroundImage,omitempty" json:"backgroundImage,omitempty"`
	FontFamily      string `bson:"fontFamily,omitempty" json:"fontFamily,omitempty"`
	BorderRadius    string `bson:"borderRadius,omitempty" json:"borderRadius,omitempty"`
}

// ContentConfigEvent - Matches ContentConfig structure
type ContentConfigEvent struct {
	WelcomeTitle   string `bson:"welcomeTitle" json:"welcomeTitle"`
	WelcomeMessage string `bson:"welcomeMessage,omitempty" json:"welcomeMessage,omitempty"`
	FooterText     string `bson:"footerText,omitempty" json:"footerText,omitempty"`
	TermsURL       string `bson:"termsUrl,omitempty" json:"termsUrl,omitempty"`
	PrivacyURL     string `bson:"privacyUrl,omitempty" json:"privacyUrl,omitempty"`
	SupportEmail   string `bson:"supportEmail,omitempty" json:"supportEmail,omitempty"`
	SupportPhone   string `bson:"supportPhone,omitempty" json:"supportPhone,omitempty"`
	SupportURL     string `bson:"supportUrl,omitempty" json:"supportUrl,omitempty"`
}

// BehaviorConfigEvent - Matches BehaviorConfig structure
type BehaviorConfigEvent struct {
	SessionDurationSeconds   int      `bson:"sessionDurationSeconds" json:"sessionDurationSeconds"`
	IdleTimeoutSeconds       int      `bson:"idleTimeoutSeconds,omitempty" json:"idleTimeoutSeconds,omitempty"`
	AllowAutoLogin           bool     `bson:"allowAutoLogin" json:"allowAutoLogin"`
	RememberDeviceDays       int      `bson:"rememberDeviceDays" json:"rememberDeviceDays"`
	AutoRedirect             bool     `bson:"autoRedirect" json:"autoRedirect"`
	DefaultRedirectURL       string   `bson:"defaultRedirectUrl,omitempty" json:"defaultRedirectUrl,omitempty"`
	RequireEmailVerification bool     `bson:"requireEmailVerification" json:"requireEmailVerification"`
	RequireMobileOtp         bool     `bson:"requireMobileOtp" json:"requireMobileOtp"`
	RequireTermsAcceptance   bool     `bson:"requireTermsAcceptance" json:"requireTermsAcceptance"`
	RequirePrivacyAcceptance bool     `bson:"requirePrivacyAcceptance" json:"requirePrivacyAcceptance"`
	AllowGuestAccess         bool     `bson:"allowGuestAccess" json:"allowGuestAccess"`
	MaxDevicesPerUser        int      `bson:"maxDevicesPerUser" json:"maxDevicesPerUser"`
	BlockedMACs              []string `bson:"blockedMacs,omitempty" json:"blockedMacs,omitempty"`
	AllowedDomains           []string `bson:"allowedDomains,omitempty" json:"allowedDomains,omitempty"`
	EnableSocialLogin        bool     `bson:"enableSocialLogin" json:"enableSocialLogin"`
	AllowedSocialMethods     []string `bson:"allowedSocialMethods,omitempty" json:"allowedSocialMethods,omitempty"`
}

// AdvancedConfigEvent - Matches AdvancedConfig structure
type AdvancedConfigEvent struct {
	CustomCSS          string   `bson:"customCss,omitempty" json:"customCss,omitempty"`
	CustomHTML         string   `bson:"customHtml,omitempty" json:"customHtml,omitempty"`
	CustomJavaScript   string   `bson:"customJavaScript,omitempty" json:"customJavaScript,omitempty"`
	DefaultLanguage    string   `bson:"defaultLanguage" json:"defaultLanguage"`
	SupportedLanguages []string `bson:"supportedLanguages,omitempty" json:"supportedLanguages,omitempty"`
	Timezone           string   `bson:"timezone" json:"timezone"`
	EnableAnalytics    bool     `bson:"enableAnalytics" json:"enableAnalytics"`
	GoogleAnalyticsID  string   `bson:"googleAnalyticsId,omitempty" json:"googleAnalyticsId,omitempty"`
	FacebookPixelID    string   `bson:"facebookPixelId,omitempty" json:"facebookPixelId,omitempty"`
}

// ComponentStylesEvent - Matches ComponentStyles structure
type ComponentStylesEvent struct {
	Header *ComponentStyleEvent `bson:"header,omitempty" json:"header,omitempty"`
	Form   *ComponentStyleEvent `bson:"form,omitempty" json:"form,omitempty"`
	Input  *ComponentStyleEvent `bson:"input,omitempty" json:"input,omitempty"`
	Button *ComponentStyleEvent `bson:"button,omitempty" json:"button,omitempty"`
}

// ComponentStyleEvent - Matches ComponentStyle structure
type ComponentStyleEvent struct {
	BackgroundColor string `bson:"backgroundColor,omitempty" json:"backgroundColor,omitempty"`
	BorderColor     string `bson:"borderColor,omitempty" json:"borderColor,omitempty"`
	BorderRadius    string `bson:"borderRadius,omitempty" json:"borderRadius,omitempty"`
	Padding         string `bson:"padding,omitempty" json:"padding,omitempty"`
	FontFamily      string `bson:"fontFamily,omitempty" json:"fontFamily,omitempty"`
}

// TypographyConfigEvent - Matches TypographyConfig structure
type TypographyConfigEvent struct {
	FontFamily       string  `bson:"fontFamily" json:"fontFamily"`
	FontSizeBase     int     `bson:"fontSizeBase" json:"fontSizeBase"`
	FontWeightNormal int     `bson:"fontWeightNormal" json:"fontWeightNormal"`
	FontWeightBold   int     `bson:"fontWeightBold" json:"fontWeightBold"`
	LineHeight       float64 `bson:"lineHeight" json:"lineHeight"`
}
