package ticketevent

// AssignmentApprovedEvent captures approval decisions.
type AssignmentApprovedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   AssignmentApprovalPayload `json:"payload" bson:"payload"`
}

// AssignmentDeniedEvent captures denial decisions.
type AssignmentDeniedEvent struct {
	BaseEvent `json:",inline" bson:",inline"`
	Payload   AssignmentApprovalPayload `json:"payload" bson:"payload"`
}

type AssignmentApprovalPayload struct {
	TicketID          string `json:"ticketId" bson:"ticketId"`
	TenantID          string `json:"tenantId" bson:"tenantId"`
	RequestedAssignee string `json:"requestedAssigneeId,omitempty" bson:"requestedAssigneeId,omitempty"`
	AssignmentGroupID string `json:"assignmentGroupId,omitempty" bson:"assignmentGroupId,omitempty"`
	ApproverID        string `json:"approverId,omitempty" bson:"approverId,omitempty"`
	Reason            string `json:"reason,omitempty" bson:"reason,omitempty"`
}
