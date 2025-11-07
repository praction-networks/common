# Tenant Event Models - Hierarchy Support Changelog

## 📅 **Date**: November 6, 2025
## 🔖 **Version**: 2.0.0 (Breaking Change - adds new fields)

---

## ✅ **Changes Summary**

### **1. Updated Existing Event Models**

#### **TenantInsertEventModel** - Added 9 Hierarchy Fields
```go
// NEW FIELDS ADDED:
Path              string   // ✅ "/isp_abc/reseller_xyz/dist_123"
PathDepth         int      // ✅ 3
Ancestors         []string // ✅ ["isp_abc", "reseller_xyz"]
Level             int      // ✅ 2
IsLeaf            bool     // ✅ true (no children)
ChildrenCount     int      // ✅ 0
ChildIDs          []string // ✅ []
AllowedChildTypes []string // ✅ ["Operator"]
MaxDepth          int      // ✅ 1
```

#### **TenantUpdateEventModel** - Added 8 Hierarchy Fields
```go
// NEW FIELDS ADDED (same as Insert, minus ChildIDs which was already there):
Path              string
PathDepth         int
Ancestors         []string
Level             int
IsLeaf            bool
ChildrenCount     int
AllowedChildTypes []string
MaxDepth          int
```

#### **TenantDeleteEventModel** - Enhanced with Hierarchy Context
```go
// NEW FIELDS ADDED:
ParentTenantID string   // ✅ Parent needs to decrement childrenCount
Path           string   // ✅ For cascade delete validation
Ancestors      []string // ✅ Notify all ancestor services
```

---

### **2. New Event Types for Hierarchy Operations**

#### **A. TenantParentChangedEventModel** ⭐
**When Published:** Parent relationship changes (reassignment)

**Fields:**
```go
TenantID       string    // Which tenant moved
OldParentID    string    // Previous parent
NewParentID    string    // New parent
OldPath        string    // Previous path
NewPath        string    // New path
OldAncestors   []string  // Previous ancestors
NewAncestors   []string  // New ancestors
OldLevel       int       // Previous level
NewLevel       int       // New level
ChangedAt      time.Time // Timestamp
ChangedBy      string    // User ID
```

**Consumer Action Required:**
- Update all entities belonging to this tenant with new path/ancestors/level
- Update cached tenant hierarchy data
- Recalculate aggregations if needed

**NATS Subject:** `tenant.parent.changed`

---

#### **B. TenantChildAddedEventModel**
**When Published:** Child tenant is added to parent

**Fields:**
```go
ParentID      string    // Parent tenant
ChildID       string    // New child
ChildType     string    // "Reseller", "Distributor", etc.
ChildPath     string    // Child's path
ParentPath    string    // Parent's path
NewChildCount int       // Parent's updated child count
AddedAt       time.Time
AddedBy       string
```

**Consumer Action Required:**
- Grant parent admins access to child's data
- Include child in parent's aggregate metrics
- Update parent's child count cache

**NATS Subject:** `tenant.child.added`

---

#### **C. TenantChildRemovedEventModel**
**When Published:** Child tenant is removed (deleted or reassigned)

**Fields:**
```go
ParentID      string    // Parent tenant
ChildID       string    // Removed child
ChildType     string    // Type of removed child
NewChildCount int       // Parent's updated child count
IsNowLeaf     bool      // True if parent has no children now
RemovedAt     time.Time
RemovedBy     string
```

**Consumer Action Required:**
- Revoke parent admins access to removed child's data
- Exclude child from parent's aggregate metrics
- Update parent's child count cache

**NATS Subject:** `tenant.child.removed`

---

#### **D. TenantHierarchyRecomputedEventModel**
**When Published:** Hierarchy fields recalculated (migrations, fixes)

**Fields:**
```go
TenantID     string    // Affected tenant
Path         string    // New computed path
PathDepth    int       // New depth
Ancestors    []string  // New ancestors
Level        int       // New level
RecomputedAt time.Time
Reason       string    // "parent_changed", "migration", "manual_fix"
```

**Consumer Action Required:**
- Sync hierarchy fields in local cache
- Update entities if path/ancestors/level changed

**NATS Subject:** `tenant.hierarchy.recomputed`

---

## 🎯 **Migration Guide for Existing Services**

### **Phase 1: Backwards-Compatible (Week 1)**
1. ✅ Add hierarchy fields to your models (mark as optional)
2. ✅ Subscribe to new event types (log but don't process yet)
3. ✅ Start populating hierarchy fields for NEW entities
4. ✅ Deploy to production

### **Phase 2: Backfill Data (Week 2)**
1. ✅ Write migration script to populate hierarchy fields for existing entities
2. ✅ Run migration in staging
3. ✅ Verify data correctness
4. ✅ Run migration in production

### **Phase 3: Start Using (Week 3)**
1. ✅ Update queries to use hierarchy fields
2. ✅ Enable event processing
3. ✅ Add hierarchy-based features (subtree queries, aggregations)
4. ✅ Monitor performance

### **Phase 4: Enforce (Week 4)**
1. ✅ Make hierarchy fields required
2. ✅ Remove old non-hierarchical queries
3. ✅ Full rollout

---

## 📊 **Example Scenarios**

### **Scenario 1: ISP Creates Reseller**

**Event Flow:**
```
1. tenant-service creates reseller
   → Publishes: tenant.created (with path="/isp_rapidnet/reseller_mumbai")
   
2. tenant-service adds child to ISP
   → Publishes: tenant.child.added (parentID="isp_rapidnet", childID="reseller_mumbai")
   
3. user-service receives events:
   - Caches reseller metadata
   - Grants ISP admin access to reseller's users
   
4. ticket-service receives events:
   - Caches reseller metadata
   - Enables ISP support team to view reseller's tickets
```

---

### **Scenario 2: Distributor Moves to New Reseller**

**Event Flow:**
```
1. tenant-service validates new parent
2. tenant-service recomputes path/ancestors/level
3. tenant-service publishes:
   a. tenant.child.removed (oldParent, distributorID)
   b. tenant.parent.changed (distributorID, oldPath, newPath)
   c. tenant.child.added (newParent, distributorID)
   
4. user-service receives:
   - Updates all distributor's users with new path/ancestors
   - Revokes old parent's access
   - Grants new parent's access
   
5. ticket-service receives:
   - Updates all distributor's tickets with new path
   - Reroutes escalations to new parent
```

---

## 🔧 **Field-by-Field Usage**

| Field | Purpose | Query Example |
|-------|---------|---------------|
| `Path` | Subtree queries | `tenantPath: /^\/isp_abc/` |
| `PathDepth` | Limit depth | `pathDepth: {$lte: 3}` |
| `Ancestors` | Bottom-up access | `tenantAncestors: "isp_abc"` |
| `Level` | Group by level | `$group: {_id: "$tenantLevel"}` |
| `IsLeaf` | Filter leaf nodes | `isLeaf: true` |
| `ChildrenCount` | Show count | Display in UI |
| `ChildIDs` | Direct children | `tenantID: {$in: childIDs}` |
| `AllowedChildTypes` | Validation | Check before operations |
| `MaxDepth` | Depth limits | Enforce in UI |

---

## 📚 **Related Documentation**

- **Integration Guide**: `HIERARCHY_INTEGRATION_GUIDE.md` (this directory)
- **Architecture**: `/docs/ARCHITECTURE_MULTI_TENANT_HIERARCHY.md`
- **Quick Reference**: `/docs/HIERARCHY_QUICK_REFERENCE.md`
- **Implementation**: `/docs/IMPLEMENTATION_GUIDE_HIERARCHY.md`

---

## 🚨 **Breaking Changes**

### **For Consumers of Tenant Events:**

**Before (v1.x):**
```go
type User struct {
    TenantID string
}

// Query all users in hierarchy - required API call to tenant-service
usersInHierarchy = getUsersForTenantAndDescendants(tenantID)
```

**After (v2.0):**
```go
type User struct {
    TenantID        string
    TenantPath      string   // ✅ NEW - enables local queries
    TenantAncestors []string // ✅ NEW - enables access control
    TenantLevel     int      // ✅ NEW - enables analytics
}

// Query all users in hierarchy - NO API calls needed!
users = db.users.find({tenantPath: /^\/isp_abc/})
```

---

## ✅ **Compatibility**

- **Backwards Compatible**: New fields are `omitempty` - old consumers won't break
- **Forward Compatible**: Services can gradually adopt hierarchy features
- **Versioned**: Use `Version` field to handle concurrent updates

---

## 🎉 **What This Enables**

### **For User Service:**
- ✅ Show all users in ISP hierarchy (including resellers, distributors)
- ✅ Count users per level (reseller vs distributor)
- ✅ Grant parent admins access to child users
- ✅ Aggregate user counts up the hierarchy

### **For Ticket Service:**
- ✅ Show all tickets in ISP hierarchy
- ✅ Route tickets based on hierarchy level
- ✅ Escalate to parent support team
- ✅ SLA based on tenant level

### **For Helpdesk Service:**
- ✅ Assign agents based on hierarchy
- ✅ Dashboard with hierarchy breakdown
- ✅ Knowledge base inheritance

### **For Analytics/Reporting:**
- ✅ Revenue by hierarchy level
- ✅ User growth per reseller
- ✅ Ticket volume by distributor
- ✅ Drill-down from ISP → Reseller → Distributor

---

## 🔮 **Future Enhancements**

Possible future additions:
- `TenantHierarchyMigratedEventModel` - Bulk hierarchy migrations
- `TenantDescendantsUpdatedEventModel` - When entire subtree changes
- `TenantAggregatesUpdatedEventModel` - When aggregate counts change
- Aggregate collection support (separate from events)

---

## 📞 **Contact**

For questions or issues:
- **Architecture**: See `/docs/ARCHITECTURE_MULTI_TENANT_HIERARCHY.md`
- **Implementation**: See `/docs/IMPLEMENTATION_GUIDE_HIERARCHY.md`
- **Examples**: See `/docs/HIERARCHY_QUICK_REFERENCE.md`

