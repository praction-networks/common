package utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// IPv4ToUint32 converts IPv4 string to uint32 (big-endian)
// Example: "192.168.1.100" -> 0xC0A80164
func IPv4ToUint32(ip string) (uint32, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return 0, fmt.Errorf("invalid IPv4: %s", ip)
	}
	ipv4 := parsed.To4()
	if ipv4 == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ip)
	}
	return binary.BigEndian.Uint32(ipv4), nil
}

// Uint32ToIPv4 converts uint32 to IPv4 string
// Example: 0xC0A80164 -> "192.168.1.100"
func Uint32ToIPv4(ip uint32) string {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, ip)
	return net.IP(bytes).String()
}

// IPv6ToBytes converts IPv6 string to 16-byte binary
// Example: "2001:db8::1" -> [0x20, 0x01, 0x0d, 0xb8, ...]
func IPv6ToBytes(ip string) ([]byte, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil, fmt.Errorf("invalid IPv6: %s", ip)
	}
	ipv6 := parsed.To16()
	if ipv6 == nil {
		return nil, fmt.Errorf("not an IPv6 address: %s", ip)
	}
	return ipv6, nil
}

// BytesToIPv6 converts 16-byte binary to IPv6 string
// Example: [0x20, 0x01, ...] -> "2001:db8::1"
func BytesToIPv6(bytes []byte) string {
	return net.IP(bytes).String()
}

// MACToBytes converts MAC string to 6-byte binary
// Example: "00:11:22:33:44:55" -> [0x00, 0x11, 0x22, 0x33, 0x44, 0x55]
func MACToBytes(mac string) ([]byte, error) {
	hw, err := net.ParseMAC(mac)
	if err != nil {
		return nil, fmt.Errorf("invalid MAC address: %w", err)
	}
	if len(hw) != 6 {
		return nil, fmt.Errorf("MAC address must be 6 bytes, got %d", len(hw))
	}
	return hw, nil
}

// BytesToMAC converts 6-byte binary to MAC string
// Example: [0x00, 0x11, 0x22, 0x33, 0x44, 0x55] -> "00:11:22:33:44:55"
func BytesToMAC(bytes []byte) string {
	return net.HardwareAddr(bytes).String()
}

// UsernameToBytes converts username string to UTF-8 bytes
// Example: "john.doe@example.com" -> []byte("john.doe@example.com")
func UsernameToBytes(username string) []byte {
	return []byte(username)
}

// BytesToUsername converts UTF-8 bytes to username string
// Example: []byte("john.doe@example.com") -> "john.doe@example.com"
func BytesToUsername(bytes []byte) string {
	return string(bytes)
}

// TimeToMillis converts time.Time to int64 (milliseconds since Unix epoch)
// Example: time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC) -> 1736941800000
func TimeToMillis(t time.Time) int64 {
	return t.UnixMilli()
}

// MillisToTime converts int64 (milliseconds) to time.Time
// Example: 1736941800000 -> time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
func MillisToTime(millis int64) time.Time {
	return time.Unix(millis/1000, (millis%1000)*1e6)
}

// ParseEventTimeFromDevice parses EventTime from device timestamp
// MUST use device timestamp, not system time
// Supports multiple timestamp formats commonly used by network devices
func ParseEventTimeFromDevice(deviceTimestamp string) (int64, error) {
	// Common timestamp formats used by network devices
	formats := []string{
		time.RFC3339,                // 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano,            // 2006-01-02T15:04:05.999999999Z07:00
		"2006-01-02T15:04:05Z07:00", // RFC3339 without nanoseconds
		"2006-01-02 15:04:05",       // MySQL format
		"2006-01-02T15:04:05",       // ISO 8601 without timezone
		"Jan 2 15:04:05 2006",       // Syslog format
		"2006-01-02T15:04:05.000Z",  // ISO 8601 with milliseconds
	}

	for _, format := range formats {
		t, err := time.Parse(format, deviceTimestamp)
		if err == nil {
			return TimeToMillis(t), nil
		}
	}

	return 0, fmt.Errorf("failed to parse device timestamp: %s (tried %d formats)", deviceTimestamp, len(formats))
}
