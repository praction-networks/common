package peerassistevent

import (
	"encoding/json"
	"testing"
)

func TestPeerAssistRequest_RoundTrip(t *testing.T) {
	in := PeerAssistRequest{
		ID:            "par_01",
		TenantID:      "tenant_xyz",
		RequesterID:   "user_a",
		Status:        "PENDING_APPROVAL",
		ReasonCode:    "NEEDS_EXTRA_CABLE",
		Note:          "ran out at site 12",
		RequestedAtMs: 1715200000000,
	}
	b, _ := json.Marshal(in)
	var got PeerAssistRequest
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("got %+v, want %+v", got, in)
	}
}

func TestPeerHandoff_RoundTrip(t *testing.T) {
	in := PeerHandoff{
		ID:         "ph_01",
		TenantID:   "tenant_xyz",
		FromUserID: "user_a",
		ToUserID:   "user_b",
		Status:     "AWAITING_CONFIRM",
		Lines: []PeerHandoffLine{
			{ID: "line_01", ConsumableID: strPtr("consumable_x"), Quantity: f64Ptr(5.0)},
		},
		StartedAtMs: 1715200000000,
	}
	b, _ := json.Marshal(in)
	var got PeerHandoff
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Lines) != 1 || got.Lines[0].ID != "line_01" {
		t.Errorf("lines mismatch: %+v", got)
	}
}

func strPtr(s string) *string   { return &s }
func f64Ptr(f float64) *float64 { return &f }
