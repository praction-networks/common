package tenantevent

import "time"

// ThemeInsertEventModel defines the model for theme created events
// Unified subject: theme.created
// Consumers filter by portalType in event payload
type ThemeInsertEventModel struct {
	ID         string                 `bson:"_id" json:"id"`
	TenantID   string                 `bson:"tenantId" json:"tenantId"`
	PortalType string                 `bson:"portalType" json:"portalType"` // "hotspot" | "tenant_admin" | "landing"
	ThemeID    string                 `bson:"themeId" json:"themeId"`
	Version    int                    `bson:"version" json:"version"`     // Version for matching/validation
	ThemeData  map[string]interface{} `bson:"themeData" json:"themeData"` // HotspotPortalTheme | TenantAdminTheme | LandingPageTheme
	CreatedAt  time.Time              `bson:"createdAt" json:"createdAt"`
}

// ThemeUpdateEventModel defines the model for theme updated events
// Unified subject: theme.updated
// Consumers filter by portalType in event payload and match version
type ThemeUpdateEventModel struct {
	ID         string                 `bson:"_id" json:"id"`
	TenantID   string                 `bson:"tenantId" json:"tenantId"`
	PortalType string                 `bson:"portalType" json:"portalType"`
	ThemeID    string                 `bson:"themeId" json:"themeId"`
	Version    int                    `bson:"version" json:"version"` // Version for matching/validation
	Changes    map[string]interface{} `bson:"changes" json:"changes"` // Only changed fields
	UpdatedAt  time.Time              `bson:"updatedAt" json:"updatedAt"`
}

// ThemeDeleteEventModel defines the model for theme deleted events
// Unified subject: theme.deleted
// Consumers filter by portalType in event payload
type ThemeDeleteEventModel struct {
	ID         string `bson:"_id" json:"id"`
	TenantID   string `bson:"tenantId" json:"tenantId"`
	PortalType string `bson:"portalType" json:"portalType"`
	ThemeID    string `bson:"themeId" json:"themeId"`
	Version    int    `bson:"version" json:"version"` // Version for matching/validation
}

// ThemeSetDefaultEventModel defines the model for theme set default events
// Unified subject: theme.set_default
// Consumers filter by portalType in event payload
type ThemeSetDefaultEventModel struct {
	ID                string `bson:"_id" json:"id"`
	TenantID          string `bson:"tenantId" json:"tenantId"`
	PortalType        string `bson:"portalType" json:"portalType"`
	ThemeID           string `bson:"themeId" json:"themeId"`
	Version           int    `bson:"version" json:"version"` // Version for matching/validation
	PreviousDefaultID string `bson:"previousDefaultId,omitempty" json:"previousDefaultId,omitempty"`
}
