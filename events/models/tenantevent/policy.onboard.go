package tenantevent

import "fmt"

// PolicyOnboard is the per-tenant onboard-wizard policy bucket.
//
// Source: backend-contract §11.2. Drives subscriber-service KYC channel
// selection. Note: enabledFeatures.core.isUserKYCEnabled lives on the
// existing TenantModel.EnabledFeatures struct and is NOT duplicated here.
type PolicyOnboard struct {
	AgreementChannel string `json:"agreementChannel,omitempty" bson:"agreementChannel,omitempty"` // "SIGNATURE" | "OTP"
	OtpChannel       string `json:"otpChannel,omitempty"       bson:"otpChannel,omitempty"`       // "SMS" | "WHATSAPP" | "BOTH"
}

func (o PolicyOnboard) Validate() error {
	switch o.AgreementChannel {
	case "", "SIGNATURE", "OTP":
		// ok
	default:
		return fmt.Errorf("AgreementChannel must be SIGNATURE or OTP (got %q)", o.AgreementChannel)
	}
	switch o.OtpChannel {
	case "", "SMS", "WHATSAPP", "BOTH":
		// ok
	default:
		return fmt.Errorf("OtpChannel must be SMS, WHATSAPP, or BOTH (got %q)", o.OtpChannel)
	}
	return nil
}
