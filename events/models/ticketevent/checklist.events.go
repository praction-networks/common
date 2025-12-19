package ticketevent

// ChecklistCreatedEvent notifies creation of technician checklist.
type ChecklistCreatedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   ChecklistCreatedPayload `json:"payload" bson:"payload"`
}

type ChecklistCreatedPayload struct {
	TicketID      string `json:"ticketId" bson:"ticketId"`
	TenantID      string `json:"tenantId" bson:"tenantId"`
	ChecklistID   string `json:"checklistId" bson:"checklistId"`
	ChecklistType string `json:"checklistType" bson:"checklistType"`
}

// ChecklistItemCompletedEvent tracks completion of checklist items.
type ChecklistItemCompletedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   ChecklistItemCompletedPayload `json:"payload" bson:"payload"`
}

type ChecklistItemCompletedPayload struct {
	TicketID    string `json:"ticketId" bson:"ticketId"`
	TenantID    string `json:"tenantId" bson:"tenantId"`
	ChecklistID string `json:"checklistId" bson:"checklistId"`
	ItemID      string `json:"itemId" bson:"itemId"`
	ItemTitle   string `json:"itemTitle" bson:"itemTitle"`
}
