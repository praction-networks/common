package inventoryevent

type DeviceInsertEventModel struct {
	ID       string         `json:"id"`
	Metadata map[string]any `json:"metadata"`
	IsActive bool           `json:"isActive"`
}
