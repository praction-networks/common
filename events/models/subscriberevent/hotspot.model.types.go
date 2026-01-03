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
	HotspotProfileTerminated   HotspotProfileStatus = "TERMINATED"
	HotspotProfileBlacklisted  HotspotProfileStatus = "BLACKLISTED"
	HotspotProfileExpired      HotspotProfileStatus = "EXPIRED"
)
