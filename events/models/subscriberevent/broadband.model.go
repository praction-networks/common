package subscriberevent

type BroadbandConnectionType string

const (
	BroadbandConnectionFTTH     BroadbandConnectionType = "FTTH"
	BroadbandConnectionWireless BroadbandConnectionType = "WIRELESS"
	BroadbandConnectionOther    BroadbandConnectionType = "OTHER"
)

type BroadbandSubscriptionStatus string

const (
	BroadbandStatusPendingInstall BroadbandSubscriptionStatus = "PENDING_INSTALL"
	BroadbandStatusActive         BroadbandSubscriptionStatus = "ACTIVE"
	BroadbandStatusSuspended      BroadbandSubscriptionStatus = "SUSPENDED"
	BroadbandStatusTerminated     BroadbandSubscriptionStatus = "TERMINATED"
	BroadbandStatusChurned        BroadbandSubscriptionStatus = "CHURNED"
	BroadbandStatusBlacklisted    BroadbandSubscriptionStatus = "BLACKLISTED"
	BroadbandStatusGrace          BroadbandSubscriptionStatus = "GRACE"
	BroadbandStatusExpired        BroadbandSubscriptionStatus = "EXPIRED"
	BroadbandStatusCancelled      BroadbandSubscriptionStatus = "CANCELLED"
)

type BroabandAuthMethod string

const (
	BroabandAuthMethodPPPoE  BroabandAuthMethod = "PPPOE"
	BroabandAuthMethodStatic BroabandAuthMethod = "STATIC_IP"
	BroabandAuthMethodDHCP   BroabandAuthMethod = "DHCP"
)

type BroadbandCPEMetadata struct {
	SerialNumber        string                `json:"serialNumber,omitempty" bson:"serialNumber,omitempty"`
	Model               string                `json:"model,omitempty" bson:"model,omitempty"`
	Manufacturer        string                `json:"manufacturer,omitempty" bson:"manufacturer,omitempty"`
	FirmwareVersion     string                `json:"firmwareVersion,omitempty" bson:"firmwareVersion,omitempty"`
	SoftwareVersion     string                `json:"softwareVersion,omitempty" bson:"softwareVersion,omitempty"`
	HardwareVersion     string                `json:"hardwareVersion,omitempty" bson:"hardwareVersion,omitempty"`
	MACAddress          string                `json:"macAddress,omitempty" bson:"macAddress,omitempty"`
	CPEWirelessMetadata []CPEWirelessMetadata `json:"cpeWirelessMetadata,omitempty" bson:"cpeWirelessMetadata,omitempty"`
}

type CPEWirelessMetadata struct {
	SSID               string `json:"ssid,omitempty" bson:"ssid,omitempty"`
	Password           string `json:"password,omitempty" bson:"password,omitempty"`
	Security           string `json:"security,omitempty" bson:"security,omitempty"`
	Channel            string `json:"channel,omitempty" bson:"channel,omitempty"`
	Frequency          string `json:"frequency,omitempty" bson:"frequency,omitempty"`
	SignalStrength     string `json:"signalStrength,omitempty" bson:"signalStrength,omitempty"`
	SignalToNoiseRatio string `json:"signalToNoiseRatio,omitempty" bson:"signalToNoiseRatio,omitempty"`
}

// PPPoEAuthConfig contains configuration for PPPoE authentication method
type PPPoEAuthConfig struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"` // PPPoE username
	Password string `json:"password,omitempty" bson:"password,omitempty"` // PPPoE password
	IP       string `json:"ip,omitempty" bson:"ip,omitempty"`             // Assigned IP address (if static)
	MAC      string `json:"mac,omitempty" bson:"mac,omitempty"`           // MAC address for authentication
}

// StaticIPAuthConfig contains configuration for Static IP authentication method
type StaticIPAuthConfig struct {
	IP         string   `json:"ip,omitempty" bson:"ip,omitempty"`                 // Static IP address
	MAC        string   `json:"mac,omitempty" bson:"mac,omitempty"`               // MAC address for authentication
	Gateway    string   `json:"gateway,omitempty" bson:"gateway,omitempty"`       // Gateway IP address
	SubnetMask string   `json:"subnetMask,omitempty" bson:"subnetMask,omitempty"` // Subnet mask
	DNS        []string `json:"dns,omitempty" bson:"dns,omitempty"`               // Primary DNS server
}

// DHCPAuthConfig contains configuration for DHCP authentication method
type DHCPAuthConfig struct {
	MAC string `json:"mac,omitempty" bson:"mac,omitempty"` // MAC address for DHCP authentication
}

// BroadbandAuthConfig is a union type that can hold any of the three authentication config types
// Based on the AuthMethod, only one of the fields should be populated
type BroadbandAuthConfig struct {
	PPPoE    *PPPoEAuthConfig    `json:"pppoe,omitempty" bson:"pppoe,omitempty"`       // PPPoE configuration (when AuthMethod is PPPOE)
	StaticIP *StaticIPAuthConfig `json:"staticIp,omitempty" bson:"staticIp,omitempty"` // Static IP configuration (when AuthMethod is STATIC_IP)
	DHCP     *DHCPAuthConfig     `json:"dhcp,omitempty" bson:"dhcp,omitempty"`         // DHCP configuration (when AuthMethod is DHCP)
}
