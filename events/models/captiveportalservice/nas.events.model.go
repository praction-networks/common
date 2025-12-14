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
	NASFieldRedirectorDomain = "redirectorDomain"
	NASFieldForwardDomain    = "forwardDomain"
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
type NASCreatedEvent struct {
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
	RedirectorDomain string   `bson:"redirectorDomain" json:"redirectorDomain"`                   // Redirector domain (e.g., "redirector.rapidnet.in")
	ForwardDomain    string   `bson:"forwardDomain" json:"forwardDomain"`                         // Forward domain (e.g., "forward.rapidnet.in")

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

	// Metadata
	Tags  []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Notes string   `bson:"notes,omitempty" json:"notes,omitempty"`

	Version int `bson:"version" json:"version"`
}

// NASDevice represents a Network Access Server (WiFi access point, router, etc.)
type NASUpdatedEvent struct {
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
	RedirectorDomain string   `bson:"redirectorDomain" json:"redirectorDomain"`                   // Redirector domain (e.g., "redirector.rapidnet.in")
	ForwardDomain    string   `bson:"forwardDomain" json:"forwardDomain"`                         // Forward domain (e.g., "forward.rapidnet.in")
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

	// Metadata
	Tags  []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Notes string   `bson:"notes,omitempty" json:"notes,omitempty"`

	Version int `bson:"version" json:"version"`
}

type NASDeletedEvent struct {
	ID    string `bson:"_id,omitempty" json:"id"`
	NASId string `bson:"nasId" json:"nasId"` // Unique NAS identifier (e.g., NAS001, AP-CAFE-01)

	Version int `bson:"version" json:"version"`
}
