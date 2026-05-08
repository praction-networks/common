package peerassistevent

// PeerHandoff is the ceremony when FE-A transfers inventory to FE-B.
// State machine: AWAITING_CONFIRM → CONFIRMED | DECLINED | RECEIPT_TIMEOUT.
type PeerHandoff struct {
	ID           string            `json:"id"                     bson:"_id"`
	TenantID     string            `json:"tenantId"               bson:"tenantId"`
	RequestID    string            `json:"requestId,omitempty"    bson:"requestId,omitempty"`
	FromUserID   string            `json:"fromUserId"             bson:"fromUserId"`
	ToUserID     string            `json:"toUserId"               bson:"toUserId"`
	Status       string            `json:"status"                 bson:"status"`
	Lines        []PeerHandoffLine `json:"lines"                  bson:"lines"`
	StartedAtMs  int64             `json:"startedAtMs"            bson:"startedAtMs"`
	ResolvedAtMs int64             `json:"resolvedAtMs,omitempty" bson:"resolvedAtMs,omitempty"`
}

// PeerHandoffLine — exactly one of AssetID or ConsumableID is set.
type PeerHandoffLine struct {
	ID           string   `json:"id"                     bson:"id"`
	AssetID      *string  `json:"assetId,omitempty"      bson:"assetId,omitempty"`
	ConsumableID *string  `json:"consumableId,omitempty" bson:"consumableId,omitempty"`
	Quantity     *float64 `json:"quantity,omitempty"     bson:"quantity,omitempty"`
}
