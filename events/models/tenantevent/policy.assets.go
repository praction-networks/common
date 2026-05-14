package tenantevent

import "fmt"

// PolicyAssets is the per-tenant asset-policy bucket.
// Source: backend-contract §11.3, Sweep 4 Q9 (dropoffLockWindowHours), Q10 (recoveryAcknowledgement).
type PolicyAssets struct {
	DropoffLockWindowHours  int    `json:"dropoffLockWindowHours,omitempty"  bson:"dropoffLockWindowHours,omitempty"`
	RecoveryAcknowledgement string `json:"recoveryAcknowledgement,omitempty" bson:"recoveryAcknowledgement,omitempty"` // NONE | SIGNATURE | OTP

	// Reconciliation & variance
	ReconcileNudgeHourLocal int `json:"reconcileNudgeHourLocal,omitempty" bson:"reconcileNudgeHourLocal,omitempty"` // hour-of-day (0–23) for server nudge job
	VarianceEscalationCount int `json:"varianceEscalationCount,omitempty" bson:"varianceEscalationCount,omitempty"` // auto ops-ticket threshold

	// Geofencing
	WarehouseGeofences      []WarehouseGeofence `json:"warehouseGeofences,omitempty"      bson:"warehouseGeofences,omitempty"`      // per-tenant geofence overrides
	ProximityThresholdM     int                 `json:"proximityThresholdM,omitempty"     bson:"proximityThresholdM,omitempty"`     // geofence radius in meters
	NearbyWarehouseThresholdKm float64          `json:"nearbyWarehouseThresholdKm,omitempty" bson:"nearbyWarehouseThresholdKm,omitempty"` // warehouse finder radius

	// Consumables
	ConsumableLowThresholds []ConsumableLowThreshold `json:"consumableLowThresholds,omitempty" bson:"consumableLowThresholds,omitempty"` // low-stock alert configs

	// Peer-handoff & courier
	PeerReceiptTimeoutMinutes int    `json:"peerReceiptTimeoutMinutes,omitempty" bson:"peerReceiptTimeoutMinutes,omitempty"` // peer-handoff timeout
	CourierTimeoutHours       int    `json:"courierTimeoutHours,omitempty"       bson:"courierTimeoutHours,omitempty"`       // peer-assist courier timeout
	PeerScope                 string `json:"peerScope,omitempty"                 bson:"peerScope,omitempty"`                 // REGION | ZONE | ALL

	// Job transfer
	JobTransferApproverSlaMinutes    int      `json:"jobTransferApproverSlaMinutes,omitempty"    bson:"jobTransferApproverSlaMinutes,omitempty"`    // approver deadline
	JobTransferPeerOfferMinutes      int      `json:"jobTransferPeerOfferMinutes,omitempty"      bson:"jobTransferPeerOfferMinutes,omitempty"`      // peer offer window
	JobTransferTotalLifetimeHours    int      `json:"jobTransferTotalLifetimeHours,omitempty"    bson:"jobTransferTotalLifetimeHours,omitempty"`    // full transfer TTL
	JobTransferReasonsAllowed        []string `json:"jobTransferReasonsAllowed,omitempty"        bson:"jobTransferReasonsAllowed,omitempty"`        // nil = all reasons allowed
	JobTransferRequiresFeManagerOnly bool     `json:"jobTransferRequiresFeManagerOnly,omitempty" bson:"jobTransferRequiresFeManagerOnly,omitempty"` // approver gating
}

// Validate checks that enumerated string fields hold only documented values
// and that numeric fields are within their documented ranges.
func (a PolicyAssets) Validate() error {
	switch a.RecoveryAcknowledgement {
	case "", "NONE", "SIGNATURE", "OTP":
		// ok
	default:
		return fmt.Errorf("RecoveryAcknowledgement must be NONE, SIGNATURE, or OTP (got %q)", a.RecoveryAcknowledgement)
	}
	switch a.PeerScope {
	case "", "REGION", "ZONE", "ALL":
		// ok
	default:
		return fmt.Errorf("PeerScope must be REGION, ZONE, or ALL (got %q)", a.PeerScope)
	}
	if a.ReconcileNudgeHourLocal < 0 || a.ReconcileNudgeHourLocal > 23 {
		return fmt.Errorf("ReconcileNudgeHourLocal must be 0–23 (got %d)", a.ReconcileNudgeHourLocal)
	}
	return nil
}

// WarehouseGeofence defines a per-warehouse geofence override.
// When set, RadiusM overrides the tenant-level ProximityThresholdM for this warehouse.
// RadiusM == 0 means use ProximityThresholdM.
type WarehouseGeofence struct {
	WarehouseID string  `json:"warehouseId" bson:"warehouseId"`
	Lat         float64 `json:"lat"         bson:"lat"`
	Lng         float64 `json:"lng"         bson:"lng"`
	RadiusM     int     `json:"radiusM"     bson:"radiusM"`
}

// ConsumableLowThreshold defines the low-stock alert threshold for one consumable type.
type ConsumableLowThreshold struct {
	ConsumableID string `json:"consumableId" bson:"consumableId"`
	Threshold    int    `json:"threshold"    bson:"threshold"`
}
