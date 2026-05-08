package inventoryevent

import (
	"encoding/json"
	"testing"
)

// TestConsumable_RoundTrip — catalogue shape; quantity is per Unit.
func TestConsumable_RoundTrip(t *testing.T) {
	in := Consumable{
		ID:       "consumable_01",
		TenantID: "tenant_xyz",
		Name:     "CAT6 Cable",
		Unit:     "metre",
		SKU:      "CAB-CAT6",
	}
	b, _ := json.Marshal(in)
	var got Consumable
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("got %+v, want %+v", got, in)
	}
}

// TestConsumableHolding_QuantityRoundTrip — holdings track decimal-aware
// metres for cable; integer counts for connectors. JSON is float64 to
// preserve fractional metres.
func TestConsumableHolding_QuantityRoundTrip(t *testing.T) {
	in := ConsumableHolding{
		ID:           "holding_01",
		TenantID:     "tenant_xyz",
		UserID:       "user_abc",
		ConsumableID: "consumable_01",
		Quantity:     12.5,
	}
	b, _ := json.Marshal(in)
	var got ConsumableHolding
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Quantity != 12.5 {
		t.Errorf("quantity: got %v, want 12.5", got.Quantity)
	}
}

// TestConsumableMovement_RoundTrip — movement is an immutable audit row.
func TestConsumableMovement_RoundTrip(t *testing.T) {
	in := ConsumableMovement{
		ID:           "move_01",
		TenantID:     "tenant_xyz",
		ConsumableID: "consumable_01",
		Kind:         "ISSUED",
		Quantity:     5.0,
		FromUserID:   "warehouse_wm",
		ToUserID:     "user_abc",
		AtMs:         1715200000000,
	}
	b, _ := json.Marshal(in)
	var got ConsumableMovement
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("got %+v, want %+v", got, in)
	}
}
