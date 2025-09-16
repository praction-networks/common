package tenantevent

type SMSProviderInsertEventModel struct {
	ID         string         `json:"id"`
	Provider   string         `json:"provider"`
	AssignedTo []string       `json:"assignedTo"`
	Metadata   map[string]any `json:"metadata"`
	IsActive   bool           `json:"isActive"`
	Version    int            `json:"version"`
}

type SMSProviderUpdateEventModel struct {
	ID         string         `json:"id"`
	Provider   string         `json:"provider"`
	AssignedTo []string       `json:"assignedTo"`
	Metadata   map[string]any `json:"metadata"`
	IsActive   bool           `json:"isActive"`
	Version    int            `json:"version"`
}

type SMSProviderDeleteEventModel struct {
	ID string `json:"id"`
}
