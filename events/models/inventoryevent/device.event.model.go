package inventoryevent

type DeviceInsertEventModel struct {
	ID       string         `bson:"_id" json:"id"`
	Metadata map[string]any `bson:"metadata" json:"metadata"`
	IsActive bool           `bson:"isActive" json:"isActive"`
}

type DeviceUpdateEventModel struct {
	ID       string         `bson:"_id" json:"id"`
	Metadata map[string]any `bson:"metadata" json:"metadata"`
	IsActive bool           `bson:"isActive" json:"isActive"`
}

type DeviceDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}
