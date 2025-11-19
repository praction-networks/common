package tenantevent

type MailServerInsertEventModel struct {
	ID              string           `bson:"_id" json:"id"`
	SortCode        string           `bson:"sort_code" json:"sort_code"`
	AssignedTo      []string         `bson:"assignedTo" json:"assignedTo"`
	SMTPConfig      *SMTPConfig      `bson:"smtp" json:"smtp"`
	SendGridConfig  *SendGridConfig  `bson:"sendgrid" json:"sendgrid"`
	MailgunConfig   *MailgunConfig   `bson:"mailgun" json:"mailgun"`
	PostalConfig    *PostalConfig    `bson:"postal" json:"postal"`
	MailchimpConfig *MailchimpConfig `bson:"mailchimp" json:"mailchimp"`
	IsActive        bool             `bson:"isActive" json:"isActive"`
	IsSystem        bool             `bson:"isSystem" json:"isSystem"`
	Version         int              `bson:"version" json:"version"`
}

type MailServerUpdateEventModel struct {
	ID              string           `bson:"_id" json:"id"`
	SortCode        string           `bson:"sort_code" json:"sort_code"`
	AssignedTo      []string         `bson:"assignedTo" json:"assignedTo"`
	SMTPConfig      *SMTPConfig      `bson:"smtp" json:"smtp"`
	SendGridConfig  *SendGridConfig  `bson:"sendgrid" json:"sendgrid"`
	MailgunConfig   *MailgunConfig   `bson:"mailgun" json:"mailgun"`
	PostalConfig    *PostalConfig    `bson:"postal" json:"postal"`
	MailchimpConfig *MailchimpConfig `bson:"mailchimp" json:"mailchimp"`
	IsActive        bool             `bson:"isActive" json:"isActive"`
	IsSystem        bool             `bson:"isSystem" json:"isSystem"`
	Version         int              `bson:"version" json:"version"`
}

type MailServerDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}

type SMTPConfig struct {
	Host       string `bson:"host" json:"host"`
	Port       int    `bson:"port" json:"port"`
	Username   string `bson:"username" json:"username"`
	Password   string `bson:"password" json:"password"`
	Encryption string `bson:"encryption" json:"encryption"`
	From       string `bson:"from" json:"from"`
}

type SendGridConfig struct {
	APIKey string `bson:"api_key" json:"api_key"`
	From   string `bson:"from" json:"from"`
}

type MailgunConfig struct {
	Domain       string `bson:"domain" json:"domain"`
	APIKey       string `bson:"api_key" json:"api_key"`
	PublicAPIKey string `bson:"public_api_key" json:"public_api_key"`
	From         string `bson:"from" json:"from"`
}

type PostalConfig struct {
	ServerURL string `bson:"server_url" json:"server_url"`
	APIKey    string `bson:"api_key" json:"api_key"`
	From      string `bson:"from" json:"from"`
}

type MailchimpConfig struct {
	APIKey       string `bson:"api_key" json:"api_key"`
	ServerPrefix string `bson:"server_prefix" json:"server_prefix"`
	From         string `bson:"from" json:"from"`
}
