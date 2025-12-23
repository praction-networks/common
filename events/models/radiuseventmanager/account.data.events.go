package radiuseventmanager

import "time"

// RadAcctEvent represents a RADIUS accounting event from FreeRADIUS
// This is the typed representation of CDC events from the radAcct table
type RadAcctEvent struct {
	RadAcctID           string     `json:"radAcctId"`
	TenantID            string     `json:"tenantId"`
	AcctSessionID       string     `json:"acctSessionId"`
	AcctUniqueID        string     `json:"acctUniqueId"`
	UserName            string     `json:"userName"`
	Realm               string     `json:"realm,omitempty"`
	NasIPAddress        string     `json:"nasIpAddress"`
	NasPortID           string     `json:"nasPortId,omitempty"`
	NasPortType         string     `json:"nasPortType,omitempty"`
	AcctStartTime       *time.Time `json:"acctStartTime,omitempty"`
	AcctUpdateTime      *time.Time `json:"acctUpdateTime,omitempty"`
	AcctStopTime        *time.Time `json:"acctStopTime,omitempty"`
	AcctInterval        *int       `json:"acctInterval,omitempty"`
	AcctSessionTime     *int       `json:"acctSessionTime,omitempty"`
	AcctAuthentic       string     `json:"acctAuthentic,omitempty"`
	ConnectInfoStart    string     `json:"connectInfoStart,omitempty"`
	ConnectInfoStop     string     `json:"connectInfoStop,omitempty"`
	AcctInputOctets     *int64     `json:"acctInputOctets,omitempty"`
	AcctOutputOctets    *int64     `json:"acctOutputOctets,omitempty"`
	CalledStationID     string     `json:"calledStationId,omitempty"`
	CallingStationID    string     `json:"callingStationId,omitempty"`
	AcctTerminateCause  string     `json:"acctTerminateCause,omitempty"`
	ServiceType         string     `json:"serviceType,omitempty"`
	FramedProtocol      string     `json:"framedProtocol,omitempty"`
	FramedIPAddress     string     `json:"framedIpAddress,omitempty"`
	FramedIPv6Address   string     `json:"framedIpv6Address,omitempty"`
	FramedIPv6Prefix    string     `json:"framedIpv6Prefix,omitempty"`
	FramedInterfaceID   string     `json:"framedInterfaceId,omitempty"`
	DelegatedIPv6Prefix string     `json:"delegatedIpv6Prefix,omitempty"`
	Version             int        `json:"version"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	Timestamp           time.Time  `json:"timestamp"` // When detected via CDC
}

// RadAcctSessionStartEvent represents a new RADIUS session start (INSERT on radAcct)
// Contains all session information when a new RADIUS session begins
type RadAcctSessionStartEvent struct {
	RadAcctEvent
}

// RadAcctSessionUpdateEvent represents a RADIUS session update (UPDATE on radAcct)
// This includes:
// 1. Interim updates: Periodic accounting updates sent by NAS during active session (usage data)
// 2. Session stops: When acctStopTime is set via UPDATE (session stopped via UPDATE)
// When IsSessionStop is false, this is an interim update with usage statistics
// When IsSessionStop is true, this represents a session stop via UPDATE (acctStopTime set)
type RadAcctSessionUpdateEvent struct {
	RadAcctEvent
	IsSessionStop bool `json:"isSessionStop"` // True if acctStopTime is set (session stopped via UPDATE)
	// When IsSessionStop is false, this is an interim update with:
	// - AcctUpdateTime: Time of this interim update
	// - AcctSessionTime: Current session duration
	// - AcctInputOctets: Current input bytes (downloaded)
	// - AcctOutputOctets: Current output bytes (uploaded)
	// - AcctInterval: Interval between updates (if configured)
}

// RadAcctSessionEndEvent represents a RADIUS session end (DELETE on radAcct)
// Published when a radAcct record is deleted, indicating the session has ended
// Contains final session statistics and termination information
type RadAcctSessionEndEvent struct {
	RadAcctID          string     `json:"radAcctId"`                    // Unique RADIUS accounting record ID
	TenantID           string     `json:"tenantId"`                    // Tenant identifier
	AcctSessionID      string     `json:"acctSessionId"`               // RADIUS accounting session ID
	AcctUniqueID       string     `json:"acctUniqueId"`                // Unique accounting identifier
	UserName           string     `json:"userName"`                    // RADIUS username
	Realm              string     `json:"realm,omitempty"`              // RADIUS realm
	NasIPAddress       string     `json:"nasIpAddress"`                // NAS IP address
	NasPortID          string     `json:"nasPortId,omitempty"`         // NAS port identifier
	NasPortType        string     `json:"nasPortType,omitempty"`       // NAS port type (e.g., "Wireless-802.11")
	AcctStartTime      *time.Time `json:"acctStartTime,omitempty"`      // Session start time
	AcctStopTime       *time.Time `json:"acctStopTime,omitempty"`      // Session stop time
	AcctSessionTime    *int       `json:"acctSessionTime,omitempty"`  // Total session duration in seconds
	AcctInputOctets    *int64     `json:"acctInputOctets,omitempty"`  // Total bytes downloaded
	AcctOutputOctets   *int64     `json:"acctOutputOctets,omitempty"` // Total bytes uploaded
	AcctTerminateCause string     `json:"acctTerminateCause,omitempty"` // Session termination cause
	CalledStationID    string     `json:"calledStationId,omitempty"`   // Called station ID (SSID/AP name)
	CallingStationID   string     `json:"callingStationId,omitempty"`   // Calling station ID (MAC address)
	FramedIPAddress    string     `json:"framedIpAddress,omitempty"`   // IP address assigned to user
	ServiceType         string     `json:"serviceType,omitempty"`      // Service type
	FramedProtocol      string     `json:"framedProtocol,omitempty"`  // Framed protocol
	Version             int        `json:"version"`                    // Record version
	Timestamp           time.Time  `json:"timestamp"`                  // When detected via CDC
}
