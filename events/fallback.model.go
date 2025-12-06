package events

import (
	"time"
)

const (
	FailedNATSFieldID         = "_id"
	FailedNATSFieldStreamName = "streamName"
	FailedNATSFieldSubject    = "subject"
	FailedNATSFieldPayload    = "payload"
	FailedNATSFieldAttempts   = "attempts"
	FailedNATSFieldTimestamp  = "timestamp"
	FailedNATSFieldLastError  = "lastError"
)

type FailedNATSEvent struct {
	ID         string    `bson:"_id" json:"_id"`
	StreamName string    `bson:"streamName" json:"streamName"`
	Subject    string    `bson:"subject" json:"subject"`
	Payload    []byte    `bson:"payload" json:"payload"`
	Attempts   int       `bson:"attempts" json:"attempts"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	LastError  string    `bson:"lastError,omitempty" json:"lastError,omitempty"`
}
