package subscriberevent

import "testing"

// TestSubscriberStatus_NewValues — these two values were added per the
// field-central onboarding-wizard contract (subscriber-prd §7.2). The
// existing SubscriberStatus type has no Valid() method (plain type+const
// pattern), so the test asserts the constants are non-empty and distinct
// from the existing values.
func TestSubscriberStatus_NewValues(t *testing.T) {
	cases := []struct {
		name string
		v    SubscriberStatus
		want string
	}{
		{"PENDING_KYC", SubscriberStatusPendingKYC, "PENDING_KYC"},
		{"PENDING_PAYMENT", SubscriberStatusPendingPayment, "PENDING_PAYMENT"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if string(tc.v) != tc.want {
				t.Errorf("got %q, want %q", tc.v, tc.want)
			}
			for _, existing := range []SubscriberStatus{SubscriberStatusActive, SubscriberStatusInactive, SubscriberStatusOnboarding} {
				if tc.v == existing {
					t.Errorf("new value %q collides with existing %q", tc.v, existing)
				}
			}
		})
	}
}
