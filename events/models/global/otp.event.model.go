package global

import (
	"time"
)

type OTPEventModel struct {
	Name           string        `json:"name,omitempty"`
	OTP            string        `json:"otp"`
	TenantID       string        `json:"tenant_id"`
	ExpiryDuration time.Duration `json:"expiry_duration" default:"30m"`
}

type OTPResendEventModel struct {
	Name           string        `json:"name"`
	OTP            string        `json:"otp"`
	TenantID       string        `json:"tenant_id"`
	Resend         bool          `json:"resend" default:"false"`
	Type           string        `json:"type" default:"sms"`
	ExpiryDuration time.Duration `json:"expiry_duration" default:"30m"`
}
