package tenantevent

type KYCProvideInsertEventModel struct {
	ID         string         `bson:"_id" json:"id"`
	Provider   string         `bson:"provider" json:"provider"`
	AssignedTo []string       `bson:"assignedTo" json:"assignedTo"`
	Metadata   map[string]any `bson:"metadata,omitempty" json:"metadata"`
	IsActive   bool           `bson:"isActive" json:"isActive"`
	Version    int            `bson:"version" json:"version"`
}

type KYCProvideUpdateEventModel struct {
	ID         string         `bson:"_id" json:"id"`
	Provider   string         `bson:"provider" json:"provider"`
	AssignedTo []string       `json:"assignedTo"`
	Metadata   map[string]any `bson:"metadata,omitempty" json:"metadata"`
	IsActive   bool           `bson:"isActive" json:"isActive"`
	Version    int            `bson:"version" json:"version"`
}

type KYCProvideDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}
