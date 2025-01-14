package domainUserEventModdel

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
