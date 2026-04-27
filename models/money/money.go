// Package money provides a typed int64 representation of monetary
// values stored in the smallest unit of the currency (paise for INR,
// cents for USD). It exists because float64 is unsafe for currency
// math: 0.1 + 0.2 != 0.3, sums of 1000+ rows accumulate rounding error
// that shows up in audit reports as "off by one paisa", and the wire
// format ambiguity ("is 1.5 fifty paise or one and a half rupees?")
// has bitten the codebase before.
//
// Conventions:
//
//   - Internal storage is int64 paise. 1 INR == 100 paise.
//   - JSON marshals as a plain integer (paise). API consumers format
//     for display on their side. This matches Stripe's convention and
//     avoids round-trip float drift.
//   - Negative values are legal — represent debits, refunds, etc.
//     Validation that a particular field must be non-negative belongs
//     in the schema layer (`gte=0`), not here.
//   - The default currency is INR. The package has helpers parameterised
//     on a Currency value so the same type can carry USD/EUR cents
//     without code changes; today every helper without an explicit
//     currency assumes INR for backwards compatibility with the
//     dataset.
//
// Use Money everywhere a money value crosses a service boundary or a
// DB write — Schemas, models, event payloads. Pure-display floats
// (e.g. "₹4,99.99" formatting) belong in the dashboard layer.
package money

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Currency is a thin wrapper around an ISO-4217 code so the type
// system can distinguish e.g. INR from USD. The current dataset is INR
// only — multi-currency arithmetic (cross-currency totals) is not
// supported by the helpers below; the type carries currency for
// future-proofing and audit clarity.
type Currency string

const (
	CurrencyINR Currency = "INR"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

// MinorUnitsPerMajor reports how many minor units (paise / cents)
// make up one major unit (rupees / dollars). Hardcoded at 100 for the
// currencies in use; revisit when adding currencies with different
// decimal places (JPY = 1, KWD = 1000).
func (c Currency) MinorUnitsPerMajor() int64 {
	return 100
}

// Money is a typed integer count of the smallest unit of the assumed
// currency (paise for INR, cents for USD). The default-zero value is
// a legitimate "zero money" — no special "unset" sentinel.
type Money int64

// FromPaise constructs a Money from a paise count.
func FromPaise(paise int64) Money { return Money(paise) }

// FromRupees converts a rupee value (e.g. 499.99) to Money. Rounds to
// the nearest paise using banker's rounding to keep aggregate sums
// consistent with downstream financial reporting.
func FromRupees(rupees float64) Money {
	return Money(int64(math.RoundToEven(rupees * 100)))
}

// MustFromString parses a "499.99" / "499" / "0.50" decimal string
// into Money and panics on bad input. Use this in seeds or tests
// where the input is hardcoded; production callers should use
// ParseString.
func MustFromString(s string) Money {
	m, err := ParseString(s)
	if err != nil {
		panic(err)
	}
	return m
}

// ParseString parses a decimal rupee string into Money. Accepts
// "499.99", "499", "0.5", "0.50", or "+/-" signed forms. Rejects
// values with more than 2 decimal places — paise is the precision
// floor, anything finer is operator error and silent rounding would
// lose audit trail.
func ParseString(s string) (Money, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("money: empty input")
	}
	negative := false
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		negative = true
		s = s[1:]
	}
	parts := strings.SplitN(s, ".", 2)
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("money: invalid integer part %q: %w", parts[0], err)
	}
	var frac int64
	if len(parts) == 2 {
		fracStr := parts[1]
		if len(fracStr) > 2 {
			return 0, fmt.Errorf("money: more than 2 decimal places in %q (paise is the smallest unit)", s)
		}
		// Pad to 2 digits — "0.5" → 50 paise, not 5.
		for len(fracStr) < 2 {
			fracStr += "0"
		}
		frac, err = strconv.ParseInt(fracStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("money: invalid fractional part %q: %w", fracStr, err)
		}
	}
	total := whole*100 + frac
	if negative {
		total = -total
	}
	return Money(total), nil
}

// Paise returns the underlying int64 count.
func (m Money) Paise() int64 { return int64(m) }

// Rupees returns the rupee equivalent as a float64 — convenience for
// display only. Don't use in arithmetic; round-tripping through float
// reintroduces the very precision loss this type exists to prevent.
func (m Money) Rupees() float64 { return float64(m) / 100.0 }

// IsZero reports whether the value is exactly zero paise.
func (m Money) IsZero() bool { return m == 0 }

// IsNegative reports whether the value is below zero. Used by schema
// validators that need to reject negative values without dropping
// down to int64.
func (m Money) IsNegative() bool { return m < 0 }

// Add returns m + other. Both values must be in the same currency —
// this package does not perform FX conversion.
func (m Money) Add(other Money) Money { return m + other }

// Sub returns m - other.
func (m Money) Sub(other Money) Money { return m - other }

// MulInt scales the amount by an integer multiplier (e.g. quantity).
// Use this rather than (m * Money(qty)) to make the unit-less scalar
// explicit at the call site.
func (m Money) MulInt(n int64) Money { return Money(int64(m) * n) }

// String renders as "<currency> <rupees>.<paise>" with two decimals,
// e.g. "INR 499.99". For UI formatting (locale-aware ₹ symbol,
// thousands separators) defer to the dashboard's display layer —
// this method exists for log lines and error messages, not customer
// invoices.
func (m Money) String() string {
	return fmt.Sprintf("INR %.2f", m.Rupees())
}

// Equal reports paise-level equality. Defined explicitly so a future
// switch to a struct-backed Money (e.g. carrying a Currency tag)
// doesn't silently change comparison semantics.
func (m Money) Equal(other Money) bool { return m == other }
