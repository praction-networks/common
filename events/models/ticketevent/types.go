package ticketevent

// EventName enumerates all ticket-domain events emitted by ticket-service.
type EventName string

const (
	TicketCreatedEventName            EventName = "ticket.created"
	TicketUpdatedEventName            EventName = "ticket.updated"
	AssignmentRequestedEventName      EventName = "ticket.assignment.requested"
	TicketCommentAddedEventName       EventName = "ticket.comment.added"
	CustomerRepliedEventName          EventName = "ticket.customer.replied"
	TicketClosedEventName             EventName = "ticket.closed"
	TicketAttachmentAddedEventName    EventName = "ticket.attachment.added"
	TechnicianStatusUpdateEventName   EventName = "ticket.technician.status.update"
	TicketStatusChangedEventName      EventName = "ticket.status.changed"
	TicketAssignedEventName           EventName = "ticket.assigned"
	AppointmentScheduledEventName     EventName = "ticket.appointment.scheduled"
	AppointmentStatusChangedEventName EventName = "ticket.appointment.status.changed"
	ChecklistItemCompletedEventName   EventName = "ticket.checklist.item.completed"
	ChecklistCreatedEventName         EventName = "ticket.checklist.created"
	SLABreachedEventName              EventName = "ticket.sla.breached"
	TicketEscalatedEventName          EventName = "ticket.escalated"
	TicketResolvedEventName           EventName = "ticket.resolved"
	SLAStateChangedEventName          EventName = "ticket.sla.changed"
	AssignmentApprovedEventName       EventName = "ticket.assignment.approved"
	AssignmentDeniedEventName         EventName = "ticket.assignment.denied"
)

// EventHeaders carries metadata accompanying every ticket-domain event.
type EventHeaders struct {
	TenantID      string `json:"tenantId" bson:"tenantId"`
	CorrelationID string `json:"correlationId,omitempty" bson:"correlationId,omitempty"`
	CausationID   string `json:"causationId,omitempty" bson:"causationId,omitempty"`
	UserID        string `json:"userId,omitempty" bson:"userId,omitempty"`
}

// BaseEvent embeds common envelope fields shared across ticket events.
type BaseEvent struct {
	ID         string       `json:"id" bson:"_id"`
	OccurredAt string       `json:"occurredAt" bson:"occurredAt"`
	Name       EventName    `json:"name" bson:"name"`
	Headers    EventHeaders `json:"headers" bson:"headers"`
}
