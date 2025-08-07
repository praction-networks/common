# I9 Networks - Custom Authentication Architecture

## 🏗️ **Architecture Overview**

### **Simplified 2-System Architecture**
```
Request → Custom Auth Service (Auth + Business) → Casbin (Authz) → Response
```

### **Service Responsibilities**

#### **🎯 Custom Auth Service**
- **Authentication**: JWT token generation/validation, password hashing, session management
- **User Management**: Extended user profiles with multi-assignment system
- **Business Logic**: Multi-zone/department assignments, context-aware permissions
- **Event Publishing**: Real-time events for other services
- **Security**: Rate limiting, account lockout, audit logging

#### **🎯 Casbin**
- **Authorization**: Policy-based access control
- **Permission Enforcement**: Context-aware permission validation
- **Role Management**: Dynamic role and permission management

## 🎯 **Multi-Zone/Department Assignment System**

### **User Access Structure**
```typescript
// User with multiple zone/department assignments
const user = {
  id: "user-123",
  firstName: "John",
  lastName: "Doe",
  email: "john@company.com",
  password: "hashed_password",
  userAccess: [
    {
      zone: "ZoneA",
      department: "FieldTeam",
      role: "field-engineer",
      isPrimary: true,
      isActive: true
    },
    {
      zone: "ZoneB",
      department: "FieldTeam", 
      role: "field-sales",
      isPrimary: false,
      isActive: true
    },
    {
      zone: "ZoneA",
      department: "SupportTeam",
      role: "support-engineer",
      isPrimary: false,
      isActive: true
    }
  ],
  extendedProfile: {
    mobile: "+1234567890",
    whatsapp: "+1234567890",
    gender: "male",
    dob: "1990-01-01",
    permanentAddress: {...},
    currentAddress: {...}
  },
  isActive: true,
  createdAt: "2025-01-01T00:00:00Z",
  updatedAt: "2025-01-01T00:00:00Z"
}
```

### **JWT Token with Multi-Tenant Context**
```json
{
  "sub": "user-123",
  "name": "John Doe",
  "email": "john@company.com",
  "iat": 1640908800,
  "exp": 1640995200,
  
  // Multi-tenant context
  "currentZone": "ZoneA",
  "currentDepartment": "FieldTeam",
  "currentRole": "field-engineer",
  "isPrimary": true,
  "allAssignments": [
    {
      "zone": "ZoneA",
      "department": "FieldTeam",
      "role": "field-engineer",
      "isPrimary": true,
      "isActive": true
    },
    {
      "zone": "ZoneB",
      "department": "FieldTeam",
      "role": "field-sales",
      "isPrimary": false,
      "isActive": true
    }
  ],
  
  // Permissions (calculated from Casbin)
  "permissions": ["user:read", "field:manage", "zone:access"],
  
  // Business context
  "context": {
    "timeRestriction": "9:00-18:00",
    "locationRestriction": "ZoneA",
    "deviceRestriction": "company-devices"
  }
}
```

## 🔄 **Authentication Flow**

### **User Login Flow**
```typescript
1. User submits credentials
POST /api/auth/login
{
  "email": "john@company.com",
  "password": "password123"
}

2. Custom Auth Service validates credentials
3. Generates JWT with multi-tenant context
4. Returns enhanced token with all assignments

Response:
{
  "user": {...},
  "token": "jwt_with_multi_tenant_context",
  "refreshToken": "refresh_token",
  "expiresIn": 3600
}
```

### **Permission Check Flow**
```typescript
1. User makes request with JWT
GET /api/zones/ZoneA/departments/FieldTeam/users
Headers: Authorization: Bearer jwt_token

2. Custom Auth Service validates JWT
3. Extracts current zone/department context
4. Checks Casbin permissions
5. Applies business rules
6. Returns response or error
```

### **Context Switching Flow**
```typescript
1. User switches to different zone/department
POST /api/auth/context/switch
{
  "zone": "ZoneB",
  "department": "FieldTeam"
}

2. Custom Auth Service validates assignment
3. Generates new JWT with updated context
4. Returns new token

Response:
{
  "token": "new_jwt_with_updated_context",
  "expiresIn": 3600
}
```

## 🎯 **API Endpoints**

### **Authentication Endpoints**
```typescript
POST /api/auth/register     // Create user with multi-assignments
POST /api/auth/login        // Login with credentials
POST /api/auth/logout       // Logout and invalidate token
POST /api/auth/refresh      // Refresh access token
POST /api/auth/context/switch // Switch zone/department context
```

### **User Management Endpoints**
```typescript
GET    /api/users           // Get all users
GET    /api/users/{id}      // Get user by ID
PUT    /api/users/{id}      // Update user profile
DELETE /api/users/{id}      // Delete user

// Assignment Management
GET    /api/users/{id}/assignments           // Get user assignments
POST   /api/users/{id}/assignments          // Add new assignment
PUT    /api/users/{id}/assignments/{assId}  // Update assignment
DELETE /api/users/{id}/assignments/{assId}  // Remove assignment
```

### **Permission Endpoints**
```typescript
GET /api/permissions/user/{id}              // Get user permissions
POST /api/permissions/check                 // Check specific permission
GET /api/permissions/context/{zone}/{dept}  // Get context permissions
```

## 🏗️ **Implementation Plan**

### **Phase 1: Core Authentication (Week 1)**
**Goal**: Build basic authentication system

**Tasks**:
- [ ] **JWT Token System**: Generate/validate tokens with multi-tenant context
- [ ] **Password Authentication**: bcrypt hashing, password policies
- [ ] **User CRUD Operations**: Create, read, update, delete users
- [ ] **Session Management**: Redis-based session storage
- [ ] **Basic Security**: Rate limiting, account lockout

**Deliverables**:
- ✅ Custom Auth Service with JWT support
- ✅ User registration and login
- ✅ Basic security features
- ✅ Session management

### **Phase 2: Multi-Tenant Features (Week 2)**
**Goal**: Implement multi-assignment system

**Tasks**:
- [ ] **Multi-Assignment System**: Zone/department/role assignments
- [ ] **Context-Aware Permissions**: Zone/department specific access
- [ ] **Business Rule Engine**: Time, location, device restrictions
- [ ] **Assignment Management**: Add, update, remove assignments
- [ ] **Context Switching**: Switch between zone/department contexts

**Deliverables**:
- ✅ Multi-assignment user system
- ✅ Context-aware permission calculation
- ✅ Business rule enforcement
- ✅ Assignment management APIs

### **Phase 3: Casbin Integration (Week 3)**
**Goal**: Integrate Casbin for authorization

**Tasks**:
- [ ] **Casbin Setup**: Policy-based access control
- [ ] **Permission Calculation**: Dynamic permission aggregation
- [ ] **Context-Aware Policies**: Zone/department specific policies
- [ ] **Role Management**: Dynamic role and permission management
- [ ] **Policy Enforcement**: Real-time permission checking

**Deliverables**:
- ✅ Casbin integration
- ✅ Context-aware policies
- ✅ Dynamic permission calculation
- ✅ Real-time authorization

### **Phase 4: Advanced Features (Week 4)**
**Goal**: Add advanced security and business features

**Tasks**:
- [ ] **Multi-Factor Authentication**: OTP, WebAuthn, FaceID
- [ ] **Audit Logging**: Comprehensive activity tracking
- [ ] **Event Publishing**: Real-time events for other services
- [ ] **Performance Optimization**: Caching, query optimization
- [ ] **Security Enhancements**: Advanced security features

**Deliverables**:
- ✅ Multi-factor authentication
- ✅ Comprehensive audit logging
- ✅ Event-driven architecture
- ✅ Performance optimization

## 🎯 **Key Benefits**

### **✅ Complete Control**
- ✅ **JWT Structure**: Custom JWT with multi-tenant context
- ✅ **Password Policies**: Custom password policies and hashing
- ✅ **Session Management**: Custom session management
- ✅ **Business Rules**: Custom business rule enforcement
- ✅ **Event Publishing**: Custom event publishing system

### **✅ Simplified Architecture**
- ✅ **Single Auth Service**: One system instead of multiple
- ✅ **Direct Casbin Integration**: No integration complexity
- ✅ **Unified JWT Tokens**: Single token with all context
- ✅ **Single Source of Truth**: Everything in one place
- ✅ **No Cross-Service Calls**: Faster auth flow

### **✅ Better Performance**
- ✅ **No External Calls**: Direct database access
- ✅ **Optimized JWT Validation**: Fast token validation
- ✅ **Cached Permissions**: Permission caching
- ✅ **In-Memory Sessions**: Fast session storage
- ✅ **Event-Driven**: Real-time updates

### **✅ Multi-Assignment Support**
- ✅ **Multiple Zones**: User can work in ZoneA, ZoneB, ZoneC
- ✅ **Multiple Departments**: User can be in FieldTeam, SalesTeam, SupportTeam
- ✅ **Multiple Roles**: User can have different roles per zone/department
- ✅ **Primary Assignment**: One primary assignment for default context
- ✅ **Active/Inactive**: Enable/disable specific assignments

This architecture provides complete control over authentication and authorization while maintaining simplicity and performance! 🎉
