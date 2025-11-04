package authevent

type RoleCreateEventModel struct {
	ID       string `bson:"_id" json:"id"`
	IsSystem bool   `bson:"isSystem" json:"isSystem"`
	IsActive bool   `bson:"isActive" json:"isActive"`
	Version  int    `bson:"version" json:"version"`
}

type RoleUpdateEventModel struct {
	ID         string   `bson:"_id" json:"id"`
	IsSystem   bool     `bson:"isSystem" json:"isSystem"`
	IsActive   bool     `bson:"isActive" json:"isActive"`
	AssignedTo []string `bson:"assignedTo" json:"assignedTo"`
	Version    int      `bson:"version" json:"version"`
}

type RoleDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}
