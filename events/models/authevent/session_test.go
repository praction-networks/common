package authevent

import (
	"encoding/json"
	"testing"
)

// TestAuthSessionRevokedEvent_RoundTrip — payload shape per design §4.
func TestAuthSessionRevokedEvent_RoundTrip(t *testing.T) {
	in := AuthSessionRevokedEvent{
		UserID:    "user_xyz",
		TenantID:  "tenant_abc",
		SessionID: "sess_123",
		Reason:    "user_initiated",
		AtMs:      1715200000000,
	}
	b, _ := json.Marshal(in)
	var got AuthSessionRevokedEvent
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, in)
	}
}

// TestAuthSessionRevokedEvent_AllSentinel — sign-out-all writes SessionID="all"
func TestAuthSessionRevokedEvent_AllSentinel(t *testing.T) {
	in := AuthSessionRevokedEvent{
		UserID:    "user_xyz",
		TenantID:  "tenant_abc",
		SessionID: "all",
		Reason:    "sign_out_all",
		AtMs:      1715200000000,
	}
	b, _ := json.Marshal(in)
	var got AuthSessionRevokedEvent
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, in)
	}
}
