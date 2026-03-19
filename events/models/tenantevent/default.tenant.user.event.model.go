package tenantevent

// DefaultTenantUserEventModel is published by tenant-service when an ISP/Enterprise
// tenant is created with default user details. The tenant-user-service consumes this
// event to auto-create the first admin user with the TenantAdmin role.
type DefaultTenantUserEventModel struct {
	TenantID  string `bson:"tenantId" json:"tenantId"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Email     string `bson:"email" json:"email"`
	Mobile    string `bson:"mobile" json:"mobile"`
	Gender    string `bson:"gender" json:"gender"`
	DOB       string `bson:"dob" json:"dob"`
}
