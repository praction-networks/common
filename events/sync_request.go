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

	// SyncBillingValidatePlanAssignment is used by subscriber-service to validate
	// that a plan can be assigned to a broadband connection via billing-service price book lookup
	SyncBillingValidatePlanAssignment Subject = "billing.sync.validate.plan.assignment"
)

// SyncReply is the response payload for sync request-reply subjects
type SyncReply struct {
	OK           bool   `json:"ok"`
	Error        string `json:"error,omitempty"`
	SubscriberID string `json:"subscriberId,omitempty"`
	ProfileID    string `json:"profileId,omitempty"`
}

// ValidatePlanAssignmentRequest is sent by subscriber-service to billing-service
// to validate that a plan can be assigned to a broadband connection
type ValidatePlanAssignmentRequest struct {
	PlanID             string  `json:"planId"`
	SellerTenantID     string  `json:"sellerTenantId"`
	BuyerTenantID      string  `json:"buyerTenantId"`
	BuyerTenantType    string  `json:"buyerTenantType"`
	PriceType          string  `json:"priceType"`          // WHOLESALE or RETAIL
	TargetTenantID     *string `json:"targetTenantId,omitempty"` // Optional specific tenant override
}

// ValidatePlanAssignmentResponse is the response from billing-service
// confirming whether a plan can be assigned and at what price
type ValidatePlanAssignmentResponse struct {
	OK           bool    `json:"ok"`
	Error        string  `json:"error,omitempty"`
	PriceBookID  string  `json:"priceBookId,omitempty"`
	Price        float64 `json:"price,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	PriceBookScope string `json:"priceBookScope,omitempty"`
}
