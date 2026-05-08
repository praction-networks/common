package tenantevent

// PolicyAssets is the per-tenant asset-policy bucket.
// Source: backend-contract §11.3, Sweep 4 Q9 (dropoffLockWindowHours), Q10 (recoveryAcknowledgement).
type PolicyAssets struct {
	DropoffLockWindowHours  int    `json:"dropoffLockWindowHours,omitempty"  bson:"dropoffLockWindowHours,omitempty"`
	RecoveryAcknowledgement string `json:"recoveryAcknowledgement,omitempty" bson:"recoveryAcknowledgement,omitempty"` // NONE | SIGNATURE | OTP
}
