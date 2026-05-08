// Package peerassistevent defines Peer Assist v4 types: PeerAssistRequest
// (FE-A asks for help), PeerHandoff (the ceremony when FE-B agrees to take
// inventory off FE-A's hands), PeerHandoffLine (per-item).
//
// Owned operationally by inventory-service. Distinct from Job Transfer v3
// which transfers a ticket between FE-techs (owned by ticket-service —
// see common/events/models/jobtransferevent/).
//
// Source: backend-contract §8.3, §8.4, assets-prd.md §6.7.
package peerassistevent
