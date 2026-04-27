package money

import "testing"

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
		{"1.5", 150, false},     // pad fractional to 2 digits
		{"1.50", 150, false},
		{"499.99", 49999, false},
		{"+499.99", 49999, false},
		{"-499.99", -49999, false},
		{"0.01", 1, false},
		{"0.50", 50, false},
		{" 100.00 ", 10000, false}, // trim
		{"", 0, true},
		{"abc", 0, true},
		{"1.234", 0, true},  // > 2 decimal places — reject, not round
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
