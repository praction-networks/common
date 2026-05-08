package tenantevent

import "fmt"

// PolicyAuth is the per-tenant authentication-policy bucket. NEW in PR 2.
// Source: backend-contract §3.2, Sweep 4 Q7.
type PolicyAuth struct {
	AccessTokenTtlMinutes int `json:"accessTokenTtlMinutes,omitempty" bson:"accessTokenTtlMinutes,omitempty"`
}

// Validate enforces tenant-tunable limits. Zero-value (0) is allowed and means
// "use the platform default" (10 min). Set values must fall within 5–15.
func (a PolicyAuth) Validate() error {
	if a.AccessTokenTtlMinutes == 0 {
		return nil
	}
	if a.AccessTokenTtlMinutes < 5 || a.AccessTokenTtlMinutes > 15 {
		return fmt.Errorf("AccessTokenTtlMinutes must be between 5 and 15 (got %d)", a.AccessTokenTtlMinutes)
	}
	return nil
}
