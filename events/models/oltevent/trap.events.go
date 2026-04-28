// Package oltevent defines the wire shapes for real-time OLT runtime
// events published by olt-manager on OLTEventStream. Sources are
// primarily SNMP traps from the OLT itself (linkDown, dying-gasp, ONT
// offline notifications) plus a smaller set of derived events from
// active polling (e.g. signal-degradation crossovers).
//
// These are *runtime* events — high-velocity, time-sensitive,
// distinct from the lifecycle events on TenantStream
// (olt.created/updated/deleted). Consumers: notification-service,
// ticket-service, subscriber-service.
package oltevent

import "time"

// EventEnvelope is the common header carried by every trap-derived
// event. Embedded into the per-event-type structs below so consumers
// can rely on a uniform "where did this come from + when" surface.
type EventEnvelope struct {
	// EventID is a stable identifier for this event instance —
	// used by the dedup window in olt-manager and as the NATS
	// Msg-Id for at-most-once dedup at the stream level.
	EventID string `bson:"eventId" json:"eventId"`

	// OLT identity. OLTID matches the document ID in the
	// tenant-service OLT collection / olt-manager projection.
	OLTID     string   `bson:"oltId" json:"oltId"`
	TenantIDs []string `bson:"tenantIds" json:"tenantIds"`
	Vendor    string   `bson:"vendor" json:"vendor"` // canonical (Huawei | ZTE | …)
	OEM       string   `bson:"oem,omitempty" json:"oem,omitempty"`

	// Times. OccurredAt is the OLT-reported timestamp when
	// available (sysUpTime extraction); ReceivedAt is when
	// olt-manager processed the trap. They diverge by the time
	// it took the trap to traverse the network + the OLT's
	// clock skew.
	OccurredAt time.Time `bson:"occurredAt" json:"occurredAt"`
	ReceivedAt time.Time `bson:"receivedAt" json:"receivedAt"`

	// Source — for forensics and debugging. Operator can
	// correlate against vendor logs.
	SourceIP string `bson:"sourceIp" json:"sourceIp"`
	TrapOID  string `bson:"trapOid" json:"trapOid"`

	// FlapCount — incremented by the dedup window when the same
	// (source, trap, ont, cause) tuple repeats inside 30s. >1
	// means the underlying condition is flapping.
	FlapCount int `bson:"flapCount,omitempty" json:"flapCount,omitempty"`
}

// PONRef identifies a PON port on a chassis/slot/port basis. Same
// shape as the polling layer uses; vendor-neutral.
type PONRef struct {
	Chassis int `bson:"chassis" json:"chassis"`
	Slot    int `bson:"slot" json:"slot"`
	Port    int `bson:"port" json:"port"`
}

// TrapVarbind preserves raw varbind data for forensic / debug paths.
// Most consumers use the typed fields on the specific event; this is
// here for the unknown-trap path and audit logs.
type TrapVarbind struct {
	OID   string `bson:"oid" json:"oid"`
	Type  string `bson:"type" json:"type"` // "OctetString" | "Integer" | "OID" | …
	Value string `bson:"value" json:"value"`
}

// ONTDownCause normalises the per-vendor cause-codes attached to ONT
// offline traps. Huawei carries them as enum integers
// (hwGponDeviceOntInfoLastDownCause); ZTE/others differ. Mapping
// happens in olt-manager's vendor-specific catalog before publish.
type ONTDownCause string

const (
	ONTDownCauseUnknown        ONTDownCause = "unknown"
	ONTDownCauseDyingGasp      ONTDownCause = "dying_gasp"      // ONT power loss — last gasp
	ONTDownCauseLOSi           ONTDownCause = "los_individual"  // loss of signal (per-ONT)
	ONTDownCauseLOSPON         ONTDownCause = "los_pon"         // PON-side LOS
	ONTDownCauseLOFi           ONTDownCause = "lof_individual"  // loss of frame
	ONTDownCauseLOAi           ONTDownCause = "loa_individual"  // loss of acknowledge
	ONTDownCauseLOAMi          ONTDownCause = "loami_individual"
	ONTDownCausePowerOff       ONTDownCause = "power_off"       // explicit shutdown
	ONTDownCauseWireDown       ONTDownCause = "wire_down"       // physical disconnect
	ONTDownCauseManualDeact    ONTDownCause = "manual_deactivate"
	ONTDownCauseRogueOnu       ONTDownCause = "rogue_onu"
	ONTDownCauseAuthFail       ONTDownCause = "auth_failure"
	ONTDownCauseDeregister     ONTDownCause = "deregister"
)

// ONTDownEvent is published on OLTEventONTDownSubject. The "cause"
// is the operator-actionable signal — dying_gasp = customer power
// outage; los_individual = customer-side fibre cut; wire_down =
// distribution-side fibre cut; manual_deactivate = ops action.
type ONTDownEvent struct {
	EventEnvelope `bson:",inline"`

	PON       PONRef       `bson:"pon" json:"pon"`
	ONTID     int          `bson:"ontId" json:"ontId"`
	ONTSerial string       `bson:"ontSerial,omitempty" json:"ontSerial,omitempty"`
	Cause     ONTDownCause `bson:"cause" json:"cause"`

	// Human-readable cause text from the trap (raw vendor string,
	// not normalised). Useful for operator debugging when our
	// cause mapping doesn't fit cleanly.
	RawCauseText string `bson:"rawCauseText,omitempty" json:"rawCauseText,omitempty"`

	// Varbinds preserved for the audit log.
	Varbinds []TrapVarbind `bson:"varbinds,omitempty" json:"varbinds,omitempty"`
}

// ONTUpEvent — registered/online notification.
type ONTUpEvent struct {
	EventEnvelope `bson:",inline"`

	PON       PONRef `bson:"pon" json:"pon"`
	ONTID     int    `bson:"ontId" json:"ontId"`
	ONTSerial string `bson:"ontSerial,omitempty" json:"ontSerial,omitempty"`
}

// DyingGaspEvent — separate from ONTDownEvent so notification-service
// can fan out a different message ("customer power loss") rather
// than the generic "ONT offline".
type DyingGaspEvent struct {
	EventEnvelope `bson:",inline"`

	PON       PONRef `bson:"pon" json:"pon"`
	ONTID     int    `bson:"ontId" json:"ontId"`
	ONTSerial string `bson:"ontSerial,omitempty" json:"ontSerial,omitempty"`
}

// AlarmSeverity normalises the per-vendor severity codes.
type AlarmSeverity string

const (
	AlarmSeverityCritical AlarmSeverity = "critical"
	AlarmSeverityMajor    AlarmSeverity = "major"
	AlarmSeverityMinor    AlarmSeverity = "minor"
	AlarmSeverityWarning  AlarmSeverity = "warning"
	AlarmSeverityInfo     AlarmSeverity = "info"
	AlarmSeverityCleared  AlarmSeverity = "cleared"
)

// AlarmEvent — covers the wide swath of OLT-level alarms (PSU, fan,
// temperature, board insert/remove, BIP errors, MAC flap, ...).
// AlarmCode is vendor-specific (e.g. Huawei alarm IDs); Resource
// localises it (board 0/2, fan 0/1, etc.). Consumers fan out by
// (Severity, AlarmCode) pairs.
type AlarmEvent struct {
	EventEnvelope `bson:",inline"`

	AlarmCode   string        `bson:"alarmCode" json:"alarmCode"`
	Severity    AlarmSeverity `bson:"severity" json:"severity"`
	Resource    string        `bson:"resource" json:"resource"`
	Description string        `bson:"description" json:"description"`

	Varbinds []TrapVarbind `bson:"varbinds,omitempty" json:"varbinds,omitempty"`
}

// LinkChangeEvent — RFC 2863 linkUp / linkDown on uplink, board, or
// PON ports.
type LinkChangeEvent struct {
	EventEnvelope `bson:",inline"`

	IfIndex     int    `bson:"ifIndex" json:"ifIndex"`
	IfDescr     string `bson:"ifDescr,omitempty" json:"ifDescr,omitempty"`
	AdminStatus string `bson:"adminStatus,omitempty" json:"adminStatus,omitempty"` // up | down | testing
	OperStatus  string `bson:"operStatus" json:"operStatus"`
}

// LOSEvent — PON-side loss-of-signal (fibre fault on the OLT side of
// the splitter, affecting an entire PON branch).
type LOSEvent struct {
	EventEnvelope `bson:",inline"`

	PON         PONRef `bson:"pon" json:"pon"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
}

// ColdStartEvent — OLT itself rebooted.
type ColdStartEvent struct {
	EventEnvelope `bson:",inline"`

	UptimeBeforeReboot int64 `bson:"uptimeBeforeReboot,omitempty" json:"uptimeBeforeReboot,omitempty"`
}

// AuthFailureEvent — repeated SNMP / CLI auth failures against the
// OLT. Useful for security-ops: someone is trying credentials.
type AuthFailureEvent struct {
	EventEnvelope `bson:",inline"`

	Description string `bson:"description,omitempty" json:"description,omitempty"`
}

// UnknownTrapEvent — fallback for traps we don't yet have a catalog
// entry for. Always carries the full raw varbinds so an operator can
// figure out what the OLT is trying to tell us, and so we can extend
// the catalog without losing the record.
type UnknownTrapEvent struct {
	EventEnvelope `bson:",inline"`

	Varbinds []TrapVarbind `bson:"varbinds" json:"varbinds"`
}
