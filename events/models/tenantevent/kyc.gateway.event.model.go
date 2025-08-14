package tenantevent

type KYCProvideInsertEventModel struct {
	ID       string         `json:"id"`
	Provider string         `json:"provider"`
	Metadata map[string]any `json:"metadata"`
	IsActive bool           `json:"isActive"`
	Version  int            `json:"version"`
}

type KYCProvideUpdateEventModel struct {
	ID         string         `json:"id"`
	Provider   string         `json:"provider"`
	AssignedTo []string       `json:"assignedTo"`
	Metadata   map[string]any `json:"metadata"`
	IsActive   bool           `json:"isActive"`
	Version    int            `json:"version"`
}

type KYCProvideDeleteEventModel struct {
	ID string `json:"id"`
}
