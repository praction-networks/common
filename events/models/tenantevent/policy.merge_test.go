package tenantevent

import (
	"reflect"
	"sort"
	"testing"
)

func TestMergeTenantPolicy_VersionBumps(t *testing.T) {
	current := Defaults() // Version 1
	patch := TenantPolicyPatch{}
	merged, _ := MergeTenantPolicy(current, patch)
	if merged.Version != 2 {
		t.Errorf("Version after merge: got %d, want 2", merged.Version)
	}
}

func TestMergeTenantPolicy_SparseBucketSet(t *testing.T) {
	current := Defaults()
	patch := TenantPolicyPatch{
		Assets: &PolicyAssets{DropoffLockWindowHours: 6},
	}
	merged, _ := MergeTenantPolicy(current, patch)
	if merged.Assets.DropoffLockWindowHours != 6 {
		t.Errorf("DropoffLockWindowHours: got %d, want 6", merged.Assets.DropoffLockWindowHours)
	}
	// Other buckets unchanged
	if merged.Auth.AccessTokenTtlMinutes != 10 {
		t.Errorf("Auth bucket should be unchanged: got %d", merged.Auth.AccessTokenTtlMinutes)
	}
}

func TestMergeTenantPolicy_BoolPointerVsZero(t *testing.T) {
	bTrue := true
	current := Defaults()
	current.Account.AllowMobileChange = &bTrue
	// Patch with a different bucket; AllowMobileChange must NOT be cleared.
	patch := TenantPolicyPatch{
		Account: &PolicyAccount{}, // empty account bucket
	}
	merged, _ := MergeTenantPolicy(current, patch)
	if merged.Account.AllowMobileChange == nil || *merged.Account.AllowMobileChange != true {
		t.Errorf("AllowMobileChange should be preserved: got %v", merged.Account.AllowMobileChange)
	}
	// Patch with explicit *bool false should override
	bFalse := false
	patch2 := TenantPolicyPatch{
		Account: &PolicyAccount{AllowMobileChange: &bFalse},
	}
	merged2, _ := MergeTenantPolicy(current, patch2)
	if merged2.Account.AllowMobileChange == nil || *merged2.Account.AllowMobileChange != false {
		t.Errorf("AllowMobileChange should be set to false: got %v", merged2.Account.AllowMobileChange)
	}
}

func TestMergeTenantPolicy_SliceReplaces(t *testing.T) {
	current := Defaults()
	patch := TenantPolicyPatch{
		Shift: &PolicyShift{BreakReasons: []string{"Custom1", "Custom2"}},
	}
	merged, _ := MergeTenantPolicy(current, patch)
	if len(merged.Shift.BreakReasons) != 2 || merged.Shift.BreakReasons[0] != "Custom1" {
		t.Errorf("BreakReasons should be replaced: got %v", merged.Shift.BreakReasons)
	}
}

func TestMergeTenantPolicy_ChangedKeysAccurate(t *testing.T) {
	current := Defaults()
	patch := TenantPolicyPatch{
		Assets: &PolicyAssets{DropoffLockWindowHours: 6}, // changed: 4 -> 6
		Auth:   &PolicyAuth{AccessTokenTtlMinutes: 10},   // unchanged: 10 -> 10
	}
	_, changedKeys := MergeTenantPolicy(current, patch)
	sort.Strings(changedKeys)
	want := []string{"assets.dropoffLockWindowHours"}
	if !reflect.DeepEqual(changedKeys, want) {
		t.Errorf("changedKeys: got %v, want %v", changedKeys, want)
	}
}
