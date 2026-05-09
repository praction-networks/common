package tenantevent

// TenantPolicyUpdatedEvent is the payload of the tenant.policy.updated NATS
// event published on TenantStream after a successful PATCH /tenant/me/policy.
//
// Subject: "tenant.policy.updated" (no .v1 suffix per common/ Wave 0
// convention; payload-versioning via the policy's own Version field).
//
// Reliability: Preferred. Late events tolerated; consumers re-fetch on
// reconnect or detect via Version compare.
//
// Source: design 2026-05-09-tenant-service-policy-design.md §5.5.
type TenantPolicyUpdatedEvent struct {
	TenantID    string   `json:"tenantId"`
	Version     int      `json:"version"`     // post-merge Version
	UpdatedAtMs int64    `json:"updatedAtMs"` // unix ms
	ChangedKeys []string `json:"changedKeys"` // dot-paths, e.g. "assets.dropoffLockWindowHours"
}
