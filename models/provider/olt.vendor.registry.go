package provider

// OLT vendor registry — single source of truth for the vendors,
// models, firmware tracks, and capabilities olt-manager understands.
// tenant-service exposes this via ``GET /api/v1/tenant/olt/registry``;
// the admin dashboard reads it once on form open and uses it to drive
// vendor / model / firmware / protocol / SNMP-version pickers.
//
// Adding a new vendor = add an entry below + ship the corresponding
// olt-manager adapter. For OEM rebrands, set the ``OEM`` field — the
// olt-manager adapter then delegates to the OEM's backend.

// OLTVendor is the canonical identifier used in events and DB docs.
type OLTVendor string

const (
	// Own-OEM vendors — they manufacture their own OLTs.
	OLTVendorHuawei    OLTVendor = "Huawei"
	OLTVendorZTE       OLTVendor = "ZTE"
	OLTVendorNokia     OLTVendor = "Nokia"
	OLTVendorFiberhome OLTVendor = "Fiberhome"
	OLTVendorVSOL      OLTVendor = "VSOL"
	OLTVendorBDCOM     OLTVendor = "BDCOM"
	OLTVendorDASAN     OLTVendor = "DASAN"
	OLTVendorGenesis   OLTVendor = "Genesis"
	OLTVendorGeneric   OLTVendor = "Generic"

	// Rebrands of VSOL — same CLI grammar and SNMP MIBs as VSOL.
	OLTVendorSyrotech OLTVendor = "Syrotech"
	OLTVendorCIG      OLTVendor = "CIG"
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
//
// OEM relationships: many small ISP brands (Syrotech, CIG, etc.) ship
// rebadged hardware from a parent OEM (typically VSOL). Set ``OEM`` to
// the parent vendor's ``Value`` so olt-manager can route to the parent's
// adapter while operators still see their nameplate vendor in the UI.
// Empty ``OEM`` means the vendor is its own OEM.
type OLTVendorInfo struct {
	Value                  string           `json:"value"`             // canonical key (matches OLTVendor)
	Label                  string           `json:"label"`             // display name
	Description            string           `json:"description"`
	OEM                    string           `json:"oem,omitempty"`     // empty = own OEM; non-empty = points to OEM's Value (e.g. "VSOL")
	Models                 []string         `json:"models"`            // curated common models; UI also allows free text
	FirmwareVersions       []string         `json:"firmwareVersions"`  // curated common firmware tracks; UI allows free text
	SupportedCLIProtocols  []OLTCLIProtocol `json:"supportedCliProtocols"`
	SupportedSNMPVersions  []OLTSNMPVersion `json:"supportedSnmpVersions"`
	SupportsEnableMode     bool             `json:"supportsEnableMode"`     // "super" / "enable" prompt
	SupportsPrivateKeyAuth bool             `json:"supportsPrivateKeyAuth"` // SSH key auth (vs password only)
	DefaultCLIPort         map[string]int   `json:"defaultCliPort"`         // protocol → port
	Notes                  string           `json:"notes,omitempty"`
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
// for the UI to surface it. Keep models / firmware lists short and
// curated; users can free-type anything outside the list.
//
// Ordering matters for UI defaults: own-OEM tier-1 vendors first, then
// own-OEM tier-2, then OEM rebrands grouped near their parent, with
// Generic last as the catch-all.
var OLTVendorRegistry = []OLTVendorInfo{
	// ------------------- Tier-1 own-OEM (global) -------------------
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
			"MA5600T",
			"MA5603T",
			"MA5608T",
			"MA5606T",
		},
		FirmwareVersions: []string{
			"V800R022C00",
			"V800R021C10",
			"V800R021C00",
			"V800R020C10",
			"V800R019C10",
			"V800R018C10",
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
		Value:       string(OLTVendorZTE),
		Label:       "ZTE",
		Description: "ZTE OLTs (C300 / C320 / C600 / C650 series).",
		Models: []string{
			"C300",
			"C320",
			"C600",
			"C610",
			"C650",
		},
		FirmwareVersions: []string{
			"V2.1.0",
			"V2.0.1",
			"V2.0.0",
			"V1.5.0",
			"V1.2.5",
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
		FirmwareVersions: []string{
			"FX0.62.05",
			"FX0.61.05",
			"FX0.60.05",
			"FX0.50.05",
		},
		SupportedCLIProtocols:  []OLTCLIProtocol{OLTCLIProtocolSSH},
		SupportedSNMPVersions:  []OLTSNMPVersion{OLTSNMPVersionV2c, OLTSNMPVersionV3},
		SupportsEnableMode:     false,
		SupportsPrivateKeyAuth: true,
		DefaultCLIPort: map[string]int{
			string(OLTCLIProtocolSSH): 22,
		},
	},
	{
		Value:       string(OLTVendorFiberhome),
		Label:       "Fiberhome",
		Description: "Fiberhome OLTs (AN5516 / AN6000 series). Widely deployed in APAC.",
		Models: []string{
			"AN5516-04",
			"AN5516-06",
			"AN5516-01",
			"AN6000-7",
			"AN6000-15",
			"AN6000-17",
		},
		FirmwareVersions: []string{
			"RP0900",
			"RP0800",
			"RP0700",
			"RP0600",
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
		Value:       string(OLTVendorDASAN),
		Label:       "DASAN",
		Description: "DASAN / DZS OLTs — V-series (Korean/Swedish merger). Carrier-grade.",
		Models: []string{
			"V8240",
			"V8242",
			"V5812",
			"V5824",
			"V5808",
		},
		FirmwareVersions: []string{
			"NOS-9.5",
			"NOS-9.0",
			"NOS-8.5",
			"NOS-8.0",
		},
		SupportedCLIProtocols:  []OLTCLIProtocol{OLTCLIProtocolSSH},
		SupportedSNMPVersions:  []OLTSNMPVersion{OLTSNMPVersionV2c, OLTSNMPVersionV3},
		SupportsEnableMode:     true,
		SupportsPrivateKeyAuth: true,
		DefaultCLIPort: map[string]int{
			string(OLTCLIProtocolSSH): 22,
		},
	},

	// ------------------- Tier-2 own-OEM (Asia / India) -------------------
	{
		Value:       string(OLTVendorVSOL),
		Label:       "VSOL",
		Description: "VSOL OLTs — own brand. Also OEMs for Syrotech, CIG, and other resold brands.",
		Models: []string{
			"V1600D",
			"V1600D4",
			"V1600G",
			"V1600G1",
			"V1600G2",
			"V2724G",
			"V3000",
		},
		FirmwareVersions: []string{
			"V2.0.06",
			"V2.0.05",
			"V2.0.04",
			"V1.9.5",
			"V1.9.0",
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
		Value:       string(OLTVendorBDCOM),
		Label:       "BDCOM",
		Description: "BDCOM OLTs (P3310 / P3608 / P3616 series). Popular with mid-tier ISPs.",
		Models: []string{
			"P3310B",
			"P3310C",
			"P3608",
			"P3608-4P",
			"P3616-2TE",
			"GP3600-08",
			"GP3600-16",
		},
		FirmwareVersions: []string{
			"10.4.0K",
			"10.3.0H",
			"10.2.0E",
			"10.1.0F",
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
		Value:       string(OLTVendorGenesis),
		Label:       "Genesis",
		Description: "Genesis OLTs.",
		Models:      []string{}, // free-type until model list is curated
		FirmwareVersions: []string{},
		SupportedCLIProtocols:  []OLTCLIProtocol{OLTCLIProtocolSSH, OLTCLIProtocolTelnet},
		SupportedSNMPVersions:  []OLTSNMPVersion{OLTSNMPVersionV2c, OLTSNMPVersionV3},
		SupportsEnableMode:     true,
		SupportsPrivateKeyAuth: true,
		DefaultCLIPort: map[string]int{
			string(OLTCLIProtocolSSH):    22,
			string(OLTCLIProtocolTelnet): 23,
		},
		Notes: "Add specific models / firmware tracks as we onboard customers running Genesis hardware.",
	},

	// ------------------- VSOL rebrands (OEM = VSOL) -------------------
	{
		Value:       string(OLTVendorSyrotech),
		Label:       "Syrotech",
		Description: "Syrotech (VSOL OEM) — Indian ISP-favourite rebrand.",
		OEM:         string(OLTVendorVSOL),
		Models: []string{
			"SY-GPON-1600D",
			"SY-GPON-1600G",
			"SY-EPON-1600D",
			"SY-XPON-04",
			"SY-XPON-08",
		},
		FirmwareVersions: []string{
			"V2.0.06",
			"V2.0.05",
			"V2.0.04",
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
		Value:       string(OLTVendorCIG),
		Label:       "CIG (Cloud Infotech)",
		Description: "CIG / Cloud Infotech (VSOL OEM) — rebranded VSOL hardware.",
		OEM:         string(OLTVendorVSOL),
		Models: []string{
			"CIG-1600D",
			"CIG-1600G",
			"CIG-2724G",
		},
		FirmwareVersions: []string{
			"V2.0.06",
			"V2.0.05",
			"V2.0.04",
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

	// ------------------- Catch-all -------------------
	{
		Value:       string(OLTVendorGeneric),
		Label:       "Generic",
		Description: "Fallback for vendors without a dedicated adapter — read-only polling, no provisioning.",
		Models:      []string{},
		FirmwareVersions: []string{},
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

// ResolveOEM returns the canonical OEM value for a vendor — the
// vendor itself if it's its own OEM, or its parent OEM otherwise.
// olt-manager uses this to pick the right vendor adapter for rebrands.
func ResolveOEM(vendor string) string {
	for _, v := range OLTVendorRegistry {
		if v.Value == vendor {
			if v.OEM != "" {
				return v.OEM
			}
			return v.Value
		}
	}
	return vendor
}

// LookupOLTVendor returns the registry entry for a given vendor value
// (or nil if unknown). Caller should treat the return as immutable.
func LookupOLTVendor(vendor string) *OLTVendorInfo {
	for i := range OLTVendorRegistry {
		if OLTVendorRegistry[i].Value == vendor {
			return &OLTVendorRegistry[i]
		}
	}
	return nil
}

// IsKnownModel reports whether ``model`` is in the curated list for
// ``vendor``. Returns false if the vendor is unknown or carries an
// empty model list — strict policy: no curated list ⇒ no valid model.
func IsKnownModel(vendor, model string) bool {
	v := LookupOLTVendor(vendor)
	if v == nil || len(v.Models) == 0 {
		return false
	}
	for _, m := range v.Models {
		if m == model {
			return true
		}
	}
	return false
}

// IsKnownFirmware reports whether ``firmware`` is in the curated list
// for ``vendor``. Same strict policy as IsKnownModel.
func IsKnownFirmware(vendor, firmware string) bool {
	v := LookupOLTVendor(vendor)
	if v == nil || len(v.FirmwareVersions) == 0 {
		return false
	}
	for _, f := range v.FirmwareVersions {
		if f == firmware {
			return true
		}
	}
	return false
}
