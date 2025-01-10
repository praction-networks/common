package events

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/praction-networks/common/logger"
)

type StreamManager struct {
	Client nats.JetStreamContext
}

func NewStreamManager(client nats.JetStreamContext) *StreamManager {
	return &StreamManager{Client: client}
}

type StreamConfig struct {
	Name         string               // Name of the stream
	Description  string               // Optional description
	Subjects     []string             // Subjects associated with the stream
	Storage      nats.StorageType     // Storage type (FileStorage or MemoryStorage)
	Retention    nats.RetentionPolicy // Retention policy
	MaxConsumers int                  // Max consumers
	MaxMsgs      int64                // Max messages
	MaxBytes     int64                // Max size in bytes
	MaxAge       time.Duration        // Max age of messages
	Discard      nats.DiscardPolicy   // Discard policy
	Replicas     int                  // Replication factor
	NoAck        bool                 // Disable acknowledgment
}

// CreateOrUpdateStream ensures the stream exists or updates it if necessary.
func (sm *StreamManager) CreateOrUpdateStream(config StreamConfig) error {
	streamInfo, err := sm.Client.StreamInfo(config.Name)
	if err == nil && streamInfo != nil {
		// Check if the existing stream contains all the necessary subjects
		missingSubjects := []string{}
		for _, subject := range config.Subjects {
			if !contains(streamInfo.Config.Subjects, subject) {
				missingSubjects = append(missingSubjects, subject)
			}
		}

		if len(missingSubjects) == 0 {
			logger.Info("Stream already exists and is valid", "streamName", config.Name)
			return nil
		}

		// Update the stream with the missing subjects
		streamInfo.Config.Subjects = append(streamInfo.Config.Subjects, missingSubjects...)
		if _, err := sm.Client.UpdateStream(&streamInfo.Config); err != nil {
			return fmt.Errorf("failed to update stream %s: %w", config.Name, err)
		}

		logger.Info("Stream updated with new subjects", "streamName", config.Name, "newSubjects", missingSubjects)
		return nil
	}

	_, err = sm.Client.AddStream(&nats.StreamConfig{
		Name:         config.Name,
		Description:  config.Description,
		Subjects:     config.Subjects,
		Storage:      config.Storage,
		Retention:    config.Retention,
		MaxConsumers: config.MaxConsumers,
		MaxMsgs:      config.MaxMsgs,
		MaxBytes:     config.MaxBytes,
		MaxAge:       config.MaxAge,
		Discard:      config.Discard,
		Replicas:     config.Replicas,
		NoAck:        config.NoAck,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update stream %s: %w", config.Name, err)
	}

	logger.Info(fmt.Sprintf("Stream %s created/updated successfully.\n", config.Name))
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
