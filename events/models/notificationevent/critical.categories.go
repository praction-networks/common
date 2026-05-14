package notificationevent

// CriticalCategoriesPlatform returns the hard-coded whitelist of
// NotificationCategory values that may be marked CRITICAL. Each call returns a
// fresh slice so callers cannot mutate the platform list. Tenant policy can
// only narrow this list, not extend it (server-enforced via
// tenantevent.PolicyNotifications.Validate).
//
// Source: notifications-prd §6.1 ("What qualifies") + backend-contract §9.5:
//
//	"Server-side, only these categories may set priority === 'CRITICAL':
//	 - SLA_BREACH
//	 - MISSED_APPOINTMENT
//	 - JOB_CANCELLED when the FE is en-route (tracked via the existing
//	   postCheckIn ARRIVAL flow)
//	 - JOB_TRANSFER_PEER_OFFERED (added in backend-contract §9.5)"
//
// JOB_CANCELLED is conditionally CRITICAL (only when the FE is en-route);
// the platform whitelist still includes it because the en-route gate is a
// runtime check, not a category-level restriction. Tenants that never use
// the en-route flow can narrow it out via PolicyNotifications.
//
// JOB_TRANSFER_PEER_OFFERED is CRITICAL because the peer must act within the
// offer window; missing it causes the transfer to expire automatically.
//
// See also: notifications-prd §4.3 (category → priority defaults), Sweep 4 Q8.
func CriticalCategoriesPlatform() []NotificationCategory {
	return []NotificationCategory{
		NotificationCategorySLABreach,
		NotificationCategoryMissedAppointment,
		NotificationCategoryJobCancelled,
		NotificationCategoryJobTransferPeerOffered,
	}
}
