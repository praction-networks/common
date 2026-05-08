package notificationevent

// NotificationDeliveryStatus tracks the lifecycle of one notification record.
// SUPPRESSED applies when tenant policy or critical-category whitelist filters
// the notification before transport (server-side gate).
type NotificationDeliveryStatus string

const (
	NotificationDeliveryStatusQueued     NotificationDeliveryStatus = "QUEUED"
	NotificationDeliveryStatusSent       NotificationDeliveryStatus = "SENT"
	NotificationDeliveryStatusDelivered  NotificationDeliveryStatus = "DELIVERED"
	NotificationDeliveryStatusRead       NotificationDeliveryStatus = "READ"
	NotificationDeliveryStatusFailed     NotificationDeliveryStatus = "FAILED"
	NotificationDeliveryStatusSuppressed NotificationDeliveryStatus = "SUPPRESSED"
)

func (s NotificationDeliveryStatus) Valid() bool {
	switch s {
	case NotificationDeliveryStatusQueued, NotificationDeliveryStatusSent,
		NotificationDeliveryStatusDelivered, NotificationDeliveryStatusRead,
		NotificationDeliveryStatusFailed, NotificationDeliveryStatusSuppressed:
		return true
	}
	return false
}
