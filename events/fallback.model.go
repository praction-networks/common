package events

import (
	"time"
)

const (
	FailedNATSFieldID         = "_id"
	FailedNATSFieldStreamName = "stream_name"
	FailedNATSFieldSubject    = "subject"
	FailedNATSFieldPayload    = "payload"
	FailedNATSFieldAttempts   = "attempts"
	FailedNATSFieldTimestamp  = "timestamp"
	FailedNATSFieldLastError  = "last_error"
)

type FailedNATSEvent struct {
	ID         string    `bson:"_id" json:"_id"`
	StreamName string    `bson:"stream_name" json:"stream_name"`
	Subject    string    `bson:"subject" json:"subject"`
	Payload    []byte    `bson:"payload" json:"payload"`
	Attempts   int       `bson:"attempts" json:"attempts"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	LastError  string    `bson:"last_error,omitempty" json:"last_error,omitempty"`
}
