package events

type StreamName string

// Stream names as constants
const (
	SeedAppStream                StreamName = "SeedAppStream"
	DomainStream                 StreamName = "DomainStream"
	PostalServerStream           StreamName = "PostalServerStream"
	DomainUserStream             StreamName = "DomainUserStream"
	DomainUserNotificationStream StreamName = "DomainUserNotificationStream"
	RolesStream                  StreamName = "RolesStream"
	PolicyStream                 StreamName = "PolicyStream"
)

// Subjects defines the NATS Subjects for different events
type Subject string

const (
	DomainCreatedSubject     Subject = "domain.created"
	DomainUpdatedSubject     Subject = "domain.updated"
	DomainDeletedSubject     Subject = "domain.deleted"
	DepartmentCreatedSubject Subject = "department.created"
	DepartmentUpdateSubject  Subject = "department.updated"
	DepartmentDeleteSubject  Subject = "department.delete"

	PostalCreatedSubject Subject = "postal.created"
	PostalUpdatedSubject Subject = "postal.updated"
	PostalDeletedSubject Subject = "postal.deleted"

	SeedDomainUserCreatedSubject Subject = "seed.domainuser.created"
	DomainUserCreatedSubject     Subject = "domainuser.created"
	DomainUserUpdatedSubject     Subject = "domainuser.updated"
	DomainUserDeletedSubject     Subject = "domainuser.deleted"

	SeedRoleCreatedSubject Subject = "seed.role.created"
	RoleCreatedSubject     Subject = "role.created"
	RoleUpdatedSubject     Subject = "role.updated"
	RoleDeletedSubject     Subject = "role.deleted"

	SeedPolicyCreatedSubject Subject = "seed.policy.created"
	PolicyCreatedSubject     Subject = "policy.created"
	PolicyDeletedSubject     Subject = "policy.deleted"

	DomainaUserCreatedNotificationSubject       Subject = "domainuser.created.notification"
	DomainUserForgetPasswordNotificationSubject Subject = "domainuser.forget.password.notification"
	DomainUserResetPasswordNotificationSubject  Subject = "domainuser.reset.password.notification"
	DomainUserChangePasswordNotificationSubject Subject = "domainuser.change.password.notification"
	DomainUserChangeEmailNotificationSubject    Subject = "domainuser.change.email.notification"
	DomainUserChangeMobileNotificationSubject   Subject = "domainuser.change.mobile.notification"
	DomainUserChangeWhatsappNotificationSubject Subject = "domainuser.change.whatsapp.notification"
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
		},
	},
	PostalServerStream: {
		Name:        PostalServerStream,
		Description: "Stream for postal server-related events",
		Subjects:    []Subject{PostalCreatedSubject, PostalUpdatedSubject, PostalDeletedSubject},
	},
	DomainUserStream: {
		Name:        DomainUserStream,
		Description: "Stream for domain user-related events",
		Subjects:    []Subject{DomainUserCreatedSubject, DomainUserUpdatedSubject, DomainUserDeletedSubject},
	},
	PolicyStream: {
		Name:        PolicyStream,
		Description: "Policy events for auth services",
		Subjects:    []Subject{PostalCreatedSubject, PolicyDeletedSubject},
	},
	SeedAppStream: {
		Name:        SeedAppStream,
		Description: "Init Stream for seed inital seeding",
		Subjects:    []Subject{SeedDomainUserCreatedSubject, SeedRoleCreatedSubject, SeedPolicyCreatedSubject},
	},
	RolesStream: {
		Name:        RolesStream,
		Description: "Role events for auth services",
		Subjects:    []Subject{RoleCreatedSubject, RoleUpdatedSubject, RoleDeletedSubject},
	},
	DomainUserNotificationStream: {
		Name:        DomainUserNotificationStream,
		Description: "Stream for domain user notification events",
		Subjects: []Subject{
			DomainaUserCreatedNotificationSubject,
			DomainUserForgetPasswordNotificationSubject,
			DomainUserResetPasswordNotificationSubject,
			DomainUserChangePasswordNotificationSubject,
			DomainUserChangeEmailNotificationSubject,
			DomainUserChangeMobileNotificationSubject,
			DomainUserChangeWhatsappNotificationSubject,
		},
	},
}
