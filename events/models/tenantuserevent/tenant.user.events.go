package tenantuserevent

import (
	"time"
)

// UserAccess represents user access permissions for a specific zone/department combination
type UserAccess struct {
	Zone       string `json:"zone" bson:"zone"`
	Department string `json:"department" bson:"department"`
	Role       string `json:"role" bson:"role"`
	IsPrimary  bool   `json:"isPrimary" bson:"isPrimary"` // Primary assignment for this user
	IsActive   bool   `json:"isActive" bson:"isActive"`   // Active assignment
}

// Domain User Events (Extended Profiles)
type DomainUserCreateEvent struct {
	ID         string       `json:"id" bson:"_id"`
	FirstName  string       `json:"firstName" bson:"firstName"`
	LastName   string       `json:"lastName" bson:"lastName"`
	Mobile     string       `json:"mobile" bson:"mobile"`
	Whatsapp   string       `json:"whatsapp" bson:"whatsapp"`
	Email      string       `json:"email" bson:"email"`
	Gender     string       `json:"gender" bson:"gender"`
	DOB        string       `json:"dob" bson:"dob"`
	UserAccess []UserAccess `json:"userAccess" bson:"userAccess"` // Multiple zone/department assignments
	OnLeave    bool         `json:"onLeave" bson:"onLeave"`
	IsActive   bool         `json:"isActive" bson:"isActive"`
	Version    int          `json:"version" bson:"version"`
}

type DomainUserUpdateEvent struct {
	ID         string       `json:"id" bson:"_id"`
	FirstName  string       `json:"firstName" bson:"firstName"`
	LastName   string       `json:"lastName" bson:"lastName"`
	Mobile     string       `json:"mobile" bson:"mobile"`
	Whatsapp   string       `json:"whatsapp" bson:"whatsapp"`
	Email      string       `json:"email" bson:"email"`
	Gender     string       `json:"gender" bson:"gender"`
	DOB        string       `json:"dob" bson:"dob"`
	UserAccess []UserAccess `json:"userAccess" bson:"userAccess"` // Updated assignments
	OnLeave    bool         `json:"onLeave" bson:"onLeave"`
	IsActive   bool         `json:"isActive" bson:"isActive"`
}

type DomainUserDeleteEvent struct {
	ID        string    `json:"id" bson:"_id"`
	DeletedAt time.Time `json:"deletedAt" bson:"deletedAt"`
}
