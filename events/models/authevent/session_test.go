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
	if got.UserID != "user_xyz" || got.SessionID != "sess_123" || got.Reason != "user_initiated" {
		t.Errorf("round-trip mismatch: %+v", got)
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
	if got.SessionID != "all" {
		t.Errorf("SessionID: got %q, want all", got.SessionID)
	}
}
