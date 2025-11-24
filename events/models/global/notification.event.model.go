package global

// NotificationModel represents a notification event that can be sent via multiple channels
// Example: For login OTP, TemplateCode="login_otp", Channel=["sms", "email"]
// The system will find templates with code="login_otp" for each channel (sms, email)
type NotificationModel struct {
	// TenantID identifies which tenant this notification belongs to
	// Used to find providers and templates assigned to this tenant
	TenantID string `json:"tenant_id"`

	// Channel specifies which communication channels to use
	// Each channel will have its own template (e.g., login_otp for SMS, login_otp for email)
	// Example: ["sms", "email"] or ["whatsapp"] or ["sms", "email", "whatsapp", "telegram"]
	Channel []string `json:"channel" default:"sms,email,whatsapp,telegram"`

	// TemplateCode identifies which template to use for this notification (REQUIRED)
	// Examples: "login_otp", "password_reset", "welcome", "invoice", etc.
	// The system will look up templates with this code for each channel
	// For each channel in Channel[], it will find: template where code=TemplateCode AND channel=channel
	TemplateCode string `json:"template_code"`

	// Variables contains the template variables to replace in the template
	// Example: {"otp": "123456", "user_name": "John Doe", "expiry_minutes": "5"}
	Variables map[string]string `json:"variables"`

	// Recipient contains recipients for SMS, WhatsApp, Telegram (phone numbers)
	// Example: ["+1234567890", "+9876543210"]
	// Used when Channel contains "sms", "whatsapp", or "telegram"
	SMSRecipients []string `json:"sms_recipients,omitempty"`

	// WhatsAppRecipients contains WhatsApp recipients
	// Used when Channel contains "whatsapp"
	WhatsAppRecipients []string `json:"whatsapp_recipients,omitempty"`

	// TelegramRecipients contains Telegram recipients
	// Used when Channel contains "telegram"
	TelegramRecipients []string `json:"telegram_recipients,omitempty"`

	// MailRecipient contains email-specific recipient information
	// Used when Channel contains "email"
	// If email is in Channel, use MailRecipient.ToMail for recipients
	MailRecipient *MailNotification `json:"mail_recipient,omitempty"`

	// NotificationType is optional metadata for categorization/routing/logging
	// Examples: "authentication", "marketing", "transactional", "system"
	// This is separate from TemplateCode which identifies the specific template
	NotificationTypes []string `json:"notification_types,omitempty"`
}

// MailNotification contains email-specific recipient and sender information
type MailNotification struct {
	// ToMail contains the primary email recipients
	ToMail []string `json:"to_mail"`

	// CCMail contains CC recipients (optional)
	CCMail []string `json:"cc_mail,omitempty"`

	// BCCMail contains BCC recipients (optional)
	BCCMail []string `json:"bcc_mail,omitempty"`

	// FromMail specifies the sender email address (optional, will use provider default if not provided)
	FromMail string `json:"from_mail,omitempty"`

	// SenderName specifies the sender display name (optional)
	SenderName string `json:"sender_name,omitempty"`
}
