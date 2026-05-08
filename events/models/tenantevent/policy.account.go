package tenantevent

// PolicyAccount is the per-tenant account-management policy bucket.
// Source: backend-contract §11.5, Sweep 4 Q2 (emailChangeAllowed).
type PolicyAccount struct {
	EmailChangeAllowed bool `json:"emailChangeAllowed,omitempty" bson:"emailChangeAllowed,omitempty"`
}
