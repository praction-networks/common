package hierarchy

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/praction-networks/common/events"
	"github.com/praction-networks/common/events/models/tenantevent"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"
)

// TenantHierarchyCache defines the interface for accessing tenant hierarchy data
type TenantHierarchyCache interface {
	// Get returns hierarchy data for a tenant
	Get(tenantID string) (*helpers.TenantHierarchyData, bool)

	// GetFresh returns hierarchy data for a tenant, falling back to the
	// configured provider when the cached entry is stale (e.g. has
	// SetupComplete=false). When no provider is set this behaves identically
	// to Get. The optional provider is wired via SetProvider at startup.
	// Concrete provider lookups are expected to be Redis-fast — the helper
	// caches the result back into the in-memory map so subsequent reads
	// stay on the hot path.
	GetFresh(ctx context.Context, tenantID string) (*helpers.TenantHierarchyData, bool)

	// SetProvider wires an optional provider used by GetFresh to refresh
	// stale entries (entries with SetupComplete=false). Pass nil to disable.
	SetProvider(provider helpers.TenantHierarchyProvider)

	// Set updates or adds a tenant to the cache
	Set(tenant *helpers.TenantHierarchyData)

	// Remove removes a tenant from the cache
	Remove(tenantID string)

	// IsChild checks if childID is a descendant of parentID (direct or indirect)
	IsChild(parentID, childID string) bool

	// GetDescendants returns all descendant tenant IDs for a given parent tenant
	// This includes all children, grandchildren, etc. (recursive)
	// Used for computing accessible tenants for query filtering
	GetDescendants(parentID string) []string

	// StartSync starts the NATS listener to keep the cache updated
	StartSync(ctx context.Context, streamManager *events.JsStreamManager) error

	// LoadInitialData populates the cache with initial data (e.g. from an API call)
	LoadInitialData(tenants []*helpers.TenantHierarchyData)
}

// InMemoryCache implements TenantHierarchyCache using a thread-safe map
type InMemoryCache struct {
	cache    map[string]*helpers.TenantHierarchyData
	mutex    sync.RWMutex
	provider helpers.TenantHierarchyProvider
}

// NewInMemoryCache creates a new empty cache
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]*helpers.TenantHierarchyData),
	}
}

// SetProvider wires an optional TenantHierarchyProvider used by GetFresh as
// a fallback when the cached entry is marked SetupComplete=false. The
// provider call is the only request path that touches a hot DB/Redis lookup,
// so callers should pass a provider that hits Redis first.
func (c *InMemoryCache) SetProvider(provider helpers.TenantHierarchyProvider) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.provider = provider
}

// GetFresh returns the cached entry; when the entry is missing or marked
// SetupComplete=false AND a provider has been wired, it consults the provider
// to confirm and writes the fresh value back into the cache. This is the
// fail-open path for the tenant-setup guard: a stale NATS event must not
// trap a user whose tenant has already completed setup.
func (c *InMemoryCache) GetFresh(ctx context.Context, tenantID string) (*helpers.TenantHierarchyData, bool) {
	data, exists := c.Get(tenantID)
	// Fast path: cache hit AND already complete — no provider call needed.
	if exists && data.SetupComplete {
		return data, true
	}

	c.mutex.RLock()
	provider := c.provider
	c.mutex.RUnlock()
	if provider == nil {
		return data, exists
	}

	fresh, err := provider.GetTenantByID(ctx, tenantID)
	if err != nil || fresh == nil {
		// Provider failure — degrade to the cached value (which may be
		// stale or missing). Logging happens at the guard layer so we
		// don't double-log for every read.
		return data, exists
	}
	c.Set(fresh)
	return fresh, true
}

func (c *InMemoryCache) Get(tenantID string) (*helpers.TenantHierarchyData, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	data, exists := c.cache[tenantID]
	if !exists {
		return nil, false
	}
	// Return copy to prevent external modification
	copy := *data
	// Deep copy ancestors slice
	if data.Ancestors != nil {
		copy.Ancestors = make([]string, len(data.Ancestors))
		for i, v := range data.Ancestors {
			copy.Ancestors[i] = v
		}
	}
	return &copy, true
}

func (c *InMemoryCache) Set(tenant *helpers.TenantHierarchyData) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[tenant.ID] = tenant
}

func (c *InMemoryCache) Remove(tenantID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cache, tenantID)
}

func (c *InMemoryCache) IsChild(parentID, childID string) bool {
	child, exists := c.Get(childID)
	if !exists {
		return false
	}
	for _, ancestorID := range child.Ancestors {
		if ancestorID == parentID {
			return true
		}
	}
	return false
}

// GetDescendants returns all descendant tenant IDs for a given parent tenant
// This is computed by iterating through all tenants and checking if parentID is in their Ancestors
func (c *InMemoryCache) GetDescendants(parentID string) []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var descendants []string
	for tenantID, data := range c.cache {
		// Skip the parent itself
		if tenantID == parentID {
			continue
		}
		// Check if parentID is in this tenant's ancestors (meaning this is a descendant)
		for _, ancestorID := range data.Ancestors {
			if ancestorID == parentID {
				descendants = append(descendants, tenantID)
				break
			}
		}
	}
	return descendants
}

func (c *InMemoryCache) LoadInitialData(tenants []*helpers.TenantHierarchyData) {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, t := range tenants {
		c.cache[t.ID] = t
	}
	logger.Info("Loaded initial tenant hierarchy data", "count", len(tenants))
}

// StartSync starts the NATS listener to consume tenant events
func (c *InMemoryCache) StartSync(ctx context.Context, streamManager *events.JsStreamManager) error {
	filterSubjects := []events.Subject{
		events.TenantCreatedSubject,
		events.TenantUpdatedSubject,
		events.TenantDeletedSubject,
	}

	listener := events.NewListener(
		events.TenantStream,
		"tenant-hierarchy-cache-consumer", // Durable name
		jetstream.DeliverNewPolicy,        // Start with new messages (we assume LoadInitialData handles history)
		jetstream.AckExplicitPolicy,
		30*time.Second,
		nil,
		filterSubjects,
		streamManager,
		c.handleEvent,
	)

	// Since NewListener uses DeliverNewPolicy by default here but we might want DeliverAll if we don't have initial load?
	// The implementation plan says "Dump & Load" first, then listen. Sinc
	// DeliverNewPolicy is appropriate if we load from DB first.
	// But if we want to be purely event driven, DeliverAll could work but might be slow startup.
	// We'll stick to DeliverLastPerSubject or DeliverNew IF we load initially.
	// Let's use DeliverNewPolicy for now as per "Event Propagation" validation.

	logger.Info("Starting Tenant Hierarchy Cache Sync Listener")
	go func() {
		if err := listener.Listen(ctx); err != nil {
			logger.Error("Failed to start hierarchy cache listener", err)
		}
	}()

	return nil
}

func (c *InMemoryCache) handleEvent(ctx context.Context, msg events.Event[json.RawMessage]) error {
	var err error
	switch msg.Subject {
	case events.TenantCreatedSubject:
		var event tenantevent.TenantInsertEventModel
		if err = json.Unmarshal(msg.Data, &event); err == nil {
			c.Set(&helpers.TenantHierarchyData{
				ID:            event.ID,
				Ancestors:     event.Ancestors,
				Level:         event.Level,
				IsSystem:      event.IsSystem,
				SetupComplete: event.SetupComplete,
			})
			logger.Debug("Cache updated from TenantCreated", "tenant_id", event.ID)
		}
	case events.TenantUpdatedSubject:
		var event tenantevent.TenantUpdateEventModel
		if err = json.Unmarshal(msg.Data, &event); err == nil {
			// Merge with existing or create new
			existing, exists := c.Get(event.ID)
			if !exists {
				existing = &helpers.TenantHierarchyData{ID: event.ID}
			}
			if event.Ancestors != nil {
				existing.Ancestors = event.Ancestors
			}
			if event.Level != 0 { // Assuming 0 is not valid or we check pointer
				existing.Level = event.Level
			}
			if event.IsSystem != nil {
				existing.IsSystem = *event.IsSystem
			}
			if event.SetupComplete != nil {
				existing.SetupComplete = *event.SetupComplete
			}
			c.Set(existing)
			logger.Debug("Cache updated from TenantUpdated", "tenant_id", event.ID)
		}
	case events.TenantDeletedSubject:
		var event tenantevent.TenantDeleteEventModel
		if err = json.Unmarshal(msg.Data, &event); err == nil {
			c.Remove(event.ID)
			logger.Debug("Cache removed from TenantDeleted", "tenant_id", event.ID)
		}
	}

	if err != nil {
		logger.Error("Failed to unmarshal tenant event for cache", err, "subject", msg.Subject)
		return nil // Return nil to ACK and skip bad message
	}
	return nil
}
