# Tenant Hierarchy Integration Guide for Microservices

## 📋 **Overview**

This guide explains how **other microservices** (User Service, Ticket Service, Helpdesk Service, etc.) should integrate with the **hierarchical multi-tenancy system** using NATS events.

---

## 🎯 **Why Hierarchy Matters for Your Service**

When tenants have parent-child relationships, your service needs to:

1. **Filter data by hierarchy** - Show data for tenant + descendants
2. **Aggregate data up the tree** - Total users/tickets for a tenant's entire subtree
3. **Inherit configurations** - Child tenants may inherit settings from parents
4. **Enforce permissions** - Parent admins may have access to child data
5. **Track relationships** - Know which tenant owns which data

---

## 📦 **New Fields in Tenant Events**

### **TenantInsertEventModel & TenantUpdateEventModel**

Both event models now include these hierarchy fields:

```go
// Materialized Path - Enables efficient subtree queries
Path              string   // "/isp_rapidnet/reseller_mumbai/dist_andheri"
PathDepth         int      // 3 (number of segments in path)

// Ancestry - For bottom-up queries and permissions
Ancestors         []string // ["isp_rapidnet", "reseller_mumbai"]
Level             int      // 2 (0=root, 1=child, 2=grandchild)

// Children Management
IsLeaf            bool     // true if no children
ChildrenCount     int      // 5 (count of direct children)
ChildIDs          []string // ["dist_andheri", "dist_borivali", ...]
AllowedChildTypes []string // ["Distributor", "Operator"]
MaxDepth          int      // 5 (max allowed depth from this tenant)
```

---

## 🔔 **New Event Types**

### **1. TenantParentChangedEventModel**
Published when tenant moves to a different parent.

```go
type TenantParentChangedEventModel struct {
    TenantID       string    // "dist_andheri"
    OldParentID    string    // "reseller_mumbai"
    NewParentID    string    // "reseller_pune"
    OldPath        string    // "/isp_rapidnet/reseller_mumbai/dist_andheri"
    NewPath        string    // "/isp_rapidnet/reseller_pune/dist_andheri"
    OldAncestors   []string  // ["isp_rapidnet", "reseller_mumbai"]
    NewAncestors   []string  // ["isp_rapidnet", "reseller_pune"]
    OldLevel       int       // 2
    NewLevel       int       // 2
    ChangedAt      time.Time
    ChangedBy      string
}
```

**Use Cases:**
- **User Service**: Recompute inherited user quotas from new parent
- **Ticket Service**: Update ticket routing to new parent's support team
- **Billing Service**: Adjust revenue attribution to new parent

---

### **2. TenantChildAddedEventModel**
Published when a child tenant is added.

```go
type TenantChildAddedEventModel struct {
    ParentID      string    // "reseller_mumbai"
    ChildID       string    // "dist_andheri"
    ChildType     string    // "Distributor"
    ChildPath     string    // "/isp_rapidnet/reseller_mumbai/dist_andheri"
    ParentPath    string    // "/isp_rapidnet/reseller_mumbai"
    NewChildCount int       // 3 (parent now has 3 children)
    AddedAt       time.Time
    AddedBy       string
}
```

**Use Cases:**
- **User Service**: Grant parent admin access to child's users
- **Ticket Service**: Enable parent to view child's tickets
- **Analytics Service**: Include child's metrics in parent's dashboard

---

### **3. TenantChildRemovedEventModel**
Published when a child tenant is removed (deleted or reassigned).

```go
type TenantChildRemovedEventModel struct {
    ParentID      string    // "reseller_mumbai"
    ChildID       string    // "dist_andheri"
    ChildType     string    // "Distributor"
    NewChildCount int       // 2 (parent now has 2 children)
    IsNowLeaf     bool      // false (still has other children)
    RemovedAt     time.Time
    RemovedBy     string
}
```

**Use Cases:**
- **User Service**: Revoke parent admin access to removed child's users
- **Ticket Service**: Update ticket visibility
- **Billing Service**: Stop revenue roll-up to old parent

---

### **4. TenantHierarchyRecomputedEventModel**
Published when hierarchy fields are recalculated (rare, usually during migrations).

```go
type TenantHierarchyRecomputedEventModel struct {
    TenantID     string    // "dist_andheri"
    Path         string    // "/isp_rapidnet/reseller_pune/dist_andheri"
    PathDepth    int       // 3
    Ancestors    []string  // ["isp_rapidnet", "reseller_pune"]
    Level        int       // 2
    RecomputedAt time.Time
    Reason       string    // "parent_changed", "migration", "manual_fix"
}
```

**Use Cases:**
- **All Services**: Sync hierarchy data to local cache/denormalized copies

---

## 🛠️ **Integration Patterns**

### **Pattern 1: Store Hierarchy Fields Locally**

Each service should store critical hierarchy fields in their own database for efficient queries.

**Example - User Service:**
```go
type User struct {
    ID                string
    Name              string
    TenantID          string
    
    // ✅ Store hierarchy fields locally
    TenantPath        string   // Copy from TenantInsertEvent
    TenantAncestors   []string // Copy from TenantInsertEvent
    TenantLevel       int      // Copy from TenantInsertEvent
}
```

**Benefits:**
- **Fast queries**: No need to call tenant-service for every query
- **Efficient filtering**: Find all users in a subtree using path regex
- **Aggregations**: Group by level, count by path prefix

---

### **Pattern 2: Query Users in Tenant Hierarchy**

**Get all users in tenant + descendants:**
```go
// In User Service
func GetUsersInHierarchy(tenantID string) ([]*User, error) {
    // Get tenant path from cache or event
    tenant := GetTenantFromCache(tenantID)
    
    // Find all users whose tenantPath starts with this tenant's path
    filter := bson.M{
        "tenantPath": bson.M{
            "$regex": fmt.Sprintf("^%s", tenant.Path),
        },
    }
    
    return userRepo.Find(filter)
}
```

**Get users only at specific level:**
```go
filter := bson.M{
    "tenantLevel": 2, // Only distributors' users
    "tenantAncestors": tenantID, // Under this tenant
}
```

---

### **Pattern 3: Aggregate Data Up the Tree**

**Count total tickets for tenant + all descendants:**
```go
// In Ticket Service
func GetTotalTicketsForHierarchy(tenantID string) (int, error) {
    tenant := GetTenantFromCache(tenantID)
    
    // Count tickets where tenantPath starts with this tenant's path
    filter := bson.M{
        "tenantPath": bson.M{
            "$regex": fmt.Sprintf("^%s", tenant.Path),
        },
        "status": "open",
    }
    
    return ticketRepo.Count(filter)
}
```

---

### **Pattern 4: Event Handling**

**Subscribe to tenant events and update local data:**

```go
// In User Service - Subscribe to NATS tenant events
func setupTenantEventHandlers() {
    // Handle new tenants
    nats.Subscribe("tenant.created", func(msg *nats.Msg) {
        var event tenantevent.TenantInsertEventModel
        json.Unmarshal(msg.Data, &event)
        
        // Store tenant metadata in local cache
        cacheTenant := LocalTenant{
            ID:        event.ID,
            Name:      event.Name,
            Path:      event.Path,
            Ancestors: event.Ancestors,
            Level:     event.Level,
        }
        tenantCache.Set(event.ID, cacheTenant)
    })
    
    // Handle tenant updates
    nats.Subscribe("tenant.updated", func(msg *nats.Msg) {
        var event tenantevent.TenantUpdateEventModel
        json.Unmarshal(msg.Data, &event)
        
        // Update cached tenant data
        tenantCache.Update(event.ID, event)
    })
    
    // Handle parent changes - CRITICAL for hierarchy-dependent logic
    nats.Subscribe("tenant.parent.changed", func(msg *nats.Msg) {
        var event tenantevent.TenantParentChangedEventModel
        json.Unmarshal(msg.Data, &event)
        
        // Update ALL users under this tenant with new path
        userRepo.UpdateMany(
            bson.M{"tenantID": event.TenantID},
            bson.M{
                "$set": bson.M{
                    "tenantPath":      event.NewPath,
                    "tenantAncestors": event.NewAncestors,
                    "tenantLevel":     event.NewLevel,
                },
            },
        )
        
        // Update cache
        tenantCache.Update(event.TenantID, TenantCacheEntry{
            Path:      event.NewPath,
            Ancestors: event.NewAncestors,
            Level:     event.NewLevel,
        })
    })
    
    // Handle child added
    nats.Subscribe("tenant.child.added", func(msg *nats.Msg) {
        var event tenantevent.TenantChildAddedEventModel
        json.Unmarshal(msg.Data, &event)
        
        // Grant parent admins access to child's data
        grantHierarchicalAccess(event.ParentID, event.ChildID)
    })
    
    // Handle child removed
    nats.Subscribe("tenant.child.removed", func(msg *nats.Msg) {
        var event tenantevent.TenantChildRemovedEventModel
        json.Unmarshal(msg.Data, &event)
        
        // Revoke parent admins access to removed child's data
        revokeHierarchicalAccess(event.ParentID, event.ChildID)
    })
    
    // Handle tenant deletion
    nats.Subscribe("tenant.deleted", func(msg *nats.Msg) {
        var event tenantevent.TenantDeleteEventModel
        json.Unmarshal(msg.Data, &event)
        
        // Clean up all data for deleted tenant
        cleanupTenantData(event.ID)
        tenantCache.Delete(event.ID)
    })
}
```

---

## 📊 **Use Cases by Service**

### **User Service**

**Data to Store:**
```go
type User struct {
    ID              string
    TenantID        string
    TenantPath      string   // ✅ For subtree queries
    TenantAncestors []string // ✅ For ancestor access checks
    TenantLevel     int      // ✅ For analytics
}
```

**Queries Enabled:**
- Get all users under ISP (including resellers, distributors)
- Count users per level (reseller vs distributor)
- Check if user belongs to parent's subtree

---

### **Ticket Service**

**Data to Store:**
```go
type Ticket struct {
    ID              string
    TenantID        string
    TenantPath      string   // ✅ For hierarchy filtering
    TenantAncestors []string // ✅ For escalation routing
    TenantLevel     int      // ✅ For SLA based on level
}
```

**Features Enabled:**
- Parent ISP can view all tickets from resellers/distributors
- Escalate ticket up the hierarchy
- Different SLA based on hierarchy level

---

### **Helpdesk Service**

**Data to Store:**
```go
type SupportAgent struct {
    ID                string
    TenantID          string
    AccessiblePaths   []string // ✅ Paths this agent can access
    AccessibleLevels  []int    // ✅ Levels this agent can handle
}
```

**Features Enabled:**
- Assign tickets based on tenant hierarchy
- Route to parent support team if child has no agents
- Dashboard showing tickets grouped by hierarchy level

---

### **Analytics/Reporting Service**

**Queries Enabled:**
```go
// Total revenue for ISP + all children
db.invoices.aggregate([
    { $match: { tenantPath: /^\/isp_rapidnet/ } },
    { $group: { _id: null, total: { $sum: "$amount" } } }
])

// Revenue breakdown by level
db.invoices.aggregate([
    { $match: { tenantPath: /^\/isp_rapidnet/ } },
    { $group: { _id: "$tenantLevel", total: { $sum: "$amount" } } }
])

// Revenue by reseller
db.invoices.aggregate([
    { $match: { 
        tenantLevel: { $gte: 1 },
        tenantAncestors: "isp_rapidnet"
    }},
    { $group: { _id: { $arrayElemAt: ["$tenantAncestors", 1] }, total: { $sum: "$amount" } } }
])
```

---

## 🚀 **Step-by-Step Integration**

### **Step 1: Add Hierarchy Fields to Your Models**

```go
// In your service's internal/models/
type YourEntity struct {
    ID        string
    TenantID  string
    
    // ✅ Add these fields
    TenantPath      string   `bson:"tenantPath" json:"tenantPath"`
    TenantAncestors []string `bson:"tenantAncestors" json:"tenantAncestors"`
    TenantLevel     int      `bson:"tenantLevel" json:"tenantLevel"`
    
    // Your other fields...
}
```

---

### **Step 2: Create Indexes**

```go
// In your database initialization
func createIndexes() {
    // Path index for subtree queries
    collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys: bson.D{{Key: "tenantPath", Value: 1}},
    })
    
    // Ancestors index for ancestor lookups
    collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys: bson.D{{Key: "tenantAncestors", Value: 1}},
    })
    
    // Level index for level-based queries
    collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys: bson.D{{Key: "tenantLevel", Value: 1}},
    })
    
    // Compound index for common queries
    collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys: bson.D{
            {Key: "tenantPath", Value: 1},
            {Key: "status", Value: 1},
            {Key: "createdAt", Value: -1},
        },
    })
}
```

---

### **Step 3: Subscribe to Tenant Events**

```go
// In your service's event handler
func initTenantEventHandlers(nc *nats.Conn) {
    // 1. New tenant created
    nc.Subscribe("tenant.created", handleTenantCreated)
    
    // 2. Tenant updated (path might change!)
    nc.Subscribe("tenant.updated", handleTenantUpdated)
    
    // 3. Parent changed (CRITICAL - path/ancestors/level all change)
    nc.Subscribe("tenant.parent.changed", handleTenantParentChanged)
    
    // 4. Child added (update aggregations)
    nc.Subscribe("tenant.child.added", handleTenantChildAdded)
    
    // 5. Child removed (update aggregations)
    nc.Subscribe("tenant.child.removed", handleTenantChildRemoved)
    
    // 6. Tenant deleted (cleanup)
    nc.Subscribe("tenant.deleted", handleTenantDeleted)
    
    // 7. Hierarchy recomputed (sync data)
    nc.Subscribe("tenant.hierarchy.recomputed", handleHierarchyRecomputed)
}
```

---

### **Step 4: Implement Event Handlers**

```go
func handleTenantCreated(msg *nats.Msg) {
    var event tenantevent.TenantInsertEventModel
    json.Unmarshal(msg.Data, &event)
    
    // Store tenant metadata in cache
    tenantCache.Set(event.ID, TenantCacheEntry{
        ID:        event.ID,
        Name:      event.Name,
        Path:      event.Path,
        Ancestors: event.Ancestors,
        Level:     event.Level,
    })
    
    logger.Info("Tenant created", 
        "tenantID", event.ID, 
        "path", event.Path, 
        "level", event.Level)
}

func handleTenantParentChanged(msg *nats.Msg) {
    var event tenantevent.TenantParentChangedEventModel
    json.Unmarshal(msg.Data, &event)
    
    // ⚠️ CRITICAL: Update ALL entities belonging to this tenant
    yourRepo.UpdateMany(
        bson.M{"tenantID": event.TenantID},
        bson.M{
            "$set": bson.M{
                "tenantPath":      event.NewPath,
                "tenantAncestors": event.NewAncestors,
                "tenantLevel":     event.NewLevel,
            },
        },
    )
    
    // Update cache
    tenantCache.Update(event.TenantID, TenantCacheEntry{
        Path:      event.NewPath,
        Ancestors: event.NewAncestors,
        Level:     event.NewLevel,
    })
    
    logger.Info("Tenant parent changed - updated all entities",
        "tenantID", event.TenantID,
        "oldPath", event.OldPath,
        "newPath", event.NewPath)
}
```

---

## 🔍 **Common Query Patterns**

### **1. Find All Entities in Subtree**

```go
// Get all tickets for tenant + descendants
func GetTicketsForHierarchy(tenantID string) ([]*Ticket, error) {
    tenant := tenantCache.Get(tenantID)
    
    filter := bson.M{
        "tenantPath": bson.M{
            "$regex": fmt.Sprintf("^%s", tenant.Path),
        },
    }
    
    return ticketRepo.Find(filter)
}
```

---

### **2. Check if Entity Belongs to Hierarchy**

```go
// Check if user belongs to tenant's hierarchy
func IsUserInHierarchy(userID, tenantID string) bool {
    user := userRepo.GetByID(userID)
    
    // Check if tenant is in user's path
    return strings.HasPrefix(user.TenantPath, tenant.Path) ||
           contains(user.TenantAncestors, tenantID) ||
           user.TenantID == tenantID
}
```

---

### **3. Aggregate Counts by Hierarchy**

```go
// Count tickets per child tenant
func GetTicketCountsByChild(parentTenantID string) (map[string]int, error) {
    parent := tenantCache.Get(parentTenantID)
    
    pipeline := []bson.M{
        {
            "$match": bson.M{
                "tenantPath": bson.M{
                    "$regex": fmt.Sprintf("^%s/", parent.Path),
                },
            },
        },
        {
            "$group": bson.M{
                "_id":   "$tenantID",
                "count": bson.M{"$sum": 1},
            },
        },
    }
    
    return ticketRepo.Aggregate(pipeline)
}
```

---

### **4. Find Direct Children Only**

```go
// Get users from direct children only (not grandchildren)
func GetUsersFromDirectChildren(parentTenantID string) ([]*User, error) {
    parent := tenantCache.Get(parentTenantID)
    
    filter := bson.M{
        "tenantLevel": parent.Level + 1, // Only next level
        "tenantAncestors": bson.M{
            "$elemMatch": bson.M{
                "$eq": parentTenantID,
            },
        },
    }
    
    return userRepo.Find(filter)
}
```

---

## 📝 **Event Subject Names**

Subscribe to these NATS subjects:

```go
const (
    // Core Events
    TenantCreated  = "tenant.created"
    TenantUpdated  = "tenant.updated"
    TenantDeleted  = "tenant.deleted"
    
    // Hierarchy Events
    TenantParentChanged      = "tenant.parent.changed"
    TenantChildAdded         = "tenant.child.added"
    TenantChildRemoved       = "tenant.child.removed"
    TenantHierarchyRecomputed = "tenant.hierarchy.recomputed"
)
```

---

## ⚠️ **Important Considerations**

### **1. Event Ordering**
Events may arrive out of order. Handle idempotently:
```go
if event.Version <= cached.Version {
    logger.Warn("Stale event ignored", "eventVersion", event.Version, "cachedVersion", cached.Version)
    return
}
```

### **2. Eventual Consistency**
There may be a brief delay between tenant-service updating and your service receiving the event. Design for eventual consistency.

### **3. Cache Invalidation**
Invalidate cached hierarchy data when receiving:
- `tenant.parent.changed`
- `tenant.hierarchy.recomputed`
- `tenant.updated` (if path/level changed)

### **4. Bulk Updates**
When `tenant.parent.changed` is received, you may need to update thousands of records. Use bulk operations:
```go
bulkWrite := []mongo.WriteModel{}
for _, userID := range affectedUsers {
    bulkWrite = append(bulkWrite, mongo.NewUpdateOneModel().
        SetFilter(bson.M{"_id": userID}).
        SetUpdate(bson.M{"$set": bson.M{
            "tenantPath": newPath,
            "tenantAncestors": newAncestors,
        }}))
}
collection.BulkWrite(context.Background(), bulkWrite)
```

---

## ✅ **Checklist for Integration**

- [ ] Add `TenantPath`, `TenantAncestors`, `TenantLevel` to your entity models
- [ ] Create database indexes on hierarchy fields
- [ ] Subscribe to all 7 tenant event types
- [ ] Implement event handlers to update local data
- [ ] Create cache for tenant hierarchy metadata
- [ ] Update queries to use path/ancestors for filtering
- [ ] Test with parent change scenarios
- [ ] Test with multi-level hierarchies (3+ levels)
- [ ] Handle event ordering and idempotency
- [ ] Implement bulk update for `parent.changed` events

---

## 🎓 **Example: User Service Complete Integration**

See full example in: `docs/EXAMPLES_USER_SERVICE_HIERARCHY.md`

---

## 🆘 **Need Help?**

- **Architecture Questions**: See `ARCHITECTURE_MULTI_TENANT_HIERARCHY.md`
- **Query Examples**: See `HIERARCHY_QUICK_REFERENCE.md`
- **Implementation Guide**: See `IMPLEMENTATION_GUIDE_HIERARCHY.md`

---

## 📌 **Summary**

**Key Takeaways:**
1. ✅ Store `path`, `ancestors`, `level` in your service's database
2. ✅ Subscribe to all 7 tenant event types (especially `parent.changed`)
3. ✅ Use path regex for subtree queries
4. ✅ Use ancestors array for bottom-up access control
5. ✅ Cache tenant hierarchy data for performance
6. ✅ Handle events idempotently (check version)
7. ✅ Use bulk updates when hierarchy changes

**This enables:**
- 🎯 Efficient hierarchical data queries
- 🎯 Automatic data aggregation up the tree
- 🎯 Proper access control across hierarchy
- 🎯 Scalable multi-tenant architecture

