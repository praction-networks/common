package shiftevent

// Shift is one closed shift session — denormalised rollup written when the
// user transitions OPEN → CLOSED (or CLOSED via auto-close on idle).
// Open shifts live as a subdoc on TenantUser; only closed shifts persist here.
type Shift struct {
	ID          string       `json:"id"                    bson:"_id"`
	UserID      string       `json:"userId"                bson:"userId"`
	TenantID    string       `json:"tenantId"              bson:"tenantId"`
	State       string       `json:"state"                 bson:"state"` // OPEN | ON_BREAK | CLOSED | ABANDONED
	OpenedAt    int64        `json:"openedAt"              bson:"openedAt"`
	ClosedAt    int64        `json:"closedAt,omitempty"    bson:"closedAt,omitempty"`
	Breaks      []ShiftBreak `json:"breaks,omitempty"      bson:"breaks,omitempty"`
	IdleFlagged bool         `json:"idleFlagged,omitempty" bson:"idleFlagged,omitempty"`
}

// ShiftBreak is a single break window inside a shift.
type ShiftBreak struct {
	StartedAt int64  `json:"startedAt"         bson:"startedAt"`
	EndedAt   int64  `json:"endedAt,omitempty" bson:"endedAt,omitempty"`
	Reason    string `json:"reason,omitempty"  bson:"reason,omitempty"` // LUNCH | TRAVEL | OTHER
}
