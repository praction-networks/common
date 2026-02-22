package events

// Sync Request/Reply Subjects — Direct NATS (NOT JetStream)
// Used by captive-portal-service to get synchronous confirmation from subscriber-service
// after publishing guest hotspot events. These run alongside the existing JetStream events.
const (
	// SyncSubscriberCreated is the subject for synchronous subscriber creation confirmation
	SyncSubscriberCreatedSubject = "captive.sync.subscriber.created"

	// SyncDeviceAdded is the subject for synchronous device addition confirmation
	SyncDeviceAddedSubject = "captive.sync.device.added"

	// SyncValidityExtended is the subject for synchronous validity extension confirmation
	SyncValidityExtendedSubject = "captive.sync.validity.extended"
)

// SyncReply is the response payload for sync request-reply subjects
type SyncReply struct {
	OK           bool   `json:"ok"`
	Error        string `json:"error,omitempty"`
	SubscriberID string `json:"subscriberId,omitempty"`
	ProfileID    string `json:"profileId,omitempty"`
}
