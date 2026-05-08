package returnevent

// ReturnLineStatus enumerates drop-off line lifecycle.
type ReturnLineStatus string

const (
	ReturnLineStatusPending  ReturnLineStatus = "PENDING_WM_ACCEPT"
	ReturnLineStatusAccepted ReturnLineStatus = "ACCEPTED"
	ReturnLineStatusRejected ReturnLineStatus = "REJECTED"
)

// Return is a drop-off batch submitted by an FE-tech. WM accepts/rejects
// per line. Auto-released after dropoffLockWindowHours (tenant policy).
type Return struct {
	ID            string       `json:"id"                          bson:"_id"`
	TenantID      string       `json:"tenantId"                    bson:"tenantId"`
	UserID        string       `json:"userId"                      bson:"userId"`
	WarehouseID   string       `json:"warehouseId"                 bson:"warehouseId"`
	RecipientWmID string       `json:"recipientWmUserId,omitempty" bson:"recipientWmUserId,omitempty"`
	Lines         []ReturnLine `json:"lines"                       bson:"lines"`
	Note          string       `json:"note,omitempty"              bson:"note,omitempty"`
	SubmittedAtMs int64        `json:"submittedAtMs"               bson:"submittedAtMs"`
}

// ReturnLine — exactly one of AssetID or ConsumableID is set.
type ReturnLine struct {
	ID             string           `json:"id"                       bson:"id"`
	AssetID        *string          `json:"assetId,omitempty"        bson:"assetId,omitempty"`
	ConsumableID   *string          `json:"consumableId,omitempty"   bson:"consumableId,omitempty"`
	Quantity       *float64         `json:"quantity,omitempty"       bson:"quantity,omitempty"`
	PhotoIDs       []string         `json:"photoIds,omitempty"       bson:"photoIds,omitempty"`
	Condition      string           `json:"condition,omitempty"      bson:"condition,omitempty"`
	RecoveryReason string           `json:"recoveryReason,omitempty" bson:"recoveryReason,omitempty"`
	Status         ReturnLineStatus `json:"status"                   bson:"status"`
	RejectReason   string           `json:"rejectReason,omitempty"   bson:"rejectReason,omitempty"`
	WmNote         string           `json:"wmNote,omitempty"         bson:"wmNote,omitempty"`
}
