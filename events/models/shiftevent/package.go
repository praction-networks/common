// Package shiftevent defines the closed-shift record (denormalised rollup;
// one document per closed shift session) and the in-shift Break sub-document.
//
// Owned operationally by user-service; types live in common/ so every consumer
// (admin-dashboard, ticket-service for streak inputs, etc.) compiles against
// the same shape.
//
// Source: backend-contract.md §4.2, account-prd.md §7.0.
package shiftevent
