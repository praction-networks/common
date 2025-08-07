package events

type StreamName string

// Stream names as constants
const (
	SeedAppStream                StreamName = "SeedAppStream"
	DomainStream                 StreamName = "DomainStream"
	PostalServerStream           StreamName = "PostalServerStream"
	DomainUserStream             StreamName = "DomainUserStream"
	DomainUserNotificationStream StreamName = "DomainUserNotificationStream"
	TenantUserStream             StreamName = "TenantUserStream"
	RolesStream                  StreamName = "RolesStream"
	PolicyStream                 StreamName = "PolicyStream"
	InventoryStream              StreamName = "InventoryStream"
)

// Subjects defines the NATS Subjects for different events
type Subject string

const (

	//Domain Service Event Initialization
	DomainCreatedSubject         Subject = "domain.created"
	DomainUpdatedSubject         Subject = "domain.updated"
	DomainDeletedSubject         Subject = "domain.deleted"
	DepartmentCreatedSubject     Subject = "department.created"
	DepartmentUpdateSubject      Subject = "department.updated"
	DepartmentDeleteSubject      Subject = "department.delete"
	AppMessengerCreateSubject    Subject = "appmessenger.created"
	AppMessengerUpdateSubject    Subject = "appmessenger.updated"
	AppMessengerDeleteSubject    Subject = "appmessenger.deleted"
	KYCGatewayCreatedSubject     Subject = "kycgateway.created"
	KYCGatewayUpdateSubject      Subject = "kycgateway.updated"
	KYCGatewayDeleteSubject      Subject = "kycgateway.deleted"
	MailServerCreatedSubject     Subject = "mailserver.created"
	MailServerUpdateSubject      Subject = "mailserver.updated"
	MailServerDeleteSubject      Subject = "mailserver.deleted"
	PaymentGatewayCreatedSubject Subject = "paymentgateway.created"
	PaymentGatewayUpdateSubject  Subject = "paymentgateway.updated"
	PaymentGatewayDeleteSubject  Subject = "paymentgateway.deleted"
	ExternalRadiusCreatedSubject Subject = "externalradius.created"
	ExternalRadiusUpdateSubject  Subject = "externalradius.updated"
	ExternalRadiusDeleteSubject  Subject = "externalradius.deleted"
	SMSGatewayCreatedSubject     Subject = "smsgateway.created"
	SMSGatewayUpdateSubject      Subject = "smsgateway.updated"
	SMSGatewayDeleteSubject      Subject = "smsgateway.deleted"

	//Domain User Service Event Initialization
	SeedDomainUserCreatedSubject Subject = "seed.domainuser.created"
	DomainUserCreatedSubject     Subject = "domainuser.created"
	DomainUserUpdatedSubject     Subject = "domainuser.updated"
	DomainUserDeletedSubject     Subject = "domainuser.deleted"

	// Tenant User Assignment Events
	TenantUserAssignmentCreatedSubject Subject = "tenantuser.assignment.created"
	TenantUserAssignmentUpdatedSubject Subject = "tenantuser.assignment.updated"
	TenantUserAssignmentDeletedSubject Subject = "tenantuser.assignment.deleted"

	// Tenant User Activity Events
	TenantUserActivityCreatedSubject Subject = "tenantuser.activity.created"
	TenantUserActivityUpdatedSubject Subject = "tenantuser.activity.updated"
	TenantUserActivityDeletedSubject Subject = "tenantuser.activity.deleted"

	// Business-Specific Notification Events
	DomainUserProfileUpdateNotificationSubject Subject = "domainuser.profile.update.notification"
	TenantUserAssignmentNotificationSubject    Subject = "tenantuser.assignment.notification"
)

// StreamMetadata defines metadata for streams
type StreamMetadata struct {
	Name        StreamName
	Description string
	Subjects    []Subject
}

// Predefined stream configurations
var Streams = map[StreamName]StreamMetadata{

	DomainStream: {
		Name:        DomainStream,
		Description: "Stream for domain-related events",
		Subjects: []Subject{DomainCreatedSubject,
			DomainUpdatedSubject,
			DomainDeletedSubject,
			DepartmentCreatedSubject,
			DepartmentUpdateSubject,
			DepartmentDeleteSubject,
			AppMessengerCreateSubject,
			AppMessengerUpdateSubject,
			AppMessengerDeleteSubject,
			KYCGatewayCreatedSubject,
			KYCGatewayUpdateSubject,
			KYCGatewayDeleteSubject,
			MailServerCreatedSubject,
			MailServerUpdateSubject,
			MailServerDeleteSubject,
			PaymentGatewayCreatedSubject,
			PaymentGatewayUpdateSubject,
			PaymentGatewayDeleteSubject,
			ExternalRadiusCreatedSubject,
			ExternalRadiusUpdateSubject,
			ExternalRadiusDeleteSubject,
			SMSGatewayCreatedSubject,
			SMSGatewayUpdateSubject,
			SMSGatewayDeleteSubject,
		},
	},

	TenantUserStream: {
		Name:        TenantUserStream,
		Description: "Stream for tenant user service events",
		Subjects: []Subject{
			// Domain User Events (Extended Profiles)
			SeedDomainUserCreatedSubject,
			DomainUserCreatedSubject,
			DomainUserUpdatedSubject,
			DomainUserDeletedSubject,

			// Tenant User Assignment Events
			TenantUserAssignmentCreatedSubject,
			TenantUserAssignmentUpdatedSubject,
			TenantUserAssignmentDeletedSubject,

			// Tenant User Activity Events
			TenantUserActivityCreatedSubject,
			TenantUserActivityUpdatedSubject,
			TenantUserActivityDeletedSubject,

			// Business-Specific Notification Events
			DomainUserProfileUpdateNotificationSubject,
			TenantUserAssignmentNotificationSubject,
		},
	},
}
