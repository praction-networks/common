package tenantevent

// ==================== TENANT BRANDING EVENTS ====================

// TenantBrandingInsertEventModel represents a tenant branding creation event
// Carries branding data needed by downstream services (captive-portal-service, etc.)
type TenantBrandingInsertEventModel struct {
	ID          string `bson:"_id" json:"id"`
	TenantID    string `bson:"tenantId" json:"tenantId"`
	CompanyName string `bson:"companyName" json:"companyName"`

	// Logo URLs (extracted from Logo.MainLogo, Logo.DarkLogo, etc.)
	LogoUrl     string `bson:"logoUrl,omitempty" json:"logoUrl,omitempty"`
	LogoDarkUrl string `bson:"logoDarkUrl,omitempty" json:"logoDarkUrl,omitempty"`
	SquareLogo  string `bson:"squareLogo,omitempty" json:"squareLogo,omitempty"`
	IconLogo    string `bson:"iconLogo,omitempty" json:"iconLogo,omitempty"`
	Favicon     string `bson:"favicon,omitempty" json:"favicon,omitempty"`

	// Brand Colors
	PrimaryColor   string `bson:"primaryColor,omitempty" json:"primaryColor,omitempty"`
	SecondaryColor string `bson:"secondaryColor,omitempty" json:"secondaryColor,omitempty"`
	AccentColor    string `bson:"accentColor,omitempty" json:"accentColor,omitempty"`

	IsActive bool `bson:"isActive" json:"isActive"`
	Version  int  `bson:"version" json:"version"`
}

// TenantBrandingUpdateEventModel represents a tenant branding update event
type TenantBrandingUpdateEventModel struct {
	ID          string `bson:"_id" json:"id"`
	TenantID    string `bson:"tenantId" json:"tenantId"`
	CompanyName string `bson:"companyName,omitempty" json:"companyName,omitempty"`

	// Logo URLs
	LogoUrl     string `bson:"logoUrl,omitempty" json:"logoUrl,omitempty"`
	LogoDarkUrl string `bson:"logoDarkUrl,omitempty" json:"logoDarkUrl,omitempty"`
	SquareLogo  string `bson:"squareLogo,omitempty" json:"squareLogo,omitempty"`
	IconLogo    string `bson:"iconLogo,omitempty" json:"iconLogo,omitempty"`
	Favicon     string `bson:"favicon,omitempty" json:"favicon,omitempty"`

	// Brand Colors
	PrimaryColor   string `bson:"primaryColor,omitempty" json:"primaryColor,omitempty"`
	SecondaryColor string `bson:"secondaryColor,omitempty" json:"secondaryColor,omitempty"`
	AccentColor    string `bson:"accentColor,omitempty" json:"accentColor,omitempty"`

	IsActive bool `bson:"isActive" json:"isActive"`
	Version  int  `bson:"version" json:"version"`
}

// TenantBrandingDeleteEventModel represents a tenant branding deletion event
type TenantBrandingDeleteEventModel struct {
	ID       string `bson:"_id" json:"id"`
	TenantID string `bson:"tenantId" json:"tenantId"`
}
