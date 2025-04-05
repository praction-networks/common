package events

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/praction-networks/common/logger"
)

type JsStreamManager struct {
	JsClient jetstream.JetStream
}

func NewStreamManager(jsClient jetstream.JetStream) *JsStreamManager {
	return &JsStreamManager{JsClient: jsClient}
}

type StreamConfig struct {
	Name         StreamName                 // Name of the stream
	Description  string                     // Optional description
	Subjects     []Subject                  // Subjects associated with the stream
	Retention    jetstream.RetentionPolicy  // Retention policy
	MaxConsumers int                        // Max consumers
	MaxMsgs      int64                      // Max number of messages
	MaxAge       time.Duration              // Max age of messages
	Discard      jetstream.DiscardPolicy    // Discard policy
	Storage      jetstream.StorageType      // Storage type (FileStorage or MemoryStorage)
	Replicas     int                        // Replication factor
	NoAck        bool                       // Disable acknowledgment
	Compression  jetstream.StoreCompression // Compression algorithm
}

// ToJetStreamSubjects converts the subjects to their string representations
func (sc *StreamConfig) ToJetStreamSubjects() []string {
	subjects := make([]string, len(sc.Subjects))
	for i, s := range sc.Subjects {
		subjects[i] = string(s)
	}
	return subjects
}

// CreateOrUpdateStream ensures the stream exists or updates it if necessary.
func (jsm *JsStreamManager) CreateOrUpdateStream(ctx context.Context, config StreamConfig) error {
	_, err := jsm.JsClient.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:         string(config.Name),
		Description:  config.Description,
		Subjects:     config.ToJetStreamSubjects(),
		Retention:    config.Retention,
		MaxConsumers: config.MaxConsumers,
		MaxMsgs:      config.MaxMsgs,
		MaxAge:       config.MaxAge,
		Discard:      config.Discard,
		Storage:      config.Storage,
		Replicas:     config.Replicas,
		NoAck:        config.NoAck,
		Compression:  config.Compression,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update stream %s: %w", config.Name, err)
	}

	logger.Info(fmt.Sprintf("Stream %s created/updated successfully.", config.Name))
	return nil
}

// Stream fetches StreamInfo and returns it.
func (jsm *JsStreamManager) Stream(ctx context.Context, streamName string) (*jetstream.StreamInfo, error) {
	stream, err := jsm.JsClient.Stream(ctx, streamName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stream %s: %w", streamName, err)
	}
	streamInfo, err := stream.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stream info for %s: %w", streamName, err)
	}
	return streamInfo, nil
}

// DeleteStream removes a stream with the given name.
func (jsm *JsStreamManager) DeleteStream(ctx context.Context, streamName string) error {
	if err := jsm.JsClient.DeleteStream(ctx, streamName); err != nil {
		return fmt.Errorf("failed to delete stream %s: %w", streamName, err)
	}
	logger.Info(fmt.Sprintf("Stream %s deleted successfully.", streamName))
	return nil
}

// StreamNameBySubject returns the name of a stream that listens on the given subject.
func (jsm *JsStreamManager) StreamNameBySubject(ctx context.Context, subject string) (string, error) {
	streamName, err := jsm.JsClient.StreamNameBySubject(ctx, subject)
	if err != nil {
		return "", fmt.Errorf("failed to fetch stream name by subject %s: %w", subject, err)
	}
	return streamName, nil
}
