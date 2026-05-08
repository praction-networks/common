package notificationevent

import "testing"

// TestNotificationCriticality_AllValid — three tiers per design: INFO, TIMELY, CRITICAL.
// CRITICAL is the gated tier (server-enforced whitelist per notification.md §6).
func TestNotificationCriticality_AllValid(t *testing.T) {
	for _, v := range []NotificationCriticality{
		NotificationCriticalityInfo,
		NotificationCriticalityTimely,
		NotificationCriticalityCritical,
	} {
		if !v.Valid() {
			t.Errorf("expected %q to be Valid()", v)
		}
	}
	if NotificationCriticality("URGENT").Valid() {
		t.Errorf("Valid() must reject unknown values")
	}
}
