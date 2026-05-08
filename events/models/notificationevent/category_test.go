package notificationevent

import "testing"

// TestNotificationCategory_AllValid — every category enumerated in
// notifications-prd §4 must be Valid(). Unknown categories must reject.
func TestNotificationCategory_AllValid(t *testing.T) {
	allKnown := []NotificationCategory{
		// Ticket / job lifecycle
		NotificationCategoryJobAssigned,
		NotificationCategoryJobReassigned,
		NotificationCategoryJobCancelled,
		NotificationCategoryJobRescheduled,
		// SLA / risk
		NotificationCategorySLAWarning,
		NotificationCategorySLABreach,
		NotificationCategoryAppointmentReminder,
		NotificationCategoryMissedAppointment,
		// Customer interaction
		NotificationCategoryCustomerArrived,
		NotificationCategoryCustomerOTPSent,
		NotificationCategoryCustomerAgreementSigned,
		// Payments
		NotificationCategoryPaymentReceived,
		NotificationCategoryPaymentFailed,
		NotificationCategoryPaymentOverdue,
		// Draft / housekeeping
		NotificationCategoryDraftExpiring,
		NotificationCategoryDraftAutoPurged,
		// System / org
		NotificationCategoryAnnouncement,
		NotificationCategoryShiftReminder,
		NotificationCategoryTenantBroadcast,
		NotificationCategorySystem,
	}
	for _, c := range allKnown {
		if !c.Valid() {
			t.Errorf("expected %q to be Valid()", c)
		}
	}
	if NotificationCategory("BOGUS_VALUE").Valid() {
		t.Errorf("Valid() must reject unknown values")
	}
}
