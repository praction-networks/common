package tenantevent

// SMSType represents supported SMS gateway providers
type SMSType string

// SMSProviderInsertEventModel defines the model for SMS provider insert events
type SMSProviderInsertEventModel struct {
	ID                 string         `bson:"_id" json:"id"`
	Name               string         `bson:"name" json:"name"`
	Provider           string         `bson:"provider" json:"provider"`
	OwnerTenantID      string         `bson:"ownerTenantId" json:"ownerTenantId"`
	OwnerTenantType    string         `bson:"ownerTenantType" json:"ownerTenantType"`
	Scope              string         `bson:"scope" json:"scope"`
	AllowedTenantTypes []string       `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string       `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	Metadata           map[string]any `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsActive           bool           `bson:"isActive" json:"isActive"`
	Version            int            `bson:"version" json:"version"`
}

// SMSProviderUpdateEventModel defines the model for SMS provider update events
type SMSProviderUpdateEventModel struct {
	ID                 string         `bson:"_id" json:"id"`
	Name               string         `bson:"name,omitempty" json:"name,omitempty"`
	Provider           string         `bson:"provider,omitempty" json:"provider,omitempty"`
	OwnerTenantID      string         `bson:"ownerTenantId,omitempty" json:"ownerTenantId,omitempty"`
	OwnerTenantType    string         `bson:"ownerTenantType,omitempty" json:"ownerTenantType,omitempty"`
	Scope              string         `bson:"scope,omitempty" json:"scope,omitempty"`
	AllowedTenantTypes []string       `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string       `bson:"explicitTenantIds,omitempty" json:"explicitTenantIds,omitempty"`
	Metadata           map[string]any `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsActive           *bool          `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Version            *int           `bson:"version,omitempty" json:"version,omitempty"`
}

// SMSProviderDeleteEventModel defines the model for SMS provider delete events
type SMSProviderDeleteEventModel struct {
	ID      string `bson:"_id" json:"id"`
	Version int    `bson:"version" json:"version"`
}
