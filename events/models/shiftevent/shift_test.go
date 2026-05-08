package shiftevent

import (
	"encoding/json"
	"testing"
)

// TestShift_RoundTrip — closed-shift document with a single break. JSON shape
// is consumed by user-service /me/performance aggregator and Account-hub history.
func TestShift_RoundTrip(t *testing.T) {
	in := Shift{
		ID:       "shift_01J9X",
		UserID:   "user_abc",
		TenantID: "tenant_xyz",
		State:    "CLOSED",
		OpenedAt: 1715200000000,
		ClosedAt: 1715235600000,
		Breaks: []ShiftBreak{
			{StartedAt: 1715210000000, EndedAt: 1715212400000, Reason: "LUNCH"},
		},
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got Shift
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID != in.ID || len(got.Breaks) != 1 || got.Breaks[0].Reason != "LUNCH" {
		t.Errorf("round-trip mismatch: got %+v", got)
	}
}
