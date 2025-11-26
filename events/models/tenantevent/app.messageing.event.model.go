package tenantevent

type AppMessaegerInsertEventModel struct {
	ID                 string                 `bson:"_id" json:"id"`
	Name               string                 `bson:"name" json:"name"`
	Provider           string                 `bson:"provider" json:"provider"`
	OwnerTenantID      string                 `bson:"ownerTenantID" json:"ownerTenantID"`
	OwnerTenantType    string                 `bson:"ownerTenantType" json:"ownerTenantType"`
	Scope              string                 `bson:"scope" json:"scope"`
	AllowedTenantTypes []string               `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string               `bson:"explicitTenantIDs,omitempty" json:"explicitTenantIDs,omitempty"`
	AssignedTo         []WhatsAppTenantConfig `bson:"assignedTo,omitempty" json:"assignedTo,omitempty"` // Deprecated: Use ExplicitTenantIDs when Scope=ExplicitTenants
	Metadata           map[string]any         `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsActive           bool                   `bson:"isActive" json:"isActive"`
	Version            int                    `bson:"version" json:"version"`
}

type AppMessengerUpdateEventModel struct {
	ID                 string                 `bson:"_id" json:"id"`
	Name               string                 `bson:"name" json:"name"`
	Provider           string                 `bson:"provider" json:"provider"`
	OwnerTenantID      string                 `bson:"ownerTenantID" json:"ownerTenantID"`
	OwnerTenantType    string                 `bson:"ownerTenantType" json:"ownerTenantType"`
	Scope              string                 `bson:"scope" json:"scope"`
	AllowedTenantTypes []string               `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string               `bson:"explicitTenantIDs,omitempty" json:"explicitTenantIDs,omitempty"`
	AssignedTo         []WhatsAppTenantConfig `bson:"assignedTo,omitempty" json:"assignedTo,omitempty"` // Deprecated: Use ExplicitTenantIDs when Scope=ExplicitTenants
	Metadata           map[string]any         `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsActive           bool                   `bson:"isActive" json:"isActive"`
	Version            int                    `bson:"version" json:"version"`
}

type AppMessengerDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}

type WhatsAppTenantConfig struct {
	TenantID             string   `json:"tenantID" bson:"tenantID" validate:"required"`
	PhoneNumberID        string   `json:"phoneNumberID" bson:"phoneNumberID" validate:"required"`
	BusinessAccountID    string   `json:"businessAccountID" bson:"businessAccountID" validate:"required"`
	AccessToken          string   `json:"accessToken" bson:"accessToken" validate:"required"`
	IsActive             bool     `json:"isActive" bson:"isActive"`
	AllowedMessagesTypes []string `json:"allowedMessagesTypes,omitempty" bson:"allowedMessagesTypes,omitempty"`
}
