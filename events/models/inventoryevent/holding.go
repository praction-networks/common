package inventoryevent

// ConsumableHolding is a per-user current possession record. Quantity is
// float64 to support fractional metres (cable). Integer-unit consumables
// store whole numbers.
type ConsumableHolding struct {
	ID           string  `json:"id"           bson:"_id"`
	TenantID     string  `json:"tenantId"     bson:"tenantId"`
	UserID       string  `json:"userId"       bson:"userId"`
	ConsumableID string  `json:"consumableId" bson:"consumableId"`
	Quantity     float64 `json:"quantity"     bson:"quantity"`
}
