package subscriberevent

import (
	"time"

	"github.com/praction-networks/common/events/models/inventoryevent"
)

// CPE assignment-saga event payloads. Published by subscriber-service
// whenever an operator requests or releases a CPE binding on a broadband
// or SMB subscription. Inventory-service consumes these to perform the
// authoritative atomic asset state transition (IN_STOCK → ASSIGNED on
// request; ASSIGNED → RETURNED on release).
//
// SubscriptionType (BROADBAND | SMB) is the discriminator that lets one
// event subject serve both subscription kinds — the inventory-side
// handler doesn't need to care which type initiated the request.
//
// Source (FIELD_APP | ADMIN_DASHBOARD) is informational; carried for
// audit/analytics so ops can tell which UI initiated the assignment.
const (
	CpeRequestSourceFieldApp       = "FIELD_APP"
	CpeRequestSourceAdminDashboard = "ADMIN_DASHBOARD"
)

// BroadbandCpeRequestedEvent is published when a CSR/operator binds a
// CPE asset to a broadband or SMB subscription. The local subscription
// row is persisted with cpeStatus=PENDING_CPE_CONFIRMATION before the
// event fires; inventory-service then attempts the atomic transition
// and emits either inventory.asset.assigned (success) or
// inventory.asset.assignment_failed (race lost / not assignable).
type BroadbandCpeRequestedEvent struct {
	BroadbandSubscriptionID string                            `json:"broadbandSubscriptionId" bson:"broadbandSubscriptionId"`
	SubscriberID            string                            `json:"subscriberId" bson:"subscriberId"`
	TenantID                string                            `json:"tenantId" bson:"tenantId"`
	SubscriptionType        inventoryevent.SubscriptionType   `json:"subscriptionType" bson:"subscriptionType"` // BROADBAND | SMB
	AssetID                 string                            `json:"assetId" bson:"assetId"`
	SerialNumber            string                            `json:"serialNumber" bson:"serialNumber"`
	Source                  string                            `json:"source,omitempty" bson:"source,omitempty"` // FIELD_APP | ADMIN_DASHBOARD
	Version                 int                               `json:"version" bson:"version"`
	RequestedAt             time.Time                         `json:"requestedAt" bson:"requestedAt"`
	RequestedBy             string                            `json:"requestedBy,omitempty" bson:"requestedBy,omitempty"`
}

// BroadbandCpeReleasedEvent is published when a CPE is unbound from a
// subscription — either because the subscription was cancelled, the CPE
// is being swapped out, or the asset was marked faulty. Inventory-service
// consumes this and transitions the asset toward IN_STOCK or FAULTY
// based on the Reason.
type BroadbandCpeReleasedEvent struct {
	BroadbandSubscriptionID string                            `json:"broadbandSubscriptionId" bson:"broadbandSubscriptionId"`
	SubscriberID            string                            `json:"subscriberId" bson:"subscriberId"`
	TenantID                string                            `json:"tenantId" bson:"tenantId"`
	SubscriptionType        inventoryevent.SubscriptionType   `json:"subscriptionType" bson:"subscriptionType"`
	PrevAssetID             string                            `json:"prevAssetId" bson:"prevAssetId"`
	Reason                  string                            `json:"reason" bson:"reason"` // FAULTY | SWAP | SUBSCRIPTION_CANCELLED
	Version                 int                               `json:"version" bson:"version"`
	ReleasedAt              time.Time                         `json:"releasedAt" bson:"releasedAt"`
	ReleasedBy              string                            `json:"releasedBy,omitempty" bson:"releasedBy,omitempty"`
}
