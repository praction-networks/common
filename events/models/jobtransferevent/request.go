package jobtransferevent

// JobTransferRequest tracks a mid-job hand-off attempt. Lifecycle:
// PENDING_APPROVAL → APPROVED | DECLINED → PEER_OFFERED →
// PEER_ACCEPTED | PEER_DECLINED | EXPIRED | WITHDRAWN | FORCE_UNASSIGNED.
type JobTransferRequest struct {
	ID                string `json:"id"                bson:"_id"`
	TenantID          string `json:"tenantId"          bson:"tenantId"`
	TicketID          string `json:"ticketId"          bson:"ticketId"`
	FromUserID        string `json:"fromUserId"        bson:"fromUserId"`
	Status            string `json:"status"            bson:"status"`
	Reason            string `json:"reason,omitempty"  bson:"reason,omitempty"`
	RequestedAtMs     int64  `json:"requestedAtMs"     bson:"requestedAtMs"`
	LifetimeExpiresAt int64  `json:"lifetimeExpiresAt" bson:"lifetimeExpiresAt"` // total-lifetime hard cap (tenant policy)
}
