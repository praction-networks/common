package provider

import (
	"context"
	"encoding/json"
)

// VerificationType enumerates supported KYC verification checks.
type VerificationType string

const (
	VerifyPAN        VerificationType = "PAN"
	VerifyDigiLocker VerificationType = "DIGILOCKER"
	VerifyPennydrop  VerificationType = "PENNYDROP"
	VerifyEStamp     VerificationType = "ESTAMP"
	VerifyESign      VerificationType = "ESIGN"
	VerifyGST        VerificationType = "GST"
	VerifyPassport   VerificationType = "PASSPORT"
)

// AllVerificationTypes returns every supported verification type.
func AllVerificationTypes() []VerificationType {
	return []VerificationType{
		VerifyPAN, VerifyDigiLocker, VerifyPennydrop,
		VerifyEStamp, VerifyESign, VerifyGST, VerifyPassport,
	}
}

// KYCRequest is the universal input for any KYC verification.
type KYCRequest struct {
	Type      VerificationType  `json:"type"`
	TenantID  string            `json:"tenantId"`
	Payload   map[string]string `json:"payload"`
	RequestID string            `json:"requestId"`
}

// KYCResult is the universal output from a KYC verification.
type KYCResult struct {
	Verified         bool            `json:"verified"`
	Status           string          `json:"status"` // "success", "invalid", "not_found", "pending"
	Provider         string          `json:"provider"`
	VerificationData map[string]any  `json:"verificationData,omitempty"` // structured, frontend-friendly extracted fields
	Data             map[string]any  `json:"data,omitempty"`             // full provider response
	RawResponse      json.RawMessage `json:"rawResponse,omitempty"`
	TraceID          string          `json:"traceId,omitempty"`
}

// KYCVerifier is the interface that every KYC provider adapter must implement.
// The router calls Supports() to check if a verification type is handled, then
// Verify() to execute the actual provider API call.
type KYCVerifier interface {
	// Name returns the provider name (e.g. "CASHFREE", "SETU").
	Name() string

	// Supports returns true if this adapter can handle the given verification type.
	Supports(vType VerificationType) bool

	// Verify executes the verification against the provider's API.
	// creds contains the tenant's provider credentials from the binding metadata.
	Verify(ctx context.Context, creds map[string]any, req KYCRequest) (*KYCResult, error)
}
