package tenantevent

// PolicyShift is the per-tenant shift-policy bucket.
//
// Source: backend-contract §11.1. Drives Account hub break/idle flows and
// the server-side end-nudge job.
type PolicyShift struct {
	FirstBreakSkipsReason         bool     `json:"firstBreakSkipsReason"                   bson:"firstBreakSkipsReason"`
	SubsequentBreaksRequireReason bool     `json:"subsequentBreaksRequireReason"           bson:"subsequentBreaksRequireReason"`
	BreakReasons                  []string `json:"breakReasons,omitempty"                  bson:"breakReasons,omitempty"`
	MaxBreakMinutes               *int     `json:"maxBreakMinutes,omitempty"               bson:"maxBreakMinutes,omitempty"`
	EndNudgeGraceMinutes          int      `json:"endNudgeGraceMinutes"                    bson:"endNudgeGraceMinutes"`
	IdleMinutesThreshold          int      `json:"idleMinutesThreshold"                    bson:"idleMinutesThreshold"`
	LocationTrackingEnabled       bool     `json:"locationTrackingEnabled"                 bson:"locationTrackingEnabled"`
}
