package tenantevent

// PolicyOnboard is the per-tenant onboard-wizard policy bucket.
//
// Source: backend-contract §11.2. Drives subscriber-service KYC channel
// selection. Note: enabledFeatures.core.isUserKYCEnabled lives on the
// existing TenantModel.EnabledFeatures struct and is NOT duplicated here.
type PolicyOnboard struct {
	AgreementChannel string `json:"agreementChannel,omitempty" bson:"agreementChannel,omitempty"` // "SIGNATURE" | "OTP"
	OtpChannel       string `json:"otpChannel,omitempty"       bson:"otpChannel,omitempty"`       // "SMS" | "WHATSAPP"
}
