package tenantevent

type RadiusProviderInsertEventModel struct {
	ID       string         `json:"id"`
	Provider string         `json:"provider"`
	Metadata map[string]any `json:"metadata"`
	IsActive bool           `json:"isActive"`
	Version  int            `json:"version"`
}

type RadiusProviderUpdateEventModel struct {
	ID         string         `json:"id"`
	Provider   string         `json:"provider"`
	AssignedTo []string       `json:"assignedTo"`
	Metadata   map[string]any `json:"metadata"`
	IsActive   bool           `json:"isActive"`
	Version    int            `json:"version"`
}

type RadiusProviderDeleteEventModel struct {
	ID string `json:"id"`
}
