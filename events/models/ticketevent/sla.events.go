package ticketevent

// SLABreachedEvent flags SLA violations.
type SLABreachedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   SLABreachedPayload `json:"payload" bson:"payload"`
}

type SLABreachedPayload struct {
	TicketID   string `json:"ticketId" bson:"ticketId"`
	TenantID   string `json:"tenantId" bson:"tenantId"`
	SLAID      string `json:"slaId,omitempty" bson:"slaId,omitempty"`
	SLAType    string `json:"slaType" bson:"slaType"`
	BreachedAt string `json:"breachedAt" bson:"breachedAt"`
}

// SLAStateChangedEvent captures SLA recalculations.
type SLAStateChangedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   SLAStateChangedPayload `json:"payload" bson:"payload"`
}

type SLAStateChangedPayload struct {
	TicketID string          `json:"ticketId" bson:"ticketId"`
	TenantID string          `json:"tenantId" bson:"tenantId"`
	SLA      SLAStateDetails `json:"sla" bson:"sla"`
	Breached bool            `json:"breached" bson:"breached"`
}

type SLAStateDetails struct {
	PolicyID        string     `json:"policyId" bson:"policyId"`
	ResponseDueAt   string     `json:"responseDueAt,omitempty" bson:"responseDueAt,omitempty"`
	ResolutionDueAt string     `json:"resolutionDueAt,omitempty" bson:"resolutionDueAt,omitempty"`
	Paused          bool       `json:"paused,omitempty" bson:"paused,omitempty"`
	Pauses          []SLAPause `json:"pauses,omitempty" bson:"pauses,omitempty"`
}

type SLAPause struct {
	From   string `json:"from" bson:"from"`
	To     string `json:"to,omitempty" bson:"to,omitempty"`
	Reason string `json:"reason,omitempty" bson:"reason,omitempty"`
}
