package tenantevent

// TenantPolicy is the root tenant-configuration struct. One bucket per domain.
// Each bucket lives in its own file; this file only composes them.
//
// Buckets are named with the `Policy` prefix to disambiguate from any other
// shared types in the tenantevent package.
type TenantPolicy struct {
	Account       PolicyAccount       `json:"account,omitempty"       bson:"account,omitempty"`
	Assets        PolicyAssets        `json:"assets,omitempty"        bson:"assets,omitempty"`
	Auth          PolicyAuth          `json:"auth,omitempty"          bson:"auth,omitempty"`
	Notifications PolicyNotifications `json:"notifications,omitempty" bson:"notifications,omitempty"`
	Onboard       PolicyOnboard       `json:"onboard,omitempty"       bson:"onboard,omitempty"`
	Shift         PolicyShift         `json:"shift,omitempty"         bson:"shift,omitempty"`
}

// Defaults returns a TenantPolicy with the baseline values documented in the
// design spec (§15.3 resolved decisions). Use at construction time; unmarshal
// of stored documents does not auto-apply defaults.
func Defaults() TenantPolicy {
	return TenantPolicy{
		Account: PolicyAccount{EmailChangeAllowed: false},
		Assets:  PolicyAssets{DropoffLockWindowHours: 4, RecoveryAcknowledgement: "NONE"},
		Auth:    PolicyAuth{AccessTokenTtlMinutes: 10},
	}
}
