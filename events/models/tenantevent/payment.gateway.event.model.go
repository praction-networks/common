package tenantevent

type PaymentGatewayInsertEventModel struct {
	ID         string         `bson:"_id" json:"id"`
	Gateway    string         `bson:"gateway" json:"gateway"`
	AssignedTo []string       `bson:"assignedTo" json:"assignedTo"`
	IsActive   bool           `bson:"isActive" json:"isActive"`
	Metadata   map[string]any `bson:"metadata,omitempty" json:"metadata"`
	Version    int            `bson:"version" json:"version"`
}

type PaymentGatewayUpdateEventModel struct {
	ID         string         `bson:"_id" json:"id"`
	Gateway    string         `bson:"gateway" json:"gateway"`
	AssignedTo []string       `bson:"assignedTo" json:"assignedTo"`
	IsActive   bool           `bson:"isActive" json:"isActive"`
	Metadata   map[string]any `bson:"metadata,omitempty" json:"metadata"`
	Version    int            `bson:"version" json:"version"`
}

type PaymentGatewayDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}
