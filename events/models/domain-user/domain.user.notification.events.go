package domainUserEventModdel

import "go.mongodb.org/mongo-driver/bson/primitive"

type DomainaUserAuthInitNotification struct {
	UUID       string       `json:"uuid,omitempty"`
	UserAccess []UserAccess `json:"userAccess"`
}

type UserAccess struct {
	Domain string `json:"domain" bson:"domain" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role   string `json:"role"  bson:"role" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type RoleInitNotification struct {
	ID       primitive.ObjectID   `json:"id" bson:"_id"`
	UUID     string               `json:"uuid" bson:"uuid"`
	Policies []primitive.ObjectID `json:"policies" bson:"policies"`
	Version  int                  `json:"version" bson:"version"`
}

type RoleUpdateNotification struct {
	ID       primitive.ObjectID   `json:"id" bson:"_id"`
	UUID     string               `json:"uuid" bson:"uuid"`
	Policies []primitive.ObjectID `json:"policies" bson:"policies"`
	Version  int                  `json:"version" bson:"version"`
}

type RoleDeleteNotification struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Version int                `json:"version" bson:"version"`
}

type PolicyInitNotification struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	UUID    string             `json:"uuid" bson:"uuid"`
	Object  string             `json:"object" bson:"object"`
	Action  string             `json:"action" bson:"action"`
	Version int                `json:"version" bson:"version"`
}

type DomainaUserCreatedNotification struct {
	UUID                string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	ReqID               string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName           string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile              string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp            string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email               string `json:"email" bson:"email" validate:"required,email"`
	PasswordCreateToken string `json:"passwordCreateToken" bson:"passwordCreateToken"`
}

type DomainaUserForgetPasswordNotification struct {
	UUID                string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	ReqID               string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName           string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile              string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp            string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email               string `json:"email" bson:"email" validate:"required,email"`
	PasswordCreateToken string `json:"passwordCreateToken" bson:"passwordCreateToken"`
}

type DomainaUserResetPasswordNotification struct {
	UUID               string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	ReqID              string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName          string `json:"firstName" bson:"firstName" validate:"required"`
	Mobile             string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp           string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email              string `json:"email" bson:"email" validate:"required,email"`
	PasswordResetToken string `json:"passwordResetToken" bson:"passwordResetToken"`
}

type DomainUserChangeEmailNotification struct {
	UUID      string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	ReqID     string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile    string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp  string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
}

type DomainUserChangeMobileNotification struct {
	UUID      string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	ReqID     string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile    string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp  string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
}

type DomainUserChangeWhatsappNotification struct {
	UUID      string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	ReqID     string `json:"reqID,omitempty" bson:"reqID,omitempty"`
	FirstName string `json:"firstName" bson:"firstName" validate:"required,oneword"`
	Mobile    string `json:"mobile" bson:"mobile" validate:"required,e164"`
	Whatsapp  string `json:"whatsapp" bson:"whatsapp" validate:"required,e164"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
}
