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
	CriticalCategoriesAllowed []string `json:"criticalCategoriesAllowed,omitempty" bson:"criticalCategoriesAllowed,omitempty"`
}

// Validate ensures CriticalCategoriesAllowed is a subset of the platform list.
func (n PolicyNotifications) Validate() error {
	if n.CriticalCategoriesAllowed == nil {
		return nil
	}
	platform := make(map[string]struct{}, len(notificationevent.CriticalCategoriesPlatform))
	for _, c := range notificationevent.CriticalCategoriesPlatform {
		platform[string(c)] = struct{}{}
	}
	for _, c := range n.CriticalCategoriesAllowed {
		if _, ok := platform[c]; !ok {
			return fmt.Errorf("CriticalCategoriesAllowed contains %q which is not in platform whitelist", c)
		}
	}
	return nil
}
