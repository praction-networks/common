package tenantevent

import (
	"fmt"

	"github.com/praction-networks/common/events/models/notificationevent"
)

// PolicyNotifications is the per-tenant notifications-policy bucket.
// Source: backend-contract §11.4 + §9.5, design notification.md §6.
type PolicyNotifications struct {
	// CriticalCategoriesAllowed narrows the platform CRITICAL whitelist.
	// nil = full platform list; empty = no CRITICAL allowed; otherwise
	// must be a subset of notificationevent.CriticalCategoriesPlatform.
	CriticalCategoriesAllowed []notificationevent.NotificationCategory `json:"criticalCategoriesAllowed,omitempty" bson:"criticalCategoriesAllowed,omitempty"`
	SupportChannels           PolicySupportChannels `json:"supportChannels,omitempty"           bson:"supportChannels,omitempty"`
}

// PolicySupportChannels contains the tenant-configured contact endpoints surfaced
// in the FE app's support sheet. Each is opt-in; empty strings are omitted.
//
// Source: backend-contract §11.4.
type PolicySupportChannels struct {
	Whatsapp string `json:"whatsapp,omitempty" bson:"whatsapp,omitempty"`
	Call     string `json:"call,omitempty"     bson:"call,omitempty"`
	Email    string `json:"email,omitempty"    bson:"email,omitempty"`
}

// Validate ensures CriticalCategoriesAllowed is a subset of the platform list.
func (n PolicyNotifications) Validate() error {
	if n.CriticalCategoriesAllowed == nil {
		return nil
	}
	platformList := notificationevent.CriticalCategoriesPlatform()
	platform := make(map[notificationevent.NotificationCategory]struct{}, len(platformList))
	for _, c := range platformList {
		platform[c] = struct{}{}
	}
	for _, c := range n.CriticalCategoriesAllowed {
		if _, ok := platform[c]; !ok {
			return fmt.Errorf("CriticalCategoriesAllowed contains %q which is not in platform whitelist", c)
		}
	}
	return nil
}
