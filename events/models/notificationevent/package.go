// Package notificationevent defines shared types for notification events:
// NotificationCategory (domain tagging), NotificationCriticality (delivery
// tier), NotificationDeliveryStatus (lifecycle), and the
// CriticalCategoriesPlatform whitelist used by tenant policy validation.
//
// Owned operationally by notification-service. Consumers: every service that
// publishes notifications (subscriber, ticket, inventory, etc.) plus
// admin-dashboard for filter chips.
//
// Source: field-central notifications-prd.md §4, §16; backend-contract §11.4.
package notificationevent
