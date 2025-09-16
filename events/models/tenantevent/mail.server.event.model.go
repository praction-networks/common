package tenantevent

type MailServerInsertEventModel struct {
	ID              string           `json:"id"`
	SortCode        string           `json:"sort_code"`
	AssignedTo      []string         `json:"assignedTo"`
	SMTPConfig      *SMTPConfig      `json:"smtp"`
	SendGridConfig  *SendGridConfig  `json:"sendgrid"`
	MailgunConfig   *MailgunConfig   `json:"mailgun"`
	PostalConfig    *PostalConfig    `json:"postal"`
	MailchimpConfig *MailchimpConfig `json:"mailchimp"`
	IsActive        bool             `json:"isActive"`
	Version         int              `json:"version"`
}

type MailServerUpdateEventModel struct {
	ID              string           `json:"id"`
	SortCode        string           `json:"sort_code"`
	AssignedTo      []string         `json:"assignedTo"`
	SMTPConfig      *SMTPConfig      `json:"smtp"`
	SendGridConfig  *SendGridConfig  `json:"sendgrid"`
	MailgunConfig   *MailgunConfig   `json:"mailgun"`
	PostalConfig    *PostalConfig    `json:"postal"`
	MailchimpConfig *MailchimpConfig `json:"mailchimp"`
	IsActive        bool             `json:"isActive"`
	Version         int              `json:"version"`
}

type MailServerDeleteEventModel struct {
	ID string `json:"id"`
}

type SMTPConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Encryption string `json:"encryption"`
	From       string `json:"from"`
}

type SendGridConfig struct {
	APIKey string `json:"api_key"`
	From   string `json:"from"`
}

type MailgunConfig struct {
	Domain       string `json:"domain"`
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"public_api_key"`
	From         string `json:"from"`
}

type PostalConfig struct {
	ServerURL string `json:"server_url"`
	APIKey    string `json:"api_key"`
	From      string `json:"from"`
}

type MailchimpConfig struct {
	APIKey       string `json:"api_key"`
	ServerPrefix string `json:"server_prefix"`
	From         string `json:"from" bson:"from"`
}
