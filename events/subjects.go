package events

type Subjects string

const (
	// Domain
	DomainCreated Subjects = "domain.created"

	// Add more as needed
	PostalServerCreated Subjects = "postalserver.created"
)
