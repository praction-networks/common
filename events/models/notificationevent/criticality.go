package notificationevent

// NotificationCriticality maps to iOS interruption-level (and Android importance).
// CRITICAL maps to iOS critical-alert (entitlement-gated) — see notifications-prd §16.
type NotificationCriticality string

const (
	NotificationCriticalityInfo     NotificationCriticality = "INFO"
	NotificationCriticalityTimely   NotificationCriticality = "TIMELY"
	NotificationCriticalityCritical NotificationCriticality = "CRITICAL"
)

func (c NotificationCriticality) Valid() bool {
	switch c {
	case NotificationCriticalityInfo, NotificationCriticalityTimely, NotificationCriticalityCritical:
		return true
	}
	return false
}
