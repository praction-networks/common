package global

type NotificationModel struct {
	TenantID         string            `json:"tenant_id"`
	Channel          []string          `json:"channel" default:"sms,email,whatsapp,telegram"`
	Variables        map[string]string `json:"variables"`
	NotificationType string            `json:"notification_type"`
	Recipient        []string          `json:"recipient"`
	MailRecipient    MailNotifcation   `json:"mail_recipient,omitempty"`
	TemplateCode     string            `json:"template_code,omitempty"`
	Language         string            `json:"language,omitempty" default:"en"`
}

type MailNotifcation struct {
	ToMail     []string `json:"to_mail"`
	CCMail     []string `json:"cc_mail,omitempty"`
	BCCMail    []string `json:"bcc_mail,omitempty"`
	FromMail   string   `json:"from_mail,omitempty"`
	SenderName string   `json:"sender_name,omitempty"`
}
