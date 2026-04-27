package money

import (
	"math"
	"testing"
)

// TestParseString — covers the inputs we expect at the JSON/operator
// boundary (signed, padded, unpadded fractional). Reject paths matter
// most; rounding silently in this layer would defeat the audit
// guarantees the type exists for.
func TestParseString(t *testing.T) {
	cases := []struct {
		in        string
		want      Money
		expectErr bool
	}{
		{"0", 0, false},
		{"0.00", 0, false},
		{"1", 100, false},
		{"1.00", 100, false},
		{"1.5", 150, false}, // pad fractional to 2 digits
		{"1.50", 150, false},
		{"499.99", 49999, false},
		{"+499.99", 49999, false},
		{"-499.99", -49999, false},
		{"0.01", 1, false},
		{"0.50", 50, false},
		{" 100.00 ", 10000, false}, // trim
		{"", 0, true},
		{"abc", 0, true},
		{"1.234", 0, true}, // > 2 decimal places — reject, not round
		{"1..2", 0, true},
		{"1.x", 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParseString(tc.in)
			if (err != nil) != tc.expectErr {
				t.Fatalf("ParseString(%q): err=%v, expectErr=%v", tc.in, err, tc.expectErr)
			}
			if !tc.expectErr && got != tc.want {
				t.Errorf("ParseString(%q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestFromRupees — banker's rounding: half-to-even keeps long sums
// stable in financial reporting (a million rows of x.x5 don't all
// round up).
func TestFromRupees(t *testing.T) {
	cases := []struct {
		name string
		in   float64
		want Money
	}{
		{"zero", 0, 0},
		{"one", 1, 100},
		{"one_fifty", 1.5, 150},
		{"neg_one_fifty", -1.5, -150},
		{"large", 499.99, 49999},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FromRupees(tc.in)
			if got != tc.want {
				t.Errorf("FromRupees(%v) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestArithmetic — sums of many small values is the case where float
// drift shows up; int paise stays exact.
func TestArithmetic(t *testing.T) {
	a := FromPaise(33)
	b := FromPaise(33)
	c := FromPaise(34)
	if got := a.Add(b).Add(c); got != Money(100) {
		t.Errorf("33 + 33 + 34 paise = %d, want 100", got)
	}
	if got := FromPaise(100).Sub(FromPaise(33)); got != Money(67) {
		t.Errorf("100 - 33 paise = %d, want 67", got)
	}
	if got := FromPaise(99).MulInt(7); got != Money(693) {
		t.Errorf("99 paise * 7 = %d, want 693", got)
	}
}

func TestStringFormat(t *testing.T) {
	cases := []struct {
		paise int64
		want  string
	}{
		{0, "INR 0.00"},
		{1, "INR 0.01"},
		{50, "INR 0.50"},
		{100, "INR 1.00"},
		{49999, "INR 499.99"},
		{-49999, "INR -499.99"},
	}
	for _, tc := range cases {
		got := FromPaise(tc.paise).String()
		if got != tc.want {
			t.Errorf("FromPaise(%d).String() = %q, want %q", tc.paise, got, tc.want)
		}
	}
}

// TestMinorUnitsPerMajor — the table is the source of truth for
// how many minor units make up one major unit. Getting this wrong
// silently corrupts every non-INR transaction, so the table is
// pinned by tests.
func TestMinorUnitsPerMajor(t *testing.T) {
	cases := []struct {
		c    Currency
		want int64
	}{
		{CurrencyINR, 100},
		{CurrencyUSD, 100},
		{CurrencyEUR, 100},
		{CurrencyJPY, 1},
		{CurrencyKWD, 1000},
		{CurrencyBHD, 1000},
		{Currency("XYZ"), 100}, // unregistered → safe default of 100
	}
	for _, tc := range cases {
		if got := tc.c.MinorUnitsPerMajor(); got != tc.want {
			t.Errorf("%s.MinorUnitsPerMajor() = %d, want %d", tc.c, got, tc.want)
		}
	}
	if !CurrencyINR.IsValid() {
		t.Error("CurrencyINR should be valid")
	}
	if Currency("XYZ").IsValid() {
		t.Error("XYZ should not be valid")
	}
}

// TestParseFor_JPY — JPY is zero-decimal: "1500" parses to 1500
// minor units; any fractional input is operator error and must be
// rejected, never silently rounded (audit requirement).
func TestParseFor_JPY(t *testing.T) {
	cases := []struct {
		in        string
		want      Money
		expectErr bool
	}{
		{"1500", 1500, false},
		{"0", 0, false},
		{"-100", -100, false},
		{"1500.50", 0, true}, // JPY has no fractional part
		{"1.0", 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParseFor(CurrencyJPY, tc.in)
			if (err != nil) != tc.expectErr {
				t.Fatalf("ParseFor(JPY, %q): err=%v, expectErr=%v", tc.in, err, tc.expectErr)
			}
			if !tc.expectErr && got != tc.want {
				t.Errorf("ParseFor(JPY, %q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestParseFor_KWD — KWD is three-decimal: "1.234" parses to 1234
// minor units. Padding short fractions and rejecting >3 decimals
// matches audit rules for the Gulf dinars.
func TestParseFor_KWD(t *testing.T) {
	cases := []struct {
		in        string
		want      Money
		expectErr bool
	}{
		{"1", 1000, false},
		{"1.234", 1234, false},
		{"1.5", 1500, false},  // pad 1 → 100
		{"1.23", 1230, false}, // pad 23 → 230
		{"0.001", 1, false},
		{"-1.234", -1234, false},
		{"1.2345", 0, true}, // > 3 decimals
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParseFor(CurrencyKWD, tc.in)
			if (err != nil) != tc.expectErr {
				t.Fatalf("ParseFor(KWD, %q): err=%v, expectErr=%v", tc.in, err, tc.expectErr)
			}
			if !tc.expectErr && got != tc.want {
				t.Errorf("ParseFor(KWD, %q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// TestFromMajor_USD — sanity check: $1.00 is 100 cents, banker's
// rounding rounds 0.005 to even cent.
func TestFromMajor_USD(t *testing.T) {
	cases := []struct {
		in   float64
		want Money
	}{
		{0, 0},
		{1, 100},
		{49.99, 4999},
		{-49.99, -4999},
	}
	for _, tc := range cases {
		got := FromMajor(CurrencyUSD, tc.in)
		if got != tc.want {
			t.Errorf("FromMajor(USD, %v) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

// TestFormatFor — render uses the currency's precision: 2 for INR,
// 0 for JPY, 3 for KWD. Logs and error messages depend on this being
// stable.
func TestFormatFor(t *testing.T) {
	cases := []struct {
		c    Currency
		m    Money
		want string
	}{
		{CurrencyINR, FromPaise(49999), "INR 499.99"},
		{CurrencyUSD, FromPaise(4999), "USD 49.99"},
		{CurrencyJPY, FromPaise(1500), "JPY 1500"},
		{CurrencyKWD, FromPaise(1234), "KWD 1.234"},
		{CurrencyKWD, FromPaise(0), "KWD 0.000"},
	}
	for _, tc := range cases {
		got := tc.m.FormatFor(tc.c)
		if got != tc.want {
			t.Errorf("FormatFor(%s, %d) = %q, want %q", tc.c, tc.m, got, tc.want)
		}
	}
}

// TestAddSubSafe — the whole point is to refuse cross-currency
// arithmetic at boundaries where currency comes from two independent
// sources. Same-currency case must work transparently.
func TestAddSubSafe(t *testing.T) {
	a, b := FromPaise(100), FromPaise(50)
	if got, err := a.AddSafe(CurrencyINR, b, CurrencyINR); err != nil || got != 150 {
		t.Errorf("AddSafe same currency: got=%d err=%v, want 150 nil", got, err)
	}
	if _, err := a.AddSafe(CurrencyINR, b, CurrencyUSD); err == nil {
		t.Error("AddSafe across currencies should error")
	}
	if got, err := a.SubSafe(CurrencyINR, b, CurrencyINR); err != nil || got != 50 {
		t.Errorf("SubSafe same currency: got=%d err=%v, want 50 nil", got, err)
	}
	if _, err := a.SubSafe(CurrencyINR, b, CurrencyUSD); err == nil {
		t.Error("SubSafe across currencies should error")
	}
}

// TestMulIntChecked — overflow detection. The unchecked path is the
// hot path; this is for caller-supplied multipliers (CSV imports,
// API quantity fields).
func TestMulIntChecked(t *testing.T) {
	if got, err := FromPaise(99).MulIntChecked(7); err != nil || got != 693 {
		t.Errorf("MulIntChecked 99*7: got=%d err=%v, want 693 nil", got, err)
	}
	if got, err := FromPaise(0).MulIntChecked(math.MaxInt64); err != nil || got != 0 {
		t.Errorf("MulIntChecked 0*max: got=%d err=%v, want 0 nil", got, err)
	}
	if _, err := FromPaise(math.MaxInt64).MulIntChecked(2); err == nil {
		t.Error("MulIntChecked maxint*2 should overflow")
	}
}

// TestFixedRateConverter — round-trip a USD value through INR at a
// fixed rate. Same-currency must be a no-op; missing rate must error,
// not silently return zero.
func TestFixedRateConverter(t *testing.T) {
	conv := FixedRateConverter{
		Rates: map[Currency]map[Currency]float64{
			CurrencyUSD: {CurrencyINR: 83.5},
			CurrencyINR: {CurrencyUSD: 1.0 / 83.5},
		},
	}
	usd := FromMajor(CurrencyUSD, 10) // 1000 cents
	inr, err := conv.Convert(usd, CurrencyUSD, CurrencyINR)
	if err != nil {
		t.Fatalf("USD→INR: %v", err)
	}
	// 10 USD * 83.5 = 835 INR = 83500 paise
	if inr != FromPaise(83500) {
		t.Errorf("10 USD → INR = %d paise, want 83500", inr.Paise())
	}
	if got, err := conv.Convert(usd, CurrencyUSD, CurrencyUSD); err != nil || got != usd {
		t.Errorf("USD→USD should be no-op: got=%d err=%v", got, err)
	}
	if _, err := conv.Convert(usd, CurrencyJPY, CurrencyINR); err == nil {
		t.Error("missing source rate should error")
	}
	if _, err := conv.Convert(usd, CurrencyUSD, CurrencyJPY); err == nil {
		t.Error("missing target rate should error")
	}
}
