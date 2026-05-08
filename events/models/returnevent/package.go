// Package returnevent defines two-step drop-off types: Return (top-level
// batch), ReturnLine (per-asset/consumable line), ReturnLineStatus (lifecycle).
//
// Owned operationally by inventory-service. Custody flips on WM accept,
// not FE submit — see backend-contract §8.2.
//
// Package directory is `returnevent` (not `return`) because `return` is a
// reserved Go keyword.
package returnevent
