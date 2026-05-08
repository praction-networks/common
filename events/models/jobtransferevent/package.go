// Package jobtransferevent defines Job Transfer v3 types: JobTransferRequest
// (the transfer attempt) and JobTransferOffer (per-peer offer inside a
// request).
//
// Owned operationally by ticket-service. Inventory-service consumes the
// associated events for clean-up CTAs. Types live in common/ because both
// services and admin-dashboard need the same shape.
//
// Source: backend-contract §8.5.
package jobtransferevent
