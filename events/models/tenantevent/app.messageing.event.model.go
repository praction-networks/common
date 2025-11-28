package tenantevent

// AppMessengerInsertEventModel defines the model for app messaging provider insert events
type AppMessengerInsertEventModel struct {
	ID                 string         `bson:"_id" json:"id"`
	Name               string         `bson:"name" json:"name" `
	Provider           string         `bson:"provider" json:"provider" `
	OwnerTenantID      string         `bson:"ownerTenantID" json:"ownerTenantID"`
	OwnerTenantType    string         `bson:"ownerTenantType" json:"ownerTenantType" `
	Scope              string         `bson:"scope" json:"scope" `
	AllowedTenantTypes []string       `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string       `bson:"explicitTenantIDs,omitempty" json:"explicitTenantIDs,omitempty" `
	Metadata           map[string]any `bson:"metadata,omitempty" json:"metadata,omitempty" `
	IsActive           bool           `bson:"isActive" json:"isActive" `
	Version            int            `bson:"version" json:"version" `
}

// AppMessengerUpdateEventModel defines the model for app messaging provider update events
type AppMessengerUpdateEventModel struct {
	ID                 string         `bson:"_id" json:"id"`
	Name               string         `bson:"name,omitempty" json:"name,omitempty"`
	Provider           string         `bson:"provider,omitempty" json:"provider,omitempty"`
	OwnerTenantID      string         `bson:"ownerTenantID,omitempty" json:"ownerTenantID,omitempty" `
	OwnerTenantType    string         `bson:"ownerTenantType,omitempty" json:"ownerTenantType,omitempty"`
	Scope              string         `bson:"scope,omitempty" json:"scope,omitempty" `
	AllowedTenantTypes []string       `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string       `bson:"explicitTenantIDs,omitempty" json:"explicitTenantIDs,omitempty" `
	Metadata           map[string]any `bson:"metadata,omitempty" json:"metadata,omitempty" `
	IsActive           *bool          `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Version            *int           `bson:"version,omitempty" json:"version,omitempty" `
}

// AppMessengerDeleteEventModel defines the model for app messaging provider delete events
type AppMessengerDeleteEventModel struct {
	ID      string `bson:"_id" json:"id"`
	Version int    `bson:"version" json:"version"`
}
