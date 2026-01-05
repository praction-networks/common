package logengineevent

import "time"

// DeviceConfig represents a network device configuration event
type DeviceConfigEvent struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenantId"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`         // BNG, CGNAT, UNIFIED
	Vendor       string    `json:"vendor"`       // CISCO_ASR, JUNIPER_MX, HUAWEI, MIKROTIK, etc.
	IPAddress    string    `json:"ipAddress"`
	Protocols    []string  `json:"protocols"`    // SYSLOG, NETFLOW_V9, IPFIX
	SyslogPort   int       `json:"syslogPort"`
	NetFlowPort  int       `json:"netflowPort"`
	TopologyMode string    `json:"topologyMode"` // UNIFIED, SEPARATE
	IPv6Enabled  bool      `json:"ipv6Enabled"`
	Enabled      bool      `json:"enabled"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// DeviceType constants
const (
	DeviceTypeBNG     = "BNG"
	DeviceTypeCGNAT   = "CGNAT"
	DeviceTypeUnified = "UNIFIED"
)

// Vendor constants
const (
	VendorCiscoASR    = "CISCO_ASR"
	VendorJuniperMX   = "JUNIPER_MX"
	VendorHuawei      = "HUAWEI"
	VendorMikroTik    = "MIKROTIK"
	VendorA10Thunder  = "A10_THUNDER"
	VendorGeneric     = "GENERIC"
)

// Protocol constants
const (
	ProtocolSyslog    = "SYSLOG"
	ProtocolNetFlowV5 = "NETFLOW_V5"
	ProtocolNetFlowV9 = "NETFLOW_V9"
	ProtocolIPFIX     = "IPFIX"
)

// TopologyMode constants
const (
	TopologyUnified  = "UNIFIED"
	TopologySeparate = "SEPARATE"
)

