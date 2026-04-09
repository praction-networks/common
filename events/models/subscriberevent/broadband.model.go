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
	RadioBand          string `json:"radioBand,omitempty" bson:"radioBand,omitempty"`         // "2.4GHz", "5GHz", "6GHz"
	NetworkType        string `json:"networkType,omitempty" bson:"networkType,omitempty"`     // "Main", "Guest", "IoT"
	InterfaceName      string `json:"interfaceName,omitempty" bson:"interfaceName,omitempty"` // e.g., "wlan0", "wlan1"
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
	IPv4     string `json:"ipv4,omitempty" bson:"ipv4,omitempty"`         // Assigned IPv4 address
	IPv6     string `json:"ipv6,omitempty" bson:"ipv6,omitempty"`         // Assigned IPv6 network/prefix
	MAC      string `json:"mac,omitempty" bson:"mac,omitempty"`           // MAC address for authentication
}

// StaticIPAuthConfig contains configuration for Static IP authentication method
type StaticIPAuthConfig struct {
	IPv4           string   `json:"ipv4,omitempty" bson:"ipv4,omitempty"`                     // Assigned IPv4 address
	IPv6           string   `json:"ipv6,omitempty" bson:"ipv6,omitempty"`                     // Assigned IPv6 network/prefix
	MAC            string   `json:"mac,omitempty" bson:"mac,omitempty"`                       // MAC address for authentication
	GatewayIPv4    string   `json:"gatewayIpv4,omitempty" bson:"gatewayIpv4,omitempty"`       // Gateway IPv4 address
	GatewayIPv6    string   `json:"gatewayIpv6,omitempty" bson:"gatewayIpv6,omitempty"`       // Gateway IPv6 address
	SubnetMaskIPv4 string   `json:"subnetMaskIpv4,omitempty" bson:"subnetMaskIpv4,omitempty"` // IPv4 subnet mask
	SubnetMaskIPv6 string   `json:"subnetMaskIpv6,omitempty" bson:"subnetMaskIpv6,omitempty"` // IPv6 subnet prefix length
	DNSIPv4        []string `json:"dnsIpv4,omitempty" bson:"dnsIpv4,omitempty"`               // IPv4 DNS servers
	DNSIPv6        []string `json:"dnsIpv6,omitempty" bson:"dnsIpv6,omitempty"`               // IPv6 DNS servers
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
