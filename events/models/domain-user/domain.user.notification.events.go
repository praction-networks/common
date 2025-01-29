package domainUserEventModdel

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserAccess struct {
	Domain string `json:"domain" bson:"domain"`
	Role   string `json:"role" bson:"role"`
}

type UserDepartment struct {
	Domain     string `json:"domain" bson:"domain"`
	Department string `json:"department" bson:"department"`
}

// Department represents a department entity with parent-child hierarchy
type Department struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UUID        string              `json:"uuid" bson:"uuid,omitempty"`
	Name        string              `json:"name" bson:"name,omitempty"`
	SystemName  string              `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID *primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                 `json:"version" bson:"version"`
}

type Domain struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UUID        string              `json:"uuid" bson:"uuid,omitempty"`
	Name        string              `json:"name" bson:"name,omitempty"`
	SystemName  string              `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID *primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                 `json:"version" bson:"version"`
}
type DomainUserCreateEvent struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UUID           string             `json:"uuid,omitempty"`
	Mobile         string             `json:"mobile"`
	Whatsapp       string             `json:"whatsapp"`
	Email          string             `json:"email"`
	UserAccess     []UserAccess       `json:"userAccess"`
	UserDepartment []UserDepartment   `json:"userDepartment"`
	OnLeave        bool               `json:"onLeave" bson:"onLeave"`
	IsActive       bool               `json:"isActive" bson:"isActive"`
	Version        int                `json:"version" bson:"version"`
}

// Domain represents a domain entity with parent-child hierarchy

type DomainUserUpdateEvent struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	UUID           string             `json:"uuid,omitempty"`
	Mobile         string             `json:"mobile"`
	Whatsapp       string             `json:"whatsapp"`
	Email          string             `json:"email"`
	UserAccess     []UserAccess       `json:"userAccess"`
	UserDepartment []UserDepartment   `json:"userDepartment"`
	OnLeave        bool               `json:"onLeave" bson:"onLeave"`
	IsActive       bool               `json:"isActive" bson:"isActive"`
	Version        int                `json:"version" bson:"version"`
}

type DomainUserDeleteEvent struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

type RoleInitEventModel struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id"`
	UUID        string                `json:"uuid" bson:"uuid"`
	Policies    []*primitive.ObjectID `json:"policies,omitempty" bson:"policies,omitempty"`
	SystemName  string                `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID *primitive.ObjectID   `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"`
	Version     int                   `json:"version" bson:"version"`
}

type RoleUpdateEventModel struct {
	ID       primitive.ObjectID    `json:"id" bson:"_id"`
	UUID     string                `json:"uuid" bson:"uuid"`
	Policies []*primitive.ObjectID `json:"policies" bson:"policies"`
	Version  int                   `json:"version" bson:"version"`
}

type RoleDeleteEvenetModel struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Version int                `json:"version" bson:"version"`
}

type PolicyInitEventModel struct {
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
