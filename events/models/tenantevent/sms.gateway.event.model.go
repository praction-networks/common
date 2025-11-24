package tenantevent

// SMSType represents supported SMS gateway providers
type SMSType string

const (
	SMSGatewayMSG91OTP SMSType = "MSG91"
)

type SMSProviderInsertEventModel struct {
	ID         string         `bson:"_id" json:"id"`
	Provider   string         `bson:"provider" json:"provider"`
	AssignedTo []string       `bson:"assignedTo" json:"assignedTo"`
	Metadata   map[string]any `bson:"metadata,omitempty" json:"metadata"`
	IsActive   bool           `bson:"isActive" json:"isActive"`
	Version    int            `bson:"version" json:"version"`
}

type SMSProviderUpdateEventModel struct {
	ID         string         `bson:"_id" json:"id"`
	Provider   string         `bson:"provider" json:"provider"`
	AssignedTo []string       `bson:"assignedTo" json:"assignedTo"`
	Metadata   map[string]any `bson:"metadata,omitempty" json:"metadata"`
	IsActive   bool           `bson:"isActive" json:"isActive"`
	Version    int            `bson:"version" json:"version"`
}

type SMSProviderDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}
