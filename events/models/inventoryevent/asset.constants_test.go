package inventoryevent

import "testing"

// TestAssetStatus_PendingReturn — added per assets-prd §11.2 two-step drop-off.
// FE submits, asset locks in PENDING_RETURN until WM accepts (advances to
// RETURNED) or rejects (asset bounces back to ASSIGNED/INSTALLED).
//
// State-machine reasoning:
//   - From ASSIGNED: FE who has the device in their kit submits a return →
//     PENDING_RETURN.
//   - From INSTALLED: FE recovering a device from a customer site submits a
//     return → PENDING_RETURN. (The recovery flow first transitions
//     INSTALLED → ASSIGNED; this transition is for the case where the FE
//     returns the unit directly to the warehouse without a kit-handoff in
//     between, which the existing INSTALLED → RETURNED edge already supports
//     but doesn't model the WM-accept gate.)
//   - From PENDING_RETURN → RETURNED: WM accepts the drop-off (custody
//     flips on accept, not on submit).
//   - From PENDING_RETURN → ASSIGNED / INSTALLED: WM rejects (or auto-release
//     after dropoffLockWindowHours expires); the asset goes back to whichever
//     state the FE held it in.
func TestAssetStatus_PendingReturn(t *testing.T) {
	if string(AssetStatusPendingReturn) != "PENDING_RETURN" {
		t.Errorf("constant value: got %q, want PENDING_RETURN", AssetStatusPendingReturn)
	}

	for _, from := range []AssetStatus{AssetStatusAssigned, AssetStatusInstalled} {
		if !IsValidTransition(from, AssetStatusPendingReturn) {
			t.Errorf("expected %s → PENDING_RETURN to be valid", from)
		}
	}

	for _, to := range []AssetStatus{AssetStatusReturned, AssetStatusAssigned, AssetStatusInstalled} {
		if !IsValidTransition(AssetStatusPendingReturn, to) {
			t.Errorf("expected PENDING_RETURN → %s to be valid", to)
		}
	}

	if _, ok := ValidAssetTransitions[AssetStatusPendingReturn]; !ok {
		t.Errorf("PENDING_RETURN missing as a key in ValidAssetTransitions")
	}
}
