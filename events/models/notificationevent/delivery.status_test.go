package notificationevent

import "testing"

// TestNotificationDeliveryStatus_AllValid — six states track the lifecycle of
// a single notification across queue, transport, device, and user-interaction.
func TestNotificationDeliveryStatus_AllValid(t *testing.T) {
	for _, v := range []NotificationDeliveryStatus{
		NotificationDeliveryStatusQueued,
		NotificationDeliveryStatusSent,
		NotificationDeliveryStatusDelivered,
		NotificationDeliveryStatusRead,
		NotificationDeliveryStatusFailed,
		NotificationDeliveryStatusSuppressed,
	} {
		if !v.Valid() {
			t.Errorf("expected %q to be Valid()", v)
		}
	}
}
