package tenantevent

// MessagingTemplateInsertEventModel represents a messaging template creation event
type MessagingTemplateInsertEventModel struct {
	ID           string                      `bson:"_id" json:"id"`
	Name         string                      `bson:"name" json:"name"`
	Code         string                      `bson:"code" json:"code"`
	Channel      string                      `bson:"channel" json:"channel"`
	Language     string                      `bson:"language" json:"language"`
	Description  string                      `bson:"description,omitempty" json:"description,omitempty"`
	IsActive     bool                        `bson:"isActive" json:"isActive"`
	Tags         []string                    `bson:"tags,omitempty" json:"tags,omitempty"`
	Variables    []MessagingTemplateVariable `bson:"variables" json:"variables"`
	AssignedTo   []string                    `bson:"assignedTo,omitempty" json:"assignedTo,omitempty"`
	AssignableBy []string                    `bson:"assignableBy,omitempty" json:"assignableBy,omitempty"`

	// Channel-specific fields
	SMSBody           string                     `bson:"sms_body,omitempty" json:"sms_body,omitempty"`
	EmailSubject      string                     `bson:"email_subject,omitempty" json:"email_subject,omitempty"`
	EmailTextBody     string                     `bson:"email_text_body,omitempty" json:"email_text_body,omitempty"`
	EmailHTMLBody     string                     `bson:"email_html_body,omitempty" json:"email_html_body,omitempty"`
	TelegramBody      string                     `bson:"telegram_body,omitempty" json:"telegram_body,omitempty"`
	TelegramParseMode string                     `bson:"telegram_parse_mode,omitempty" json:"telegram_parse_mode,omitempty"`
	WhatsApp          *MessagingWhatsAppTemplate `bson:"whatsapp,omitempty" json:"whatsapp,omitempty"`

	Version int `bson:"version" json:"version"`
}

// MessagingTemplateUpdateEventModel represents a messaging template update event
type MessagingTemplateUpdateEventModel struct {
	ID           string                      `bson:"_id" json:"id"`
	Name         string                      `bson:"name" json:"name"`
	Code         string                      `bson:"code" json:"code"`
	Channel      string                      `bson:"channel" json:"channel"`
	Language     string                      `bson:"language" json:"language"`
	Description  string                      `bson:"description,omitempty" json:"description,omitempty"`
	IsActive     bool                        `bson:"isActive" json:"isActive"`
	Tags         []string                    `bson:"tags,omitempty" json:"tags,omitempty"`
	Variables    []MessagingTemplateVariable `bson:"variables" json:"variables"`
	AssignedTo   []string                    `bson:"assignedTo,omitempty" json:"assignedTo,omitempty"`
	AssignableBy []string                    `bson:"assignableBy,omitempty" json:"assignableBy,omitempty"`

	// Channel-specific fields
	SMSBody           string                     `bson:"sms_body,omitempty" json:"sms_body,omitempty"`
	EmailSubject      string                     `bson:"email_subject,omitempty" json:"email_subject,omitempty"`
	EmailTextBody     string                     `bson:"email_text_body,omitempty" json:"email_text_body,omitempty"`
	EmailHTMLBody     string                     `bson:"email_html_body,omitempty" json:"email_html_body,omitempty"`
	TelegramBody      string                     `bson:"telegram_body,omitempty" json:"telegram_body,omitempty"`
	TelegramParseMode string                     `bson:"telegram_parse_mode,omitempty" json:"telegram_parse_mode,omitempty"`
	WhatsApp          *MessagingWhatsAppTemplate `bson:"whatsapp,omitempty" json:"whatsapp,omitempty"`

	Version int `bson:"version" json:"version"`
}

// MessagingTemplateDeleteEventModel represents a messaging template deletion event
type MessagingTemplateDeleteEventModel struct {
	ID string `bson:"_id" json:"id"`
}

// MessagingTemplateVariable represents a template variable definition
type MessagingTemplateVariable struct {
	Key          string      `bson:"key" json:"key"`
	Placeholder  *int        `bson:"placeholder,omitempty" json:"placeholder,omitempty"`
	Type         string      `bson:"type" json:"type"`
	Required     bool        `bson:"required" json:"required"`
	DefaultValue interface{} `bson:"default_value,omitempty" json:"default_value,omitempty"`
	Description  string      `bson:"description,omitempty" json:"description,omitempty"`
}

// MessagingWhatsAppTemplate represents WhatsApp-specific template structure
type MessagingWhatsAppTemplate struct {
	Namespace    string                       `bson:"namespace,omitempty" json:"namespace,omitempty"`
	Name         string                       `bson:"name" json:"name"`
	Category     string                       `bson:"category" json:"category"`
	LanguageCode string                       `bson:"language_code" json:"language_code"`
	Status       string                       `bson:"status" json:"status"`
	Components   []MessagingWhatsAppComponent `bson:"components" json:"components"`
	ProviderMeta map[string]any               `bson:"provider_meta,omitempty" json:"provider_meta,omitempty"`
}

// MessagingWhatsAppComponent represents a WhatsApp template component
type MessagingWhatsAppComponent struct {
	Type    string                    `bson:"type" json:"type"`
	SubType string                    `bson:"sub_type,omitempty" json:"sub_type,omitempty"`
	Text    string                    `bson:"text,omitempty" json:"text,omitempty"`
	Format  string                    `bson:"format,omitempty" json:"format,omitempty"`
	Index   *int                      `bson:"index,omitempty" json:"index,omitempty"`
	Buttons []MessagingWhatsAppButton `bson:"buttons,omitempty" json:"buttons,omitempty"`
	Example *MessagingWhatsAppExample `bson:"example,omitempty" json:"example,omitempty"`
}

// MessagingWhatsAppButton represents a WhatsApp template button
type MessagingWhatsAppButton struct {
	Type  string `bson:"type" json:"type"`
	Text  string `bson:"text" json:"text"`
	Url   string `bson:"url,omitempty" json:"url,omitempty"`
	Phone string `bson:"phone,omitempty" json:"phone,omitempty"`
	Index int    `bson:"index" json:"index"`
}

// MessagingWhatsAppExample represents WhatsApp template examples
type MessagingWhatsAppExample struct {
	HeaderTextExamples [][]string `bson:"header_text_examples,omitempty" json:"header_text_examples,omitempty"`
	BodyTextExamples   [][]string `bson:"body_text_examples,omitempty" json:"body_text_examples,omitempty"`
	UrlButtonExamples  [][]string `bson:"url_button_examples,omitempty" json:"url_button_examples,omitempty"`
}
