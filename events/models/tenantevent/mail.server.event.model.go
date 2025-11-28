package tenantevent

// MailServerInsertEventModel defines the model for mail server insert events
type MailServerInsertEventModel struct {
	ID                 string           `bson:"_id" json:"id"`
	Name               string           `bson:"name" json:"name" `
	SortCode           string           `bson:"sortCode" json:"sortCode" `
	OwnerTenantID      string           `bson:"ownerTenantID" json:"ownerTenantID" `
	OwnerTenantType    string           `bson:"ownerTenantType" json:"ownerTenantType" `
	Scope              string           `bson:"scope" json:"scope" `
	AllowedTenantTypes []string         `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string         `bson:"explicitTenantIDs,omitempty" json:"explicitTenantIDs,omitempty" `
	SMTPConfig         *SMTPConfig      `bson:"smtp,omitempty" json:"smtp,omitempty"`
	SendGridConfig     *SendGridConfig  `bson:"sendgrid,omitempty" json:"sendgrid,omitempty"`
	MailgunConfig      *MailgunConfig   `bson:"mailgun,omitempty" json:"mailgun,omitempty"`
	PostalConfig       *PostalConfig    `bson:"postal,omitempty" json:"postal,omitempty"`
	MailchimpConfig    *MailchimpConfig `bson:"mailchimp,omitempty" json:"mailchimp,omitempty"`
	IsActive           bool             `bson:"isActive" json:"isActive" `
	Version            int              `bson:"version" json:"version" `
}

// MailServerUpdateEventModel defines the model for mail server update events
type MailServerUpdateEventModel struct {
	ID                 string           `bson:"_id" json:"id"`
	Name               string           `bson:"name,omitempty" json:"name,omitempty" `
	SortCode           string           `bson:"sortCode,omitempty" json:"sortCode,omitempty" `
	OwnerTenantID      string           `bson:"ownerTenantID,omitempty" json:"ownerTenantID,omitempty" `
	OwnerTenantType    string           `bson:"ownerTenantType,omitempty" json:"ownerTenantType,omitempty"`
	Scope              string           `bson:"scope,omitempty" json:"scope,omitempty" `
	AllowedTenantTypes []string         `bson:"allowedTenantTypes,omitempty" json:"allowedTenantTypes,omitempty"`
	ExplicitTenantIDs  []string         `bson:"explicitTenantIDs,omitempty" json:"explicitTenantIDs,omitempty" `
	SMTPConfig         *SMTPConfig      `bson:"smtp,omitempty" json:"smtp,omitempty"`
	SendGridConfig     *SendGridConfig  `bson:"sendgrid,omitempty" json:"sendgrid,omitempty"`
	MailgunConfig      *MailgunConfig   `bson:"mailgun,omitempty" json:"mailgun,omitempty"`
	PostalConfig       *PostalConfig    `bson:"postal,omitempty" json:"postal,omitempty"`
	MailchimpConfig    *MailchimpConfig `bson:"mailchimp,omitempty" json:"mailchimp,omitempty"`
	IsActive           *bool            `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Version            *int             `bson:"version,omitempty" json:"version,omitempty" `
}

// MailServerDeleteEventModel defines the model for mail server delete events
type MailServerDeleteEventModel struct {
	ID      string `bson:"_id" json:"id"`
	Version int    `bson:"version" json:"version"`
}

// SMTPConfig schema for SMTP configuration
type SMTPConfig struct {
	Host       string `bson:"host" json:"host"`
	Port       int    `bson:"port" json:"port"`
	Username   string `bson:"username" json:"username" `
	Password   string `bson:"password" json:"password"`
	Encryption string `bson:"encryption" json:"encryption" `
	From       string `bson:"from" json:"from" `
}

// SendGridConfig schema for SendGrid configuration
type SendGridConfig struct {
	APIKey string `bson:"api_key" json:"api_key" `
	From   string `bson:"from" json:"from" `
}

// MailgunConfig schema for Mailgun configuration
type MailgunConfig struct {
	Domain       string `bson:"domain" json:"domain"`
	APIKey       string `bson:"api_key" json:"api_key" `
	PublicAPIKey string `bson:"public_api_key" json:"public_api_key" `
	From         string `bson:"from" json:"from" `
}

// PostalConfig schema for Postal configuration
type PostalConfig struct {
	ServerURL string `bson:"server_url" json:"server_url" `
	APIKey    string `bson:"api_key" json:"api_key" `
	From      string `bson:"from" json:"from" `
}

// MailchimpConfig schema for Mailchimp configuration
type MailchimpConfig struct {
	APIKey       string `bson:"api_key" json:"api_key" `
	ServerPrefix string `bson:"server_prefix" json:"server_prefix" `
	From         string `bson:"from" json:"from" `
}
