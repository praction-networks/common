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
	Data      map[string]interface{} `json:"data"`     // Deleted row values
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

