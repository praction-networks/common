package ticketevent

// TicketCreatedEvent captures initial creation of a ticket.
type TicketCreatedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketCreatedPayload `json:"payload" bson:"payload"`
}

type TicketCreatedPayload struct {
	TicketID     string `json:"ticketId" bson:"ticketId"`
	TenantID     string `json:"tenantId" bson:"tenantId"`
	Number       string `json:"number" bson:"number"`
	StatusID     string `json:"statusId,omitempty" bson:"statusId,omitempty"`
	StatusName   string `json:"statusName,omitempty" bson:"statusName,omitempty"`
	PriorityID   string `json:"priorityId,omitempty" bson:"priorityId,omitempty"`
	PriorityName string `json:"priorityName,omitempty" bson:"priorityName,omitempty"`
}

// TicketUpdatedEvent records arbitrary ticket field changes.
type TicketUpdatedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketUpdatedPayload `json:"payload" bson:"payload"`
}

type TicketUpdatedPayload struct {
	TicketID string                 `json:"ticketId" bson:"ticketId"`
	TenantID string                 `json:"tenantId" bson:"tenantId"`
	Changes  map[string]interface{} `json:"changes" bson:"changes"`
}

// TicketClosedEvent signals ticket closure.
type TicketClosedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketClosedPayload `json:"payload" bson:"payload"`
}

type TicketClosedPayload struct {
	TicketID string `json:"ticketId" bson:"ticketId"`
	TenantID string `json:"tenantId" bson:"tenantId"`
}

// TicketResolvedEvent marks resolution prior to closure.
type TicketResolvedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketResolvedPayload `json:"payload" bson:"payload"`
}

type TicketResolvedPayload struct {
	TicketID string `json:"ticketId" bson:"ticketId"`
	TenantID string `json:"tenantId" bson:"tenantId"`
}

// TicketStatusChangedEvent captures lifecycle transitions.
type TicketStatusChangedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketStatusChangedPayload `json:"payload" bson:"payload"`
}

type TicketStatusChangedPayload struct {
	TicketID       string `json:"ticketId" bson:"ticketId"`
	TenantID       string `json:"tenantId" bson:"tenantId"`
	FromStatusID   string `json:"fromStatusId,omitempty" bson:"fromStatusId,omitempty"`
	ToStatusID     string `json:"toStatusId,omitempty" bson:"toStatusId,omitempty"`
	FromStatusName string `json:"fromStatusName,omitempty" bson:"fromStatusName,omitempty"`
	ToStatusName   string `json:"toStatusName,omitempty" bson:"toStatusName,omitempty"`
	Reason         string `json:"reason,omitempty" bson:"reason,omitempty"`
}

// TicketAssignedEvent reflects final assignment decisions.
type TicketAssignedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketAssignedPayload `json:"payload" bson:"payload"`
}

type TicketAssignedPayload struct {
	TicketID         string   `json:"ticketId" bson:"ticketId"`
	TenantID         string   `json:"tenantId" bson:"tenantId"`
	AssignedTo       string   `json:"assignedTo,omitempty" bson:"assignedTo,omitempty"`
	AssignedGroups   []string `json:"assignedGroups,omitempty" bson:"assignedGroups,omitempty"`
	PreviousAssignee string   `json:"previousAssignee,omitempty" bson:"previousAssignee,omitempty"`
	Reason           string   `json:"reason,omitempty" bson:"reason,omitempty"`
}

// AssignmentRequestedEvent tracks requests for assignment.
type AssignmentRequestedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   AssignmentRequestedPayload `json:"payload" bson:"payload"`
}

type AssignmentRequestedPayload struct {
	TicketID          string `json:"ticketId" bson:"ticketId"`
	TenantID          string `json:"tenantId" bson:"tenantId"`
	RequestedAssignee string `json:"requestedAssigneeId,omitempty" bson:"requestedAssigneeId,omitempty"`
	AssignmentGroupID string `json:"assignmentGroupId,omitempty" bson:"assignmentGroupId,omitempty"`
	Reason            string `json:"reason,omitempty" bson:"reason,omitempty"`
}

// TicketEscalatedEvent captures escalation triggers.
type TicketEscalatedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketEscalatedPayload `json:"payload" bson:"payload"`
}

type TicketEscalatedPayload struct {
	TicketID        string `json:"ticketId" bson:"ticketId"`
	TenantID        string `json:"tenantId" bson:"tenantId"`
	EscalationLevel int    `json:"escalationLevel" bson:"escalationLevel"`
	Reason          string `json:"reason" bson:"reason"`
}
