package returnevent

import (
	"encoding/json"
	"testing"
)

// TestReturnLineStatus_Constants — three states cover the lifecycle.
func TestReturnLineStatus_Constants(t *testing.T) {
	for _, s := range []ReturnLineStatus{
		ReturnLineStatusPending, ReturnLineStatusAccepted, ReturnLineStatusRejected,
	} {
		if s == "" {
			t.Errorf("status constant must not be empty")
		}
	}
}

// TestReturn_RoundTrip — drop-off batch with two lines.
func TestReturn_RoundTrip(t *testing.T) {
	in := Return{
		ID:            "return_01",
		TenantID:      "tenant_xyz",
		UserID:        "user_abc",
		WarehouseID:   "warehouse_main",
		Lines: []ReturnLine{
			{ID: "line_01", AssetID: strPtr("asset_aa"), Status: ReturnLineStatusPending},
			{ID: "line_02", ConsumableID: strPtr("consumable_xx"), Quantity: f64Ptr(3.5), Status: ReturnLineStatusPending},
		},
		SubmittedAtMs: 1715200000000,
	}
	b, _ := json.Marshal(in)
	var got Return
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Lines) != 2 || got.Lines[0].Status != ReturnLineStatusPending {
		t.Errorf("lines mismatch: %+v", got)
	}
}

func strPtr(s string) *string   { return &s }
func f64Ptr(f float64) *float64 { return &f }
