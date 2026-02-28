package tenantuserevent

// UserAccess represents user access permissions for a specific zone/department combination
type UserAccess struct {
	TenantID  string `json:"tenantId" bson:"tenantId"`
	RoleID    string `json:"roleId" bson:"roleId"`
	IsPrimary bool   `json:"isPrimary" bson:"isPrimary"` // Primary assignment for this user
	IsActive  bool   `json:"isActive" bson:"isActive"`   // Active assignment
}

// Domain User Events (Extended Profiles)
type TenantUserCreateEvent struct {
	ID                       string       `json:"id" bson:"_id"`
	FirstName                string       `json:"firstName" bson:"firstName"`
	MiddleName               string       `json:"middleName" bson:"middleName"`
	LastName                 string       `json:"lastName" bson:"lastName"`
	Mobile                   string       `json:"mobile" bson:"mobile"`
	IsMobileVerified         bool         `json:"isMobileVerified" bson:"isMobileVerified"`
	OfficialMobile           string       `json:"officialMobile" bson:"officialMobile"`
	IsOfficialMobileVerified bool         `json:"isOfficialMobileVerified" bson:"isOfficialMobileVerified"`
	Whatsapp                 string       `json:"whatsapp" bson:"whatsapp"`
	IsWhatsappVerified       bool         `json:"isWhatsappVerified" bson:"isWhatsappVerified"`
	Email                    string       `json:"email" bson:"email"`
	IsEmailVerified          bool         `json:"isEmailVerified" bson:"isEmailVerified"`
	OfficialEmail            string       `json:"officialEmail" bson:"officialEmail"`
	IsOfficialEmailVerified  bool         `json:"isOfficialEmailVerified" bson:"isOfficialEmailVerified"`
	Gender                   string       `json:"gender" bson:"gender"`
	DOB                      string       `json:"dob" bson:"dob"`
	UserAccess               []UserAccess `json:"userAccess" bson:"userAccess"` // Multiple zone/department assignments
	OnLeave                  bool         `json:"onLeave" bson:"onLeave"`
	IsActive                 bool         `json:"isActive" bson:"isActive"`
	IsSystem                 bool         `json:"isSystem" bson:"isSystem"`
	Version                  int          `json:"version" bson:"version"`

	// Organizational Hierarchy
	Designation   string   `json:"designation,omitempty" bson:"designation,omitempty"`
	Department    string   `json:"department,omitempty" bson:"department,omitempty"`
	ReportsTo     string   `json:"reportsTo,omitempty" bson:"reportsTo,omitempty"`
	DirectReports []string `json:"directReports,omitempty" bson:"directReports,omitempty"`
	OrgLevel      int      `json:"orgLevel" bson:"orgLevel"`
}

type TenantUserUpdateEvent struct {
	ID                       string       `json:"id" bson:"_id"`
	FirstName                string       `json:"firstName" bson:"firstName"`
	MiddleName               string       `json:"middleName" bson:"middleName"`
	LastName                 string       `json:"lastName" bson:"lastName"`
	Mobile                   string       `json:"mobile" bson:"mobile"`
	IsMobileVerified         bool         `json:"isMobileVerified" bson:"isMobileVerified"`
	OfficialMobile           string       `json:"officialMobile" bson:"officialMobile"`
	IsOfficialMobileVerified bool         `json:"isOfficialMobileVerified" bson:"isOfficialMobileVerified"`
	Whatsapp                 string       `json:"whatsapp" bson:"whatsapp"`
	IsWhatsappVerified       bool         `json:"isWhatsappVerified" bson:"isWhatsappVerified"`
	Email                    string       `json:"email" bson:"email"`
	IsEmailVerified          bool         `json:"isEmailVerified" bson:"isEmailVerified"`
	OfficialEmail            string       `json:"officialEmail" bson:"officialEmail"`
	IsOfficialEmailVerified  bool         `json:"isOfficialEmailVerified" bson:"isOfficialEmailVerified"`
	Gender                   string       `json:"gender" bson:"gender"`
	DOB                      string       `json:"dob" bson:"dob"`
	UserAccess               []UserAccess `json:"userAccess" bson:"userAccess"` // Updated assignments
	IsSystem                 bool         `json:"isSystem" bson:"isSystem"`
	OnLeave                  bool         `json:"onLeave" bson:"onLeave"`
	IsActive                 bool         `json:"isActive" bson:"isActive"`
	Version                  int          `json:"version" bson:"version"`

	// Organizational Hierarchy
	Designation   string   `json:"designation,omitempty" bson:"designation,omitempty"`
	Department    string   `json:"department,omitempty" bson:"department,omitempty"`
	ReportsTo     string   `json:"reportsTo,omitempty" bson:"reportsTo,omitempty"`
	DirectReports []string `json:"directReports,omitempty" bson:"directReports,omitempty"`
	OrgLevel      int      `json:"orgLevel" bson:"orgLevel"`
}

type TenantUserDeleteEvent struct {
	ID         string       `json:"id" bson:"_id"`
	UserAccess []UserAccess `json:"userAccess" bson:"userAccess"` // Updated assignments
	Version    int          `json:"version" bson:"version"`
}

// TenantUserPasswordSetEvent is published when a new user sets their initial password
// after OTP verification during onboarding. This event is consumed by auth-service
// to create the credential password record securely via internal NATS communication.
type TenantUserPasswordSetEvent struct {
	UserID       string `json:"userId"`       // User ID (cuid2)
	PasswordHash string `json:"passwordHash"` // Bcrypt hashed password
	Timestamp    string `json:"timestamp"`    // ISO 8601 timestamp
}

// UserPreferences represents user UI preferences (shared with tenant-user-service model)
type UserPreferences struct {
	SidebarMenu     *SidebarMenuPreferences     `json:"sidebarMenu,omitempty"`
	TenantFavorites *TenantFavoritesPreferences `json:"tenantFavorites,omitempty"`
	Theme           *UserThemePreferences       `json:"theme,omitempty"`
}

// SidebarMenuPreferences represents sidebar menu organization preferences
type SidebarMenuPreferences struct {
	GroupOrder        []string            `json:"groupOrder,omitempty"`        // Ordered group titles
	ItemOrder         map[string][]string `json:"itemOrder,omitempty"`         // groupTitle -> ordered item hrefs
	Favorites         []string            `json:"favorites,omitempty"`         // Starred item hrefs
	UseFavoritesOrder bool                `json:"useFavoritesOrder,omitempty"` // If true, use favorites-based order
}

// TenantFavoritesPreferences represents tenant favorites and ordering preferences
type TenantFavoritesPreferences struct {
	FavoriteTenantIds []string `json:"favoriteTenantIds,omitempty"` // Starred tenant IDs
	TenantOrder       []string `json:"tenantOrder,omitempty"`       // Custom ordered tenant IDs
	UseFavoritesOrder bool     `json:"useFavoritesOrder,omitempty"` // If true, favorites first
}

// UserThemePreferences represents user-specific theme overrides
type UserThemePreferences struct {
	OverrideTenantTheme bool   `json:"overrideTenantTheme,omitempty"`
	PrimaryColor        string `json:"primaryColor,omitempty"`
	TextOnPrimary       string `json:"textOnPrimary,omitempty"`
	AccentColor         string `json:"accentColor,omitempty"`
	BackgroundColor     string `json:"backgroundColor,omitempty"`
	SurfaceColor        string `json:"surfaceColor,omitempty"`
	BorderColor         string `json:"borderColor,omitempty"`
	PreferredLanguage   string `json:"preferredLanguage,omitempty"`
	// Note: Timezone is tenant-specific only, not user-specific
}

// TenantUserPreferencesUpdatedEvent is published when user preferences are updated
// This event is consumed by auth-service to update the cached preferences
type TenantUserPreferencesUpdatedEvent struct {
	UserID      string           `json:"userId"`      // User ID (cuid2)
	Preferences *UserPreferences `json:"preferences"` // Updated preferences
	Timestamp   string           `json:"timestamp"`   // ISO 8601 timestamp
}
