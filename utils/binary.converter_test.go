package utils

import (
	"testing"
	"time"
)

func TestIPv4ToUint32(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    uint32
		wantErr bool
	}{
		{
			name:    "valid IPv4",
			ip:      "192.168.1.100",
			want:    0xC0A80164,
			wantErr: false,
		},
		{
			name:    "valid IPv4 - public IP",
			ip:      "103.83.131.238",
			want:    0x675383EE, // 103*256^3 + 83*256^2 + 131*256 + 238
			wantErr: false,
		},
		{
			name:    "invalid IP",
			ip:      "invalid",
			want:    0,
			wantErr: true,
		},
		{
			name:    "IPv6 address",
			ip:      "2001:db8::1",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPv4ToUint32(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPv4ToUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IPv4ToUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint32ToIPv4(t *testing.T) {
	tests := []struct {
		name string
		ip   uint32
		want string
	}{
		{
			name: "valid IPv4",
			ip:   0xC0A80164,
			want: "192.168.1.100",
		},
		{
			name: "valid IPv4 - public IP",
			ip:   0x675383EE,
			want: "103.83.131.238",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Uint32ToIPv4(tt.ip)
			if got != tt.want {
				t.Errorf("Uint32ToIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4RoundTrip(t *testing.T) {
	ips := []string{
		"192.168.1.100",
		"103.83.131.238",
		"142.250.192.78",
		"10.0.0.1",
		"172.16.0.1",
	}

	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			uint32Val, err := IPv4ToUint32(ip)
			if err != nil {
				t.Fatalf("IPv4ToUint32() error = %v", err)
			}
			got := Uint32ToIPv4(uint32Val)
			if got != ip {
				t.Errorf("Round trip failed: got %v, want %v", got, ip)
			}
		})
	}
}

func TestIPv6ToBytes(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantLen int
		wantErr bool
	}{
		{
			name:    "valid IPv6",
			ip:      "2001:db8::1",
			wantLen: 16,
			wantErr: false,
		},
		{
			name:    "valid IPv6 - full format",
			ip:      "2001:0db8:0000:0000:0000:0000:0000:0001",
			wantLen: 16,
			wantErr: false,
		},
		{
			name:    "invalid IP",
			ip:      "invalid",
			wantLen: 0,
			wantErr: true,
		},
		{
			name:    "IPv4-mapped IPv6 (valid IPv6 representation)",
			ip:      "192.168.1.1",
			wantLen: 16, // IPv4 can be represented as IPv6 (IPv4-mapped)
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPv6ToBytes(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPv6ToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("IPv6ToBytes() length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestIPv6RoundTrip(t *testing.T) {
	ips := []string{
		"2001:db8::1",
		"::1",
		"2001:db8:100::1",
	}

	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			bytes, err := IPv6ToBytes(ip)
			if err != nil {
				t.Fatalf("IPv6ToBytes() error = %v", err)
			}
			got := BytesToIPv6(bytes)
			// Note: Go's net.IP.String() normalizes IPv6 (removes leading zeros)
			// So "2001:0db8:0000:0000:0000:0000:0000:0001" becomes "2001:db8::1"
			// We need to convert both to compare
			gotNormalized, _ := IPv6ToBytes(got)
			wantNormalized, _ := IPv6ToBytes(ip)
			if len(gotNormalized) != len(wantNormalized) {
				t.Errorf("Round trip failed: got %v, want %v", got, ip)
			}
			for i := range gotNormalized {
				if gotNormalized[i] != wantNormalized[i] {
					t.Errorf("Round trip failed: got %v, want %v", got, ip)
					break
				}
			}
		})
	}
}

func TestMACToBytes(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		wantLen int
		wantErr bool
	}{
		{
			name:    "valid MAC - colon format",
			mac:     "00:11:22:33:44:55",
			wantLen: 6,
			wantErr: false,
		},
		{
			name:    "valid MAC - hyphen format",
			mac:     "00-11-22-33-44-55",
			wantLen: 6,
			wantErr: false,
		},
		{
			name:    "invalid MAC",
			mac:     "invalid",
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MACToBytes(tt.mac)
			if (err != nil) != tt.wantErr {
				t.Errorf("MACToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("MACToBytes() length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestMACRoundTrip(t *testing.T) {
	macs := []string{
		"00:11:22:33:44:55",
		"aa:bb:cc:dd:ee:ff",
		"00-11-22-33-44-55",
	}

	for _, mac := range macs {
		t.Run(mac, func(t *testing.T) {
			bytes, err := MACToBytes(mac)
			if err != nil {
				t.Fatalf("MACToBytes() error = %v", err)
			}
			got := BytesToMAC(bytes)
			// Note: net.HardwareAddr.String() always uses colon format
			expected := "00:11:22:33:44:55"
			if mac == "00-11-22-33-44-55" {
				expected = "00:11:22:33:44:55"
			}
			if got != expected && got != mac {
				t.Errorf("Round trip: got %v, want %v or %v", got, mac, expected)
			}
		})
	}
}

func TestUsernameRoundTrip(t *testing.T) {
	usernames := []string{
		"john.doe@example.com",
		"user123",
		"test_user",
		"admin",
	}

	for _, username := range usernames {
		t.Run(username, func(t *testing.T) {
			bytes := UsernameToBytes(username)
			got := BytesToUsername(bytes)
			if got != username {
				t.Errorf("Round trip failed: got %v, want %v", got, username)
			}
		})
	}
}

func TestTimeToMillis(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want int64
	}{
		{
			name: "Unix epoch",
			t:    time.Unix(0, 0),
			want: 0,
		},
		{
			name: "2025-01-15 10:30:00 UTC",
			t:    time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
			want: 1736937000000, // Calculated: Jan 15, 2025 10:30:00 UTC in milliseconds
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimeToMillis(tt.t)
			if got != tt.want {
				t.Errorf("TimeToMillis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMillisToTime(t *testing.T) {
	tests := []struct {
		name    string
		millis  int64
		want    time.Time
		wantStr string
	}{
		{
			name:    "Unix epoch",
			millis:  0,
			want:    time.Unix(0, 0),
			wantStr: "1970-01-01T00:00:00Z",
		},
		{
			name:    "2025-01-15 10:30:00 UTC",
			millis:  1736937000000, // Jan 15, 2025 10:30:00 UTC in milliseconds
			want:    time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
			wantStr: "2025-01-15T10:30:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MillisToTime(tt.millis)
			if !got.Equal(tt.want) {
				t.Errorf("MillisToTime() = %v, want %v", got, tt.want)
			}
			if got.Format(time.RFC3339) != tt.wantStr {
				t.Errorf("MillisToTime() string = %v, want %v", got.Format(time.RFC3339), tt.wantStr)
			}
		})
	}
}

func TestTimeRoundTrip(t *testing.T) {
	times := []time.Time{
		time.Now(),
		time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
		time.Date(2020, 6, 1, 12, 0, 0, 0, time.UTC),
		time.Unix(0, 0),
	}

	for _, tTime := range times {
		t.Run(tTime.Format(time.RFC3339), func(t *testing.T) {
			millis := TimeToMillis(tTime)
			got := MillisToTime(millis)
			// Compare up to millisecond precision
			if got.Unix() != tTime.Unix() || got.Nanosecond()/1e6 != tTime.Nanosecond()/1e6 {
				t.Errorf("Round trip failed: got %v, want %v", got, tTime)
			}
		})
	}
}

func TestParseEventTimeFromDevice(t *testing.T) {
	tests := []struct {
		name    string
		ts      string
		wantErr bool
	}{
		{
			name:    "RFC3339 format",
			ts:      "2025-01-15T10:30:00Z",
			wantErr: false,
		},
		{
			name:    "RFC3339 with timezone",
			ts:      "2025-01-15T10:30:00+05:30",
			wantErr: false,
		},
		{
			name:    "MySQL format",
			ts:      "2025-01-15 10:30:00",
			wantErr: false,
		},
		{
			name:    "Syslog format",
			ts:      "Jan 15 10:30:00 2025",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			ts:      "invalid timestamp",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseEventTimeFromDevice(tt.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEventTimeFromDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
