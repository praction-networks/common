package events

type StreamName string

// Stream names as constants
const (
	SeedAppStream          StreamName = "SeedAppStream"
	TenantStream           StreamName = "TenantStream"
	TenantUserStream       StreamName = "TenantUserStream"
	TenantUserRoleStream   StreamName = "TenantUserRoleStream"
	TenantUserPolicyStream StreamName = "TenantUserPolicyStream"
	InventoryStream        StreamName = "InventoryStream"
)

// Subjects defines the NATS Subjects for different events
type Subject string

const (

	//Domain Service Event Initialization
	TenantCreatedSubject         Subject = "tenant.created"
	TenantUpdatedSubject         Subject = "tenant.updated"
	TenantDeletedSubject         Subject = "tenant.deleted"
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
	SeedTenantUserCreatedSubject Subject = "seed.tenantuser.created"
	TenantUserCreatedSubject     Subject = "tenantuser.created"
	TenantUserUpdatedSubject     Subject = "tenantuser.updated"
	TenantUserDeletedSubject     Subject = "tenantuser.deleted"

	// Tenant User Role Service Event Initialization
	SeedTenantUserRoleCreatedSubject Subject = "seed.tenantuserrole.created"
	TenantUserRoleCreatedSubject     Subject = "tenantuserrole.created"
	TenantUserRoleUpdatedSubject     Subject = "tenantuserrole.updated"
	TenantUserRoleDeletedSubject     Subject = "tenantuserrole.deleted"

	// Tenant User Policy Service Event Initialization
	SeedTenantUserPolicyCreatedSubject Subject = "seed.tenantuserpolicy.created"
	TenantUserPolicyCreatedSubject     Subject = "tenantuserpolicy.created"
	TenantUserPolicyUpdatedSubject     Subject = "tenantuserpolicy.updated"
	TenantUserPolicyDeletedSubject     Subject = "tenantuserpolicy.deleted"
)

// StreamMetadata defines metadata for streams
type StreamMetadata struct {
	Name        StreamName
	Description string
	Subjects    []Subject
}

// Predefined stream configurations
var Streams = map[StreamName]StreamMetadata{

	TenantStream: {
		Name:        TenantStream,
		Description: "Stream for domain-related events",
		Subjects: []Subject{TenantCreatedSubject,
			TenantUpdatedSubject,
			TenantDeletedSubject,
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
			// Tenant User Events (Extended Profiles)
			SeedTenantUserCreatedSubject,
			TenantUserCreatedSubject,
			TenantUserUpdatedSubject,
			TenantUserDeletedSubject,

			// Tenant User Role Events
			SeedTenantUserRoleCreatedSubject,
			TenantUserRoleCreatedSubject,
			TenantUserRoleUpdatedSubject,
			TenantUserRoleDeletedSubject,

			// Tenant User Policy Events
			SeedTenantUserPolicyCreatedSubject,
			TenantUserPolicyCreatedSubject,
			TenantUserPolicyUpdatedSubject,
			TenantUserPolicyDeletedSubject,
		},
	},
}
