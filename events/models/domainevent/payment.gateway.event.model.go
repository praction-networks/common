package domainevent

type PaymentGatewayInsertEventModel struct {
	ID       string         `json:"id"`
	Gateway  string         `json:"gateway"`
	IsActive bool           `json:"isActive"`
	Metadata map[string]any `json:"metadata"`
	Version  int            `json:"version"`
}

type PaymentGatewayUpdateEventModel struct {
	ID         string         `json:"id"`
	Gateway    string         `json:"gateway"`
	AssignedTo []string       `json:"assignedTo"`
	IsActive   bool           `json:"isActive"`
	Metadata   map[string]any `json:"metadata"`
	Version    int            `json:"version"`
}

type PaymentGatewayDeleteEventModel struct {
	ID string `json:"id"`
}
