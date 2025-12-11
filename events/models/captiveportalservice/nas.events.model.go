package captiveportalservice

import (
	"time"
)

// NAS Device field constants
const (
	NASFieldID              = "_id"
	NASFieldNASId           = "nasId"
	NASFieldTenantIDs       = "tenantIds"
	NASFieldLocationAddress = "locationAddress"
	NASFieldName            = "name"
	NASFieldIPAddress       = "ipAddress"
	NASFieldIsActive        = "isActive"
	NASFieldSSIDs           = "ssids"
	NASFieldCreatedAt       = "createdAt"
	NASFieldUpdatedAt       = "updatedAt"
	NASFieldCreatedBy       = "createdBy"
	NASFieldUpdatedBy       = "updatedBy"
	NASFieldVersion         = "version"
	NASFieldTags            = "tags"
	NASFieldNotes           = "notes"
	NASFieldStats           = "stats"
	NASFieldLastSeen        = "lastSeen"
)

// NASDevice represents a Network Access Server (WiFi access point, router, etc.)
type NASDevice struct {
	ID        string   `bson:"_id,omitempty" json:"id"`
	NASId     string   `bson:"nasId" json:"nasId"`         // Unique NAS identifier (e.g., NAS001, AP-CAFE-01)
	TenantIDs []string `bson:"tenantIds" json:"tenantIds"` // Which tenants own this NAS

	// Device information
	Name         string `bson:"name" json:"name"`                                     // Friendly name (e.g., "Cool Cafe Main AP")
	Description  string `bson:"description,omitempty" json:"description,omitempty"`   // Additional details
	Manufacturer string `bson:"manufacturer,omitempty" json:"manufacturer,omitempty"` // MikroTik, Ubiquiti, pfSense
	Model        string `bson:"model,omitempty" json:"model,omitempty"`               // Device model

	// Network configuration
	IPAddress       string   `bson:"ipAddress" json:"ipAddress"`                                 // NAS IP address
	AllowedIPRanges []string `bson:"allowedIPRanges,omitempty" json:"allowedIPRanges,omitempty"` // Client IP ranges (e.g., ["10.0.0.0/24"])
	SSIDs           []string `bson:"ssids" json:"ssids"`                                         // WiFi SSIDs broadcasted by this NAS

	// Security credentials
	Secret       string `bson:"secret" json:"-"`                                // For HMAC/JWT signing and API authentication (bcrypt hashed, never sent to client)
	RADIUSSecret string `bson:"radiusSecret" json:"-"`                          // For RADIUS protocol communication (bcrypt hashed, never sent to client)
	PublicKey    string `bson:"publicKey,omitempty" json:"publicKey,omitempty"` // For asymmetric crypto (optional)

	// Security settings
	TokenExpirySeconds int  `bson:"tokenExpirySeconds" json:"tokenExpirySeconds"` // How long NAS tokens are valid (default: 300)
	RequireMutualTLS   bool `bson:"requireMutualTLS" json:"requireMutualTLS"`     // Require mutual TLS for NAS calls
	EnableRateLimit    bool `bson:"enableRateLimit" json:"enableRateLimit"`       // Enable rate limiting

	// Rate limiting settings
	MaxSessionsPerHour  int `bson:"maxSessionsPerHour" json:"maxSessionsPerHour"`   // Max new sessions per hour
	MaxDevicesPerMAC    int `bson:"maxDevicesPerMAC" json:"maxDevicesPerMAC"`       // Max devices per MAC (prevent sharing)
	MaxRetriesPerDevice int `bson:"maxRetriesPerDevice" json:"maxRetriesPerDevice"` // Max failed auth attempts

	// Status
	IsActive bool      `bson:"isActive" json:"isActive"`                     // Is this NAS currently active?
	LastSeen time.Time `bson:"lastSeen,omitempty" json:"lastSeen,omitempty"` // Last health check

	// Statistics (denormalized)
	Stats NASStats `bson:"stats,omitempty" json:"stats,omitempty"`

	// Metadata
	Tags  []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Notes string   `bson:"notes,omitempty" json:"notes,omitempty"`

	// Audit fields
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	CreatedBy string    `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedBy string    `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	Version   int       `bson:"version" json:"version"`
}

// NASStats contains denormalized statistics for NAS device
type NASStats struct {
	TotalSessions      int       `bson:"totalSessions" json:"totalSessions"`
	ActiveSessions     int       `bson:"activeSessions" json:"activeSessions"`
	TotalDevices       int       `bson:"totalDevices" json:"totalDevices"`
	FailedAuthAttempts int       `bson:"failedAuthAttempts" json:"failedAuthAttempts"`
	LastSessionAt      time.Time `bson:"lastSessionAt,omitempty" json:"lastSessionAt,omitempty"`
}
