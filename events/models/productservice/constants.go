package productservice

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
	StatusDraft    Status = "DRAFT"
	StautsAssigned Status = "ASSIGNED"
)