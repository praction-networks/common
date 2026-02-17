package ticketevent

// TicketReopenedEvent captures ticket reopen action.
type TicketReopenedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketReopenedPayload `json:"payload" bson:"payload"`
}

type TicketReopenedPayload struct {
	TicketID string `json:"ticketId" bson:"ticketId"`
	TenantID string `json:"tenantId" bson:"tenantId"`
	Reason   string `json:"reason,omitempty" bson:"reason,omitempty"`
}

// TicketMergedEvent captures ticket merge action.
type TicketMergedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketMergedPayload `json:"payload" bson:"payload"`
}

type TicketMergedPayload struct {
	SourceTicketID string `json:"sourceTicketId" bson:"sourceTicketId"`
	TargetTicketID string `json:"targetTicketId" bson:"targetTicketId"`
	TenantID       string `json:"tenantId" bson:"tenantId"`
	Reason         string `json:"reason,omitempty" bson:"reason,omitempty"`
}

// TicketSplitEvent captures ticket split action.
type TicketSplitEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketSplitPayload `json:"payload" bson:"payload"`
}

type TicketSplitPayload struct {
	SourceTicketID string   `json:"sourceTicketId" bson:"sourceTicketId"`
	NewTicketID    string   `json:"newTicketId" bson:"newTicketId"`
	TenantID       string   `json:"tenantId" bson:"tenantId"`
	MessageIDs     []string `json:"messageIds" bson:"messageIds"`
	Reason         string   `json:"reason,omitempty" bson:"reason,omitempty"`
}

// TicketClassifiedEvent captures NLP classification completion.
type TicketClassifiedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketClassifiedPayload `json:"payload" bson:"payload"`
}

type TicketClassifiedPayload struct {
	TicketID    string   `json:"ticketId" bson:"ticketId"`
	TenantID    string   `json:"tenantId" bson:"tenantId"`
	CategoryKey *string  `json:"categoryKey,omitempty" bson:"categoryKey,omitempty"`
	PriorityKey *string  `json:"priorityKey,omitempty" bson:"priorityKey,omitempty"`
	Tags        []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Confidence  float64  `json:"confidence,omitempty" bson:"confidence,omitempty"`
	Provider    string   `json:"provider,omitempty" bson:"provider,omitempty"` // e.g., "openai", "google"
}

// TicketAutoLinkedEvent captures requester auto-linking to customer/subscriber.
type TicketAutoLinkedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketAutoLinkedPayload `json:"payload" bson:"payload"`
}

type TicketAutoLinkedPayload struct {
	TicketID      string  `json:"ticketId" bson:"ticketId"`
	TenantID      string  `json:"tenantId" bson:"tenantId"`
	CustomerID    *string `json:"customerId,omitempty" bson:"customerId,omitempty"`
	ContactID     *string `json:"contactId,omitempty" bson:"contactId,omitempty"`
	SubscriberID  *string `json:"subscriberId,omitempty" bson:"subscriberId,omitempty"`
	MatchType     string  `json:"matchType" bson:"matchType"`         // "exact", "fuzzy", "manual"
	MatchedEntity string  `json:"matchedEntity" bson:"matchedEntity"` // "customer", "subscriber", "tenant-user"
}

// TicketCreatedFromEmailEvent captures ticket creation from inbound email.
type TicketCreatedFromEmailEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   TicketCreatedFromEmailPayload `json:"payload" bson:"payload"`
}

type TicketCreatedFromEmailPayload struct {
	TicketID     string `json:"ticketId" bson:"ticketId"`
	TenantID     string `json:"tenantId" bson:"tenantId"`
	EmailAddress string `json:"emailAddress" bson:"emailAddress"`
	InboxID      string `json:"inboxId" bson:"inboxId"`
	Provider     string `json:"provider" bson:"provider"` // "gmail", "outlook"
}
