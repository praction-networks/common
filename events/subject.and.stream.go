package events

// Stream names as constants
const (
	DomainStream       = "DomainStream"
	PostalServerStream = "PostalServerStream"
	DomainUserStream   = "DomainUserStream"
)

// Subjects defines the NATS subjects for different events

const (
	DomainSubject       = "domain.*"
	PostalServerSubject = "postalserver.*"
	DomainUserSubject   = "domainuser.*"
)

// StreamMetadata defines metadata for streams
type StreamMetadata struct {
	Name        string
	Description string
	Subjects    string
}

// Predefined stream configurations
var Streams = map[string]StreamMetadata{
	DomainStream: {
		Name:        DomainStream,
		Description: "Stream for domain-related events",
		Subjects:    DomainSubject,
	},
	PostalServerStream: {
		Name:        PostalServerStream,
		Description: "Stream for postal server-related events",
		Subjects:    PostalServerSubject,
	},
	DomainUserStream: {
		Name:        DomainUserStream,
		Description: "Stream for domain user-related events",
		Subjects:    DomainUserSubject,
	},
}
