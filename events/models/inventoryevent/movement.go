package inventoryevent

// ConsumableMovement is an immutable audit row. Kind enumerates the cause:
// ISSUED, CONSUMED, RETURNED, RECOVERED, ADJUSTED.
type ConsumableMovement struct {
	ID           string  `json:"id"                   bson:"_id"`
	TenantID     string  `json:"tenantId"             bson:"tenantId"`
	ConsumableID string  `json:"consumableId"         bson:"consumableId"`
	Kind         string  `json:"kind"                 bson:"kind"`
	Quantity     float64 `json:"quantity"             bson:"quantity"`
	FromUserID   string  `json:"fromUserId,omitempty" bson:"fromUserId,omitempty"`
	ToUserID     string  `json:"toUserId,omitempty"   bson:"toUserId,omitempty"`
	AtMs         int64   `json:"atMs"                 bson:"atMs"`
}
