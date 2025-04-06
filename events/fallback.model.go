package events

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FailedNATSEvent struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	StreamName string             `bson:"stream_name"`
	Subject    string             `bson:"subject"`
	MsgID      string             `bson:"msg_id"`
	Payload    []byte             `bson:"payload"`
	Attempts   int                `bson:"attempts"`
	Timestamp  time.Time          `bson:"timestamp"`
}
