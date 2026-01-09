package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// HashSourceIP creates a SHA256 hash of the source IP and encodes it as base64url (URL-safe)
// This is used for file organization to avoid issues with special characters in IPs
// Returns a 44-character base64url-encoded hash
func HashSourceIP(sourceIP string) string {
	if sourceIP == "" {
		return ""
	}

	// Create SHA256 hash
	hash := sha256.Sum256([]byte(sourceIP))

	// Encode as base64url (URL-safe, no padding)
	// base64url uses - and _ instead of + and /, and omits padding =
	hashStr := base64.RawURLEncoding.EncodeToString(hash[:])

	return hashStr
}

