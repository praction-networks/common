package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
)

// IPPrefixCode represents a compressed prefix code for file-level indexing
type IPPrefixCode struct {
	Code   string // "AA", "AB", etc.
	Prefix string // "103.83.131.0/24"
	CIDR   int    // 24
	Family string // "IPv4" or "IPv6"
}

// CalculateIPv4PrefixCodes calculates prefix codes for IPv4 addresses
// Returns /16 and /24 prefixes
func CalculateIPv4PrefixCodes(ipv4 uint32) []IPPrefixCode {
	codes := []IPPrefixCode{}

	// Extract IP string
	ipStr := Uint32ToIPv4(ipv4)
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return codes
	}

	// Calculate /16 prefix (first 2 octets)
	prefix16 := net.IPNet{
		IP:   net.IPv4(ip[0], ip[1], 0, 0),
		Mask: net.CIDRMask(16, 32),
	}
	code16 := generatePrefixCode(prefix16.String(), "IPv4", 16)
	codes = append(codes, IPPrefixCode{
		Code:   code16,
		Prefix: prefix16.String(),
		CIDR:   16,
		Family: "IPv4",
	})

	// Calculate /24 prefix (first 3 octets)
	prefix24 := net.IPNet{
		IP:   net.IPv4(ip[0], ip[1], ip[2], 0),
		Mask: net.CIDRMask(24, 32),
	}
	code24 := generatePrefixCode(prefix24.String(), "IPv4", 24)
	codes = append(codes, IPPrefixCode{
		Code:   code24,
		Prefix: prefix24.String(),
		CIDR:   24,
		Family: "IPv4",
	})

	return codes
}

// CalculateIPv6PrefixCodes calculates prefix codes for IPv6 addresses
// Returns /32 and /64 prefixes
func CalculateIPv6PrefixCodes(ipv6 []byte) []IPPrefixCode {
	codes := []IPPrefixCode{}

	if len(ipv6) != 16 {
		return codes
	}

	ip := net.IP(ipv6)

	// Calculate /32 prefix (first 32 bits)
	prefix32IP := make(net.IP, 16)
	copy(prefix32IP, ip)
	for i := 4; i < 16; i++ {
		prefix32IP[i] = 0
	}
	prefix32 := net.IPNet{
		IP:   prefix32IP,
		Mask: net.CIDRMask(32, 128),
	}
	code32 := generatePrefixCode(prefix32.String(), "IPv6", 32)
	codes = append(codes, IPPrefixCode{
		Code:   code32,
		Prefix: prefix32.String(),
		CIDR:   32,
		Family: "IPv6",
	})

	// Calculate /64 prefix (first 64 bits)
	prefix64IP := make(net.IP, 16)
	copy(prefix64IP, ip)
	for i := 8; i < 16; i++ {
		prefix64IP[i] = 0
	}
	prefix64 := net.IPNet{
		IP:   prefix64IP,
		Mask: net.CIDRMask(64, 128),
	}
	code64 := generatePrefixCode(prefix64.String(), "IPv6", 64)
	codes = append(codes, IPPrefixCode{
		Code:   code64,
		Prefix: prefix64.String(),
		CIDR:   64,
		Family: "IPv6",
	})

	return codes
}

// generatePrefixCode generates a short code for a prefix (2-character code)
func generatePrefixCode(prefix string, family string, cidr int) string {
	// Use prefix string + CIDR as key for consistent hashing
	key := fmt.Sprintf("%s:%s:%d", family, prefix, cidr)

	// Generate short code (2 characters) using SHA-256 hash
	hash := sha256.Sum256([]byte(key))
	// Use first 2 bytes of hash, encode as hex, take first 2 characters
	hashCode := hex.EncodeToString(hash[:2])[:2]

	// Return uppercase 2-character code (no family prefix - keep it simple)
	return string(hashCode[0]) + string(hashCode[1])
}

// ExtractUniquePrefixCodes extracts unique prefix codes from a list
func ExtractUniquePrefixCodes(prefixCodes []IPPrefixCode) []string {
	seen := make(map[string]bool)
	unique := []string{}

	for _, pc := range prefixCodes {
		if !seen[pc.Code] {
			seen[pc.Code] = true
			unique = append(unique, pc.Code)
		}
	}

	return unique
}

// GeneratePrefixCodeFromIP generates prefix code from IP string (for backwards compatibility)
func GeneratePrefixCodeFromIP(ipStr string) []IPPrefixCode {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		// Convert to uint32
		ipv4Uint32 := uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3])
		return CalculateIPv4PrefixCodes(ipv4Uint32)
	}

	// IPv6
	return CalculateIPv6PrefixCodes(ip.To16())
}
