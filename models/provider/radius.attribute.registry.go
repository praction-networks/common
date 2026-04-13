package provider

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// RadiusAttributeDataType defines the expected data type of a RADIUS attribute value.
type RadiusAttributeDataType string

const (
	RadiusDataTypeString  RadiusAttributeDataType = "string"
	RadiusDataTypeNumber  RadiusAttributeDataType = "number"
	RadiusDataTypeIP      RadiusAttributeDataType = "ip_address"
	RadiusDataTypeBoolean RadiusAttributeDataType = "boolean"
	RadiusDataTypeEnum    RadiusAttributeDataType = "enum"
	RadiusDataTypeInteger RadiusAttributeDataType = "integer"
	RadiusDataTypeIPv6    RadiusAttributeDataType = "ipv6_address"
)

// RadiusAttributeCategory groups attributes by functional area.
type RadiusAttributeCategory string

const (
	RadiusCategoryAuthentication RadiusAttributeCategory = "authentication"
	RadiusCategoryAuthorization  RadiusAttributeCategory = "authorization"
	RadiusCategoryAccounting     RadiusAttributeCategory = "accounting"
	RadiusCategoryQoS            RadiusAttributeCategory = "qos"
	RadiusCategoryFiltering      RadiusAttributeCategory = "filtering"
	RadiusCategoryRouting        RadiusAttributeCategory = "routing"
	RadiusCategoryTunneling      RadiusAttributeCategory = "tunneling"
	RadiusCategoryVendorSpecific RadiusAttributeCategory = "vendor_specific"
	RadiusCategoryGeneral        RadiusAttributeCategory = "general"
)



// RadiusAttributeSchema defines a single RADIUS attribute and its validation rules.
type RadiusAttributeSchema struct {
	Key         string                  `json:"key"`
	AttributeID int                     `json:"attributeId,omitempty"` // Numeric RADIUS attribute ID
	Label       string                  `json:"label"`
	Description string                  `json:"description"`
	DataType    RadiusAttributeDataType `json:"dataType"`
	Category    RadiusAttributeCategory `json:"category,omitempty"`
	Placeholder string                  `json:"placeholder,omitempty"`
	Required    bool                    `json:"required,omitempty"`
	Options     []FieldOption           `json:"options,omitempty"` // Used if DataType is 'enum'
	Pattern     string                  `json:"pattern,omitempty"` // Regex pattern for validation
	MinValue    *float64                `json:"minValue,omitempty"`
	MaxValue    *float64                `json:"maxValue,omitempty"`
	MinLength   *int                    `json:"minLength,omitempty"`
	MaxLength   *int                    `json:"maxLength,omitempty"`
	Examples    []string                `json:"examples,omitempty"` // Example values
	Deprecated  bool                    `json:"deprecated,omitempty"`
}

// VendorDictionaryInfo holds a collection of RADIUS attributes for a specific vendor dictionary.
type VendorDictionaryInfo struct {
	Value       string                  `json:"value"` // e.g., "MIKROTIK"
	Label       string                  `json:"label"` // e.g., "MikroTik"
	Description string                  `json:"description"`
	VendorID    int                     `json:"vendorId,omitempty"` // Vendor-Specific attribute ID
	Attributes  []RadiusAttributeSchema `json:"attributes"`
}

// RadiusDictionaryRegistry is the single source of truth for all supported RADIUS dictionaries.
var RadiusDictionaryRegistry = map[string]VendorDictionaryInfo{
	"STANDARD": {
		Value:       "STANDARD",
		Label:       "Standard RADIUS",
		Description: "Standard RADIUS parameters defined by RFCs (RFC 2865, RFC 2866, RFC 2869).",
		Attributes: []RadiusAttributeSchema{
			{
				Key:         "User-Name",
				AttributeID: 1,
				Label:       "User Name",
				Description: "The name of the user being authenticated",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    true,
				MinLength:   intPtr(1),
				MaxLength:   intPtr(253),
				Examples:    []string{"john.doe@example.com", "user123"},
			},
			{
				Key:         "User-Password",
				AttributeID: 2,
				Label:       "User Password",
				Description: "The password of the user (encrypted)",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinLength:   intPtr(1),
				MaxLength:   intPtr(128),
			},
			{
				Key:         "CHAP-Password",
				AttributeID: 3,
				Label:       "CHAP Password",
				Description: "The CHAP password for CHAP authentication",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinLength:   intPtr(17),
				MaxLength:   intPtr(17),
			},
			{
				Key:         "NAS-IP-Address",
				AttributeID: 4,
				Label:       "NAS IP Address",
				Description: "The IP address of the Network Access Server",
				DataType:    RadiusDataTypeIP,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Examples:    []string{"192.168.1.1", "10.0.0.1"},
			},
			{
				Key:         "NAS-Port",
				AttributeID: 5,
				Label:       "NAS Port",
				Description: "The physical port number of the NAS",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				MinValue:    float64Ptr(0),
				MaxValue:    float64Ptr(65535),
				Examples:    []string{"0", "5000"},
			},
			{
				Key:         "Service-Type",
				AttributeID: 6,
				Label:       "Service Type",
				Description: "The type of service the user is requesting",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "1", Label: "Login"},
					{Value: "2", Label: "Framed"},
					{Value: "3", Label: "Callback Login"},
					{Value: "4", Label: "Callback Framed"},
					{Value: "5", Label: "Outbound"},
					{Value: "6", Label: "Administrative"},
					{Value: "7", Label: "NAS Prompt"},
					{Value: "8", Label: "Authenticate Only"},
					{Value: "9", Label: "Callback NAS Prompt"},
					{Value: "10", Label: "Call Check"},
					{Value: "11", Label: "Callback Administrative"},
					{Value: "12", Label: "Framed-Management"},
				},
			},
			{
				Key:         "Framed-Protocol",
				AttributeID: 7,
				Label:       "Framed Protocol",
				Description: "The framing protocol to be used",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "1", Label: "PPP"},
					{Value: "2", Label: "SLIP"},
					{Value: "3", Label: "ARAP"},
					{Value: "4", Label: "Gandalf SLIP"},
					{Value: "5", Label: "Xylogics IPX SLIP"},
					{Value: "6", Label: "X.75 Synchronous"},
				},
			},
			{
				Key:         "Framed-IP-Address",
				AttributeID: 8,
				Label:       "Framed IP Address",
				Description: "Address to be configured for the user (255.255.255.254 indicates NAS should select)",
				DataType:    RadiusDataTypeIP,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"192.168.1.100", "10.0.0.50"},
			},
			{
				Key:         "Framed-IP-Netmask",
				AttributeID: 9,
				Label:       "Framed IP Netmask",
				Description: "Netmask to be configured for the user",
				DataType:    RadiusDataTypeIP,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"255.255.255.0", "255.255.0.0"},
			},
			{
				Key:         "Framed-Routing",
				AttributeID: 10,
				Label:       "Framed Routing",
				Description: "Routing method for the user",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryRouting,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "None"},
					{Value: "1", Label: "Send routing packets"},
					{Value: "2", Label: "Listen for routing packets"},
					{Value: "3", Label: "Send and Listen"},
				},
			},
			{
				Key:         "Filter-Id",
				AttributeID: 11,
				Label:       "Filter ID",
				Description: "Filter list to be applied to the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryFiltering,
				Required:    false,
				Examples:    []string{"filter.in", "filter.out"},
			},
			{
				Key:         "Framed-MTU",
				AttributeID: 12,
				Label:       "Framed MTU",
				Description: "Maximum Transmission Unit for the user",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(64),
				MaxValue:    float64Ptr(65535),
				Examples:    []string{"1500", "1492"},
			},
			{
				Key:         "Framed-Compression",
				AttributeID: 13,
				Label:       "Framed Compression",
				Description: "Compression protocol to be used",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "None"},
					{Value: "1", Label: "VJ TCP/IP header compression"},
					{Value: "2", Label: "IPX header compression"},
					{Value: "3", Label: "Stac-LZS compression"},
				},
			},
			{
				Key:         "Login-IP-Host",
				AttributeID: 14,
				Label:       "Login IP Host",
				Description: "System with which to connect the user",
				DataType:    RadiusDataTypeIP,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"192.168.1.10"},
			},
			{
				Key:         "Login-Service",
				AttributeID: 15,
				Label:       "Login Service",
				Description: "Service to use for connecting the user",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "Telnet"},
					{Value: "1", Label: "Rlogin"},
					{Value: "2", Label: "TCP-Clear"},
					{Value: "3", Label: "PortMaster"},
					{Value: "4", Label: "LAT"},
					{Value: "5", Label: "X25-PAD"},
					{Value: "6", Label: "X25-T3POS"},
					{Value: "7", Label: "Unassigned"},
					{Value: "8", Label: "TCP"},
				},
			},
			{
				Key:         "Login-TCP-Port",
				AttributeID: 16,
				Label:       "Login TCP Port",
				Description: "TCP port to which to connect the user",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				MaxValue:    float64Ptr(65535),
				Examples:    []string{"23", "22", "80"},
			},
			{
				Key:         "Reply-Message",
				AttributeID: 18,
				Label:       "Reply Message",
				Description: "Message to be displayed to the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				MaxLength:   intPtr(253),
				Examples:    []string{"Welcome to the network!", "Access granted"},
			},
			{
				Key:         "Callback-Number",
				AttributeID: 19,
				Label:       "Callback Number",
				Description: "Phone number to call back",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"1234567890"},
			},
			{
				Key:         "Callback-Id",
				AttributeID: 20,
				Label:       "Callback ID",
				Description: "Callback identifier",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Framed-Route",
				AttributeID: 22,
				Label:       "Framed Route",
				Description: "Routing information to be configured for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryRouting,
				Required:    false,
				Pattern:     `^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/\d{1,2}(\s+[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(\s+\d+)?)?$`,
				Examples:    []string{"192.168.1.0/24", "10.0.0.0/8 192.168.1.1 1"},
			},
			{
				Key:         "Framed-IPX-Network",
				AttributeID: 23,
				Label:       "Framed IPX Network",
				Description: "IPX network number to be configured",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryRouting,
				Required:    false,
				Pattern:     `^[0-9A-Fa-f]{1,8}$`,
				Examples:    []string{"ABCDEF12", "12345678"},
			},
			{
				Key:         "State",
				AttributeID: 24,
				Label:       "State",
				Description: "State information to be maintained between client and server",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Class",
				AttributeID: 25,
				Label:       "Class",
				Description: "Class information to be maintained between client and server",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Vendor-Specific",
				AttributeID: 26,
				Label:       "Vendor Specific",
				Description: "Vendor-specific attributes",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryVendorSpecific,
				Required:    false,
			},
			{
				Key:         "Session-Timeout",
				AttributeID: 27,
				Label:       "Session Timeout",
				Description: "Maximum number of seconds of service to be provided",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"3600", "7200"},
			},
			{
				Key:         "Idle-Timeout",
				AttributeID: 28,
				Label:       "Idle Timeout",
				Description: "Maximum number of consecutive seconds of idle connection allowed",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"300", "600"},
			},
			{
				Key:         "Termination-Action",
				AttributeID: 29,
				Label:       "Termination Action",
				Description: "Action to take when service is terminated",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "Default"},
					{Value: "1", Label: "RADIUS-Request"},
				},
			},
			{
				Key:         "Called-Station-Id",
				AttributeID: 30,
				Label:       "Called Station ID",
				Description: "Phone number or station ID that was called",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Examples:    []string{"1234567890", "00:11:22:33:44:55"},
			},
			{
				Key:         "Calling-Station-Id",
				AttributeID: 31,
				Label:       "Calling Station ID",
				Description: "Phone number or station ID that is calling",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Examples:    []string{"0987654321", "AA:BB:CC:DD:EE:FF"},
			},
			{
				Key:         "NAS-Identifier",
				AttributeID: 32,
				Label:       "NAS Identifier",
				Description: "String identifying the NAS",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Examples:    []string{"router1", "nas.example.com"},
			},
			{
				Key:         "Proxy-State",
				AttributeID: 33,
				Label:       "Proxy State",
				Description: "State information to be maintained between proxy and server",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Login-LAT-Service",
				AttributeID: 34,
				Label:       "Login LAT Service",
				Description: "System with which the user is to be connected by LAT",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Login-LAT-Node",
				AttributeID: 35,
				Label:       "Login LAT Node",
				Description: "Node with which the user is to be connected by LAT",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Login-LAT-Group",
				AttributeID: 36,
				Label:       "Login LAT Group",
				Description: "LAT group codes which this user is authorized to use",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Framed-AppleTalk-Link",
				AttributeID: 37,
				Label:       "Framed AppleTalk Link",
				Description: "AppleTalk network number for the user",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				MaxValue:    float64Ptr(65535),
			},
			{
				Key:         "Framed-AppleTalk-Network",
				AttributeID: 38,
				Label:       "Framed AppleTalk Network",
				Description: "AppleTalk network number for the user",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Framed-AppleTalk-Zone",
				AttributeID: 39,
				Label:       "Framed AppleTalk Zone",
				Description: "AppleTalk default zone for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Acct-Status-Type",
				AttributeID: 40,
				Label:       "Acct Status Type",
				Description: "Type of accounting message",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				Options: []FieldOption{
					{Value: "1", Label: "Start"},
					{Value: "2", Label: "Stop"},
					{Value: "3", Label: "Interim-Update"},
					{Value: "4", Label: "On"},
					{Value: "5", Label: "Off"},
					{Value: "6", Label: "Tunnel-Start"},
					{Value: "7", Label: "Tunnel-Stop"},
					{Value: "8", Label: "Tunnel-Reject"},
					{Value: "9", Label: "Tunnel-Link-Start"},
					{Value: "10", Label: "Tunnel-Link-Stop"},
					{Value: "11", Label: "Tunnel-Link-Reject"},
					{Value: "12", Label: "Failed"},
				},
			},
			{
				Key:         "Acct-Delay-Time",
				AttributeID: 41,
				Label:       "Acct Delay Time",
				Description: "Time in seconds that the client has been trying to send this request",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Acct-Input-Octets",
				AttributeID: 42,
				Label:       "Acct Input Octets",
				Description: "Number of octets received from the port",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Acct-Output-Octets",
				AttributeID: 43,
				Label:       "Acct Output Octets",
				Description: "Number of octets sent to the port",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Acct-Session-Id",
				AttributeID: 44,
				Label:       "Acct Session ID",
				Description: "Unique accounting session ID",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAccounting,
				Required:    false,
			},
			{
				Key:         "Acct-Authentic",
				AttributeID: 45,
				Label:       "Acct Authentic",
				Description: "How the user was authenticated",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				Options: []FieldOption{
					{Value: "1", Label: "RADIUS"},
					{Value: "2", Label: "Local"},
					{Value: "3", Label: "Remote"},
					{Value: "4", Label: "Diameter"},
				},
			},
			{
				Key:         "Acct-Session-Time",
				AttributeID: 46,
				Label:       "Acct Session Time",
				Description: "Number of seconds the user has been logged in",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Acct-Input-Packets",
				AttributeID: 47,
				Label:       "Acct Input Packets",
				Description: "Number of packets received from the port",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Acct-Output-Packets",
				AttributeID: 48,
				Label:       "Acct Output Packets",
				Description: "Number of packets sent to the port",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Acct-Terminate-Cause",
				AttributeID: 49,
				Label:       "Acct Terminate Cause",
				Description: "How the session was terminated",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				Options: []FieldOption{
					{Value: "1", Label: "User Request"},
					{Value: "2", Label: "Lost Carrier"},
					{Value: "3", Label: "Lost Service"},
					{Value: "4", Label: "Idle Timeout"},
					{Value: "5", Label: "Session Timeout"},
					{Value: "6", Label: "Admin Reset"},
					{Value: "7", Label: "Admin Reboot"},
					{Value: "8", Label: "Port Error"},
					{Value: "9", Label: "NAS Error"},
					{Value: "10", Label: "NAS Request"},
					{Value: "11", Label: "NAS Reboot"},
					{Value: "12", Label: "Port Unneeded"},
					{Value: "13", Label: "Port Preempted"},
					{Value: "14", Label: "Port Suspended"},
					{Value: "15", Label: "Service Unavailable"},
					{Value: "16", Label: "Callback"},
					{Value: "17", Label: "User Error"},
					{Value: "18", Label: "Host Request"},
				},
			},
			{
				Key:         "Acct-Multi-Session-Id",
				AttributeID: 50,
				Label:       "Acct Multi Session ID",
				Description: "Unique ID to link together multiple sessions",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAccounting,
				Required:    false,
			},
			{
				Key:         "Acct-Link-Count",
				AttributeID: 51,
				Label:       "Acct Link Count",
				Description: "Count of links in a multilink session",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Acct-Input-Gigawords",
				AttributeID: 52,
				Label:       "Acct Input Gigawords",
				Description: "Number of times the Acct-Input-Octets counter has wrapped around 2^32",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Acct-Output-Gigawords",
				AttributeID: 53,
				Label:       "Acct Output Gigawords",
				Description: "Number of times the Acct-Output-Octets counter has wrapped around 2^32",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Event-Timestamp",
				AttributeID: 55,
				Label:       "Event Timestamp",
				Description: "Timestamp when the event occurred",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "NAS-Port-Type",
				AttributeID: 61,
				Label:       "NAS Port Type",
				Description: "Type of physical port on the NAS",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "Async"},
					{Value: "1", Label: "Sync"},
					{Value: "2", Label: "ISDN Sync"},
					{Value: "3", Label: "ISDN Async V.120"},
					{Value: "4", Label: "ISDN Async V.110"},
					{Value: "5", Label: "Virtual"},
					{Value: "6", Label: "PIAFS"},
					{Value: "7", Label: "HDLC Clear Channel"},
					{Value: "8", Label: "X.25"},
					{Value: "9", Label: "X.75"},
					{Value: "10", Label: "G.3 Fax"},
					{Value: "11", Label: "SDSL - Symmetric DSL"},
					{Value: "12", Label: "ADSL-CAP - ADSL Carrierless Amplitude Phase"},
					{Value: "13", Label: "ADSL-DMT - ADSL Discrete Multi-Tone"},
					{Value: "14", Label: "IDSL - ISDN Digital Subscriber Line"},
					{Value: "15", Label: "Ethernet"},
					{Value: "16", Label: "xDSL - Digital Subscriber Line of unknown type"},
					{Value: "17", Label: "Cable"},
					{Value: "18", Label: "Wireless - Other"},
					{Value: "19", Label: "Wireless - IEEE 802.11"},
					{Value: "20", Label: "Token-Ring"},
					{Value: "21", Label: "FDDI"},
					{Value: "22", Label: "Wireless - CDMA2000"},
					{Value: "23", Label: "Wireless - UMTS"},
					{Value: "24", Label: "Wireless - 1X-EV"},
					{Value: "25", Label: "IAPP"},
					{Value: "26", Label: "FTTP - Fiber to the Premises"},
					{Value: "27", Label: "Wireless - IEEE 802.16"},
					{Value: "28", Label: "Wireless - IEEE 802.20"},
					{Value: "29", Label: "Wireless - IEEE 802.22"},
					{Value: "30", Label: "PPPoA - PPP over ATM"},
					{Value: "31", Label: "PPPoEoA - PPP over Ethernet over ATM"},
					{Value: "32", Label: "PPPoEoE - PPP over Ethernet over Ethernet"},
					{Value: "33", Label: "PPPoEoVLAN - PPP over Ethernet over VLAN"},
					{Value: "34", Label: "PPPoEoQinQ - PPP over Ethernet over QinQ"},
					{Value: "35", Label: "xPON - Passive Optical Network"},
					{Value: "36", Label: "Wireless - XGP"},
					{Value: "37", Label: "WiMAX Pre-Release 8"},
					{Value: "38", Label: "WIMAX Release 8"},
					{Value: "39", Label: "WIMAX 2"},
					{Value: "40", Label: "CAPWAP"},
					{Value: "41", Label: "WIMAX Relay"},
				},
			},
			{
				Key:         "Connect-Info",
				AttributeID: 77,
				Label:       "Connect Info",
				Description: "Information about the connection",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Examples:    []string{"100Mbps", "1Gbps Full Duplex"},
			},
			{
				Key:         "NAS-Port-Id",
				AttributeID: 87,
				Label:       "NAS Port ID",
				Description: "Text identifier for the NAS port",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Examples:    []string{"eth0", "ge-0/0/0", "FastEthernet0/1"},
			},
			{
				Key:         "Framed-Pool",
				AttributeID: 88,
				Label:       "Framed Pool",
				Description: "IP pool from which to assign an IP address",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"pool1", "default-pool"},
			},
			{
				Key:         "NAS-IPv6-Address",
				AttributeID: 95,
				Label:       "NAS IPv6 Address",
				Description: "The IPv6 address of the Network Access Server",
				DataType:    RadiusDataTypeIPv6,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Examples:    []string{"2001:db8::1"},
			},
			{
				Key:         "Framed-Interface-Id",
				AttributeID: 96,
				Label:       "Framed Interface ID",
				Description: "IPv6 interface identifier for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Pattern:     `^[0-9a-fA-F:]+$`,
				Examples:    []string{"0:0:0:1"},
			},
			{
				Key:         "Framed-IPv6-Prefix",
				AttributeID: 97,
				Label:       "Framed IPv6 Prefix",
				Description: "IPv6 prefix to be configured for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Pattern:     `^[0-9a-fA-F:]+/\d{1,3}$`,
				Examples:    []string{"2001:db8::/64"},
			},
			{
				Key:         "Login-IPv6-Host",
				AttributeID: 98,
				Label:       "Login IPv6 Host",
				Description: "IPv6 system with which to connect the user",
				DataType:    RadiusDataTypeIPv6,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"2001:db8::10"},
			},
			{
				Key:         "Framed-IPv6-Route",
				AttributeID: 99,
				Label:       "Framed IPv6 Route",
				Description: "IPv6 routing information to be configured for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryRouting,
				Required:    false,
				Pattern:     `^[0-9a-fA-F:]+/\d{1,3}(\s+[0-9a-fA-F:]+(\s+\d+)?)?$`,
				Examples:    []string{"2001:db8:1::/48", "2001:db8:2::/48 2001:db8::1 1"},
			},
			{
				Key:         "Framed-IPv6-Pool",
				AttributeID: 100,
				Label:       "Framed IPv6 Pool",
				Description: "IPv6 pool from which to assign a prefix",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"ipv6-pool1"},
			},
			{
				Key:         "Delegated-IPv6-Prefix",
				AttributeID: 123,
				Label:       "Delegated IPv6 Prefix",
				Description: "IPv6 prefix to delegate to user's local network (DHCPv6-PD)",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Pattern:     `^[0-9a-fA-F:]+/\d{1,3}$`,
				Examples:    []string{"2001:db8:a::/48"},
			},
			{
				Key:         "Tunnel-Type",
				AttributeID: 64,
				Label:       "Tunnel Type",
				Description: "The tunneling protocol being used",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryTunneling,
				Required:    false,
				Options: []FieldOption{
					{Value: "1", Label: "PPTP"},
					{Value: "2", Label: "L2TP"},
					{Value: "3", Label: "L2F"},
					{Value: "4", Label: "ATMP"},
					{Value: "5", Label: "VLAN"},
					{Value: "6", Label: "IP"},
					{Value: "7", Label: "IP-in-IP"},
					{Value: "8", Label: "GRE"},
					{Value: "9", Label: "IPSec"},
					{Value: "10", Label: "VXLAN"},
					{Value: "11", Label: "NVGRE"},
					{Value: "12", Label: "MPLS"},
					{Value: "13", Label: "MPLS-in-IP"},
					{Value: "14", Label: "MPLS-in-GRE"},
					{Value: "15", Label: "IPSec-in-GRE"},
					{Value: "16", Label: "Wireless (WLAN)"},
					{Value: "17", Label: "Wireless (WPAN)"},
				},
			},
			{
				Key:         "Tunnel-Medium-Type",
				AttributeID: 65,
				Label:       "Tunnel Medium Type",
				Description: "The transport medium for the tunnel",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryTunneling,
				Required:    false,
				Options: []FieldOption{
					{Value: "1", Label: "IPv4"},
					{Value: "2", Label: "IPv6"},
					{Value: "3", Label: "NSAP"},
					{Value: "4", Label: "HDLC"},
					{Value: "5", Label: "BBN 1822"},
					{Value: "6", Label: "802"},
					{Value: "7", Label: "E.163"},
					{Value: "8", Label: "E.164"},
					{Value: "9", Label: "F.69"},
					{Value: "10", Label: "X.121"},
					{Value: "11", Label: "IPX"},
					{Value: "12", Label: "Appletalk"},
					{Value: "13", Label: "Decnet IV"},
					{Value: "14", Label: "Banyan Vines"},
					{Value: "15", Label: "E.164 with NSAP subaddress"},
				},
			},
			{
				Key:         "Tunnel-Client-Endpoint",
				AttributeID: 66,
				Label:       "Tunnel Client Endpoint",
				Description: "IP address of the tunnel client",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryTunneling,
				Required:    false,
				Examples:    []string{"192.168.1.1", "2001:db8::1"},
			},
			{
				Key:         "Tunnel-Server-Endpoint",
				AttributeID: 67,
				Label:       "Tunnel Server Endpoint",
				Description: "IP address of the tunnel server",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryTunneling,
				Required:    false,
				Examples:    []string{"10.0.0.1", "2001:db8::2"},
			},
			{
				Key:         "Tunnel-Password",
				AttributeID: 69,
				Label:       "Tunnel Password",
				Description: "Password for authenticating to the tunnel server",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryTunneling,
				Required:    false,
			},
			{
				Key:         "Tunnel-Private-Group-Id",
				AttributeID: 81,
				Label:       "Tunnel Private Group ID",
				Description: "Group ID for the tunnel",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryTunneling,
				Required:    false,
				Examples:    []string{"100", "VLAN100"},
			},
		},
	},

	"MIKROTIK": {
		Value:       "MIKROTIK",
		Label:       "MikroTik",
		Description: "MikroTik RouterOS specific RADIUS attributes (Vendor ID: 14988).",
		VendorID:    14988,
		Attributes: []RadiusAttributeSchema{
			{
				Key:         "Mikrotik-Rate-Limit",
				Label:       "Rate Limit",
				Description: "Bandwidth limit for the user. Format: rx-rate[/tx-rate] [rx-burst-rate[/tx-burst-rate] [rx-burst-threshold[/tx-burst-threshold] [rx-burst-time[/tx-burst-time] [priority] [rx-rate-min[/tx-rate-min]]]]",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryQoS,
				Required:    false,
				Pattern:     `^\d+[kKmMgG](\/\d+[kKmMgG])?(\s+\d+[kKmMgG](\/\d+[kKmMgG])?(\s+\d+[kKmMgG](\/\d+[kKmMgG])?(\s+\d+[kKmMgG](\/\d+[kKmMgG])?(\s+\d+(\s+\d+[kKmMgG](\/\d+[kKmMgG])?)?)?)?)?$`,
				Examples:    []string{"10M/10M", "5M/10M 10M/20M 5M/10M 10s 5", "1M/2M 2M/4M 1M/2M 10s 8 512k/1M"},
			},
			{
				Key:         "Mikrotik-Address-List",
				Label:       "Address List",
				Description: "Firewall address list to which user is added",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryFiltering,
				Required:    false,
				Examples:    []string{"allowed", "blocked", "premium"},
			},
			{
				Key:         "Mikrotik-Group",
				Label:       "User Group",
				Description: "MikroTik User group name",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"default", "premium", "trial"},
			},
			{
				Key:         "Mikrotik-Wireless-Enc-Algo",
				Label:       "Wireless Encryption Algorithm",
				Description: "Wireless encryption algorithm for the user",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				Options: []FieldOption{
					{Value: "none", Label: "None"},
					{Value: "40bit-wep", Label: "40-bit WEP"},
					{Value: "104bit-wep", Label: "104-bit WEP"},
					{Value: "aes-ccm", Label: "AES-CCM"},
					{Value: "tkip", Label: "TKIP"},
				},
			},
			{
				Key:         "Mikrotik-Wireless-Enc-Key",
				Label:       "Wireless Encryption Key",
				Description: "Wireless encryption key for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinLength:   intPtr(1),
				MaxLength:   intPtr(64),
			},
			{
				Key:         "Mikrotik-Wireless-Psk",
				Label:       "Wireless PSK",
				Description: "Wireless Pre-Shared Key for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinLength:   intPtr(8),
				MaxLength:   intPtr(64),
			},
			{
				Key:         "Mikrotik-Wireless-Vlan-Id",
				Label:       "Wireless VLAN ID",
				Description: "VLAN ID to assign to wireless client",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryTunneling,
				Required:    false,
				MinValue:    float64Ptr(1),
				MaxValue:    float64Ptr(4094),
				Examples:    []string{"100", "200"},
			},
			{
				Key:         "Mikrotik-Wireless-Vlan-Psk",
				Label:       "Wireless VLAN PSK",
				Description: "PSK for VLAN wireless client",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinLength:   intPtr(8),
				MaxLength:   intPtr(64),
			},
			{
				Key:         "Mikrotik-Wireless-Vlan-Enc-Algo",
				Label:       "Wireless VLAN Encryption Algorithm",
				Description: "Encryption algorithm for VLAN wireless client",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				Options: []FieldOption{
					{Value: "none", Label: "None"},
					{Value: "aes-ccm", Label: "AES-CCM"},
					{Value: "tkip", Label: "TKIP"},
				},
			},
			{
				Key:         "Mikrotik-Wireless-Vlan-Enc-Key",
				Label:       "Wireless VLAN Encryption Key",
				Description: "Encryption key for VLAN wireless client",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinLength:   intPtr(1),
				MaxLength:   intPtr(64),
			},
			{
				Key:         "Mikrotik-Wireless-Multicast-Key",
				Label:       "Wireless Multicast Key",
				Description: "Multicast encryption key for wireless",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinLength:   intPtr(1),
				MaxLength:   intPtr(64),
			},
			{
				Key:         "Mikrotik-Wireless-Multicast-Encryption",
				Label:       "Wireless Multicast Encryption",
				Description: "Multicast encryption type for wireless",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				Options: []FieldOption{
					{Value: "none", Label: "None"},
					{Value: "aes-ccm", Label: "AES-CCM"},
					{Value: "tkip", Label: "TKIP"},
				},
			},
			{
				Key:         "Mikrotik-Wireless-Fast-Transition",
				Label:       "Wireless Fast Transition",
				Description: "Enable 802.11r Fast Transition for the client",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
			},
			{
				Key:         "Mikrotik-Wireless-Pmk-Caching",
				Label:       "Wireless PMK Caching",
				Description: "Enable PMK caching for the client",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
			},
			{
				Key:         "Mikrotik-Wireless-Group-Key-Update",
				Label:       "Wireless Group Key Update",
				Description: "Group key update interval in seconds",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinValue:    float64Ptr(60),
				MaxValue:    float64Ptr(86400),
				Examples:    []string{"3600", "7200"},
			},
			{
				Key:         "Mikrotik-Wireless-Preauth-Timeout",
				Label:       "Wireless Preauth Timeout",
				Description: "Preauthentication timeout in seconds",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinValue:    float64Ptr(1),
				MaxValue:    float64Ptr(3600),
				Examples:    []string{"30", "60"},
			},
			{
				Key:         "Mikrotik-Realm",
				Label:       "Realm",
				Description: "Realm name for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Examples:    []string{"example.com", "isp.net"},
			},
			{
				Key:         "Mikrotik-Recv-Limit",
				Label:       "Receive Limit",
				Description: "Total receive limit in bytes (0 = unlimited)",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryQoS,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"10737418240", "5368709120"},
			},
			{
				Key:         "Mikrotik-Xmit-Limit",
				Label:       "Transmit Limit",
				Description: "Total transmit limit in bytes (0 = unlimited)",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryQoS,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"10737418240", "5368709120"},
			},
			{
				Key:         "Mikrotik-Recv-Limit-Gigawords",
				Label:       "Receive Limit Gigawords",
				Description: "High 32 bits of receive limit",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryQoS,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Mikrotik-Xmit-Limit-Gigawords",
				Label:       "Transmit Limit Gigawords",
				Description: "High 32 bits of transmit limit",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryQoS,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Mikrotik-Total-Limit",
				Label:       "Total Limit",
				Description: "Total (tx+rx) limit in bytes (0 = unlimited)",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryQoS,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"21474836480", "10737418240"},
			},
			{
				Key:         "Mikrotik-Total-Limit-Gigawords",
				Label:       "Total Limit Gigawords",
				Description: "High 32 bits of total limit",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryQoS,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Mikrotik-Session-Timeout",
				Label:       "Session Timeout",
				Description: "Maximum session time in seconds",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"3600", "7200"},
			},
			{
				Key:         "Mikrotik-Idle-Timeout",
				Label:       "Idle Timeout",
				Description: "Maximum idle time in seconds",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"300", "600"},
			},
			{
				Key:         "Mikrotik-Wireless-Forward",
				Label:       "Wireless Forward",
				Description: "Enable wireless forwarding for the client",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Mikrotik-Wireless-Client-To-Client-Forwarding",
				Label:       "Wireless Client-to-Client Forwarding",
				Description: "Enable client-to-client forwarding",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Mikrotik-Wireless-Encap-Session-Timeout",
				Label:       "Wireless Encap Session Timeout",
				Description: "Encapsulation session timeout in seconds",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Mikrotik-Wireless-Encap-Type",
				Label:       "Wireless Encap Type",
				Description: "Wireless encapsulation type",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "None"},
					{Value: "1", Label: "WPA"},
					{Value: "2", Label: "WPA2"},
				},
			},
			{
				Key:         "Mikrotik-Wireless-Mac-Format",
				Label:       "Wireless MAC Format",
				Description: "Wireless MAC address format",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "Default"},
					{Value: "1", Label: "XX:XX:XX:XX:XX:XX"},
					{Value: "2", Label: "XX-XX-XX-XX-XX-XX"},
				},
			},
			{
				Key:         "Mikrotik-Wireless-Skip-802.1x",
				Label:       "Wireless Skip 802.1x",
				Description: "Skip 802.1x authentication for this client",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
			},
			{
				Key:         "Mikrotik-Wireless-Algo-Preauth",
				Label:       "Wireless Algorithm Preauth",
				Description: "Enable algorithm preauthentication",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
			},
			{
				Key:         "Mikrotik-Wireless-Protected-Frame",
				Label:       "Wireless Protected Frame",
				Description: "Enable protected frames (802.11w)",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
			},
			{
				Key:         "Mikrotik-Wireless-Sa-Query-Timeout",
				Label:       "Wireless SA Query Timeout",
				Description: "SA Query timeout in seconds",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinValue:    float64Ptr(1),
				MaxValue:    float64Ptr(100),
				Examples:    []string{"1", "2"},
			},
			{
				Key:         "Mikrotik-Wireless-Sa-Query-Retry",
				Label:       "Wireless SA Query Retry",
				Description: "SA Query retry count",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinValue:    float64Ptr(0),
				MaxValue:    float64Ptr(10),
				Examples:    []string{"2", "4"},
			},
			{
				Key:         "Mikrotik-Advertise-URL",
				Label:       "Advertise URL",
				Description: "URL to advertise to the client",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				Pattern:     `^https?://[^\s]+$`,
				Examples:    []string{"https://example.com", "http://portal.example.com"},
			},
			{
				Key:         "Mikrotik-Advertise-Interval",
				Label:       "Advertise Interval",
				Description: "Advertisement interval in seconds",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"300", "600"},
			},
			{
				Key:         "Mikrotik-Delegated-IPv6-Pool",
				Label:       "Delegated IPv6 Pool",
				Description: "IPv6 pool for prefix delegation",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"ipv6-pool1", "pd-pool"},
			},
			{
				Key:         "Mikrotik-Delegated-IPv6-Prefix",
				Label:       "Delegated IPv6 Prefix",
				Description: "IPv6 prefix to delegate",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Pattern:     `^[0-9a-fA-F:]+/\d{1,3}$`,
				Examples:    []string{"2001:db8::/64", "fd00::/48"},
			},
			{
				Key:         "Mikrotik-IPsec-Policy",
				Label:       "IPsec Policy",
				Description: "IPsec policy for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"policy1", "ipsec-policy"},
			},
			{
				Key:         "Mikrotik-IPsec-Secret",
				Label:       "IPsec Secret",
				Description: "IPsec pre-shared key",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinLength:   intPtr(1),
				MaxLength:   intPtr(128),
			},
			{
				Key:         "Mikrotik-PPPoe-Service-Name",
				Label:       "PPPoE Service Name",
				Description: "PPPoE service name",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"isp-service", "pppoe"},
			},
			{
				Key:         "Mikrotik-PPTP-Encoding",
				Label:       "PPTP Encoding",
				Description: "PPTP encoding type",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "None"},
					{Value: "1", Label: "MPPE 40-bit"},
					{Value: "2", Label: "MPPE 128-bit"},
				},
			},
			{
				Key:         "Mikrotik-L2TP-Encoding",
				Label:       "L2TP Encoding",
				Description: "L2TP encoding type",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "None"},
					{Value: "1", Label: "MPPE 40-bit"},
					{Value: "2", Label: "MPPE 128-bit"},
				},
			},
			{
				Key:         "Mikrotik-Wireless-Comment",
				Label:       "Wireless Comment",
				Description: "Comment for wireless client",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				MaxLength:   intPtr(255),
			},
		},
	},

	"CISCO": {
		Value:       "CISCO",
		Label:       "Cisco",
		Description: "Cisco specific RADIUS attributes (Vendor ID: 9).",
		VendorID:    9,
		Attributes: []RadiusAttributeSchema{
			{
				Key:         "Cisco-AVPair",
				Label:       "Cisco AVPair",
				Description: "Arbitrary AV-Pairs used for advanced Cisco configs. Format: attribute=value",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryVendorSpecific,
				Required:    false,
				Pattern:     `^[a-zA-Z0-9_-]+=.+$`,
				Examples:    []string{"client-mac-address=00:11:22:33:44:55", "ip:inacl#1=permit ip any any", "tunnel-type=1"},
			},
			{
				Key:         "Cisco-NAS-Port",
				Label:       "Cisco NAS Port",
				Description: "Cisco specific NAS port format",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Cisco-Multilink-ID",
				Label:       "Multilink ID",
				Description: "Multilink bundle ID",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Cisco-Multilink-Maxlinks",
				Label:       "Multilink Max Links",
				Description: "Maximum number of links in multilink bundle",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryGeneral,
				Required:    false,
				MinValue:    float64Ptr(1),
				MaxValue:    float64Ptr(255),
				Examples:    []string{"2", "4", "8"},
			},
			{
				Key:         "Cisco-Password-Expiry",
				Label:       "Password Expiry",
				Description: "Days until password expires",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"30", "60", "90"},
			},
			{
				Key:         "Cisco-Service-Info",
				Label:       "Service Info",
				Description: "Service information for the user",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"service1", "premium"},
			},
			{
				Key:         "Cisco-Command-Code",
				Label:       "Command Code",
				Description: "Command code to execute",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Cisco-Script-Path",
				Label:       "Script Path",
				Description: "Path to script to execute",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Cisco-Disconnect-Cause",
				Label:       "Disconnect Cause",
				Description: "Reason for disconnection",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Cisco-Data-Rate",
				Label:       "Data Rate",
				Description: "Data rate in bps",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryQoS,
				Required:    false,
				MinValue:    float64Ptr(0),
				Examples:    []string{"10000000", "100000000"},
			},
			{
				Key:         "Cisco-Presence",
				Label:       "Presence",
				Description: "Presence attribute",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Cisco-Call-Filter",
				Label:       "Call Filter",
				Description: "Call filter to apply",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryFiltering,
				Required:    false,
			},
			{
				Key:         "Cisco-Idlepat",
				Label:       "Idle Pattern",
				Description: "Idle pattern detection",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Cisco-Tunnel-Private-Group-ID",
				Label:       "Tunnel Private Group ID",
				Description: "Private group ID for tunneling",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryTunneling,
				Required:    false,
				Examples:    []string{"100", "VLAN100"},
			},
			{
				Key:         "Cisco-Tunnel-Preference",
				Label:       "Tunnel Preference",
				Description: "Tunnel preference value",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryTunneling,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Cisco-IP-Direct",
				Label:       "IP Direct",
				Description: "Enable IP direct routing",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryRouting,
				Required:    false,
			},
			{
				Key:         "Cisco-PPP-VJ-Slot-Comp",
				Label:       "PPP VJ Slot Compression",
				Description: "Enable Van Jacobson TCP/IP header compression",
				DataType:    RadiusDataTypeBoolean,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
			},
			{
				Key:         "Cisco-PPP-Async-Map",
				Label:       "PPP Async Map",
				Description: "PPP async control character map",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				MaxValue:    float64Ptr(4294967295),
			},
			{
				Key:         "Cisco-Packet-Type",
				Label:       "Packet Type",
				Description: "Packet type identifier",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Cisco-Num-In-Sessions",
				Label:       "Number of In Sessions",
				Description: "Number of incoming sessions",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAccounting,
				Required:    false,
				MinValue:    float64Ptr(0),
			},
			{
				Key:         "Cisco-Proxy-Acl",
				Label:       "Proxy ACL",
				Description: "Proxy access control list",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryFiltering,
				Required:    false,
			},
			{
				Key:         "Cisco-Authen-Method",
				Label:       "Authentication Method",
				Description: "Authentication method used",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "Not Set"},
					{Value: "1", Label: "None"},
					{Value: "2", Label: "KRB5"},
					{Value: "3", Label: "Line"},
					{Value: "4", Label: "Enable"},
					{Value: "5", Label: "Local"},
					{Value: "6", Label: "Tacacs+"},
					{Value: "7", Label: "RADIUS"},
				},
			},
			{
				Key:         "Cisco-Authen-Type",
				Label:       "Authentication Type",
				Description: "Type of authentication",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "Not Set"},
					{Value: "1", Label: "ASCII"},
					{Value: "2", Label: "PAP"},
					{Value: "3", Label: "CHAP"},
					{Value: "4", Label: "ARAP"},
					{Value: "5", Label: "MS-CHAP"},
					{Value: "6", Label: "MS-CHAPv2"},
					{Value: "7", Label: "EAP"},
				},
			},
			{
				Key:         "Cisco-Service-Type",
				Label:       "Service Type",
				Description: "Service type for authorization",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "None"},
					{Value: "1", Label: "Shell"},
					{Value: "2", Label: "TTY"},
					{Value: "3", Label: "Network"},
					{Value: "4", Label: "Fwd"},
					{Value: "5", Label: "Outbound"},
					{Value: "6", Label: "System"},
				},
			},
			{
				Key:         "Cisco-User-Privilege-Level",
				Label:       "User Privilege Level",
				Description: "Cisco privilege level (0-15)",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				MaxValue:    float64Ptr(15),
				Examples:    []string{"1", "15", "0"},
			},
			{
				Key:         "Cisco-IP-Pool-Definition",
				Label:       "IP Pool Definition",
				Description: "IP pool name to use",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Examples:    []string{"POOL1", "default-pool"},
			},
			{
				Key:         "Cisco-Assign-IP-Pool",
				Label:       "Assign IP Pool",
				Description: "IP pool assignment method",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "Local Pool"},
					{Value: "1", Label: "DHCP Proxy"},
				},
			},
			{
				Key:         "Cisco-Link-Compression",
				Label:       "Link Compression",
				Description: "Enable link compression",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "None"},
					{Value: "1", Label: "Stac"},
					{Value: "2", Label: "MPPC"},
					{Value: "3", Label: "Stac-LZS"},
				},
			},
			{
				Key:         "Cisco-Target-Utility",
				Label:       "Target Utility",
				Description: "Target utility for compression",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(0),
				MaxValue:    float64Ptr(255),
			},
			{
				Key:         "Cisco-Maximum-Channels",
				Label:       "Maximum Channels",
				Description: "Maximum number of channels",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryAuthorization,
				Required:    false,
				MinValue:    float64Ptr(1),
				Examples:    []string{"1", "2", "4"},
			},
			{
				Key:         "Cisco-Encryption-Type",
				Label:       "Encryption Type",
				Description: "Encryption type to use",
				DataType:    RadiusDataTypeEnum,
				Category:    RadiusCategoryAuthentication,
				Required:    false,
				Options: []FieldOption{
					{Value: "0", Label: "None"},
					{Value: "1", Label: "MD5"},
					{Value: "2", Label: "DES"},
				},
			},
			{
				Key:         "Cisco-Event",
				Label:       "Event",
				Description: "Event type",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
			{
				Key:         "Cisco-Account-Info",
				Label:       "Account Info",
				Description: "Accounting information",
				DataType:    RadiusDataTypeString,
				Category:    RadiusCategoryAccounting,
				Required:    false,
			},
			{
				Key:         "Cisco-Progress",
				Label:       "Progress",
				Description: "Progress indicator",
				DataType:    RadiusDataTypeInteger,
				Category:    RadiusCategoryGeneral,
				Required:    false,
			},
		},
	},
}

// GetRadiusDictionaryConfig returns the full list of supported vendors and their attribute schemas for frontend rendering.
func GetRadiusDictionaryConfig() []VendorDictionaryInfo {
	dictionaries := make([]VendorDictionaryInfo, 0, len(RadiusDictionaryRegistry))
	for _, info := range RadiusDictionaryRegistry {
		dictionaries = append(dictionaries, info)
	}
	return dictionaries
}

// ValidateRadiusAttributes validates a map of key-value attributes against the known registry schemas.
func ValidateRadiusAttributes(attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}

	// Build a flat lookup map for quick validation of any supported key
	schemaMap := make(map[string]RadiusAttributeSchema)
	for _, vendor := range RadiusDictionaryRegistry {
		for _, attr := range vendor.Attributes {
			schemaMap[attr.Key] = attr
		}
	}

	for key, rawValue := range attributes {
		schema, exists := schemaMap[key]
		if !exists {
			return fmt.Errorf("RADIUS attribute '%s' is not defined in any supported vendor dictionary", key)
		}

		if err := validateAttributeValue(key, rawValue, schema); err != nil {
			return err
		}
	}

	return nil
}

func validateAttributeValue(key string, rawValue interface{}, schema RadiusAttributeSchema) error {
	if rawValue == nil {
		if schema.Required {
			return fmt.Errorf("RADIUS attribute '%s' is required but missing", key)
		}
		return nil
	}

	switch schema.DataType {
	case RadiusDataTypeString, RadiusDataTypeIP, RadiusDataTypeIPv6:
		val, ok := rawValue.(string)
		if !ok {
			return fmt.Errorf("RADIUS attribute '%s' must be a string", key)
		}

		if schema.MinLength != nil && len(val) < *schema.MinLength {
			return fmt.Errorf("RADIUS attribute '%s' must be at least %d characters", key, *schema.MinLength)
		}
		if schema.MaxLength != nil && len(val) > *schema.MaxLength {
			return fmt.Errorf("RADIUS attribute '%s' must not exceed %d characters", key, *schema.MaxLength)
		}

		if schema.Pattern != "" {
			regex, err := regexp.Compile(schema.Pattern)
			if err == nil && !regex.MatchString(val) {
				return fmt.Errorf("RADIUS attribute '%s' value '%s' does not match required pattern", key, val)
			}
		}

		switch schema.DataType {
		case RadiusDataTypeIP:
			if net.ParseIP(val) == nil || strings.Contains(val, ":") {
				return fmt.Errorf("RADIUS attribute '%s' must be a valid IPv4 address", key)
			}
		case RadiusDataTypeIPv6:
			if net.ParseIP(val) == nil || !strings.Contains(val, ":") {
				return fmt.Errorf("RADIUS attribute '%s' must be a valid IPv6 address", key)
			}
		}

	case RadiusDataTypeInteger, RadiusDataTypeNumber:
		var val float64
		switch v := rawValue.(type) {
		case float64:
			val = v
		case int:
			val = float64(v)
		case int32:
			val = float64(v)
		case int64:
			val = float64(v)
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("RADIUS attribute '%s' must be a number", key)
			}
			val = parsed
		default:
			return fmt.Errorf("RADIUS attribute '%s' must be a number", key)
		}

		if schema.DataType == RadiusDataTypeInteger && val != float64(int64(val)) {
			return fmt.Errorf("RADIUS attribute '%s' must be an integer", key)
		}

		if schema.MinValue != nil && val < *schema.MinValue {
			return fmt.Errorf("RADIUS attribute '%s' must be at least %v", key, *schema.MinValue)
		}
		if schema.MaxValue != nil && val > *schema.MaxValue {
			return fmt.Errorf("RADIUS attribute '%s' must be at most %v", key, *schema.MaxValue)
		}

	case RadiusDataTypeBoolean:
		_, ok := rawValue.(bool)
		if !ok {
			valStr, okStr := rawValue.(string)
			if okStr && (strings.ToLower(valStr) == "true" || strings.ToLower(valStr) == "false") {
				return nil
			}
			return fmt.Errorf("RADIUS attribute '%s' must be a boolean", key)
		}

	case RadiusDataTypeEnum:
		valStr, ok := rawValue.(string)
		if !ok {
			return fmt.Errorf("RADIUS attribute '%s' must be a string (enum option)", key)
		}
		validOption := false
		for _, opt := range schema.Options {
			if opt.Value == valStr {
				validOption = true
				break
			}
		}
		if !validOption {
			return fmt.Errorf("RADIUS attribute '%s' value '%s' is not a valid option", key, valStr)
		}
	}

	return nil
}

// Helpers for schema definitions
func float64Ptr(v float64) *float64 { return &v }
func intPtr(v int) *int { return &v }
