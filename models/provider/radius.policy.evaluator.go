package provider

import (
	"fmt"
	"net"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------
// RADIUS Policy Evaluator
//
// Evaluates a set of RadiusNetworkPolicies against an incoming RADIUS
// Access-Request attribute map and returns a merged, priority-ordered list
// of reply attributes to inject into the Access-Accept.
//
// Also provides helpers for resolving RadiusReplyRules (plan-level
// multi-vendor reply templates) filtered by the NAS vendor.
// ---------------------------------------------------------------------------

// PolicyConditionOp mirrors the DB constant so the evaluator has no DB import.
type PolicyConditionOp string

const (
	PolicyCondOpEq      PolicyConditionOp = "=="
	PolicyCondOpNeq     PolicyConditionOp = "!="
	PolicyCondOpIn      PolicyConditionOp = "in"
	PolicyCondOpNotIn   PolicyConditionOp = "not_in"
	PolicyCondOpMatches PolicyConditionOp = "matches"
	PolicyCondOpRange   PolicyConditionOp = "range"
)

// PolicyReplyOp mirrors the DB constant.
type PolicyReplyOp string

const (
	PolicyReplyOpOverride PolicyReplyOp = ":="
	PolicyReplyOpAddOnce  PolicyReplyOp = "="
	PolicyReplyOpAppend   PolicyReplyOp = "+="
)

// ---------------------------------------------------------------------------
// Input types (plain structs — no DB/GORM imports needed here)
// ---------------------------------------------------------------------------

// PolicyCondition is one condition belonging to a network policy.
type PolicyCondition struct {
	Attr  string            // RADIUS attribute name e.g. "Tunnel-Private-Group-Id"
	Op    PolicyConditionOp // comparison operator
	Value string            // comparison target
}

// PolicyReplyAttr is one reply attribute belonging to a network policy.
type PolicyReplyAttr struct {
	Type   string        // "reply" | "check" | "config"
	Vendor string        // "" or registry key e.g. "CISCO"
	Attr   string        // attribute name
	Op     PolicyReplyOp // ":=" | "=" | "+="
	Value  string
}

// NetworkPolicy is the evaluatable form of a RadiusNetworkPolicy.
// The repository layer maps DB rows into this struct before calling Evaluate.
type NetworkPolicy struct {
	ID         string
	Name       string
	Priority   int // lower = higher priority
	Conditions []PolicyCondition
	ReplyAttrs []PolicyReplyAttr
}

// ReplyRule is the evaluatable form of a RadiusReplyRule row.
type ReplyRule struct {
	SortOrder int
	Type      string
	Vendor    string // "" = applies to all vendors
	Attr      string
	Op        string
	Value     string
}

// ResolvedReplyAttr is a single reply attribute after all policies are merged.
type ResolvedReplyAttr struct {
	Type             string
	Vendor           string
	Attr             string
	Op               PolicyReplyOp
	Value            string
	SourcePolicyID   string // which policy produced this (for audit)
	SourcePolicyName string
}

// ---------------------------------------------------------------------------
// Network Policy Evaluator
// ---------------------------------------------------------------------------

// EvaluateNetworkPolicies matches all provided policies against requestAttrs,
// merges the reply attributes from every matching policy (ordered by Priority),
// and returns the final set of attributes to inject into the Access-Accept.
//
// Merge semantics:
//   :=  (Override) — first match per Attr wins; later policies cannot override it.
//   =   (AddOnce)  — same as :=, first occurrence per Attr is kept.
//   +=  (Append)   — always added regardless; produces multi-value attrs.
//
// requestAttrs maps RADIUS attribute names (e.g. "Tunnel-Private-Group-Id")
// to their string values as received in the Access-Request packet.
func EvaluateNetworkPolicies(
	requestAttrs map[string]string,
	policies []NetworkPolicy,
) []ResolvedReplyAttr {
	if len(policies) == 0 || len(requestAttrs) == 0 {
		return nil
	}

	// Sort policies by priority ascending (lower number = evaluated first).
	sorted := make([]NetworkPolicy, len(policies))
	copy(sorted, policies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority < sorted[j].Priority
	})

	// seen tracks which (type+attr) keys have been set by := or =
	// so that later lower-priority policies cannot override them.
	type attrKey struct{ typ, attr string }
	seen := make(map[attrKey]bool)

	var result []ResolvedReplyAttr

	for _, policy := range sorted {
		if !matchesAllConditions(requestAttrs, policy.Conditions) {
			continue
		}
		for _, ra := range policy.ReplyAttrs {
			k := attrKey{ra.Type, ra.Attr}
			switch ra.Op {
			case PolicyReplyOpAppend: // += always appends, no duplicate check
				result = append(result, ResolvedReplyAttr{
					Type:             ra.Type,
					Vendor:           ra.Vendor,
					Attr:             ra.Attr,
					Op:               ra.Op,
					Value:            ra.Value,
					SourcePolicyID:   policy.ID,
					SourcePolicyName: policy.Name,
				})
			default: // := and = — first-write wins
				if !seen[k] {
					seen[k] = true
					result = append(result, ResolvedReplyAttr{
						Type:             ra.Type,
						Vendor:           ra.Vendor,
						Attr:             ra.Attr,
						Op:               ra.Op,
						Value:            ra.Value,
						SourcePolicyID:   policy.ID,
						SourcePolicyName: policy.Name,
					})
				}
			}
		}
	}

	return result
}

// matchesAllConditions returns true only when every condition in the list
// evaluates to true against requestAttrs (AND logic).
func matchesAllConditions(requestAttrs map[string]string, conditions []PolicyCondition) bool {
	for _, c := range conditions {
		requestVal, exists := requestAttrs[c.Attr]
		if !exists {
			// Attribute not present in request — condition cannot match.
			return false
		}
		if !evaluateCondition(requestVal, c.Op, c.Value) {
			return false
		}
	}
	return true
}

// evaluateCondition compares a single request value against a condition.
func evaluateCondition(requestVal string, op PolicyConditionOp, condVal string) bool {
	switch op {
	case PolicyCondOpEq:
		return strings.EqualFold(requestVal, condVal)

	case PolicyCondOpNeq:
		return !strings.EqualFold(requestVal, condVal)

	case PolicyCondOpIn:
		for _, item := range strings.Split(condVal, ",") {
			if strings.EqualFold(strings.TrimSpace(item), requestVal) {
				return true
			}
		}
		return false

	case PolicyCondOpNotIn:
		for _, item := range strings.Split(condVal, ",") {
			if strings.EqualFold(strings.TrimSpace(item), requestVal) {
				return false
			}
		}
		return true

	case PolicyCondOpMatches:
		rx, err := regexp.Compile(condVal)
		if err != nil {
			return false
		}
		return rx.MatchString(requestVal)

	case PolicyCondOpRange:
		// condVal format: "low-high" e.g. "500-600"
		return inNumericRange(requestVal, condVal)

	default:
		return false
	}
}

// inNumericRange checks whether requestVal (parsed as float64) falls within
// the "low-high" range specified in rangeVal.
func inNumericRange(requestVal, rangeVal string) bool {
	parts := strings.SplitN(rangeVal, "-", 2)
	if len(parts) != 2 {
		return false
	}
	lo, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	hi, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	rv, err3 := strconv.ParseFloat(strings.TrimSpace(requestVal), 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}
	return rv >= lo && rv <= hi
}

// ---------------------------------------------------------------------------
// Plan-level Reply Rule resolver
// ---------------------------------------------------------------------------

// ResolveReplyRules filters a plan's RadiusReplyRules to those applicable
// for the given NAS vendor and returns them sorted by SortOrder.
//
// Rules with an empty Vendor apply to ALL NAS vendors.
// Rules with a non-empty Vendor apply only when nasVendor matches (case-insensitive).
//
// nasVendor should be the registry key e.g. "MIKROTIK", "CISCO", "ERX".
// Pass "" to retrieve only vendor-agnostic (standard) rules.
func ResolveReplyRules(rules []ReplyRule, nasVendor string) []ReplyRule {
	nasVendor = strings.ToUpper(strings.TrimSpace(nasVendor))

	var result []ReplyRule
	for _, r := range rules {
		v := strings.ToUpper(strings.TrimSpace(r.Vendor))
		if v == "" || v == nasVendor {
			result = append(result, r)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SortOrder < result[j].SortOrder
	})
	return result
}

// ---------------------------------------------------------------------------
// Policy validation helpers (used before saving to DB)
// ---------------------------------------------------------------------------

// ValidateNetworkPolicy validates a NetworkPolicy before it is persisted.
// Returns a slice of all validation errors found (never stops at first).
func ValidateNetworkPolicy(p NetworkPolicy) []error {
	var errs []error

	if strings.TrimSpace(p.Name) == "" {
		errs = append(errs, fmt.Errorf("policy name is required"))
	}

	for i, c := range p.Conditions {
		if strings.TrimSpace(c.Attr) == "" {
			errs = append(errs, fmt.Errorf("condition[%d]: attr is required", i))
		}
		if strings.TrimSpace(c.Value) == "" {
			errs = append(errs, fmt.Errorf("condition[%d]: value is required", i))
		}
		switch c.Op {
		case PolicyCondOpEq, PolicyCondOpNeq, PolicyCondOpIn, PolicyCondOpNotIn,
			PolicyCondOpMatches, PolicyCondOpRange:
		default:
			errs = append(errs, fmt.Errorf("condition[%d]: unknown operator %q", i, c.Op))
		}
		if c.Op == PolicyCondOpMatches {
			if _, err := regexp.Compile(c.Value); err != nil {
				errs = append(errs, fmt.Errorf("condition[%d]: invalid regex %q: %v", i, c.Value, err))
			}
		}
		if c.Op == PolicyCondOpRange {
			parts := strings.SplitN(c.Value, "-", 2)
			if len(parts) != 2 {
				errs = append(errs, fmt.Errorf("condition[%d]: range value must be 'low-high', got %q", i, c.Value))
			} else {
				if _, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64); err != nil {
					errs = append(errs, fmt.Errorf("condition[%d]: range low bound is not a number", i))
				}
				if _, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err != nil {
					errs = append(errs, fmt.Errorf("condition[%d]: range high bound is not a number", i))
				}
			}
		}
	}

	for i, ra := range p.ReplyAttrs {
		if strings.TrimSpace(ra.Attr) == "" {
			errs = append(errs, fmt.Errorf("replyAttr[%d]: attr is required", i))
			continue
		}
		switch ra.Type {
		case "reply", "check", "config", "":
		default:
			errs = append(errs, fmt.Errorf("replyAttr[%d]: unknown type %q (use reply/check/config)", i, ra.Type))
		}
		switch ra.Op {
		case PolicyReplyOpOverride, PolicyReplyOpAddOnce, PolicyReplyOpAppend:
		default:
			errs = append(errs, fmt.Errorf("replyAttr[%d]: unknown op %q (use :=, =, +=)", i, ra.Op))
		}

		// Validate attribute exists in registry (vendor-scoped when vendor is set).
		vendor := strings.TrimSpace(ra.Vendor)
		if vendor != "" {
			if !IsValidKeyForVendor(vendor, ra.Attr) {
				errs = append(errs, fmt.Errorf("replyAttr[%d]: attribute %q is not defined for vendor %q",
					i, ra.Attr, vendor))
			}
		} else {
			if !IsValidKey(ra.Attr) {
				errs = append(errs, fmt.Errorf("replyAttr[%d]: attribute %q is not defined in any vendor dictionary",
					i, ra.Attr))
			}
		}
	}

	return errs
}

// ValidateReplyRule validates a single ReplyRule row before it is persisted.
func ValidateReplyRule(r ReplyRule) []error {
	var errs []error

	if strings.TrimSpace(r.Attr) == "" {
		errs = append(errs, fmt.Errorf("attr is required"))
	}
	switch r.Op {
	case string(PolicyReplyOpOverride), string(PolicyReplyOpAddOnce), string(PolicyReplyOpAppend):
	default:
		errs = append(errs, fmt.Errorf("unknown op %q (use :=, =, +=)", r.Op))
	}
	switch r.Type {
	case "reply", "check", "config", "":
	default:
		errs = append(errs, fmt.Errorf("unknown type %q (use reply/check/config)", r.Type))
	}

	vendor := strings.ToUpper(strings.TrimSpace(r.Vendor))
	if r.Attr != "" {
		if vendor != "" {
			if !IsValidKeyForVendor(vendor, r.Attr) {
				errs = append(errs, fmt.Errorf("attribute %q is not defined for vendor %q", r.Attr, vendor))
			}
		} else {
			if !IsValidKey(r.Attr) {
				errs = append(errs, fmt.Errorf("attribute %q is not defined in any vendor dictionary", r.Attr))
			}
		}
	}

	return errs
}

// ---------------------------------------------------------------------------
// Convenience: extract common condition fields from a RADIUS request map.
// The caller populates requestAttrs from the actual packet attributes.
// These helpers translate semantic names into standard RADIUS attribute names.
// ---------------------------------------------------------------------------

// RequestAttrVLAN returns the canonical RADIUS attribute name used to carry
// the VLAN ID in an Access-Request. Most BNG/NAS vendors use
// Tunnel-Private-Group-Id; some use NAS-Port-Id with a "vlan:" prefix.
const RequestAttrVLAN = "Tunnel-Private-Group-Id"

// RequestAttrNASIP is the standard attribute carrying the NAS IP address.
const RequestAttrNASIP = "NAS-IP-Address"

// RequestAttrCalledStation is the Called-Station-Id attribute (SSID, port, DSLAM info).
const RequestAttrCalledStation = "Called-Station-Id"

// RequestAttrCallingStation is the Calling-Station-Id attribute (client MAC address).
const RequestAttrCallingStation = "Calling-Station-Id"

// NASIPInCIDR reports whether the NAS-IP-Address in requestAttrs falls within
// the given CIDR block. Useful for building NAS-subnet based conditions.
func NASIPInCIDR(requestAttrs map[string]string, cidr string) (bool, error) {
	nasIP, ok := requestAttrs[RequestAttrNASIP]
	if !ok {
		return false, nil
	}
	ip := net.ParseIP(strings.TrimSpace(nasIP))
	if ip == nil {
		return false, fmt.Errorf("invalid NAS-IP-Address %q", nasIP)
	}
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, fmt.Errorf("invalid CIDR %q: %v", cidr, err)
	}
	return network.Contains(ip), nil
}
