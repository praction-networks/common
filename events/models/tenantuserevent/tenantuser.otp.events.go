package tenantuserevent

import (
	"time"
)

type TenantUserOTPSentEvent struct {
	ID        string    `json:"id" bson:"_id"`
	Mobile    string    `json:"mobile" bson:"mobile"`
	Email     string    `json:"email" bson:"email"`
	OTP       string    `json:"otp" bson:"otp"`
	Type      string    `json:"type" bson:"type"` // e.g., "SMS", "Email"
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Version   int       `json:"version" bson:"version"`
}

type TenantUserOTPVerifiedEvent struct {
	ID         string    `json:"id" bson:"_id"`
	Mobile     string    `json:"mobile" bson:"mobile"`
	Email      string    `json:"email" bson:"email"`
	OTP        string    `json:"otp" bson:"otp"`
	VerifiedAt time.Time `json:"verifiedAt" bson:"verifiedAt"`
	Version    int       `json:"version" bson:"version"`
}

type TenantUserOTPExpiredEvent struct {
	ID        string    `json:"id" bson:"_id"`
	Mobile    string    `json:"mobile" bson:"mobile"`
	Email     string    `json:"email" bson:"email"`
	OTP       string    `json:"otp" bson:"otp"`
	ExpiredAt time.Time `json:"expiredAt" bson:"expiredAt"`
	Version   int       `json:"version" bson:"version"`
}
