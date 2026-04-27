package utils

import "testing"

// TestCurrencyFromCountry_Names — the bug this exists to prevent: an
// AddressModel stores "India" (not "IN"), so a code-only lookup fell
// back to USD for every Indian tenant. Names must resolve correctly.
func TestCurrencyFromCountry_Names(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"India", "INR"},
		{"india", "INR"},
		{"  India ", "INR"},
		{"United States", "USD"},
		{"USA", "USD"},
		{"United Kingdom", "GBP"},
		{"UK", "GBP"},
		{"United Arab Emirates", "AED"},
		{"UAE", "AED"},
		{"Japan", "JPY"},
		{"Kuwait", "KWD"},
	}
	for _, tc := range cases {
		if got := CurrencyFromCountry(tc.in); got != tc.want {
			t.Errorf("CurrencyFromCountry(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestCurrencyFromCountry_Alpha2 — alpha-2 still works for any caller
// that wires up ISO codes (RBI feeds, third-party gateway responses).
func TestCurrencyFromCountry_Alpha2(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"IN", "INR"},
		{"in", "INR"},
		{" IN ", "INR"},
		{"US", "USD"},
		{"GB", "GBP"},
		{"AE", "AED"},
		{"JP", "JPY"},
		{"KW", "KWD"},
	}
	for _, tc := range cases {
		if got := CurrencyFromCountry(tc.in); got != tc.want {
			t.Errorf("CurrencyFromCountry(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestCurrencyFromCountry_Fallback — empty / unknown defaults to INR
// (India-dominant dataset). Don't change to USD without auditing the
// resolver warning logs first.
func TestCurrencyFromCountry_Fallback(t *testing.T) {
	if got := CurrencyFromCountry(""); got != "INR" {
		t.Errorf("empty = %q, want INR", got)
	}
	if got := CurrencyFromCountry("Atlantis"); got != "INR" {
		t.Errorf("unknown = %q, want INR", got)
	}
	if got := CurrencyFromCountry("XX"); got != "INR" {
		t.Errorf("unknown alpha2 = %q, want INR", got)
	}
}

func TestIsKnownCountry(t *testing.T) {
	if !IsKnownCountry("India") {
		t.Error("India should be known")
	}
	if !IsKnownCountry("IN") {
		t.Error("IN should be known")
	}
	if IsKnownCountry("Atlantis") {
		t.Error("Atlantis should be unknown")
	}
	if IsKnownCountry("") {
		t.Error("empty should be unknown")
	}
}
