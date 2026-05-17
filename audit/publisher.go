package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/praction-networks/common/events"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"
)

// subjectMap maps resource domains to their NATS audit subjects on
// AuditGlobalStream. The same eight subjects every service emits to.
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

// Publisher emits canonical AuditEvents to AuditGlobalStream.
//
// Takes *events.JsStreamManager (the modern JetStream wrapper every
// service already provides via Wire) rather than the legacy
// nats.JetStreamContext, so callers no longer need GetConn() helpers
// on their natsconnection clients. Resource→subject routing happens
// at publish time via subjectMap.
type Publisher struct {
	streamManager *events.JsStreamManager
	serviceName   string
}

// NewPublisher creates a new audit event publisher for the given service.
// streamManager is typically wired from *natsconnection.NATSJetStreamClient
// via events.ProvideStreamManager or equivalent per-service provider.
func NewPublisher(streamManager *events.JsStreamManager, serviceName string) *Publisher {
	return &Publisher{
		streamManager: streamManager,
		serviceName:   serviceName,
	}
}

// Publish sends an audit event synchronously. Returns an error if
// marshal or JetStream publish fails. Audit publish should rarely block
// callers — prefer PublishAsync from request paths.
func (p *Publisher) Publish(ctx context.Context, event AuditEvent) error {
	if p == nil || p.streamManager == nil || p.streamManager.JsClient == nil {
		return fmt.Errorf("audit publisher not initialized")
	}

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
	if event.UserRole == "" {
		event.UserRole = helpers.GetUserRole(ctx)
	}
	// Note: UserName intentionally NOT denormalized into the audit event.
	// The audit log keeps IDs only — frontend resolves userId → name
	// against tenant-user-service at render time. Keeps the audit
	// store immutable + small, avoids stale snapshots when users rename.

	data, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal audit event", err, "resource", event.Resource, "action", string(event.Action))
		return fmt.Errorf("failed to marshal audit event: %w", err)
	}

	subject := p.resolveSubject(event.Resource)

	if _, err = p.streamManager.JsClient.Publish(ctx, string(subject), data); err != nil {
		logger.Error("Failed to publish audit event", err, "subject", string(subject), "resource", event.Resource)
		return fmt.Errorf("failed to publish audit event: %w", err)
	}

	logger.Debug("Audit event published", "subject", string(subject), "action", string(event.Action), "resource", event.Resource, "resourceId", event.ResourceID)
	return nil
}

// PublishAsync sends an audit event in a goroutine with its own short
// timeout so a slow / unavailable NATS never blocks the request path.
// Errors are logged but not returned.
func (p *Publisher) PublishAsync(ctx context.Context, event AuditEvent) {
	if p == nil {
		return
	}
	// Snapshot ctx values that affect enrichment before detaching, since
	// the caller's ctx may be cancelled by the time the goroutine runs.
	if event.UserID == "" {
		event.UserID = helpers.GetUserID(ctx)
	}
	if event.TenantID == "" {
		event.TenantID = helpers.GetTenantID(ctx)
	}
	if event.UserRole == "" {
		event.UserRole = helpers.GetUserRole(ctx)
	}

	go func() {
		publishCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := p.Publish(publishCtx, event); err != nil {
			logger.Error("Async audit event publish failed", err)
		}
	}()
}

// resolveSubject maps a resource name to the correct NATS subject.
// Falls back to AuditSystemActionSubject for unrecognised resources so
// events are never dropped — they just land on the catch-all subject.
func (p *Publisher) resolveSubject(resource string) events.Subject {
	if subject, ok := subjectMap[resource]; ok {
		return subject
	}
	return events.AuditSystemActionSubject
}
