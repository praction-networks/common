package accountdataevent

import "time"

// AccountDataCreatedEvent represents an account data creation event from PostgreSQL CDC
type AccountDataCreatedEvent struct {
	TableName string                 `json:"tableName"` // e.g., "user_profiles", "radcheck", "radusergroup"
	Data      map[string]interface{} `json:"data"`      // All column values from the row
	Timestamp time.Time              `json:"timestamp"` // When the change was detected
}

// AccountDataUpdatedEvent represents an account data update event from PostgreSQL CDC
type AccountDataUpdatedEvent struct {
	TableName string                 `json:"tableName"` // e.g., "user_profiles", "radcheck", "radusergroup"
	NewData   map[string]interface{} `json:"newData"`   // New row values
	OldData   map[string]interface{} `json:"oldData"`   // Old row values (may be nil)
	Timestamp time.Time              `json:"timestamp"` // When the change was detected
}

// AccountDataDeletedEvent represents an account data deletion event from PostgreSQL CDC
type AccountDataDeletedEvent struct {
	TableName string                 `json:"tableName"` // e.g., "user_profiles", "radcheck", "radusergroup"
	Data      map[string]interface{} `json:"data"`      // Deleted row values
	Timestamp time.Time              `json:"timestamp"` // When the change was detected
}

// UserProfileCreatedEvent represents a user profile creation event
type UserProfileCreatedEvent struct {
	ID                   string     `json:"id"`
	TenantID             string     `json:"tenantId"`
	Username             string     `json:"username"`
	Email                string     `json:"email,omitempty"`
	Phone                string     `json:"phone,omitempty"`
	FullName             string     `json:"fullName,omitempty"`
	Status               string     `json:"status"`
	UserType             string     `json:"userType"`
	SimultaneousUseLimit int        `json:"simultaneousUseLimit"`
	MaxMacAddresses      *int       `json:"maxMacAddresses,omitempty"`
	ExpiresAt            *time.Time `json:"expiresAt,omitempty"`
	Notes                string     `json:"notes,omitempty"`
	CreatedBy            string     `json:"createdBy"`
	CreatedAt            time.Time  `json:"createdAt"`
	Timestamp            time.Time  `json:"timestamp"` // When detected via CDC
}

// UserProfileUpdatedEvent represents a user profile update event
type UserProfileUpdatedEvent struct {
	ID                   string     `json:"id"`
	TenantID             string     `json:"tenantId"`
	Username             *string    `json:"username,omitempty"`
	Email                *string    `json:"email,omitempty"`
	Phone                *string    `json:"phone,omitempty"`
	FullName             *string    `json:"fullName,omitempty"`
	Status               *string    `json:"status,omitempty"`
	UserType             *string    `json:"userType,omitempty"`
	SimultaneousUseLimit *int       `json:"simultaneousUseLimit,omitempty"`
	MaxMacAddresses      *int       `json:"maxMacAddresses,omitempty"`
	ExpiresAt            *time.Time `json:"expiresAt,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
	UpdatedAt            *time.Time `json:"updatedAt,omitempty"`
	Timestamp            time.Time  `json:"timestamp"` // When detected via CDC
}

// UserProfileDeletedEvent represents a user profile deletion event
type UserProfileDeletedEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Username  string    `json:"username,omitempty"`
	Timestamp time.Time `json:"timestamp"` // When detected via CDC
}

// RadcheckCreatedEvent represents a radcheck entry creation event
type RadcheckCreatedEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Username  string    `json:"username"`
	Attribute string    `json:"attribute"`
	Op        string    `json:"op"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"` // When detected via CDC
}

// RadcheckUpdatedEvent represents a radcheck entry update event
type RadcheckUpdatedEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Username  string    `json:"username"`
	Attribute *string   `json:"attribute,omitempty"`
	Op        *string   `json:"op,omitempty"`
	Value     *string   `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp"` // When detected via CDC
}

// RadcheckDeletedEvent represents a radcheck entry deletion event
type RadcheckDeletedEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Username  string    `json:"username"`
	Attribute string    `json:"attribute,omitempty"`
	Timestamp time.Time `json:"timestamp"` // When detected via CDC
}

// RadusergroupCreatedEvent represents a radusergroup entry creation event
type RadusergroupCreatedEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Username  string    `json:"username"`
	Groupname string    `json:"groupname"`
	Priority  int       `json:"priority"`
	Timestamp time.Time `json:"timestamp"` // When detected via CDC
}

// RadusergroupUpdatedEvent represents a radusergroup entry update event
type RadusergroupUpdatedEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Username  string    `json:"username"`
	Groupname *string   `json:"groupname,omitempty"`
	Priority  *int      `json:"priority,omitempty"`
	Timestamp time.Time `json:"timestamp"` // When detected via CDC
}

// RadusergroupDeletedEvent represents a radusergroup entry deletion event
type RadusergroupDeletedEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Username  string    `json:"username"`
	Groupname string    `json:"groupname,omitempty"`
	Timestamp time.Time `json:"timestamp"` // When detected via CDC
}

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
type RadAcctSessionStartEvent struct {
	RadAcctEvent
}

// RadAcctSessionUpdateEvent represents a RADIUS session update (UPDATE on radAcct)
// This includes interim updates and session stops
type RadAcctSessionUpdateEvent struct {
	RadAcctEvent
	IsSessionStop bool `json:"isSessionStop"` // True if acctStopTime is set
}

// RadAcctSessionEndEvent represents a RADIUS session end (DELETE on radAcct or session with stop time)
type RadAcctSessionEndEvent struct {
	RadAcctID          string    `json:"radAcctId"`
	TenantID           string    `json:"tenantId"`
	AcctSessionID      string    `json:"acctSessionId"`
	UserName           string    `json:"userName"`
	AcctSessionTime    *int      `json:"acctSessionTime,omitempty"`  // Total session duration in seconds
	AcctInputOctets    *int64    `json:"acctInputOctets,omitempty"`  // Total bytes downloaded
	AcctOutputOctets   *int64    `json:"acctOutputOctets,omitempty"` // Total bytes uploaded
	AcctTerminateCause string    `json:"acctTerminateCause,omitempty"`
	Timestamp          time.Time `json:"timestamp"` // When detected via CDC
}
