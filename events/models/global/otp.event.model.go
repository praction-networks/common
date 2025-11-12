package global

import (
	"time"
)

type OTPEventModel struct {
	Name           string        `json:"name,omitempty"`
	OTP            string        `json:"otp"`
	TenantID       string        `json:"tenant_id"`
	CreatedAt      time.Time     `json:"created_at"`
	ExpiryDuration time.Duration `json:"expiry_duration" default:"30m"`
}
