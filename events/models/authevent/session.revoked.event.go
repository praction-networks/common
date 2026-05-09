package authevent

// AuthSessionRevokedEvent is the payload of the auth.session.revoked NATS
// event published on AuthStream after a session is revoked (per-session DELETE,
// sign-out-all, or logout).
//
// Subject: "auth.session.revoked"
// Reliability: BestEffort
type AuthSessionRevokedEvent struct {
	UserID    string `json:"userId"`
	TenantID  string `json:"tenantId"`
	SessionID string `json:"sessionId"` // session row id, OR "all" for sign-out-all
	Reason    string `json:"reason"`    // "user_initiated", "sign_out_all", "logout"
	AtMs      int64  `json:"atMs"`      // unix ms
}
