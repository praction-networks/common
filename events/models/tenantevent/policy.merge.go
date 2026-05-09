package tenantevent

import "fmt"

// MergeTenantPolicy returns a new TenantPolicy with each non-nil bucket from
// patch merged over current. Bumps Version on the returned merged policy.
//
// Field-level rules within each bucket:
//   - *bool / *int / *float64 set if non-nil; nil = unchanged.
//   - Plain string / int set if non-zero (omitempty semantics; combine with
//     pointer types where false-y values are meaningful).
//   - []string / []T replaced wholesale (PATCH must send full intended list).
//   - Nested struct fields recursively merged via the same rules.
//
// Returns changedKeys as dot-paths (e.g. "assets.dropoffLockWindowHours")
// — only keys whose value actually differed from current. Empty changedKeys
// is valid (no-op patch); Version still bumps so audit trail is consistent.
func MergeTenantPolicy(current TenantPolicy, patch TenantPolicyPatch) (TenantPolicy, []string) {
	merged := current
	merged.Version = current.Version + 1
	var changed []string

	if patch.Account != nil {
		merged.Account = mergeAccount(current.Account, *patch.Account, "account", &changed)
	}
	if patch.Assets != nil {
		merged.Assets = mergeAssets(current.Assets, *patch.Assets, "assets", &changed)
	}
	if patch.Auth != nil {
		merged.Auth = mergeAuth(current.Auth, *patch.Auth, "auth", &changed)
	}
	if patch.Notifications != nil {
		merged.Notifications = mergeNotifications(current.Notifications, *patch.Notifications, "notifications", &changed)
	}
	if patch.Onboard != nil {
		merged.Onboard = mergeOnboard(current.Onboard, *patch.Onboard, "onboard", &changed)
	}
	if patch.Shift != nil {
		merged.Shift = mergeShift(current.Shift, *patch.Shift, "shift", &changed)
	}
	return merged, changed
}

func mergeAccount(cur, p PolicyAccount, prefix string, changed *[]string) PolicyAccount {
	out := cur
	track := func(key string) {
		*changed = append(*changed, fmt.Sprintf("%s.%s", prefix, key))
	}
	if p.EmailChangeAllowed != nil && (cur.EmailChangeAllowed == nil || *cur.EmailChangeAllowed != *p.EmailChangeAllowed) {
		out.EmailChangeAllowed = p.EmailChangeAllowed
		track("emailChangeAllowed")
	}
	if p.AllowMobileChange != nil && (cur.AllowMobileChange == nil || *cur.AllowMobileChange != *p.AllowMobileChange) {
		out.AllowMobileChange = p.AllowMobileChange
		track("allowMobileChange")
	}
	if p.AllowWhatsappChange != nil && (cur.AllowWhatsappChange == nil || *cur.AllowWhatsappChange != *p.AllowWhatsappChange) {
		out.AllowWhatsappChange = p.AllowWhatsappChange
		track("allowWhatsappChange")
	}
	if p.FuelMode != "" && cur.FuelMode != p.FuelMode {
		out.FuelMode = p.FuelMode
		track("fuelMode")
	}
	if p.FuelRatePerKm != nil && (cur.FuelRatePerKm == nil || *cur.FuelRatePerKm != *p.FuelRatePerKm) {
		out.FuelRatePerKm = p.FuelRatePerKm
		track("fuelRatePerKm")
	}
	if p.LeaveCategories != nil && !stringSliceEqual(cur.LeaveCategories, p.LeaveCategories) {
		out.LeaveCategories = p.LeaveCategories
		track("leaveCategories")
	}
	if p.VehicleRequired != nil && (cur.VehicleRequired == nil || *cur.VehicleRequired != *p.VehicleRequired) {
		out.VehicleRequired = p.VehicleRequired
		track("vehicleRequired")
	}
	if p.DocumentVaultEnabled != nil && (cur.DocumentVaultEnabled == nil || *cur.DocumentVaultEnabled != *p.DocumentVaultEnabled) {
		out.DocumentVaultEnabled = p.DocumentVaultEnabled
		track("documentVaultEnabled")
	}
	return out
}

func mergeAssets(cur, p PolicyAssets, prefix string, changed *[]string) PolicyAssets {
	out := cur
	if p.DropoffLockWindowHours != 0 && cur.DropoffLockWindowHours != p.DropoffLockWindowHours {
		out.DropoffLockWindowHours = p.DropoffLockWindowHours
		*changed = append(*changed, prefix+".dropoffLockWindowHours")
	}
	if p.RecoveryAcknowledgement != "" && cur.RecoveryAcknowledgement != p.RecoveryAcknowledgement {
		out.RecoveryAcknowledgement = p.RecoveryAcknowledgement
		*changed = append(*changed, prefix+".recoveryAcknowledgement")
	}
	return out
}

func mergeAuth(cur, p PolicyAuth, prefix string, changed *[]string) PolicyAuth {
	out := cur
	if p.AccessTokenTtlMinutes != 0 && cur.AccessTokenTtlMinutes != p.AccessTokenTtlMinutes {
		out.AccessTokenTtlMinutes = p.AccessTokenTtlMinutes
		*changed = append(*changed, prefix+".accessTokenTtlMinutes")
	}
	return out
}

func mergeNotifications(cur, p PolicyNotifications, prefix string, changed *[]string) PolicyNotifications {
	out := cur
	if p.CriticalCategoriesAllowed != nil && !stringSliceEqual(cur.CriticalCategoriesAllowed, p.CriticalCategoriesAllowed) {
		out.CriticalCategoriesAllowed = p.CriticalCategoriesAllowed
		*changed = append(*changed, prefix+".criticalCategoriesAllowed")
	}
	if p.SupportChannels.Whatsapp != "" && cur.SupportChannels.Whatsapp != p.SupportChannels.Whatsapp {
		out.SupportChannels.Whatsapp = p.SupportChannels.Whatsapp
		*changed = append(*changed, prefix+".supportChannels.whatsapp")
	}
	if p.SupportChannels.Call != "" && cur.SupportChannels.Call != p.SupportChannels.Call {
		out.SupportChannels.Call = p.SupportChannels.Call
		*changed = append(*changed, prefix+".supportChannels.call")
	}
	if p.SupportChannels.Email != "" && cur.SupportChannels.Email != p.SupportChannels.Email {
		out.SupportChannels.Email = p.SupportChannels.Email
		*changed = append(*changed, prefix+".supportChannels.email")
	}
	return out
}

func mergeOnboard(cur, p PolicyOnboard, prefix string, changed *[]string) PolicyOnboard {
	out := cur
	if p.AgreementChannel != "" && cur.AgreementChannel != p.AgreementChannel {
		out.AgreementChannel = p.AgreementChannel
		*changed = append(*changed, prefix+".agreementChannel")
	}
	if p.OtpChannel != "" && cur.OtpChannel != p.OtpChannel {
		out.OtpChannel = p.OtpChannel
		*changed = append(*changed, prefix+".otpChannel")
	}
	return out
}

func mergeShift(cur, p PolicyShift, prefix string, changed *[]string) PolicyShift {
	out := cur
	if p.MaxBreakMinutes != nil && (cur.MaxBreakMinutes == nil || *cur.MaxBreakMinutes != *p.MaxBreakMinutes) {
		out.MaxBreakMinutes = p.MaxBreakMinutes
		*changed = append(*changed, prefix+".maxBreakMinutes")
	}
	if p.EndNudgeGraceMinutes != 0 && cur.EndNudgeGraceMinutes != p.EndNudgeGraceMinutes {
		out.EndNudgeGraceMinutes = p.EndNudgeGraceMinutes
		*changed = append(*changed, prefix+".endNudgeGraceMinutes")
	}
	if p.IdleMinutesThreshold != 0 && cur.IdleMinutesThreshold != p.IdleMinutesThreshold {
		out.IdleMinutesThreshold = p.IdleMinutesThreshold
		*changed = append(*changed, prefix+".idleMinutesThreshold")
	}
	if p.BreakReasons != nil && !stringSliceEqual(cur.BreakReasons, p.BreakReasons) {
		out.BreakReasons = p.BreakReasons
		*changed = append(*changed, prefix+".breakReasons")
	}
	// Plain bools: cannot distinguish "set to false" from "didn't send" — use false-detection only.
	// Per spec §6.1, PolicyShift's bool fields are not pointer-typed because admin sends full policy
	// for shift bucket; this is acceptable for V1. If finer control needed, migrate to *bool later.
	if p.FirstBreakSkipsReason != cur.FirstBreakSkipsReason {
		out.FirstBreakSkipsReason = p.FirstBreakSkipsReason
		*changed = append(*changed, prefix+".firstBreakSkipsReason")
	}
	if p.SubsequentBreaksRequireReason != cur.SubsequentBreaksRequireReason {
		out.SubsequentBreaksRequireReason = p.SubsequentBreaksRequireReason
		*changed = append(*changed, prefix+".subsequentBreaksRequireReason")
	}
	if p.LocationTrackingEnabled != cur.LocationTrackingEnabled {
		out.LocationTrackingEnabled = p.LocationTrackingEnabled
		*changed = append(*changed, prefix+".locationTrackingEnabled")
	}
	return out
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
