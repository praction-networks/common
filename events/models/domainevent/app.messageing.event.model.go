package domainevent

type AppMessangerInsertEventModel struct {
	ID       string         `json:"id"`
	Provider string         `json:"provider"`
	Metadata map[string]any `json:"metadata,omitempty"`
	IsActive bool           `json:"isActive"`
	Version  int            `json:"version"`
}

type AppMessangerUpdateEventModel struct {
	ID         string         `json:"id"`
	Provider   string         `json:"provider"`
	AssignedTo []string       `json:"assignedTo"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	IsActive   bool           `json:"isActive"`
	Version    int            `json:"version"`
}

type AppMessangerDeleteEventModel struct {
	ID string `json:"id"`
}
