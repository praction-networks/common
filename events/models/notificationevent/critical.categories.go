package notificationevent

// CriticalCategoriesPlatform is the hard-coded whitelist of NotificationCategory
// values that may be marked CRITICAL. Tenant policy can only narrow this list,
// not extend it (server-enforced via tenantevent.PolicyNotifications.Validate).
//
// Source: notifications-prd §6.1 ("What qualifies"):
//
//	"Server-side, only these categories may set priority === 'CRITICAL':
//	 - SLA_BREACH
//	 - MISSED_APPOINTMENT
//	 - JOB_CANCELLED when the FE is en-route (tracked via the existing
//	   postCheckIn ARRIVAL flow)"
//
// JOB_CANCELLED is conditionally CRITICAL (only when the FE is en-route);
// the platform whitelist still includes it because the en-route gate is a
// runtime check, not a category-level restriction. Tenants that never use
// the en-route flow can narrow it out via PolicyNotifications.
//
// See also: notifications-prd §4.3 (category → priority defaults), Sweep 4 Q8.
var CriticalCategoriesPlatform = []NotificationCategory{
	NotificationCategorySLABreach,
	NotificationCategoryMissedAppointment,
	NotificationCategoryJobCancelled,
}
