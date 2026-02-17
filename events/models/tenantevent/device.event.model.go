package tenantevent

// DeviceInsertEventModel defines the model for device insert events
type DeviceInsertEventModel struct {
	ID                  string   `bson:"_id" json:"id"`
	NASId               string   `bson:"nasId" json:"nasId"`
	TenantIDs           []string `bson:"tenantIds" json:"tenantIds"`
	Name                string   `bson:"name" json:"name"`
	Description         string   `bson:"description,omitempty" json:"description,omitempty"`
	DeviceTypes         []string `bson:"deviceTypes" json:"deviceTypes"`
	Manufacturer        string   `bson:"manufacturer,omitempty" json:"manufacturer,omitempty"`
	Model               string   `bson:"model,omitempty" json:"model,omitempty"`
	IPAddress           string   `bson:"ipAddress" json:"ipAddress"`
	AllowedIPRanges     []string `bson:"allowedIPRanges,omitempty" json:"allowedIPRanges,omitempty"`
	SSIDs               []string `bson:"ssids" json:"ssids"`
	RedirectorDomain    string   `bson:"redirectorDomain" json:"redirectorDomain"`
	ForwardDomain       string   `bson:"forwardDomain" json:"forwardDomain"`
	RadiusIP            []string `bson:"radiusIp,omitempty" json:"radiusIp,omitempty"`
	Secret              string   `bson:"secret" json:"secret"`             // For HMAC/JWT signing and API authentication
	RADIUSSecret        string   `bson:"radiusSecret" json:"radiusSecret"` // For RADIUS protocol communication
	TokenExpirySeconds  int      `bson:"tokenExpirySeconds" json:"tokenExpirySeconds"`
	RequireMutualTLS    bool     `bson:"requireMutualTLS" json:"requireMutualTLS"`
	EnableRateLimit     bool     `bson:"enableRateLimit" json:"enableRateLimit"`
	MaxSessionsPerHour  int      `bson:"maxSessionsPerHour" json:"maxSessionsPerHour"`
	MaxDevicesPerMAC    int      `bson:"maxDevicesPerMAC" json:"maxDevicesPerMAC"`
	MaxRetriesPerDevice int      `bson:"maxRetriesPerDevice" json:"maxRetriesPerDevice"`
	IsActive            bool     `bson:"isActive" json:"isActive"`
	Tags                []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Notes               string   `bson:"notes,omitempty" json:"notes,omitempty"`
	Version             int      `bson:"version" json:"version"`
}

// DeviceUpdateEventModel defines the model for device update events
type DeviceUpdateEventModel struct {
	ID                  string   `bson:"_id" json:"id"`
	NASId               string   `bson:"nasId,omitempty" json:"nasId,omitempty"`
	TenantIDs           []string `bson:"tenantIds,omitempty" json:"tenantIds,omitempty"`
	Name                string   `bson:"name,omitempty" json:"name,omitempty"`
	Description         string   `bson:"description,omitempty" json:"description,omitempty"`
	DeviceTypes         []string `bson:"deviceTypes,omitempty" json:"deviceTypes,omitempty"`
	Manufacturer        string   `bson:"manufacturer,omitempty" json:"manufacturer,omitempty"`
	Model               string   `bson:"model,omitempty" json:"model,omitempty"`
	IPAddress           string   `bson:"ipAddress,omitempty" json:"ipAddress,omitempty"`
	AllowedIPRanges     []string `bson:"allowedIPRanges,omitempty" json:"allowedIPRanges,omitempty"`
	SSIDs               []string `bson:"ssids,omitempty" json:"ssids,omitempty"`
	RedirectorDomain    string   `bson:"redirectorDomain,omitempty" json:"redirectorDomain,omitempty"`
	ForwardDomain       string   `bson:"forwardDomain,omitempty" json:"forwardDomain,omitempty"`
	RadiusIP            []string `bson:"radiusIp,omitempty" json:"radiusIp,omitempty"`
	Secret              string   `bson:"secret,omitempty" json:"secret,omitempty"`             // For HMAC/JWT signing and API authentication
	RADIUSSecret        string   `bson:"radiusSecret,omitempty" json:"radiusSecret,omitempty"` // For RADIUS protocol communication
	TokenExpirySeconds  *int     `bson:"tokenExpirySeconds,omitempty" json:"tokenExpirySeconds,omitempty"`
	RequireMutualTLS    *bool    `bson:"requireMutualTLS,omitempty" json:"requireMutualTLS,omitempty"`
	EnableRateLimit     *bool    `bson:"enableRateLimit,omitempty" json:"enableRateLimit,omitempty"`
	MaxSessionsPerHour  *int     `bson:"maxSessionsPerHour,omitempty" json:"maxSessionsPerHour,omitempty"`
	MaxDevicesPerMAC    *int     `bson:"maxDevicesPerMAC,omitempty" json:"maxDevicesPerMAC,omitempty"`
	MaxRetriesPerDevice *int     `bson:"maxRetriesPerDevice,omitempty" json:"maxRetriesPerDevice,omitempty"`
	IsActive            *bool    `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Tags                []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Notes               string   `bson:"notes,omitempty" json:"notes,omitempty"`
	Version             int      `bson:"version" json:"version"`
}

// DeviceDeleteEventModel defines the model for device delete events
type DeviceDeleteEventModel struct {
	ID      string `bson:"_id" json:"id"`
	NASId   string `bson:"nasId" json:"nasId"`
	Version int    `bson:"version" json:"version"`
}
