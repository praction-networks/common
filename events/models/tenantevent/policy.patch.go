package tenantevent

// TenantPolicyPatch is the wire shape for sparse PATCH /tenant/me/policy.
//
// Each bucket is *Bucket so a nil bucket means "client did not send this
// bucket"; an empty {} bucket means "client sent the bucket but no fields
// inside changed". Used as input to MergeTenantPolicy.
//
// Within each bucket, fields use *bool / *int / *float64 (per the bucket
// definitions in policy.account.go, policy.shift.go, etc.) so absent fields
// are likewise distinguishable from zero-values.
type TenantPolicyPatch struct {
	Account       *PolicyAccount       `json:"account,omitempty"       bson:"account,omitempty"`
	Assets        *PolicyAssets        `json:"assets,omitempty"        bson:"assets,omitempty"`
	Auth          *PolicyAuth          `json:"auth,omitempty"          bson:"auth,omitempty"`
	Notifications *PolicyNotifications `json:"notifications,omitempty" bson:"notifications,omitempty"`
	Onboard       *PolicyOnboard       `json:"onboard,omitempty"       bson:"onboard,omitempty"`
	Shift         *PolicyShift         `json:"shift,omitempty"         bson:"shift,omitempty"`
}

// IsEmpty reports whether every bucket is nil. The PATCH handler rejects
// empty bodies as 400 EMPTY_PATCH (see tenant.policy.validator.go).
func (p TenantPolicyPatch) IsEmpty() bool {
	return p.Account == nil && p.Assets == nil && p.Auth == nil &&
		p.Notifications == nil && p.Onboard == nil && p.Shift == nil
}
