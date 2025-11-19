package tenantevent

type AppMessaegerInsertEventModel struct {
	ID         string                 `bson:"_id" json:"id"`
	Provider   string                 `bson:"provider" json:"provider"`
	AssignedTo []WhatsAppTenantConfig `bson:"assignedTo" json:"assignedTo"`
	Metadata   map[string]any         `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsActive   bool                   `bson:"isActive" json:"isActive"`
	IsSystem   bool                   `bson:"isSystem" json:"isSystem"`
	Version    int                    `bson:"version" json:"version"`
}

type AppMessengerUpdateEventModel struct {
	ID         string                 `bson:"_id" json:"id"`
	Provider   string                 `bson:"provider" json:"provider"`
	AssignedTo []WhatsAppTenantConfig `bson:"assignedTo" json:"assignedTo"`
	Metadata   map[string]any         `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsActive   bool                   `bson:"isActive" json:"isActive"`
	IsSystem   bool                   `bson:"isSystem" json:"isSystem"`
	Version    int                    `bson:"version" json:"version"`
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
