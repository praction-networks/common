package peerassistevent

// PeerAssistRequest is FE-A's request for help. WM approves/declines.
type PeerAssistRequest struct {
	ID            string `json:"id"                     bson:"_id"`
	TenantID      string `json:"tenantId"               bson:"tenantId"`
	RequesterID   string `json:"requesterId"            bson:"requesterId"`
	Status        string `json:"status"                 bson:"status"` // PENDING_APPROVAL | APPROVED | DECLINED | CANCELLED
	ReasonCode    string `json:"reasonCode,omitempty"   bson:"reasonCode,omitempty"`
	Note          string `json:"note,omitempty"         bson:"note,omitempty"`
	RequestedAtMs int64  `json:"requestedAtMs"          bson:"requestedAtMs"`
	ResolvedAtMs  int64  `json:"resolvedAtMs,omitempty" bson:"resolvedAtMs,omitempty"`
}
