package subscriberevent

type HotspotAuthMethod string

const (
	HotspotAuthOTP        HotspotAuthMethod = "OTP"
	HotspotAuthPassword   HotspotAuthMethod = "PASSWORD"
	HotspotAuthVoucher    HotspotAuthMethod = "VOUCHER"
	HotspotAuthSocial     HotspotAuthMethod = "SOCIAL"
	HotspotAuthMAC        HotspotAuthMethod = "MAC"
	HotSpotOneTap         HotspotAuthMethod = "ONE_TAP"
	HotSpotQRCode         HotspotAuthMethod = "QR_CODE"
	HotSpotAuthSSO        HotspotAuthMethod = "SSO"
	HostSpotAuthMagicLink HotspotAuthMethod = "MAGIC_LINK"
)

type HotspotProfileStatus string

const (
	HotspotProfilePending      HotspotProfileStatus = "PENDING"
	HotspotProfileActive       HotspotProfileStatus = "ACTIVE"
	HotspotProfileSuspended    HotspotProfileStatus = "SUSPENDED"
	HotspotProfileTerminated   HotspotProfileStatus = "TERMINATED"
	HotspotProfileDisconnected HotspotProfileStatus = "DISCONNECTED"
	HotspotProfileChurned      HotspotProfileStatus = "CHURNED"
	HotspotProfileBlacklisted  HotspotProfileStatus = "BLACKLISTED"
)
