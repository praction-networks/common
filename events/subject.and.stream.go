package events

type StreamName string

// Stream names as constants
const (
	DomainStream       StreamName = "DomainStream"
	PostalServerStream StreamName = "PostalServerStream"
	DomainUserStream   StreamName = "DomainUserStream"
)

// Subjects defines the NATS Subjects for different events
type Subject string

const (
	DomainCreatedSubject     Subject = "domain.created"
	DomainUpdatedSubject     Subject = "domain.updated"
	DomainDeletedSubject     Subject = "domain.deleted"
	PostalCreatedSubject     Subject = "postal.created"
	PostalUpdatedSubject     Subject = "postal.updated"
	PostalDeletedSubject     Subject = "postal.deleted"
	DomainUserCreatedSubject Subject = "domainuser.created"
	DomainUserUpdatedSubject Subject = "domainuser.updated"
	DomainUserDeletedSubject Subject = "domainuser.deleted"
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
		Subjects:    []Subject{DomainCreatedSubject, DomainUpdatedSubject, DomainDeletedSubject},
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
}
