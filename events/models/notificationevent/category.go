package notificationevent

// NotificationCategory tags the domain a notification belongs to.
// Source of truth: notifications-prd §4.1 (NotificationCategory union type).
//
// Naming convention: Go const NotificationCategory<PascalCase> = "<UPPER_SNAKE>"
// where the wire value matches the TypeScript literal verbatim so server,
// dashboard, and mobile share one vocabulary.
type NotificationCategory string

const (
	// ── Ticket / job lifecycle ───────────────────────────────────────────

	// NotificationCategoryJobAssigned — a ticket has been assigned to the
	// recipient FE. Default priority HIGH; carries Accept/Decline actions.
	NotificationCategoryJobAssigned NotificationCategory = "JOB_ASSIGNED"
	// NotificationCategoryJobReassigned — an existing assignment was moved
	// to a different FE (recipient is the new owner).
	NotificationCategoryJobReassigned NotificationCategory = "JOB_REASSIGNED"
	// NotificationCategoryJobCancelled — a previously assigned job was
	// cancelled by dispatch or the customer. Default priority CRITICAL.
	NotificationCategoryJobCancelled NotificationCategory = "JOB_CANCELLED"
	// NotificationCategoryJobRescheduled — appointment window for an
	// assigned job was changed.
	NotificationCategoryJobRescheduled NotificationCategory = "JOB_RESCHEDULED"

	// ── SLA / risk ───────────────────────────────────────────────────────

	// NotificationCategorySLAWarning — an SLA timer is about to breach.
	// Default priority HIGH.
	NotificationCategorySLAWarning NotificationCategory = "SLA_WARNING"
	// NotificationCategorySLABreach — an SLA timer has already breached.
	// Default priority CRITICAL.
	NotificationCategorySLABreach NotificationCategory = "SLA_BREACH"
	// NotificationCategoryAppointmentReminder — scheduled reminder ahead of
	// an upcoming visit window.
	NotificationCategoryAppointmentReminder NotificationCategory = "APPOINTMENT_REMINDER"
	// NotificationCategoryMissedAppointment — FE did not arrive within the
	// committed window. Default priority CRITICAL.
	NotificationCategoryMissedAppointment NotificationCategory = "MISSED_APPOINTMENT"

	// ── Customer interaction ─────────────────────────────────────────────

	// NotificationCategoryCustomerArrived — the customer has arrived at the
	// venue / signalled readiness for the visit. Default priority HIGH.
	NotificationCategoryCustomerArrived NotificationCategory = "CUSTOMER_ARRIVED"
	// NotificationCategoryCustomerOTPSent — an OTP has been dispatched to
	// the customer (used to confirm device-side flows).
	NotificationCategoryCustomerOTPSent NotificationCategory = "CUSTOMER_OTP_SENT"
	// NotificationCategoryCustomerAgreementSigned — the customer has signed
	// the agreement document; FE may proceed with provisioning.
	NotificationCategoryCustomerAgreementSigned NotificationCategory = "CUSTOMER_AGREEMENT_SIGNED"

	// ── Payments ─────────────────────────────────────────────────────────

	// NotificationCategoryPaymentReceived — a payment was successfully
	// captured against a subscriber/ticket.
	NotificationCategoryPaymentReceived NotificationCategory = "PAYMENT_RECEIVED"
	// NotificationCategoryPaymentFailed — payment attempt failed; usually
	// requires FE retry or escalation. Default priority HIGH.
	NotificationCategoryPaymentFailed NotificationCategory = "PAYMENT_FAILED"
	// NotificationCategoryPaymentOverdue — invoice has crossed its due
	// date; informational nudge.
	NotificationCategoryPaymentOverdue NotificationCategory = "PAYMENT_OVERDUE"

	// ── Draft / housekeeping ─────────────────────────────────────────────

	// NotificationCategoryDraftExpiring — an in-progress draft will be
	// auto-purged soon; offers Resume/Snooze/Discard.
	NotificationCategoryDraftExpiring NotificationCategory = "DRAFT_EXPIRING"
	// NotificationCategoryDraftAutoPurged — a draft was removed by the
	// retention policy; informational only.
	NotificationCategoryDraftAutoPurged NotificationCategory = "DRAFT_AUTO_PURGED"

	// ── System / org ─────────────────────────────────────────────────────

	// NotificationCategoryAnnouncement — generic in-app announcement (org
	// or product team).
	NotificationCategoryAnnouncement NotificationCategory = "ANNOUNCEMENT"
	// NotificationCategoryShiftReminder — scheduled reminder around the
	// FE's shift start/end.
	NotificationCategoryShiftReminder NotificationCategory = "SHIFT_REMINDER"
	// NotificationCategoryTenantBroadcast — tenant-wide message from an
	// admin (e.g. weather advisory, mass outage).
	NotificationCategoryTenantBroadcast NotificationCategory = "TENANT_BROADCAST"
	// NotificationCategorySystem — fallback bucket for platform-emitted
	// notifications that don't fit a more specific category.
	NotificationCategorySystem NotificationCategory = "SYSTEM"
)

// Valid reports whether c is a recognized NotificationCategory.
func (c NotificationCategory) Valid() bool {
	switch c {
	case NotificationCategoryJobAssigned,
		NotificationCategoryJobReassigned,
		NotificationCategoryJobCancelled,
		NotificationCategoryJobRescheduled,
		NotificationCategorySLAWarning,
		NotificationCategorySLABreach,
		NotificationCategoryAppointmentReminder,
		NotificationCategoryMissedAppointment,
		NotificationCategoryCustomerArrived,
		NotificationCategoryCustomerOTPSent,
		NotificationCategoryCustomerAgreementSigned,
		NotificationCategoryPaymentReceived,
		NotificationCategoryPaymentFailed,
		NotificationCategoryPaymentOverdue,
		NotificationCategoryDraftExpiring,
		NotificationCategoryDraftAutoPurged,
		NotificationCategoryAnnouncement,
		NotificationCategoryShiftReminder,
		NotificationCategoryTenantBroadcast,
		NotificationCategorySystem:
		return true
	}
	return false
}
