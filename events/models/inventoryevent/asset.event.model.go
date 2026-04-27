package inventoryevent

import "time"

// Asset event models published on InventoryStream by inventory-service.
// These are consumed by subscriber-service (and any other service that
// needs to project asset state) for keeping local cache projections in
// sync with the authoritative inventory record.
//
// All payloads carry Version for the version-guard pattern: consumers
// fetch existing local state by AssetID and skip events where
// existing.LastEventVersion >= event.Version. See feedback memory
// "feedback_event_architecture.md" → "Version guards on Update events"
// for placement guidance.

// AssetCreatedEvent fires when a new asset is inserted into inventory
// (typically via inward posting of serial-tracked items).
type AssetCreatedEvent struct {
	AssetID       string            `json:"assetId" bson:"assetId"`
	OwnerTenantID string            `json:"ownerTenantId" bson:"ownerTenantId"`
	TemplateID    string            `json:"templateId" bson:"templateId"`
	TemplateName  string            `json:"templateName,omitempty" bson:"templateName,omitempty"`
	SerialNumber  string            `json:"serialNumber" bson:"serialNumber"`
	MACAddress    string            `json:"macAddress,omitempty" bson:"macAddress,omitempty"`
	Manufacturer  string            `json:"manufacturer,omitempty" bson:"manufacturer,omitempty"`
	Model         string            `json:"model,omitempty" bson:"model,omitempty"`
	SKU           string            `json:"sku,omitempty" bson:"sku,omitempty"`
	Status        string            `json:"status" bson:"status"`       // e.g. "IN_STOCK"
	Condition     string            `json:"condition" bson:"condition"` // e.g. "NEW"
	Identifiers   map[string]string `json:"identifiers,omitempty" bson:"identifiers,omitempty"`
	Version       int               `json:"version" bson:"version"`
	CreatedAt     time.Time         `json:"createdAt" bson:"createdAt"`
	CreatedBy     string            `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
}

// AssetTenantAssignedEvent fires when stock-tenant ownership transfers
// (e.g. parent ISP → franchisee). This is distinct from subscriber
// assignment — it changes WHO OWNS THE STOCK, not who uses the unit.
type AssetTenantAssignedEvent struct {
	AssetID       string    `json:"assetId" bson:"assetId"`
	FromTenantID  string    `json:"fromTenantId,omitempty" bson:"fromTenantId,omitempty"`
	ToTenantID    string    `json:"toTenantId" bson:"toTenantId"`
	Version       int       `json:"version" bson:"version"`
	TransferredAt time.Time `json:"transferredAt" bson:"transferredAt"`
	TransferredBy string    `json:"transferredBy,omitempty" bson:"transferredBy,omitempty"`
	Reason        string    `json:"reason,omitempty" bson:"reason,omitempty"`
}

// AssetConditionChangedEvent fires when QC re-grades an asset (e.g. a
// RETURNED unit gets graded REFURBISHED, or a working unit gets
// downgraded to DEGRADED). Drives the IsAssignable flag in
// subscriber-service's available_assets projection.
type AssetConditionChangedEvent struct {
	AssetID      string    `json:"assetId" bson:"assetId"`
	OwnerTenantID string   `json:"ownerTenantId" bson:"ownerTenantId"`
	OldCondition string    `json:"oldCondition,omitempty" bson:"oldCondition,omitempty"`
	NewCondition string    `json:"newCondition" bson:"newCondition"`
	Version      int       `json:"version" bson:"version"`
	ChangedAt    time.Time `json:"changedAt" bson:"changedAt"`
	ChangedBy    string    `json:"changedBy,omitempty" bson:"changedBy,omitempty"`
	Reason       string    `json:"reason,omitempty" bson:"reason,omitempty"`
}

// AssetAssignedEvent fires when an asset is bound to a subscriber after
// inventory-service successfully processes a subscriber.broadband.cpe_requested
// event (atomic Status: IN_STOCK → ASSIGNED).
type AssetAssignedEvent struct {
	AssetID                 string            `json:"assetId" bson:"assetId"`
	OwnerTenantID           string            `json:"ownerTenantId" bson:"ownerTenantId"`
	SubscriberID            string            `json:"subscriberId" bson:"subscriberId"`
	BroadbandSubscriptionID string            `json:"broadbandSubscriptionId,omitempty" bson:"broadbandSubscriptionId,omitempty"`
	SubscriptionType        string            `json:"subscriptionType,omitempty" bson:"subscriptionType,omitempty"` // "BROADBAND" | "SMB"
	SerialNumber            string            `json:"serialNumber" bson:"serialNumber"`
	MACAddress              string            `json:"macAddress,omitempty" bson:"macAddress,omitempty"`
	Manufacturer            string            `json:"manufacturer,omitempty" bson:"manufacturer,omitempty"`
	Model                   string            `json:"model,omitempty" bson:"model,omitempty"`
	Condition               string            `json:"condition" bson:"condition"`
	Identifiers             map[string]string `json:"identifiers,omitempty" bson:"identifiers,omitempty"`
	Version                 int               `json:"version" bson:"version"`
	AssignedAt              time.Time         `json:"assignedAt" bson:"assignedAt"`
	AssignedBy              string            `json:"assignedBy,omitempty" bson:"assignedBy,omitempty"`
}

// AssetAssignmentFailedEvent fires when subscriber-service requested an
// asset assignment via subscriber.broadband.cpe_requested but the asset
// could not be transitioned to ASSIGNED (e.g. already taken by another
// race-winner, condition no longer assignable, tenant scope mismatch).
//
// Subscriber-service consumes this and marks the relevant broadband/SMB
// subscription's cpeStatus = ASSIGNMENT_FAILED so the operator can pick
// another CPE.
type AssetAssignmentFailedEvent struct {
	AssetID                 string    `json:"assetId" bson:"assetId"`
	SubscriberID            string    `json:"subscriberId" bson:"subscriberId"`
	BroadbandSubscriptionID string    `json:"broadbandSubscriptionId,omitempty" bson:"broadbandSubscriptionId,omitempty"`
	SubscriptionType        string    `json:"subscriptionType,omitempty" bson:"subscriptionType,omitempty"`
	Reason                  string    `json:"reason" bson:"reason"`         // ALREADY_ASSIGNED | NOT_ASSIGNABLE | WRONG_TENANT | NOT_FOUND | RACE_LOST
	CurrentStatus           string    `json:"currentStatus,omitempty" bson:"currentStatus,omitempty"`
	CurrentCondition        string    `json:"currentCondition,omitempty" bson:"currentCondition,omitempty"`
	AttemptedAt             time.Time `json:"attemptedAt" bson:"attemptedAt"`
}

// AssetInstalledEvent fires when a field tech marks an asset as
// physically installed at the customer premises (Status ASSIGNED → INSTALLED).
type AssetInstalledEvent struct {
	AssetID                 string    `json:"assetId" bson:"assetId"`
	OwnerTenantID           string    `json:"ownerTenantId" bson:"ownerTenantId"`
	SubscriberID            string    `json:"subscriberId,omitempty" bson:"subscriberId,omitempty"`
	BroadbandSubscriptionID string    `json:"broadbandSubscriptionId,omitempty" bson:"broadbandSubscriptionId,omitempty"`
	Version                 int       `json:"version" bson:"version"`
	InstalledAt             time.Time `json:"installedAt" bson:"installedAt"`
	InstalledBy             string    `json:"installedBy,omitempty" bson:"installedBy,omitempty"`
}

// AssetReturnedEvent fires when an asset is returned (subscription
// cancelled, swap-out, etc.). Asset transitions toward IN_STOCK after
// QC re-grades the Condition.
type AssetReturnedEvent struct {
	AssetID                 string    `json:"assetId" bson:"assetId"`
	OwnerTenantID           string    `json:"ownerTenantId" bson:"ownerTenantId"`
	PrevSubscriberID        string    `json:"prevSubscriberId,omitempty" bson:"prevSubscriberId,omitempty"`
	BroadbandSubscriptionID string    `json:"broadbandSubscriptionId,omitempty" bson:"broadbandSubscriptionId,omitempty"`
	// Reason carries an AssetReleaseReason value (string-coded). Kept
	// as `string` for now; promoting to the typed alias is a separate
	// nicety that can land on the next common release without changing
	// the wire format.
	Reason     string    `json:"reason,omitempty" bson:"reason,omitempty"`
	Version    int       `json:"version" bson:"version"`
	ReturnedAt time.Time `json:"returnedAt" bson:"returnedAt"`
	ReturnedBy string    `json:"returnedBy,omitempty" bson:"returnedBy,omitempty"`
}

// AssetFaultyEvent fires when an asset is marked non-functional. The
// subscription it was bound to (if any) should transition to
// AWAITING_REPLACEMENT until a new CPE is assigned.
type AssetFaultyEvent struct {
	AssetID                 string `json:"assetId" bson:"assetId"`
	OwnerTenantID           string `json:"ownerTenantId" bson:"ownerTenantId"`
	PrevSubscriberID        string `json:"prevSubscriberId,omitempty" bson:"prevSubscriberId,omitempty"`
	BroadbandSubscriptionID string `json:"broadbandSubscriptionId,omitempty" bson:"broadbandSubscriptionId,omitempty"`
	// FaultReason carries an AssetReleaseReason value (string-coded).
	FaultReason string    `json:"faultReason,omitempty" bson:"faultReason,omitempty"`
	Version     int       `json:"version" bson:"version"`
	MarkedAt    time.Time `json:"markedAt" bson:"markedAt"`
	MarkedBy    string    `json:"markedBy,omitempty" bson:"markedBy,omitempty"`
}

// AssetRMAEvent fires when a faulty asset enters the vendor RMA pipeline.
type AssetRMAEvent struct {
	AssetID       string    `json:"assetId" bson:"assetId"`
	OwnerTenantID string    `json:"ownerTenantId" bson:"ownerTenantId"`
	VendorID      string    `json:"vendorId,omitempty" bson:"vendorId,omitempty"`
	RMANumber     string    `json:"rmaNumber,omitempty" bson:"rmaNumber,omitempty"`
	Version       int       `json:"version" bson:"version"`
	SentAt        time.Time `json:"sentAt" bson:"sentAt"`
	SentBy        string    `json:"sentBy,omitempty" bson:"sentBy,omitempty"`
}

// AssetScrappedEvent fires when an asset reaches terminal state.
// Subscriber-side projection should drop the row from available_assets.
type AssetScrappedEvent struct {
	AssetID       string `json:"assetId" bson:"assetId"`
	OwnerTenantID string `json:"ownerTenantId" bson:"ownerTenantId"`
	// Reason carries an AssetReleaseReason value (string-coded).
	Reason     string    `json:"reason,omitempty" bson:"reason,omitempty"`
	Version    int       `json:"version" bson:"version"`
	ScrappedAt time.Time `json:"scrappedAt" bson:"scrappedAt"`
	ScrappedBy string    `json:"scrappedBy,omitempty" bson:"scrappedBy,omitempty"`
}
