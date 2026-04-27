package tenantevent

import "time"

// OLTCLIConfig carries the credentials olt-manager needs to drive the
// vendor CLI (Huawei MA5800, ZTE C320, Nokia FX, etc.) for ONT
// provisioning, board/port queries, and command execution.
//
// Either Password or PrivateKey must be present — enforced at the
// tenant-service request-schema layer, not here, since the event model
// is the wire shape and may legitimately omit the unused alternative.
type OLTCLIConfig struct {
	Protocol         string `bson:"protocol" json:"protocol"`                                 // ssh | telnet
	Port             int    `bson:"port" json:"port"`                                         // default 22 (ssh) / 23 (telnet)
	Username         string `bson:"username" json:"username"`
	Password         string `bson:"password,omitempty" json:"password,omitempty"`
	PrivateKey       string `bson:"privateKey,omitempty" json:"privateKey,omitempty"`         // SSH key alt to password
	PrivateKeyPhrase string `bson:"privateKeyPhrase,omitempty" json:"privateKeyPhrase,omitempty"`
	EnablePassword   string `bson:"enablePassword,omitempty" json:"enablePassword,omitempty"` // vendor-specific (super/enable mode)
}

// OLTSNMPConfig carries the credentials olt-manager needs for SNMP
// polling — link state, board/port stats, ONT counters. Field set
// switches on Version: v1/v2c → ReadCommunity (+ optional WriteCommunity);
// v3 → username + auth/priv.
//
// SNMP v1/v2c distinguishes read (GET/GETBULK — high frequency status
// polling) from write (SET — config push). Many ISPs use a separate
// write community to reduce the blast radius of credential leaks; if
// the operator only configures one community, leave WriteCommunity
// blank and olt-manager will fall back to ReadCommunity for SET ops.
type OLTSNMPConfig struct {
	Version        string `bson:"version" json:"version"`                                   // v1 | v2c | v3
	Port           int    `bson:"port" json:"port"`                                         // default 161
	ReadCommunity  string `bson:"readCommunity,omitempty" json:"readCommunity,omitempty"`   // v1/v2c — required
	WriteCommunity string `bson:"writeCommunity,omitempty" json:"writeCommunity,omitempty"` // v1/v2c — optional; falls back to readCommunity
	Username       string `bson:"username,omitempty" json:"username,omitempty"`             // v3
	SecurityLevel  string `bson:"securityLevel,omitempty" json:"securityLevel,omitempty"`   // v3: noAuthNoPriv | authNoPriv | authPriv
	AuthProtocol   string `bson:"authProtocol,omitempty" json:"authProtocol,omitempty"`     // v3: MD5 | SHA | SHA224 | SHA256 | SHA384 | SHA512
	AuthPassword   string `bson:"authPassword,omitempty" json:"authPassword,omitempty"`
	PrivProtocol   string `bson:"privProtocol,omitempty" json:"privProtocol,omitempty"`     // v3: DES | 3DES | AES128 | AES192 | AES256
	PrivPassword   string `bson:"privPassword,omitempty" json:"privPassword,omitempty"`
}

// OLTInsertEventModel is published on TenantStream / olt.created when
// tenant-service creates a new OLT. olt-manager projects the full
// document into its local oltDevices collection, keyed by ID.
type OLTInsertEventModel struct {
	ID           string   `bson:"_id" json:"id"`
	TenantIDs    []string `bson:"tenantIds" json:"tenantIds"`
	AssignableBy []string `bson:"assignableBy,omitempty" json:"assignableBy,omitempty"`

	Name            string `bson:"name" json:"name"`
	Description     string `bson:"description,omitempty" json:"description,omitempty"`
	Vendor          string `bson:"vendor" json:"vendor"` // Huawei | ZTE | Nokia | Generic
	Model           string `bson:"model" json:"model"`
	FirmwareVersion string `bson:"firmwareVersion,omitempty" json:"firmwareVersion,omitempty"`
	SerialNumber    string `bson:"serialNumber,omitempty" json:"serialNumber,omitempty"`

	MgmtIP string `bson:"mgmtIp" json:"mgmtIp"`

	CLI  OLTCLIConfig  `bson:"cli" json:"cli"`
	SNMP OLTSNMPConfig `bson:"snmp" json:"snmp"`

	Tags     []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Status   string   `bson:"status" json:"status"` // provisioned | active | suspended
	IsActive bool     `bson:"isActive" json:"isActive"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	CreatedBy string    `bson:"createdBy" json:"createdBy"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy string    `bson:"updatedBy" json:"updatedBy"`
	Version   int       `bson:"version" json:"version"`
}

// OLTUpdateEventModel mirrors the device update model: pointer types
// for scalars where the absence of a field must be distinguishable from
// the zero value (IsActive true/false/unset, Port 0/unset).
type OLTUpdateEventModel struct {
	ID           string   `bson:"_id" json:"id"`
	TenantIDs    []string `bson:"tenantIds,omitempty" json:"tenantIds,omitempty"`
	AssignableBy []string `bson:"assignableBy,omitempty" json:"assignableBy,omitempty"`

	Name            string `bson:"name,omitempty" json:"name,omitempty"`
	Description     string `bson:"description,omitempty" json:"description,omitempty"`
	Vendor          string `bson:"vendor,omitempty" json:"vendor,omitempty"`
	Model           string `bson:"model,omitempty" json:"model,omitempty"`
	FirmwareVersion string `bson:"firmwareVersion,omitempty" json:"firmwareVersion,omitempty"`
	SerialNumber    string `bson:"serialNumber,omitempty" json:"serialNumber,omitempty"`

	MgmtIP string `bson:"mgmtIp,omitempty" json:"mgmtIp,omitempty"`

	CLI  *OLTCLIConfig  `bson:"cli,omitempty" json:"cli,omitempty"`
	SNMP *OLTSNMPConfig `bson:"snmp,omitempty" json:"snmp,omitempty"`

	Tags     []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Status   string   `bson:"status,omitempty" json:"status,omitempty"`
	IsActive *bool    `bson:"isActive,omitempty" json:"isActive,omitempty"`

	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy string    `bson:"updatedBy" json:"updatedBy"`
	Version   int       `bson:"version" json:"version"`
}

// OLTDeleteEventModel — minimal payload, mirrors DeviceDeleteEventModel.
type OLTDeleteEventModel struct {
	ID        string    `bson:"_id" json:"id"`
	TenantIDs []string  `bson:"tenantIds,omitempty" json:"tenantIds,omitempty"`
	Version   int       `bson:"version" json:"version"`
	DeletedAt time.Time `bson:"deletedAt" json:"deletedAt"`
	DeletedBy string    `bson:"deletedBy,omitempty" json:"deletedBy,omitempty"`
}
