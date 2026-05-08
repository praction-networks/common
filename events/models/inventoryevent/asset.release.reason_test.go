package inventoryevent

import "testing"

// TestAssetReleaseReason_ChurnRecovery — CHURN_RECOVERY is the only net-new
// value in the field-central rollout. UPGRADE_SWAP / LOST_BY_CUSTOMER /
// END_OF_LEASE reuse existing SWAP / LOST / END_OF_LIFE.
//
// Coverage:
//   - constant exists and equals "CHURN_RECOVERY"
//   - present in validAssetReleaseReasons map (so NormalizeAssetReleaseReason
//     accepts it and AssetReleaseReasonOneOf will include it for oneof= tags)
//   - normalize round-trips the canonical and lowercase forms
func TestAssetReleaseReason_ChurnRecovery(t *testing.T) {
	if string(AssetReleaseReasonChurnRecovery) != "CHURN_RECOVERY" {
		t.Errorf("constant value: got %q, want CHURN_RECOVERY", AssetReleaseReasonChurnRecovery)
	}
	if !validAssetReleaseReasons[AssetReleaseReasonChurnRecovery] {
		t.Errorf("CHURN_RECOVERY missing from validAssetReleaseReasons map")
	}
	for _, in := range []string{"CHURN_RECOVERY", "churn_recovery", "  churn_recovery  "} {
		if got := NormalizeAssetReleaseReason(in); got != AssetReleaseReasonChurnRecovery {
			t.Errorf("Normalize(%q): got %q, want CHURN_RECOVERY", in, got)
		}
	}
}
