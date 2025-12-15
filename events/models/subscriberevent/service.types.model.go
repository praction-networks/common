package subscriberevent

type ServiceType string

const (
	// Core Connectivity Services
	ServiceTypeBroadband   ServiceType = "BROADBAND"    // FTTH / Wireless / DSL
	ServiceTypePublicIP    ServiceType = "PUBLIC_IP"    // Static Public IP, /30, /29 etc.
	ServiceTypeVPN         ServiceType = "VPN"          // L2TP, IPsec, MPLS, SD-WAN
	ServiceTypeManagedWiFi ServiceType = "MANAGED_WIFI" // Enterprise Managed Wi-Fi

	// Communication Services
	ServiceTypeVOIP ServiceType = "VOIP" // Voice over IP (SIP, Hosted PBX)

	// Entertainment / Media
	ServiceTypeIPTV ServiceType = "IPTV" // Live TV
	ServiceTypeOTT  ServiceType = "OTT"  // Disney+, SonyLiv, JioCinema, Amazon Prime, Hotstar etc.

	// Smart Home / IoT
	ServiceTypeSmartHome    ServiceType = "SMART_HOME"    // Home automation bundles
	ServiceTypeSmartCamera  ServiceType = "SMART_CAMERA"  // CCTV/Cloud Camera monitoring
	ServiceTypeHomeSecurity ServiceType = "HOME_SECURITY" // Sensors, alarms, monitoring

	// Hotspot / Guest WiFi
	ServiceTypeHotspot ServiceType = "HOTSPOT" // Captive portal access

	// Enterprise / Special Services
	ServiceTypeColocation ServiceType = "COLOCATION" // Datacenter racks, hosting
	ServiceTypeCloudPBX   ServiceType = "CLOUD_PBX"  // Enterprise phone systems
	ServiceTypeILL        ServiceType = "ILL"        // Internet Leased Line
	ServiceTypeMPLS       ServiceType = "MPLS"       // Layer-3 MPLS
	ServiceTypeSDWAN      ServiceType = "SDWAN"      // SD-WAN subscription

	// Catch-All
	ServiceTypeOther ServiceType = "OTHER" // Anything not yet categorized
)
