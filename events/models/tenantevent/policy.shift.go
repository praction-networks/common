package tenantevent

// PolicyShift is the per-tenant shift-policy bucket.
//
// Source: backend-contract §11.1. Drives Account hub break/idle flows and
// the server-side end-nudge job.
type PolicyShift struct {
	FirstBreakSkipsReason         bool     `json:"firstBreakSkipsReason,omitempty"         bson:"firstBreakSkipsReason,omitempty"`
	SubsequentBreaksRequireReason bool     `json:"subsequentBreaksRequireReason,omitempty" bson:"subsequentBreaksRequireReason,omitempty"`
	BreakReasons                  []string `json:"breakReasons,omitempty"                  bson:"breakReasons,omitempty"`
	MaxBreakMinutes               *int     `json:"maxBreakMinutes,omitempty"               bson:"maxBreakMinutes,omitempty"`
	EndNudgeGraceMinutes          int      `json:"endNudgeGraceMinutes,omitempty"          bson:"endNudgeGraceMinutes,omitempty"`
	IdleMinutesThreshold          int      `json:"idleMinutesThreshold,omitempty"          bson:"idleMinutesThreshold,omitempty"`
	LocationTrackingEnabled       bool     `json:"locationTrackingEnabled,omitempty"       bson:"locationTrackingEnabled,omitempty"`
}
