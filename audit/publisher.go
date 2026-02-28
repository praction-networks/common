package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/praction-networks/common/events"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"
)

// subjectMap maps resource domains to their NATS audit subjects
var subjectMap = map[string]events.Subject{
	"tenant-user": events.AuditUserActionSubject,
	"auth":        events.AuditAuthActionSubject,
	"tenant":      events.AuditTenantActionSubject,
	"subscriber":  events.AuditSubscriberActionSubject,
	"plan":        events.AuditPlanActionSubject,
	"inventory":   events.AuditInventoryActionSubject,
	"ticket":      events.AuditTicketActionSubject,
	"system":      events.AuditSystemActionSubject,
}

// Publisher provides a thin interface for publishing audit events to NATS JetStream.
// Each service creates one Publisher instance via Wire and injects it into handlers.
type Publisher struct {
	js          nats.JetStreamContext
	serviceName string
}

// NewPublisher creates a new audit event publisher for the given service.
func NewPublisher(js nats.JetStreamContext, serviceName string) *Publisher {
	return &Publisher{
		js:          js,
		serviceName: serviceName,
	}
}

// Publish sends an audit event synchronously. Returns an error if publishing fails.
func (p *Publisher) Publish(ctx context.Context, event AuditEvent) error {
	event.ID = uuid.New().String()
	event.Service = p.serviceName

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Enrich from context if not already set
	if event.UserID == "" {
		event.UserID = helpers.GetUserID(ctx)
	}
	if event.TenantID == "" {
		event.TenantID = helpers.GetTenantID(ctx)
	}

	data, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal audit event", err, "resource", event.Resource, "action", string(event.Action))
		return fmt.Errorf("failed to marshal audit event: %w", err)
	}

	subject := p.resolveSubject(event.Resource)

	_, err = p.js.Publish(string(subject), data)
	if err != nil {
		logger.Error("Failed to publish audit event", err, "subject", string(subject), "resource", event.Resource)
		return fmt.Errorf("failed to publish audit event: %w", err)
	}

	logger.Debug("Audit event published", "subject", string(subject), "action", string(event.Action), "resource", event.Resource, "resourceId", event.ResourceID)
	return nil
}

// PublishAsync sends an audit event asynchronously in a goroutine.
// Errors are logged but not returned — audit should never block the request.
func (p *Publisher) PublishAsync(ctx context.Context, event AuditEvent) {
	go func() {
		if err := p.Publish(context.Background(), event); err != nil {
			logger.Error("Async audit event publish failed", err)
		}
	}()
}

// resolveSubject maps a resource name to the correct NATS subject.
// Falls back to system action if the resource is not mapped.
func (p *Publisher) resolveSubject(resource string) events.Subject {
	if subject, ok := subjectMap[resource]; ok {
		return subject
	}
	return events.AuditSystemActionSubject
}
