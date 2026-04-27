package provider

// OLT vendor registry — single source of truth for the vendors,
// models, and capabilities olt-manager understands. tenant-service
// exposes this via ``GET /api/v1/tenant/olt/registry``; the admin
// dashboard reads it once on form open and uses it to drive the
// vendor / model / protocol / SNMP-version pickers.
//
// Adding a new vendor = add an entry below + ship the corresponding
// olt-manager adapter. Frontend auto-renders the new options without
// a code change.

// OLTVendor is the canonical identifier used in events and DB docs.
type OLTVendor string

const (
	OLTVendorHuawei  OLTVendor = "Huawei"
	OLTVendorZTE     OLTVendor = "ZTE"
	OLTVendorNokia   OLTVendor = "Nokia"
	OLTVendorGeneric OLTVendor = "Generic"
)

// OLTCLIProtocol — vendor-CLI transport.
type OLTCLIProtocol string

const (
	OLTCLIProtocolSSH    OLTCLIProtocol = "ssh"
	OLTCLIProtocolTelnet OLTCLIProtocol = "telnet"
)

// OLTSNMPVersion — supported SNMP variants.
type OLTSNMPVersion string

const (
	OLTSNMPVersionV1  OLTSNMPVersion = "v1"
	OLTSNMPVersionV2c OLTSNMPVersion = "v2c"
	OLTSNMPVersionV3  OLTSNMPVersion = "v3"
)

// OLTSNMPSecurityLevel — v3 security mode.
type OLTSNMPSecurityLevel string

const (
	OLTSNMPSecurityNoAuthNoPriv OLTSNMPSecurityLevel = "noAuthNoPriv"
	OLTSNMPSecurityAuthNoPriv   OLTSNMPSecurityLevel = "authNoPriv"
	OLTSNMPSecurityAuthPriv     OLTSNMPSecurityLevel = "authPriv"
)

// OLTSNMPAuthProtocol — v3 authentication digest.
type OLTSNMPAuthProtocol string

const (
	OLTSNMPAuthMD5    OLTSNMPAuthProtocol = "MD5"
	OLTSNMPAuthSHA    OLTSNMPAuthProtocol = "SHA"
	OLTSNMPAuthSHA224 OLTSNMPAuthProtocol = "SHA224"
	OLTSNMPAuthSHA256 OLTSNMPAuthProtocol = "SHA256"
	OLTSNMPAuthSHA384 OLTSNMPAuthProtocol = "SHA384"
	OLTSNMPAuthSHA512 OLTSNMPAuthProtocol = "SHA512"
)

// OLTSNMPPrivProtocol — v3 encryption cipher.
type OLTSNMPPrivProtocol string

const (
	OLTSNMPPrivDES    OLTSNMPPrivProtocol = "DES"
	OLTSNMPPriv3DES   OLTSNMPPrivProtocol = "3DES"
	OLTSNMPPrivAES128 OLTSNMPPrivProtocol = "AES128"
	OLTSNMPPrivAES192 OLTSNMPPrivProtocol = "AES192"
	OLTSNMPPrivAES256 OLTSNMPPrivProtocol = "AES256"
)

// OLTStatus — operational lifecycle stage.
type OLTStatus string

const (
	OLTStatusProvisioned OLTStatus = "provisioned"
	OLTStatusActive      OLTStatus = "active"
	OLTStatusSuspended   OLTStatus = "suspended"
)

// OLTVendorInfo describes one vendor's capability surface. The frontend
// uses these flags to gate which inputs are visible and which values
// are selectable for a given vendor.
type OLTVendorInfo struct {
	Value                  string             `json:"value"`                  // canonical key (matches OLTVendor)
	Label                  string             `json:"label"`                  // display name
	Description            string             `json:"description"`
	Models                 []string           `json:"models"`                 // curated common models; UI also allows free text
	SupportedCLIProtocols  []OLTCLIProtocol   `json:"supportedCliProtocols"`
	SupportedSNMPVersions  []OLTSNMPVersion   `json:"supportedSnmpVersions"`
	SupportsEnableMode     bool               `json:"supportsEnableMode"`     // "super" / "enable" prompt
	SupportsPrivateKeyAuth bool               `json:"supportsPrivateKeyAuth"` // SSH key auth (vs password only)
	DefaultCLIPort         map[string]int     `json:"defaultCliPort"`         // protocol → port
	Notes                  string             `json:"notes,omitempty"`
}

// OLTRegistry is the wire shape returned from
// ``GET /api/v1/tenant/olt/registry``. Keep additive — frontends pin
// to current keys and ignore unknowns.
type OLTRegistry struct {
	Vendors            []OLTVendorInfo        `json:"vendors"`
	CLIProtocols       []OLTCLIProtocol       `json:"cliProtocols"`
	SNMPVersions       []OLTSNMPVersion       `json:"snmpVersions"`
	SNMPSecurityLevels []OLTSNMPSecurityLevel `json:"snmpSecurityLevels"`
	SNMPAuthProtocols  []OLTSNMPAuthProtocol  `json:"snmpAuthProtocols"`
	SNMPPrivProtocols  []OLTSNMPPrivProtocol  `json:"snmpPrivProtocols"`
	Statuses           []OLTStatus            `json:"statuses"`
}

// OLTVendorRegistry — adding a vendor here is the only change needed
// for the UI to surface it. Keep models lists short and curated; users
// can free-type if their model isn't listed.
var OLTVendorRegistry = []OLTVendorInfo{
	{
		Value:       string(OLTVendorHuawei),
		Label:       "Huawei",
		Description: "Huawei OLTs (MA5800 / MA5600T series). SmartAX line.",
		Models: []string{
			"MA5800-X2",
			"MA5800-X7",
			"MA5800-X15",
			"MA5800-X17",
			"MA5683T",
			"MA5680T",
		},
		SupportedCLIProtocols:  []OLTCLIProtocol{OLTCLIProtocolSSH, OLTCLIProtocolTelnet},
		SupportedSNMPVersions:  []OLTSNMPVersion{OLTSNMPVersionV2c, OLTSNMPVersionV3},
		SupportsEnableMode:     true, // "enable" / "super" privileged mode
		SupportsPrivateKeyAuth: true,
		DefaultCLIPort: map[string]int{
			string(OLTCLIProtocolSSH):    22,
			string(OLTCLIProtocolTelnet): 23,
		},
	},
	{
		Value:       string(OLTVendorZTE),
		Label:       "ZTE",
		Description: "ZTE OLTs (C300 / C320 / C600 series).",
		Models: []string{
			"C300",
			"C320",
			"C600",
			"C610",
			"C650",
		},
		SupportedCLIProtocols:  []OLTCLIProtocol{OLTCLIProtocolSSH, OLTCLIProtocolTelnet},
		SupportedSNMPVersions:  []OLTSNMPVersion{OLTSNMPVersionV2c, OLTSNMPVersionV3},
		SupportsEnableMode:     true,
		SupportsPrivateKeyAuth: true,
		DefaultCLIPort: map[string]int{
			string(OLTCLIProtocolSSH):    22,
			string(OLTCLIProtocolTelnet): 23,
		},
	},
	{
		Value:       string(OLTVendorNokia),
		Label:       "Nokia",
		Description: "Nokia ISAM / FX-series PON OLTs.",
		Models: []string{
			"ISAM 7360 FX-4",
			"ISAM 7360 FX-8",
			"ISAM 7360 FX-16",
			"ISAM 7360 FX-30",
			"FX-4 Lite",
		},
		SupportedCLIProtocols:  []OLTCLIProtocol{OLTCLIProtocolSSH},
		SupportedSNMPVersions:  []OLTSNMPVersion{OLTSNMPVersionV2c, OLTSNMPVersionV3},
		SupportsEnableMode:     false, // Nokia uses single-shell auth
		SupportsPrivateKeyAuth: true,
		DefaultCLIPort: map[string]int{
			string(OLTCLIProtocolSSH): 22,
		},
	},
	{
		Value:       string(OLTVendorGeneric),
		Label:       "Generic",
		Description: "Fallback for vendors without a dedicated adapter — read-only polling, no provisioning.",
		Models:      []string{}, // free-type only
		SupportedCLIProtocols: []OLTCLIProtocol{
			OLTCLIProtocolSSH,
			OLTCLIProtocolTelnet,
		},
		SupportedSNMPVersions: []OLTSNMPVersion{
			OLTSNMPVersionV1,
			OLTSNMPVersionV2c,
			OLTSNMPVersionV3,
		},
		SupportsEnableMode:     false,
		SupportsPrivateKeyAuth: true,
		DefaultCLIPort: map[string]int{
			string(OLTCLIProtocolSSH):    22,
			string(OLTCLIProtocolTelnet): 23,
		},
	},
}

// GetOLTRegistry returns the full registry for the OLT integration
// form. Caller should treat the return as immutable.
func GetOLTRegistry() OLTRegistry {
	return OLTRegistry{
		Vendors:      OLTVendorRegistry,
		CLIProtocols: []OLTCLIProtocol{OLTCLIProtocolSSH, OLTCLIProtocolTelnet},
		SNMPVersions: []OLTSNMPVersion{
			OLTSNMPVersionV1,
			OLTSNMPVersionV2c,
			OLTSNMPVersionV3,
		},
		SNMPSecurityLevels: []OLTSNMPSecurityLevel{
			OLTSNMPSecurityNoAuthNoPriv,
			OLTSNMPSecurityAuthNoPriv,
			OLTSNMPSecurityAuthPriv,
		},
		SNMPAuthProtocols: []OLTSNMPAuthProtocol{
			OLTSNMPAuthMD5,
			OLTSNMPAuthSHA,
			OLTSNMPAuthSHA224,
			OLTSNMPAuthSHA256,
			OLTSNMPAuthSHA384,
			OLTSNMPAuthSHA512,
		},
		SNMPPrivProtocols: []OLTSNMPPrivProtocol{
			OLTSNMPPrivDES,
			OLTSNMPPriv3DES,
			OLTSNMPPrivAES128,
			OLTSNMPPrivAES192,
			OLTSNMPPrivAES256,
		},
		Statuses: []OLTStatus{
			OLTStatusProvisioned,
			OLTStatusActive,
			OLTStatusSuspended,
		},
	}
}
