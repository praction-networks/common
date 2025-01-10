package events

type Subjects string

const (
	// Domain-related subjects
	DomainEvenet       Subjects = "domain.*"
	PostalServerEvenet Subjects = "postalserver.*"
	DomainUserEveent   Subjects = "domainuser.*"
)

// StreamMetadata defines metadata for streams
type StreamMetadata struct {
	Name        string
	Description string
	Subjects    Subjects
}

// Predefined stream configurations
var Streams = map[string]StreamMetadata{
	"DomainStream": {
		Name:        "DomainStream",
		Description: "Stream for domain related events",
		Subjects:    Subjects(DomainEvenet),
	},
	"DomainUserStream": {
		Name:        "DomainUserStream",
		Description: "Stream for domain user related events",
		Subjects:    Subjects(DomainUserEveent),
	},
	"PostalStream": {
		Name:        "PostalStream",
		Description: "Stream for postal server related events",
		Subjects:    Subjects(PostalServerEvenet),
	},
}
