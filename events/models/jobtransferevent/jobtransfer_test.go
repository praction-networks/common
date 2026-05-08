package jobtransferevent

import (
	"encoding/json"
	"testing"
)

func TestJobTransferRequest_RoundTrip(t *testing.T) {
	in := JobTransferRequest{
		ID:                "jtr_01",
		TenantID:          "tenant_xyz",
		TicketID:          "ticket_abc",
		FromUserID:        "user_a",
		Status:            "PENDING_APPROVAL",
		Reason:            "EMERGENCY",
		RequestedAtMs:     1715200000000,
		LifetimeExpiresAt: 1715214400000,
	}
	b, _ := json.Marshal(in)
	var got JobTransferRequest
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("got %+v, want %+v", got, in)
	}
}

func TestJobTransferOffer_RoundTrip(t *testing.T) {
	in := JobTransferOffer{
		ID:          "jto_01",
		RequestID:   "jtr_01",
		ToUserID:    "user_b",
		Status:      "OFFERED",
		OfferedAtMs: 1715200300000,
		ExpiresAtMs: 1715200900000,
	}
	b, _ := json.Marshal(in)
	var got JobTransferOffer
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != in {
		t.Errorf("got %+v, want %+v", got, in)
	}
}
