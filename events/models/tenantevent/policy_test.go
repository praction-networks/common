package tenantevent

import (
	"encoding/json"
	"testing"

	"github.com/praction-networks/common/events/models/notificationevent"
)

// TestDefaults — Defaults() returns a struct with the documented baseline
// values: emailChangeAllowed=false, dropoffLockWindowHours=4,
// recoveryAcknowledgement=NONE, accessTokenTtlMinutes=10.
func TestDefaults(t *testing.T) {
	d := Defaults()
	if d.Account.EmailChangeAllowed {
		t.Errorf("Account.EmailChangeAllowed default: got true, want false")
	}
	if d.Assets.DropoffLockWindowHours != 4 {
		t.Errorf("Assets.DropoffLockWindowHours default: got %d, want 4", d.Assets.DropoffLockWindowHours)
	}
	if d.Assets.RecoveryAcknowledgement != "NONE" {
		t.Errorf("Assets.RecoveryAcknowledgement default: got %q, want NONE", d.Assets.RecoveryAcknowledgement)
	}
	if d.Auth.AccessTokenTtlMinutes != 10 {
		t.Errorf("Auth.AccessTokenTtlMinutes default: got %d, want 10", d.Auth.AccessTokenTtlMinutes)
	}
}

// TestBackcompatUnmarshal — old policy documents (without the new buckets)
// must hydrate without error; missing fields fall back to zero-values.
func TestBackcompatUnmarshal(t *testing.T) {
	old := []byte(`{"account":{}, "assets":{}}`)
	var p TenantPolicy
	if err := json.Unmarshal(old, &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Auth.AccessTokenTtlMinutes != 0 {
		t.Errorf("expected zero-value 0 (defaults applied at construction): got %d", p.Auth.AccessTokenTtlMinutes)
	}
}

// TestPolicyAccount_RoundTrip — EmailChangeAllowed gates the email-change endpoint.
func TestPolicyAccount_RoundTrip(t *testing.T) {
	in := PolicyAccount{EmailChangeAllowed: true}
	b, _ := json.Marshal(in)
	var got PolicyAccount
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("got %+v, want %+v", got, in)
	}
}

// TestPolicyAssets_RoundTrip — three new keys: dropoffLockWindowHours, recoveryAcknowledgement.
func TestPolicyAssets_RoundTrip(t *testing.T) {
	in := PolicyAssets{
		DropoffLockWindowHours:  6,
		RecoveryAcknowledgement: "SIGNATURE",
	}
	b, _ := json.Marshal(in)
	var got PolicyAssets
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("got %+v, want %+v", got, in)
	}
}

// TestPolicyAuth_RoundTrip — first key: AccessTokenTtlMinutes.
func TestPolicyAuth_RoundTrip(t *testing.T) {
	in := PolicyAuth{AccessTokenTtlMinutes: 12}
	b, _ := json.Marshal(in)
	var got PolicyAuth
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.AccessTokenTtlMinutes != 12 {
		t.Errorf("got %d, want 12", got.AccessTokenTtlMinutes)
	}
}

func TestPolicyAuth_Validate(t *testing.T) {
	cases := []struct {
		v       int
		wantErr bool
	}{
		{0, false},
		{5, false},
		{10, false},
		{15, false},
		{4, true},
		{16, true},
		{60, true},
	}
	for _, tc := range cases {
		err := PolicyAuth{AccessTokenTtlMinutes: tc.v}.Validate()
		if (err != nil) != tc.wantErr {
			t.Errorf("Validate(%d) err=%v, wantErr=%v", tc.v, err, tc.wantErr)
		}
	}
}

// TestPolicyNotifications_Validate_Narrowing — tenant may narrow but not expand.
func TestPolicyNotifications_Validate_Narrowing(t *testing.T) {
	platform := notificationevent.CriticalCategoriesPlatform
	if len(platform) == 0 {
		t.Skip("CriticalCategoriesPlatform empty — populate before validating")
	}

	cases := []struct {
		name    string
		allowed []string
		wantErr bool
	}{
		{"nil means full platform list", nil, false},
		{"empty means no critical allowed", []string{}, false},
		{"subset OK", []string{string(platform[0])}, false},
		{"superset rejected", append([]string{"BOGUS_CATEGORY"}, string(platform[0])), true},
	}
	for _, tc := range cases {
		err := PolicyNotifications{CriticalCategoriesAllowed: tc.allowed}.Validate()
		if (err != nil) != tc.wantErr {
			t.Errorf("%s: err=%v, wantErr=%v", tc.name, err, tc.wantErr)
		}
	}
}
