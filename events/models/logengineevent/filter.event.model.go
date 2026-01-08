package logengineevent

import "time"

// PortExclusionConfigEvent represents a port exclusion configuration event
type PortExclusionConfigEvent struct {
	ID                  string              `json:"id"`
	TenantID            string              `json:"tenantId"`
	ExcludedPorts       []PortExclusionRule `json:"excludedPorts"`
	ApplySystemDefaults bool                `json:"applySystemDefaults"`
	Enabled             bool                `json:"enabled"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
}

// PortExclusionRule represents a single port exclusion rule
type PortExclusionRule struct {
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"` // TCP, UDP, ANY
	Reason   string `json:"reason"`
	Enabled  bool   `json:"enabled"`
}

// Protocol filter constants
const (
	ProtocolFilterTCP = "TCP"
	ProtocolFilterUDP = "UDP"
	ProtocolFilterAny = "ANY"
)
