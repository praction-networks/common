// Package compliance enforces the second layer of the compliance-retain
// rule documented in reference_compliance_retain_policy.md (memory).
// The first layer is RBAC tag-based exclusion in the auth-service seed —
// every non-SuperAdmin role's PolicySelectors carries
// ExcludeTags:["compliance-retain"], blocking 14 of 15 roles.
//
// This package is the second layer: every handler that owns a
// compliance-retain DELETE endpoint MUST call CheckRetention before
// executing the delete. Even SuperAdmin gets a 409 if the record is
// younger than the regulatory retention period.
//
// Indian regulatory basis:
//   - TRAI consumer protection + DOT UASL clause 39 → 2 yr subscriber/billing
//   - PMLA → 10 yr KYC docs
//   - CERT-In Directive (Apr 2022) → 180 d log retention (audit, IPDR, DNS)
//   - PMLA + IT Act → KYC + transaction logs
//
// All defaults are conservative (longest of overlapping minimums).
// Per-tenant overrides can be wired via TenantRetentionConfig (TBD).
package compliance

import (
	"fmt"
	"net/http"
	"time"

	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"
)

// RetentionResource identifies a class of compliance-retain protected
// record. Each constant maps to a default MinRetention period.
type RetentionResource string

const (
	// ResourceSubscriber — subscriber master record (PII + billing).
	// TRAI consumer protection + DOT UASL clause 39.
	ResourceSubscriber RetentionResource = "subscriber"

	// ResourceKYCDocument — KYC document (Aadhaar/PAN/proof).
	// PMLA Act, longest of the bunch.
	ResourceKYCDocument RetentionResource = "kyc-document"

	// ResourceAuditLog — internal admin audit trail.
	// CERT-In Directive (Apr 2022) §X for ICT providers.
	ResourceAuditLog RetentionResource = "audit-log"

	// ResourceLogReport — generated extract of CGNAT/IPDR data, often
	// for LEA / TRAI / DOT request responses.
	ResourceLogReport RetentionResource = "log-report"

	// ResourceDNSQueryLog — subscriber DNS browsing pattern (PII).
	// CERT-In Directive.
	ResourceDNSQueryLog RetentionResource = "dns-query-log"
)

// defaultRetention is the regulatory floor per resource. Real
// deployments may extend via per-tenant config — never shrink below
// these defaults without legal sign-off.
var defaultRetention = map[RetentionResource]time.Duration{
	ResourceSubscriber:  2 * 365 * 24 * time.Hour,  // 2 yr
	ResourceKYCDocument: 10 * 365 * 24 * time.Hour, // 10 yr PMLA
	ResourceAuditLog:    2 * 365 * 24 * time.Hour,  // 2 yr practical (180 d minimum)
	ResourceLogReport:   2 * 365 * 24 * time.Hour,  // 2 yr practical
	ResourceDNSQueryLog: 180 * 24 * time.Hour,      // 180 d CERT-In floor
}

// RetentionViolation is the structured error returned when a delete
// is blocked by retention rules. The 409 response body carries every
// field so audit + UX can surface specifics.
type RetentionViolation struct {
	Resource          RetentionResource `json:"resource"`
	RecordID          string            `json:"recordId"`
	RecordLastUpdated time.Time         `json:"recordLastUpdated"`
	MinRetention      time.Duration     `json:"minRetentionDuration"`
	EarliestDeleteAt  time.Time         `json:"earliestDeleteAt"`
	RegulatoryBasis   string            `json:"regulatoryBasis"`
}

func (v RetentionViolation) Error() string {
	return fmt.Sprintf(
		"compliance retention: %s record %s is under retention. Last updated %s; earliest delete at %s (%s).",
		v.Resource,
		v.RecordID,
		v.RecordLastUpdated.Format(time.RFC3339),
		v.EarliestDeleteAt.Format(time.RFC3339),
		v.RegulatoryBasis,
	)
}

// basisFor returns the regulatory citation for the given resource.
// Kept in code so the 409 body explains itself without forcing the
// caller to read memory docs.
func basisFor(res RetentionResource) string {
	switch res {
	case ResourceSubscriber:
		return "TRAI consumer protection + DOT UASL clause 39"
	case ResourceKYCDocument:
		return "PMLA (Prevention of Money Laundering Act) §12"
	case ResourceAuditLog:
		return "CERT-In Directive (Apr 2022)"
	case ResourceLogReport:
		return "CERT-In Directive + DOT UASL (CGNAT/IPDR extracts)"
	case ResourceDNSQueryLog:
		return "CERT-In Directive (subscriber browsing PII)"
	default:
		return "regulatory retention policy"
	}
}

// CheckRetention is the canonical gate every compliance-retain delete
// handler must call. Returns nil when the record is old enough to
// purge; returns *RetentionViolation otherwise (handler MUST emit 409
// with that violation as the response body).
//
// recordLastUpdated should be the most-recent meaningful state change:
//   - subscriber: terminationDate or updatedAt (post-cessation)
//   - kyc doc: rejectedAt or verifiedAt (whichever later)
//   - audit log: createdAt (logs never mutate)
//   - log report: createdAt
//   - dns query log batch: log row timestamp
func CheckRetention(res RetentionResource, recordID string, recordLastUpdated time.Time) error {
	minAge, ok := defaultRetention[res]
	if !ok {
		// Unknown resource — fail closed to be safe.
		minAge = 2 * 365 * 24 * time.Hour
	}
	age := time.Since(recordLastUpdated)
	if age >= minAge {
		return nil
	}
	return &RetentionViolation{
		Resource:          res,
		RecordID:          recordID,
		RecordLastUpdated: recordLastUpdated,
		MinRetention:      minAge,
		EarliestDeleteAt:  recordLastUpdated.Add(minAge),
		RegulatoryBasis:   basisFor(res),
	}
}

// EnforceRetention is a convenience wrapper. Calls CheckRetention; on
// violation, writes a 409 Conflict response with the violation body
// and emits a structured logger.Warn so the rejection is auditable.
// Returns true when the delete should proceed, false when blocked.
//
// Usage:
//
//	if !compliance.EnforceRetention(w, r, compliance.ResourceSubscriber,
//	    id, sub.UpdatedAt) {
//	    return
//	}
//	// safe to delete
func EnforceRetention(
	w http.ResponseWriter,
	r *http.Request,
	res RetentionResource,
	recordID string,
	recordLastUpdated time.Time,
) bool {
	if err := CheckRetention(res, recordID, recordLastUpdated); err != nil {
		v := err.(*RetentionViolation)
		logger.Warn(
			"Compliance retention blocked delete",
			nil,
			"resource", string(res),
			"recordId", recordID,
			"recordLastUpdated", recordLastUpdated.Format(time.RFC3339),
			"earliestDeleteAt", v.EarliestDeleteAt.Format(time.RFC3339),
			"regulatoryBasis", v.RegulatoryBasis,
			"actingUserId", helpers.GetUserID(r.Context()),
			"actingTenantId", helpers.GetTenantID(r.Context()),
		)
		helpers.HandleAppError(w, appError.New(
			appError.ResourceConflict,
			v.Error(),
			http.StatusConflict,
			v,
		))
		return false
	}
	return true
}
