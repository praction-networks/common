package events

import (
	"time"
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
