package authevent

type RoleCreateEventModel struct {
	ID       string `json:"id"`
	IsSystem bool   `json:"isSystem"`
	IsActive bool   `json:"isActive"`
	Version  int    `json:"version"`
}

type RoleUpdateEventModel struct {
	ID         string   `json:"id"`
	IsSystem   bool     `json:"isSystem"`
	IsActive   bool     `json:"isActive"`
	AssignedTo []string `json:"assignedTo"`
	Version    int      `json:"version"`
}

type RoleDeleteEventModel struct {
	ID string `json:"id"`
}
