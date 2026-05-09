package tenantevent

// TenantPolicy is the root tenant-configuration struct. One bucket per domain.
// Each bucket lives in its own file; this file only composes them.
//
// Buckets are named with the `Policy` prefix to disambiguate from any other
// shared types in the tenantevent package.
type TenantPolicy struct {
	Version       int                 `json:"version,omitempty"       bson:"version,omitempty"`
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
		Version: 1,
		Account: PolicyAccount{},
		Assets:  PolicyAssets{DropoffLockWindowHours: 4, RecoveryAcknowledgement: "NONE"},
		Auth:    PolicyAuth{AccessTokenTtlMinutes: 10},
		Shift: PolicyShift{
			FirstBreakSkipsReason:         true,
			SubsequentBreaksRequireReason: true,
			BreakReasons:                  []string{"Lunch", "Fuel", "Customer", "Personal", "Other"},
			EndNudgeGraceMinutes:          30,
			IdleMinutesThreshold:          5,
			LocationTrackingEnabled:       true,
		},
	}
}
