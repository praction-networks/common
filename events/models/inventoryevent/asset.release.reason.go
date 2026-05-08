package inventoryevent

import "strings"

// AssetReleaseReason is the typed operator/system reason carried on
// asset transition events that take an asset OUT OF a bound or in-stock
// state — RETURNED, FAULTY, SCRAPPED. Keeping this in common ensures
// inventory-service (the producer), subscriber-service (which projects
// the reason onto cpe_binding_history), and any audit/analytics consumer
// all agree on the same string set.
//
// Not every value is meaningful for every transition (e.g. END_OF_LIFE
// is really a scrap path), but the type is unified so consumers can
// treat the reason as a single comparable field. Use the
// NormalizeAssetReleaseReason helper at every ingress to defend against
// case-mangled or whitespace-padded inputs.
type AssetReleaseReason string

const (
	// AssetReleaseReasonFaulty — unit failed in service.
	AssetReleaseReasonFaulty AssetReleaseReason = "FAULTY"
	// AssetReleaseReasonSwap — operator swapped the unit (e.g. customer
	// upgraded plan, capacity bump, preventative replacement).
	AssetReleaseReasonSwap AssetReleaseReason = "SWAP"
	// AssetReleaseReasonSubscriptionCancelled — sub was cancelled and
	// the unit is being returned to stock.
	AssetReleaseReasonSubscriptionCancelled AssetReleaseReason = "SUBSCRIPTION_CANCELLED"
	// AssetReleaseReasonDOA — dead-on-arrival; never functioned in service.
	AssetReleaseReasonDOA AssetReleaseReason = "DOA"
	// AssetReleaseReasonEndOfLife — typically a scrap-path reason; the
	// unit reached end of useful life and is being terminally retired.
	AssetReleaseReasonEndOfLife AssetReleaseReason = "END_OF_LIFE"
	// AssetReleaseReasonDamaged — physical damage (customer-caused,
	// transport, etc.). Distinct from FAULTY which implies functional
	// failure of an otherwise-intact unit.
	AssetReleaseReasonDamaged AssetReleaseReason = "DAMAGED"
	// AssetReleaseReasonLost — unit cannot be recovered from the field.
	AssetReleaseReasonLost AssetReleaseReason = "LOST"
	// AssetReleaseReasonOther — catch-all; pair with a free-text note.
	AssetReleaseReasonOther AssetReleaseReason = "OTHER"
	// AssetReleaseReasonChurnRecovery — operator recovered the unit from a
	// churning customer (subscription not yet cancelled at the moment of
	// recovery; distinct from SUBSCRIPTION_CANCELLED which presumes the
	// cancellation already landed).
	// Source: field-central assets-prd.md §6.4.
	AssetReleaseReasonChurnRecovery AssetReleaseReason = "CHURN_RECOVERY"
	// AssetReleaseReasonUnknown — sentinel for unrecognized / unparseable
	// inputs. Storing UNKNOWN is preferable to dropping the row, as it
	// surfaces the bad-reason via analytics queries.
	AssetReleaseReasonUnknown AssetReleaseReason = "UNKNOWN"
)

// validAssetReleaseReasons is a lookup of recognized values for fast
// membership checks during normalization.
var validAssetReleaseReasons = map[AssetReleaseReason]bool{
	AssetReleaseReasonFaulty:                true,
	AssetReleaseReasonSwap:                  true,
	AssetReleaseReasonSubscriptionCancelled: true,
	AssetReleaseReasonDOA:                   true,
	AssetReleaseReasonEndOfLife:             true,
	AssetReleaseReasonDamaged:               true,
	AssetReleaseReasonLost:                  true,
	AssetReleaseReasonOther:                 true,
	AssetReleaseReasonChurnRecovery:         true,
}

// NormalizeAssetReleaseReason trims whitespace, uppercases the input,
// and returns the matching enum value. Empty or unrecognized inputs map
// to AssetReleaseReasonUnknown — callers should warn-log when input was
// non-empty but normalized to UNKNOWN, so bad reasons are visible to ops
// rather than silently corrupted.
func NormalizeAssetReleaseReason(raw string) AssetReleaseReason {
	cleaned := AssetReleaseReason(strings.ToUpper(strings.TrimSpace(raw)))
	if cleaned == "" {
		return AssetReleaseReasonUnknown
	}
	if validAssetReleaseReasons[cleaned] {
		return cleaned
	}
	return AssetReleaseReasonUnknown
}

// AssetReleaseReasonOneOf is the space-separated list of canonical
// values usable directly in struct-tag `oneof=...` validation. Kept in
// sync with validAssetReleaseReasons.
const AssetReleaseReasonOneOf = "FAULTY SWAP SUBSCRIPTION_CANCELLED DOA END_OF_LIFE DAMAGED LOST OTHER CHURN_RECOVERY"
