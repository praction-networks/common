package tenantevent

import (
	"encoding/json"
	"testing"

	"github.com/praction-networks/common/events/models/notificationevent"
)

// TestDefaults — Defaults() returns a struct with the documented baseline
// values: emailChangeAllowed=nil (not set), dropoffLockWindowHours=4,
// recoveryAcknowledgement=NONE, accessTokenTtlMinutes=10.
//
// EmailChangeAllowed is now *bool; nil means "not configured" (sparse PATCH
// semantics). Defaults() no longer initialises it so callers can distinguish
// "server default" from "explicitly set to false".
func TestDefaults(t *testing.T) {
	d := Defaults()
	if d.Account.EmailChangeAllowed != nil {
		t.Errorf("Account.EmailChangeAllowed default: got %v, want nil", d.Account.EmailChangeAllowed)
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
// Now *bool: pointer-to-true must survive a JSON round-trip.
func TestPolicyAccount_RoundTrip(t *testing.T) {
	bTrue := true
	in := PolicyAccount{EmailChangeAllowed: &bTrue}
	b, _ := json.Marshal(in)
	var got PolicyAccount
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.EmailChangeAllowed == nil || *got.EmailChangeAllowed != true {
		t.Errorf("got %v, want pointer to true", got.EmailChangeAllowed)
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

// TestTenantPolicy_VersionField — Version is bumped server-side on every
// successful PATCH and surfaces to consumers via the tenant.policy.updated
// event so caches can detect stale state via integer compare.
func TestTenantPolicy_VersionField(t *testing.T) {
	d := Defaults()
	if d.Version != 1 {
		t.Errorf("Defaults().Version: got %d, want 1", d.Version)
	}
	in := TenantPolicy{Version: 7, Account: PolicyAccount{}}
	b, _ := json.Marshal(in)
	var got TenantPolicy
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Version != 7 {
		t.Errorf("round-trip Version: got %d, want 7", got.Version)
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

// TestPolicyAccount_DeferredKeys_RoundTrip — the 7 §11.5 deferred-but-listed
// keys plus the migrated EmailChangeAllowed (now *bool to disambiguate
// "set to false" from "not sent" under sparse PATCH semantics).
func TestPolicyAccount_DeferredKeys_RoundTrip(t *testing.T) {
	bTrue, bFalse := true, false
	rate := 7.5
	in := PolicyAccount{
		EmailChangeAllowed:   &bFalse,
		AllowMobileChange:    &bTrue,
		AllowWhatsappChange:  &bFalse,
		FuelMode:             "PER_KM",
		FuelRatePerKm:        &rate,
		LeaveCategories:      []string{"PAID", "SICK"},
		VehicleRequired:      &bTrue,
		DocumentVaultEnabled: &bTrue,
	}
	b, _ := json.Marshal(in)
	var got PolicyAccount
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.EmailChangeAllowed == nil || *got.EmailChangeAllowed != false {
		t.Errorf("EmailChangeAllowed: got %v", got.EmailChangeAllowed)
	}
	if got.FuelMode != "PER_KM" || got.FuelRatePerKm == nil || *got.FuelRatePerKm != 7.5 {
		t.Errorf("fuel: got mode=%q rate=%v", got.FuelMode, got.FuelRatePerKm)
	}
	if len(got.LeaveCategories) != 2 || got.LeaveCategories[0] != "PAID" {
		t.Errorf("leave categories: got %v", got.LeaveCategories)
	}
}

// TestPolicyAccount_NilVsZero — nil pointer (omitted) and *bool to false must
// be distinguishable in the wire form so MergeTenantPolicy can preserve the
// "client did not send" semantic.
func TestPolicyAccount_NilVsZero(t *testing.T) {
	bFalse := false
	a := PolicyAccount{}
	b, _ := json.Marshal(a)
	if string(b) != "{}" {
		t.Errorf("nil pointers must omit; got %s", string(b))
	}
	a2 := PolicyAccount{EmailChangeAllowed: &bFalse}
	b2, _ := json.Marshal(a2)
	if string(b2) != `{"emailChangeAllowed":false}` {
		t.Errorf("set-to-false must encode; got %s", string(b2))
	}
}

// TestPolicyShift_RoundTrip — 7 fields per backend-contract §11.1, used by
// the Account hub break/idle/heartbeat flows and by server-side end-nudge.
func TestPolicyShift_RoundTrip(t *testing.T) {
	max := 45
	in := PolicyShift{
		FirstBreakSkipsReason:         true,
		SubsequentBreaksRequireReason: true,
		BreakReasons:                  []string{"Lunch", "Customer"},
		MaxBreakMinutes:               &max,
		EndNudgeGraceMinutes:          30,
		IdleMinutesThreshold:          5,
		LocationTrackingEnabled:       true,
	}
	b, _ := json.Marshal(in)
	var got PolicyShift
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !got.FirstBreakSkipsReason || got.MaxBreakMinutes == nil || *got.MaxBreakMinutes != 45 {
		t.Errorf("round-trip mismatch: %+v", got)
	}
	if len(got.BreakReasons) != 2 || got.BreakReasons[0] != "Lunch" {
		t.Errorf("break reasons: %v", got.BreakReasons)
	}
}

// TestPolicyOnboard_RoundTrip — 2 fields per backend-contract §11.2.
// enabledFeatures.core.isUserKYCEnabled lives on TenantModel.EnabledFeatures
// (existing) and is NOT duplicated here.
func TestPolicyOnboard_RoundTrip(t *testing.T) {
	in := PolicyOnboard{
		AgreementChannel: "OTP",
		OtpChannel:       "WHATSAPP",
	}
	b, _ := json.Marshal(in)
	var got PolicyOnboard
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("got %+v, want %+v", got, in)
	}
}

// TestPolicyNotifications_SupportChannels_RoundTrip — §11.4 support channels.
// Strings are URLs/numbers/emails per tenant choice.
func TestPolicyNotifications_SupportChannels_RoundTrip(t *testing.T) {
	in := PolicyNotifications{
		SupportChannels: PolicySupportChannels{
			Whatsapp: "+919876543210",
			Call:     "+918888888888",
			Email:    "support@example.com",
		},
	}
	b, _ := json.Marshal(in)
	var got PolicyNotifications
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SupportChannels.Whatsapp != "+919876543210" || got.SupportChannels.Email != "support@example.com" {
		t.Errorf("support channels: got %+v", got.SupportChannels)
	}
}

// TestTenantPolicyPatch_NilBucketsIgnored — nil buckets in the wire form
// represent "client did not send this bucket"; only non-nil buckets are
// considered for merge.
func TestTenantPolicyPatch_NilBucketsIgnored(t *testing.T) {
	body := []byte(`{"assets":{"dropoffLockWindowHours":6}}`)
	var p TenantPolicyPatch
	if err := json.Unmarshal(body, &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Account != nil || p.Auth != nil || p.Shift != nil {
		t.Errorf("absent buckets must be nil; got Account=%v Auth=%v Shift=%v", p.Account, p.Auth, p.Shift)
	}
	if p.Assets == nil {
		t.Fatalf("Assets must be non-nil when sent")
	}
	if p.Assets.DropoffLockWindowHours != 6 {
		t.Errorf("DropoffLockWindowHours: got %d", p.Assets.DropoffLockWindowHours)
	}
}

// TestTenantPolicyPatch_EmptyBucketRetained — an empty {} bucket is
// distinguishable from nil; means "client sent the bucket but no fields
// inside changed". Unusual but legal.
func TestTenantPolicyPatch_EmptyBucketRetained(t *testing.T) {
	body := []byte(`{"assets":{}}`)
	var p TenantPolicyPatch
	if err := json.Unmarshal(body, &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Assets == nil {
		t.Errorf("empty bucket must round-trip as non-nil pointer")
	}
}

// TestTenantPolicyUpdatedEvent_RoundTrip — payload shape per design §5.5.
// Subject is "tenant.policy.updated" (no .v1 suffix per Wave 0 convention).
func TestTenantPolicyUpdatedEvent_RoundTrip(t *testing.T) {
	in := TenantPolicyUpdatedEvent{
		TenantID:    "tenant_xyz",
		Version:     4,
		UpdatedAtMs: 1715200000000,
		ChangedKeys: []string{"assets.dropoffLockWindowHours"},
	}
	b, _ := json.Marshal(in)
	var got TenantPolicyUpdatedEvent
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.TenantID != "tenant_xyz" || got.Version != 4 || len(got.ChangedKeys) != 1 {
		t.Errorf("round-trip mismatch: got %+v", got)
	}
}

// TestPolicyShift_Defaults — Defaults() lands the documented baseline values.
func TestPolicyShift_Defaults(t *testing.T) {
	d := Defaults()
	if !d.Shift.FirstBreakSkipsReason {
		t.Errorf("FirstBreakSkipsReason default: want true")
	}
	if !d.Shift.SubsequentBreaksRequireReason {
		t.Errorf("SubsequentBreaksRequireReason default: want true")
	}
	wantReasons := []string{"Lunch", "Fuel", "Customer", "Personal", "Other"}
	if len(d.Shift.BreakReasons) != 5 || d.Shift.BreakReasons[0] != wantReasons[0] {
		t.Errorf("BreakReasons default: got %v, want %v", d.Shift.BreakReasons, wantReasons)
	}
	if d.Shift.EndNudgeGraceMinutes != 30 {
		t.Errorf("EndNudgeGraceMinutes default: got %d, want 30", d.Shift.EndNudgeGraceMinutes)
	}
	if d.Shift.IdleMinutesThreshold != 5 {
		t.Errorf("IdleMinutesThreshold default: got %d, want 5", d.Shift.IdleMinutesThreshold)
	}
	if !d.Shift.LocationTrackingEnabled {
		t.Errorf("LocationTrackingEnabled default: want true")
	}
}
