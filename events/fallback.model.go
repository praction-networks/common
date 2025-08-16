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
	ID         string    `bson:"_id"`
	StreamName string    `bson:"stream_name"`
	Subject    string    `bson:"subject"`
	Payload    []byte    `bson:"payload"`
	Attempts   int       `bson:"attempts"`
	Timestamp  time.Time `bson:"timestamp"`
	LastError  string    `bson:"last_error,omitempty"`
}
