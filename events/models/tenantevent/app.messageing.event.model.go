package tenantevent

type AppMessaegerInsertEventModel struct {
	ID         string         `bson:"_id" json:"id"`
	Provider   string         `bson:"provider" json:"provider"`
	AssignedTo []string       `bson:"assignedTo" json:"assignedTo"`
	Metadata   map[string]any `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsActive   bool           `bson:"isActive" json:"isActive"`
	Version    int            `bson:"version" json:"version"`
}

type AppMessengerUpdateEventModel struct {
	ID         string         `bson:"_id" json:"id"`
	Provider   string         `bson:"provider" json:"provider"`
	AssignedTo []string       `bson:"assignedTo" json:"assignedTo"`
	Metadata   map[string]any `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IsActive   bool           `bson:"isActive" json:"isActive"`
	Version    int            `bson:"version" json:"version"`
}

type AppMessengerDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}
