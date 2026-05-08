package inventoryevent

// Asset state model — shared types and helpers for any service that
// projects inventory asset state (inventory-service is the authority,
// subscriber-service maintains a read projection, etc.).
//
// Two orthogonal dimensions describe an asset:
//   - Status:    WHERE the asset is in its lifecycle (IN_STOCK / ASSIGNED / ...).
//   - Condition: PHYSICAL grade (NEW / REFURBISHED / USED_OK / DEGRADED / FAULTY).
//
// Both must qualify for a subscriber-side assignment (see
// IsAssignableForBroadband). Keeping these in common ensures every
// service compares against identical strings without copy-paste drift.

// ──────────────────────────────────────────────
// Asset Lifecycle States (Status)
// ──────────────────────────────────────────────

// AssetStatus tracks WHERE an asset is in its lifecycle.
type AssetStatus string

const (
	AssetStatusInStock   AssetStatus = "IN_STOCK"
	AssetStatusReserved  AssetStatus = "RESERVED"
	AssetStatusAssigned  AssetStatus = "ASSIGNED"
	AssetStatusInstalled AssetStatus = "INSTALLED"
	AssetStatusFaulty    AssetStatus = "FAULTY"
	AssetStatusRMA       AssetStatus = "RMA"
	AssetStatusReturned  AssetStatus = "RETURNED"
	AssetStatusScrapped  AssetStatus = "SCRAPPED"
	// AssetStatusPendingReturn — FE has submitted a drop-off batch that
	// includes this asset but the WM has not yet accepted/rejected the line.
	// Custody still nominally with the FE; auto-released back to ASSIGNED/
	// INSTALLED after dropoffLockWindowHours (tenant policy) elapses.
	// Source: field-central assets-prd.md §11.2 (two-step drop-off).
	AssetStatusPendingReturn AssetStatus = "PENDING_RETURN"
)

// ValidAssetTransitions defines the allowed state transitions for assets.
var ValidAssetTransitions = map[AssetStatus][]AssetStatus{
	AssetStatusInStock:       {AssetStatusReserved, AssetStatusAssigned, AssetStatusScrapped},
	AssetStatusReserved:      {AssetStatusInStock, AssetStatusAssigned},
	AssetStatusAssigned:      {AssetStatusInstalled, AssetStatusReturned, AssetStatusFaulty, AssetStatusPendingReturn},
	AssetStatusInstalled:     {AssetStatusFaulty, AssetStatusReturned, AssetStatusPendingReturn},
	AssetStatusFaulty:        {AssetStatusRMA, AssetStatusScrapped, AssetStatusReturned},
	AssetStatusRMA:           {AssetStatusReturned, AssetStatusScrapped},
	AssetStatusReturned:      {AssetStatusInStock},
	AssetStatusScrapped:      {}, // terminal state
	AssetStatusPendingReturn: {AssetStatusReturned, AssetStatusAssigned, AssetStatusInstalled},
}

// IsValidTransition checks if a state transition is allowed.
func IsValidTransition(from, to AssetStatus) bool {
	allowed, ok := ValidAssetTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// ──────────────────────────────────────────────
// Asset Condition (physical grade)
// ──────────────────────────────────────────────

// AssetCondition describes the PHYSICAL grade of an asset, orthogonal to
// AssetStatus (which tracks lifecycle position). Both must be considered
// when deciding whether an asset can be assigned to a subscriber.
type AssetCondition string

const (
	// AssetConditionNew — fresh from manufacturer / inward, never deployed.
	AssetConditionNew AssetCondition = "NEW"
	// AssetConditionRefurbished — returned via RMA pipeline, cleaned + tested OK.
	AssetConditionRefurbished AssetCondition = "REFURBISHED"
	// AssetConditionUsedOK — returned from a customer, tested, fully functional.
	AssetConditionUsedOK AssetCondition = "USED_OK"
	// AssetConditionDegraded — works but reduced spec / reliability. NEVER assignable
	// via the standard subscriber-flow; operator override required.
	AssetConditionDegraded AssetCondition = "DEGRADED"
	// AssetConditionFaulty — non-functional; awaiting RMA or scrap decision.
	AssetConditionFaulty AssetCondition = "FAULTY"
)

// AssignableConditions are the conditions under which an IN_STOCK asset
// can be assigned to a subscriber via the standard subscriber-service flow.
// DEGRADED is intentionally excluded — operators can override per-case but
// the default-deny rule keeps service quality consistent.
var AssignableConditions = map[AssetCondition]bool{
	AssetConditionNew:         true,
	AssetConditionRefurbished: true,
	AssetConditionUsedOK:      true,
}

// IsAssignableForBroadband reports whether an asset can be assigned to a
// broadband (or SMB) subscriber. Both Status and Condition must qualify.
//
// Status check: IN_STOCK only — RESERVED/ASSIGNED/INSTALLED/FAULTY/etc.
// already represent a binding or non-functional state.
//
// Condition check: NEW / REFURBISHED / USED_OK only — DEGRADED and FAULTY
// are excluded by AssignableConditions.
func IsAssignableForBroadband(status AssetStatus, condition AssetCondition) bool {
	if status != AssetStatusInStock {
		return false
	}
	return AssignableConditions[condition]
}

// SubscriptionType discriminates which subscription kind an asset is
// bound to in subscriber-service. Used by the cpe-request saga events.
type SubscriptionType string

const (
	SubscriptionTypeBroadband SubscriptionType = "BROADBAND"
	SubscriptionTypeSMB       SubscriptionType = "SMB"
)

// AssignmentFailureReason enumerates why inventory-service rejects a
// subscriber.broadband.cpe_requested attempt. Carried in
// AssetAssignmentFailedEvent.Reason. Subscriber-side reconciliation maps
// these to user-facing messages.
const (
	AssignmentFailureNotFound        = "NOT_FOUND"
	AssignmentFailureNotAssignable   = "NOT_ASSIGNABLE"
	AssignmentFailureWrongTenant     = "WRONG_TENANT"
	AssignmentFailureAlreadyAssigned = "ALREADY_ASSIGNED"
	AssignmentFailureRaceLost        = "RACE_LOST"
)
