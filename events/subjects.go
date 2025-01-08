package events

type Subjects string

const (
	// Domain
	DomainCreated Subjects = "domain.created"

	// Postal
	PostalServerCreated Subjects = "postalserver.created"
)
