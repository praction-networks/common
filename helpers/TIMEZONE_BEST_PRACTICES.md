# Multi-Tenant Timezone Management - Best Practices

## Overview
This document explains how multi-tenant and global systems handle datetime across different timezones.

---

## Core Principles

### 1. **Store Everything in UTC**
```go
// ✅ CORRECT - Always store in UTC
createdAt := time.Now().UTC()
subscriber.CreatedAt = createdAt

// ❌ WRONG - Never store local time
createdAt := time.Now() // Depends on server timezone!
```

**Why UTC?**
- Universal reference point
- No ambiguity during daylight saving time transitions
- Consistent database queries and sorting
- Easy conversion to any timezone

### 2. **Convert to User Timezone Only for Display**
```go
// Store in database (UTC)
subscriber.CreatedAt = time.Now().UTC()

// Display to user (convert to their timezone)
userTimezone := tenant.Timezone // e.g., "Asia/Kolkata"
displayTime, _ := timeProvider.FormatTimeForDisplay(subscriber.CreatedAt, userTimezone)
// Output: "2024-12-25 20:30:45 IST"
```

### 3. **Store Timezone at Tenant/User Level**
```go
type Tenant struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Timezone     string    `json:"timezone"`     // "Asia/Kolkata"
    // ... other fields
}

type User struct {
    ID           string    `json:"id"`
    TenantID     string    `json:"tenantId"`
    Timezone     string    `json:"timezone,omitempty"` // Optional: user preference
    // ... other fields
}
```

---

## Multi-Tenant Architecture Patterns

### Pattern 1: Tenant-Level Timezone (Recommended)
Each tenant has a default timezone. All users under that tenant see times in the tenant's timezone unless they have a personal preference.

```go
func (s *SubscriberService) GetSubscriber(ctx context.Context, id string) (*schema.SubscriberResponse, error) {
    subscriber, err := s.repo.GetSubscriberByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Get tenant timezone
    tenant, _ := s.tenantRepo.GetTenantByID(ctx, subscriber.TenantID)
    tenantTimezone := tenant.Timezone // e.g., "Asia/Kolkata"
    
    // Convert all timestamps for response
    createdAtDisplay, _ := helpers.FormatTimeForDisplay(subscriber.CreatedAt, tenantTimezone)
    
    return &schema.SubscriberResponse{
        ID:        subscriber.ID,
        CreatedAt: createdAtDisplay,
        // ... other fields
    }, nil
}
```

### Pattern 2: User-Level Timezone Override
Users can override the tenant timezone with their personal preference.

```go
func GetEffectiveTimezone(user *models.User, tenant *models.Tenant) string {
    // User preference takes priority
    if user.Timezone != "" && helpers.ValidateTimezone(user.Timezone) {
        return user.Timezone
    }
    
    // Fallback to tenant timezone
    if tenant.Timezone != "" && helpers.ValidateTimezone(tenant.Timezone) {
        return tenant.Timezone
    }
    
    // Default fallback
    return helpers.TimezoneUTC
}
```

---

## Database Design

### Example Schema with Timezone Support

```sql
-- Tenants table
CREATE TABLE tenants (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',  -- IANA timezone
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),    -- Stored in UTC
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()     -- Stored in UTC
);

-- Users table
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    timezone VARCHAR(100),                          -- Optional user preference
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),    -- Always UTC
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Business records (e.g., subscribers, sessions)
CREATE TABLE subscribers (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),    -- Always UTC
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),    -- Always UTC
    valid_from TIMESTAMP,                           -- Always UTC
    valid_until TIMESTAMP,                          -- Always UTC
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
```

### MongoDB Example

```go
type Tenant struct {
    ID        string    `bson:"_id" json:"id"`
    Name      string    `bson:"name" json:"name"`
    Timezone  string    `bson:"timezone" json:"timezone"`  // "Asia/Kolkata"
    CreatedAt time.Time `bson:"createdAt" json:"createdAt"` // Stored as UTC
    UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"` // Stored as UTC
}

type Subscriber struct {
    ID         string    `bson:"_id" json:"id"`
    TenantID   string    `bson:"tenantId" json:"tenantId"`
    FullName   string    `bson:"fullName" json:"fullName"`
    CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`  // Always UTC
    ValidFrom  time.Time `bson:"validFrom" json:"validFrom"`  // Always UTC
    ValidUntil time.Time `bson:"validUntil" json:"validUntil"` // Always UTC
}
```

---

## API Design

### Request: Accept Timezone-Aware Input

```go
type CreateSubscriberRequest struct {
    FullName   string    `json:"fullName"`
    ValidFrom  string    `json:"validFrom"`  // ISO 8601 with timezone: "2024-12-25T10:30:00+05:30"
    ValidUntil string    `json:"validUntil"` // ISO 8601 with timezone
}

func (s *SubscriberService) CreateSubscriber(ctx context.Context, req *CreateSubscriberRequest) error {
    // Parse ISO 8601 timestamp (automatically handles timezone)
    validFrom, err := time.Parse(time.RFC3339, req.ValidFrom)
    if err != nil {
        return fmt.Errorf("invalid validFrom format: %w", err)
    }
    
    // Convert to UTC before storing
    validFromUTC := validFrom.UTC()
    
    subscriber := &models.Subscriber{
        FullName:   req.FullName,
        ValidFrom:  validFromUTC,  // ✅ Stored in UTC
        CreatedAt:  time.Now().UTC(),
    }
    
    return s.repo.Create(ctx, subscriber)
}
```

### Response: Convert to User's Timezone

```go
type SubscriberResponse struct {
    ID         string `json:"id"`
    FullName   string `json:"fullName"`
    CreatedAt  string `json:"createdAt"`  // ISO 8601 in user's timezone
    ValidFrom  string `json:"validFrom"`  // ISO 8601 in user's timezone
    ValidUntil string `json:"validUntil"` // ISO 8601 in user's timezone
}

func (s *SubscriberService) GetSubscriber(ctx context.Context, id string, userTimezone string) (*SubscriberResponse, error) {
    subscriber, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Convert to user's timezone for display
    loc, _ := time.LoadLocation(userTimezone)
    
    return &SubscriberResponse{
        ID:         subscriber.ID,
        FullName:   subscriber.FullName,
        CreatedAt:  subscriber.CreatedAt.In(loc).Format(time.RFC3339),  // ISO 8601 with timezone
        ValidFrom:  subscriber.ValidFrom.In(loc).Format(time.RFC3339),
        ValidUntil: subscriber.ValidUntil.In(loc).Format(time.RFC3339),
    }, nil
}
```

---

## Common Scenarios

### 1. Scheduled Tasks / Cron Jobs
```go
// Business requirement: "Send report every day at 9:00 AM tenant time"

func ScheduleReportForTenant(tenant *models.Tenant) {
    // Get tenant's timezone
    loc, _ := time.LoadLocation(tenant.Timezone)
    
    // Calculate next 9:00 AM in tenant's timezone
    now := time.Now().In(loc)
    next9AM := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, loc)
    
    if next9AM.Before(now) {
        // If 9 AM already passed today, schedule for tomorrow
        next9AM = next9AM.Add(24 * time.Hour)
    }
    
    // Convert to UTC for job scheduler
    next9AMUTC := next9AM.UTC()
    
    // Schedule job
    scheduler.ScheduleAt(next9AMUTC, func() {
        SendReport(tenant.ID)
    })
}
```

### 2. Date Range Queries
```go
// User wants: "Show me all sessions from Dec 25 to Dec 26 in my timezone"

func GetSessionsByDateRange(ctx context.Context, tenantID string, startDate, endDate string, timezone string) ([]*models.Session, error) {
    // Parse dates in user's timezone
    loc, _ := time.LoadLocation(timezone)
    
    // Start of day in user's timezone
    startTime, _ := time.ParseInLocation("2006-01-02", startDate, loc)
    // End of day in user's timezone
    endTime, _ := time.ParseInLocation("2006-01-02", endDate, loc)
    endTime = endTime.Add(24 * time.Hour).Add(-1 * time.Second) // End of day: 23:59:59
    
    // Convert to UTC for database query
    startUTC := startTime.UTC()
    endUTC := endTime.UTC()
    
    // Query database (using UTC)
    return repo.FindSessionsByDateRange(ctx, tenantID, startUTC, endUTC)
}
```

### 3. Expiry/Validity Checks
```go
// Check if subscription is valid right now

func IsSubscriptionValid(subscriber *models.Subscriber) bool {
    now := time.Now().UTC()  // Current time in UTC
    
    // Compare UTC times
    return now.After(subscriber.ValidFrom) && now.Before(subscriber.ValidUntil)
}
```

---

## Frontend Integration

### Send Timezone from Frontend
```javascript
// JavaScript/TypeScript example

// Option 1: Send user's timezone
const userTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
// Output: "Asia/Kolkata", "America/New_York", etc.

// Option 2: Send dates in ISO 8601 with timezone
const date = new Date('2024-12-25T10:30:00');
const isoString = date.toISOString(); // "2024-12-25T10:30:00.000Z"

// Include in API request
fetch('/api/subscribers', {
    method: 'POST',
    headers: { 
        'Content-Type': 'application/json',
        'X-User-Timezone': userTimezone  // Custom header
    },
    body: JSON.stringify({
        fullName: 'John Doe',
        validFrom: isoString
    })
});
```

### Display Times in User's Timezone
```javascript
// Backend returns: "2024-12-25T15:00:00Z" (UTC)
// Frontend converts to user's timezone

const utcTime = "2024-12-25T15:00:00Z";
const date = new Date(utcTime);

// Automatically displays in user's browser timezone
console.log(date.toLocaleString()); 
// India: "25/12/2024, 8:30:00 PM"
// USA EST: "25/12/2024, 10:00:00 AM"
```

---

## Testing Considerations

### Test Across Timezones
```go
func TestSubscriberValidityAcrossTimezones(t *testing.T) {
    // Create subscriber valid from Dec 25, 2024 00:00 IST to Dec 26, 2024 23:59 IST
    istLoc, _ := time.LoadLocation("Asia/Kolkata")
    
    validFrom := time.Date(2024, 12, 25, 0, 0, 0, 0, istLoc).UTC()
    validUntil := time.Date(2024, 12, 26, 23, 59, 59, 0, istLoc).UTC()
    
    subscriber := &models.Subscriber{
        ValidFrom:  validFrom,
        ValidUntil: validUntil,
    }
    
    // Test: Dec 25, 2024 12:00 IST (should be valid)
    testTime := time.Date(2024, 12, 25, 12, 0, 0, 0, istLoc).UTC()
    assert.True(t, IsSubscriptionValid(subscriber, testTime))
    
    // Test: Dec 27, 2024 00:00 IST (should be invalid)
    testTime = time.Date(2024, 12, 27, 0, 0, 0, 0, istLoc).UTC()
    assert.False(t, IsSubscriptionValid(subscriber, testTime))
}
```

---

## Common Pitfalls to Avoid

### ❌ Don't: Use Server's Local Time
```go
// BAD - Depends on server location!
subscriber.CreatedAt = time.Now()
```

### ❌ Don't: Store Timezone Offset Instead of Name
```go
// BAD - Doesn't handle DST changes
type Tenant struct {
    TimezoneOffset int  // +5:30 for India
}

// GOOD - Use IANA timezone name
type Tenant struct {
    Timezone string  // "Asia/Kolkata"
}
```

### ❌ Don't: Convert Times in Database Queries
```go
// BAD - Complex, slow, error-prone
SELECT * FROM subscribers 
WHERE CONVERT_TZ(created_at, 'UTC', 'Asia/Kolkata') > '2024-12-25 00:00:00';

// GOOD - Query in UTC, convert in application
SELECT * FROM subscribers 
WHERE created_at > '2024-12-24 18:30:00';  // UTC equivalent
```

---

## Summary Checklist

✅ **Store all timestamps in UTC**  
✅ **Store tenant/user timezone preference as IANA name** (e.g., "Asia/Kolkata")  
✅ **Convert to user timezone only for display**  
✅ **Accept ISO 8601 timestamps in API requests**  
✅ **Use timezone-aware date parsing when needed**  
✅ **Test with multiple timezones including DST transitions**  
✅ **Handle timezone validation gracefully**  
✅ **Document timezone handling in API documentation**  

---

## Additional Resources

- [IANA Time Zone Database](https://www.iana.org/time-zones)
- [ISO 8601 DateTime Format](https://en.wikipedia.org/wiki/ISO_8601)
- [Go time package documentation](https://pkg.go.dev/time)

