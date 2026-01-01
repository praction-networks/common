package subscriberevent

import "time"

// HotspotProfileCreatedEvent represents a hotspot profile creation event
type HotspotProfileCreatedEvent struct {
	ID                      string               `json:"id" bson:"id"`
	SubscriberID            string               `json:"subscriberId" bson:"subscriberId"`
	TenantID                string               `json:"tenantId" bson:"tenantId"`
	LocationIDs             []string             `bson:"locationIds,omitempty" json:"locationIds,omitempty"`
	Status                  HotspotProfileStatus `json:"status" bson:"status"`
	Username                string               `json:"username,omitempty" bson:"username,omitempty"` // RADIUS username
	Password                string               `json:"password,omitempty" bson:"password,omitempty"` // RADIUS password
	MacAddresses            []string             `bson:"macAddresses,omitempty" json:"macAddresses,omitempty"`
	BlacklistedMacAddresses []string             `bson:"blacklistedMacAddresses,omitempty" json:"blacklistedMacAddresses,omitempty"` // MAC addresses blocked from using this profile
	DefaultAuthMethod       HotspotAuthMethod    `json:"defaultAuthMethod" bson:"defaultAuthMethod"`
	AllowedAuthMethods      []HotspotAuthMethod  `bson:"allowedAuthMethods,omitempty" json:"allowedAuthMethods,omitempty"`
	MaxConcurrentSessions   int                  `json:"maxConcurrentSessions,omitempty" bson:"maxConcurrentSessions,omitempty"` // max concurrent logins for this profile
	MaxDevices               int                  `json:"maxDevices,omitempty" bson:"maxDevices,omitempty"`             // max devices allowed (policy)
	RememberDeviceDays      int                  `json:"rememberDeviceDays,omitempty" bson:"rememberDeviceDays,omitempty"`       // days to remember each device for auto-login
	CreationSource          string               `json:"creationSource,omitempty" bson:"creationSource,omitempty"`               // "CAPTIVE_PORTAL", "ADMIN", "API", "IMPORT"
	ValidFrom               time.Time            `json:"validFrom,omitempty" bson:"validFrom,omitempty"`                         // When this profile becomes valid
	ValidUntil              *time.Time           `json:"validUntil,omitempty" bson:"validUntil,omitempty"`                       // When this profile expires (nil = no expiration)
	PlanCode                string               `json:"planCode,omitempty" bson:"planCode,omitempty"`                           // Plan code used for this hotspot profile
	Version                 int                  `json:"version" bson:"version"`
}

// HotspotProfileUpdatedEvent represents a hotspot profile update event
type HotspotProfileUpdatedEvent struct {
	ID                      string               `json:"id" bson:"id"`
	SubscriberID            string               `json:"subscriberId" bson:"subscriberId"`
	TenantID                string               `json:"tenantId" bson:"tenantId"`
	LocationIDs             []string             `bson:"locationIds,omitempty" json:"locationIds,omitempty"`
	Status                  HotspotProfileStatus `json:"status,omitempty" bson:"status,omitempty"`
	UserID                  string               `json:"userId,omitempty" bson:"userId,omitempty"`
	Username                string               `json:"username,omitempty" bson:"username,omitempty"` // RADIUS username
	Password                string               `json:"password,omitempty" bson:"password,omitempty"` // RADIUS password
	MacAddresses            []string             `bson:"macAddresses,omitempty" json:"macAddresses,omitempty"`
	BlacklistedMacAddresses []string             `bson:"blacklistedMacAddresses,omitempty" json:"blacklistedMacAddresses,omitempty"` // MAC addresses blocked from using this profile
	DefaultAuthMethod       HotspotAuthMethod    `json:"defaultAuthMethod,omitempty" bson:"defaultAuthMethod,omitempty"`
	AllowedAuthMethods      []HotspotAuthMethod  `bson:"allowedAuthMethods,omitempty" json:"allowedAuthMethods,omitempty"`
	MaxConcurrentSessions   int                  `json:"maxConcurrentSessions,omitempty" bson:"maxConcurrentSessions,omitempty"` // max concurrent logins for this profile
	MaxDevices               int                  `json:"maxDevices,omitempty" bson:"maxDevices,omitempty"`             // max devices allowed (policy)
	RememberDeviceDays      int                  `json:"rememberDeviceDays,omitempty" bson:"rememberDeviceDays,omitempty"`       // days to remember each device for auto-login
	FirstLoginAt            *time.Time           `json:"firstLoginAt,omitempty" bson:"firstLoginAt,omitempty"`                 // When was the first successful login
	ValidFrom               *time.Time           `json:"validFrom,omitempty" bson:"validFrom,omitempty"`                         // When this profile becomes valid (optional update)
	ValidUntil              *time.Time           `json:"validUntil,omitempty" bson:"validUntil,omitempty"`                       // When this profile expires (optional update)
	PlanCode                string               `json:"planCode,omitempty" bson:"planCode,omitempty"`                           // Plan code used for this hotspot profile (optional update)
	Version                 int                  `json:"version" bson:"version"`
}

// HotspotProfileDeletedEvent represents a hotspot profile deletion event
type HotspotProfileDeletedEvent struct {
	ID           string `json:"id" bson:"id"`
	SubscriberID string `json:"subscriberId" bson:"subscriberId"`
	TenantID     string `json:"tenantId" bson:"tenantId"`
	Version      int    `json:"version" bson:"version"`
}

