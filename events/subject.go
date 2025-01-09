package events

type Subjects string

const (
	// Domain-related subjects
	DomainCreated  Subjects = "domain.created"
	DomainUpdated  Subjects = "domain.updated"
	DomainUDeleted Subjects = "domain.deleted"

	// Postal-related subjects
	PostalServerCreated Subjects = "postalserver.created"
	PostalServerUpdated Subjects = "postalserver.updated"
	PostalServerDeleted Subjects = "postalserver.deleted"

	//Domain-User related subjects

	DomainUserCreated Subjects = "domainuser.created"
	DomainUserUpdated Subjects = "domainuser.updated"
	DomainUserDeleted Subjects = "domainuser.deleted"
)

// StreamMetadata defines metadata for streams
type StreamMetadata struct {
	Name        string
	Description string
	Subjects    []Subjects
}

// Predefined stream configurations
var Streams = map[string]StreamMetadata{
	"DomainStream": {
		Name:        "DomainStream",
		Description: "Stream for domain related events",
		Subjects:    []Subjects{DomainCreated, DomainUpdated, DomainUDeleted},
	},
	"DomainUserStream": {
		Name:        "DomainUserStream",
		Description: "Stream for domain user related events",
		Subjects:    []Subjects{DomainUserCreated, DomainUserUpdated, DomainUserDeleted},
	},
	"PostalStream": {
		Name:        "PostalStream",
		Description: "Stream for postal server related events",
		Subjects:    []Subjects{PostalServerCreated, PostalServerUpdated, PostalServerDeleted},
	},
}
