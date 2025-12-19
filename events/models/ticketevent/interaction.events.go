package ticketevent

// TicketCommentAddedEvent captures comments or notes.
type TicketCommentAddedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketCommentAddedPayload `json:"payload" bson:"payload"`
}

type TicketCommentAddedPayload struct {
	TicketID  string `json:"ticketId" bson:"ticketId"`
	TenantID  string `json:"tenantId" bson:"tenantId"`
	CommentID string `json:"commentId" bson:"commentId"`
	AuthorID  string `json:"authorId,omitempty" bson:"authorId,omitempty"`
	Type      string `json:"type" bson:"type"`
	Text      string `json:"text" bson:"text"`
}

// CustomerRepliedEvent follows inbound replies on tickets.
type CustomerRepliedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   CustomerRepliedPayload `json:"payload" bson:"payload"`
}

type CustomerRepliedPayload struct {
	TicketID  string `json:"ticketId" bson:"ticketId"`
	TenantID  string `json:"tenantId" bson:"tenantId"`
	MessageID string `json:"messageId" bson:"messageId"`
}

// TicketAttachmentAddedEvent signals file uploads.
type TicketAttachmentAddedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketAttachmentAddedPayload `json:"payload" bson:"payload"`
}

type TicketAttachmentAddedPayload struct {
	TicketID     string `json:"ticketId" bson:"ticketId"`
	TenantID     string `json:"tenantId" bson:"tenantId"`
	AttachmentID string `json:"attachmentId" bson:"attachmentId"`
	Filename     string `json:"filename" bson:"filename"`
	URL          string `json:"url,omitempty" bson:"url,omitempty"`
}

// TechnicianStatusUpdateEvent covers field technician updates.
type TechnicianStatusUpdateEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TechnicianStatusUpdatePayload `json:"payload" bson:"payload"`
}

type TechnicianStatusUpdatePayload struct {
	TicketID     string `json:"ticketId" bson:"ticketId"`
	TenantID     string `json:"tenantId" bson:"tenantId"`
	TechnicianID string `json:"technicianId" bson:"technicianId"`
	Status       string `json:"status" bson:"status"`
	Note         string `json:"note,omitempty" bson:"note,omitempty"`
}
