package subscriberevent

type BroadbandAccessType string

const (
	BroadbandAccessFTTH     BroadbandAccessType = "FTTH"
	BroadbandAccessWireless BroadbandAccessType = "WIRELESS"
	BroadbandAccessDSL      BroadbandAccessType = "DSL"
	BroadbandAccessCable    BroadbandAccessType = "CABLE"
	BroadbandAccessVDSL     BroadbandAccessType = "VDSL"
	BroadbandAccessOther    BroadbandAccessType = "OTHER"
)

type BroadbandSubscriptionStatus string

const (
	BroadbandStatusPendingInstall BroadbandSubscriptionStatus = "PENDING_INSTALL"
	BroadbandStatusActive         BroadbandSubscriptionStatus = "ACTIVE"
	BroadbandStatusSuspended      BroadbandSubscriptionStatus = "SUSPENDED"
	BroadbandStatusDisconnected   BroadbandSubscriptionStatus = "DISCONNECTED"
	BroadbandStatusTerminated     BroadbandSubscriptionStatus = "TERMINATED"
	BroadbandStatusChurned        BroadbandSubscriptionStatus = "CHURNED"
	BroadbandStatusBlacklisted    BroadbandSubscriptionStatus = "BLACKLISTED"
	BroadbandStatusGrace          BroadbandSubscriptionStatus = "GRACE"
)

type IPAssignmentType string

const (
	IPAssignmentPPPoE  IPAssignmentType = "PPPOE"
	IPAssignmentStatic IPAssignmentType = "STATIC_IP"
	IPAssignmentDHCP   IPAssignmentType = "DHCP"
)
