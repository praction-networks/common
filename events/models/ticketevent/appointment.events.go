package ticketevent

// AppointmentScheduledEvent tracks appointment creation.
type AppointmentScheduledEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   AppointmentScheduledPayload `json:"payload" bson:"payload"`
}

type AppointmentScheduledPayload struct {
	TicketID      string `json:"ticketId" bson:"ticketId"`
	TenantID      string `json:"tenantId" bson:"tenantId"`
	AppointmentID string `json:"appointmentId" bson:"appointmentId"`
	ScheduledAt   string `json:"scheduledAt" bson:"scheduledAt"`
	Status        string `json:"status" bson:"status"`
}

// AppointmentStatusChangedEvent records appointment status transitions.
type AppointmentStatusChangedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   AppointmentStatusChangedPayload `json:"payload" bson:"payload"`
}

type AppointmentStatusChangedPayload struct {
	TicketID      string `json:"ticketId" bson:"ticketId"`
	TenantID      string `json:"tenantId" bson:"tenantId"`
	AppointmentID string `json:"appointmentId" bson:"appointmentId"`
	FromStatus    string `json:"fromStatus" bson:"fromStatus"`
	ToStatus      string `json:"toStatus" bson:"toStatus"`
}
