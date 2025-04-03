package domainadminevent

import (
	"time"
)

type UserAccess struct {
	Domain     string `json:"domain" bson:"domain"`
	Department string `json:"department" bson:"department"`
	Role       string `json:"role" bson:"role"`
}

type DomainUserCreateEvent struct {
	ID                string       `json:"id" bson:"_id"`
	Mobile            string       `json:"mobile"`
	Whatsapp          string       `json:"whatsapp"`
	Email             string       `json:"email"`
	PasswordChangedAt *time.Time   `json:"passwordChangedAt,omitempty" bson:"passwordChangedAt,omitempty"`
	UserAccess        []UserAccess `json:"userAccess"`
	OnLeave           bool         `json:"onLeave" bson:"onLeave"`
	IsActive          bool         `json:"isActive" bson:"isActive"`
	Version           int          `json:"version" bson:"version"`
}

// Domain represents a domain entity with parent-child hierarchy

type DomainUserUpdateEvent struct {
	ID                string       `json:"id" bson:"_id"`
	Mobile            string       `json:"mobile"`
	Whatsapp          string       `json:"whatsapp"`
	Email             string       `json:"email"`
	PasswordChangedAt *time.Time   `json:"passwordChangedAt,omitempty" bson:"passwordChangedAt,omitempty"`
	UserAccess        []UserAccess `json:"userAccess"`
	OnLeave           bool         `json:"onLeave" bson:"onLeave"`
	IsActive          bool         `json:"isActive" bson:"isActive"`
	Version           int          `json:"version" bson:"version"`
}

type DomainUserDeleteEvent struct {
	ID string `json:"id" bson:"_id"`
}

type RoleInitEventModel struct {
	ID          string   `json:"id" bson:"_id"`
	Policies    []string `json:"policies,omitempty" bson:"policies,omitempty"`
	SystemName  string   `json:"systemName" bson:"systemName,omitempty"`
	CreatedBy   string   `json:"createdBy" bson:"createdBy"`
	ParentRefID string   `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"`
	Version     int      `json:"version" bson:"version"`
}

type RoleUpdateEventModel struct {
	ID          string   `json:"id" bson:"_id"`
	Policies    []string `json:"policies" bson:"policies"`
	SystemName  string   `json:"systemName" bson:"systemName,omitempty"`
	CreatedBy   string   `json:"createdBy" bson:"createdBy"`
	ParentRefID string   `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"`
	Version     int      `json:"version" bson:"version"`
}

type RoleDeleteEvenetModel struct {
	ID      string `json:"id" bson:"_id"`
	Version int    `json:"version" bson:"version"`
}

type PolicyInitEventModel struct {
	ID      string `json:"id" bson:"_id"`
	Object  string `json:"object" bson:"object"`
	Action  string `json:"action" bson:"action"`
	Version int    `json:"version" bson:"version"`
}

type DomainaUserCreatedNotification struct {
	ID                  string `json:"id" bson:"_id"`
	ReqID               string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName           string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile              string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp            string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email               string `json:"email" bson:"email" validate:"required,email"`
	PasswordCreateToken string `json:"passwordCreateToken" bson:"passwordCreateToken"`
}

type DomainaUserForgetPasswordNotification struct {
	ID                  string `json:"id" bson:"_id"`
	ReqID               string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName           string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile              string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp            string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email               string `json:"email" bson:"email" validate:"required,email"`
	PasswordCreateToken string `json:"passwordCreateToken" bson:"passwordCreateToken"`
}

type DomainaUserResetPasswordNotification struct {
	ID                 string `json:"id" bson:"_id"`
	ReqID              string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName          string `json:"firstName" bson:"firstName" validate:"required"`
	Mobile             string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp           string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email              string `json:"email" bson:"email" validate:"required,email"`
	PasswordResetToken string `json:"passwordResetToken" bson:"passwordResetToken"`
}

type DomainUserChangeEmailNotification struct {
	ID        string `json:"id" bson:"_id"`
	ReqID     string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile    string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp  string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
}

type DomainUserChangeMobileNotification struct {
	ID        string `json:"id" bson:"_id"`
	ReqID     string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile    string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp  string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
}

type DomainUserChangeWhatsappNotification struct {
	ID        string `json:"id" bson:"_id"`
	ReqID     string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile    string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp  string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
}
