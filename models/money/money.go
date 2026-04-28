// Package money provides a typed int64 representation of monetary
// values stored in the smallest unit of the currency (paise for INR,
// cents for USD, yen for JPY, fils for KWD/BHD/OMR). It exists because
// float64 is unsafe for currency math: 0.1 + 0.2 != 0.3, sums of 1000+
// rows accumulate rounding error that shows up in audit reports as
// "off by one paisa", and the wire format ambiguity ("is 1.5 fifty
// paise or one and a half rupees?") has bitten the codebase before.
//
// Conventions:
//
//   - Internal storage is int64 minor units. Exact factor depends on
//     currency: 1 INR == 100 paise, 1 JPY == 1 yen, 1 KWD == 1000 fils.
//   - JSON marshals as a plain integer (minor units). API consumers
//     format for display on their side. This matches Stripe's
//     convention and avoids round-trip float drift.
//   - Negative values are legal — represent debits, refunds, etc.
//     Validation that a particular field must be non-negative belongs
//     in the schema layer (`gte=0`), not here.
//   - Currency travels in a sibling field on the parent model
//     (Invoice.Currency, Payment.Currency, …). Money itself stays a
//     bare int64 to keep wire format and DB columns stable across
//     services. Helpers that need precision use a Currency parameter.
//
// Use Money everywhere a money value crosses a service boundary or a
// DB write — Schemas, models, event payloads. Pure-display floats
// (e.g. "₹4,99.99" formatting) belong in the dashboard layer.
package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Currency is a thin wrapper around an ISO-4217 code so the type
// system can distinguish e.g. INR from USD. Currency-aware helpers
// (FromMajor, Major, ParseFor, FormatFor) honour per-currency
// precision; the legacy INR-named helpers (FromRupees, Rupees,
// ParseString, String) delegate to the currency-aware ones with INR.
type Currency string

const (
	CurrencyINR Currency = "INR"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencySGD Currency = "SGD"
	CurrencyAUD Currency = "AUD"
	CurrencyCAD Currency = "CAD"
	CurrencyAED Currency = "AED"
	CurrencyJPY Currency = "JPY" // zero-decimal
	CurrencyKWD Currency = "KWD" // three-decimal
	CurrencyBHD Currency = "BHD" // three-decimal
	CurrencyOMR Currency = "OMR" // three-decimal
	CurrencyTND Currency = "TND" // three-decimal
)

// minorUnitsByCurrency maps an ISO-4217 code to the number of minor
// units that make up one major unit. Two-decimal is the common case;
// JPY is zero-decimal; the Gulf dinars and TND are three-decimal. Add
// new currencies here before going live in a new region.
var minorUnitsByCurrency = map[Currency]int64{
	CurrencyINR: 100,
	CurrencyUSD: 100,
	CurrencyEUR: 100,
	CurrencyGBP: 100,
	CurrencySGD: 100,
	CurrencyAUD: 100,
	CurrencyCAD: 100,
	CurrencyAED: 100,
	CurrencyJPY: 1,
	CurrencyKWD: 1000,
	CurrencyBHD: 1000,
	CurrencyOMR: 1000,
	CurrencyTND: 1000,
}

// MinorUnitsPerMajor reports how many minor units make up one major
// unit for this currency. Returns 100 (the dominant case) for any
// unregistered code rather than panicking — but that's a fallback,
// not a feature: register new currencies in minorUnitsByCurrency
// before the first real transaction.
func (c Currency) MinorUnitsPerMajor() int64 {
	if v, ok := minorUnitsByCurrency[c]; ok {
		return v
	}
	return 100
}

// IsValid reports whether the currency is registered with a known
// minor-unit precision. Use at the API/schema boundary to reject
// codes the system can't safely arithmetic on.
func (c Currency) IsValid() bool {
	_, ok := minorUnitsByCurrency[c]
	return ok
}

// Money is a typed integer count of the smallest unit of the assumed
// currency. The default-zero value is a legitimate "zero money" — no
// special "unset" sentinel.
type Money int64

// FromPaise constructs a Money from a minor-unit count. Named
// "Paise" historically; works for any minor unit (cents, yen, fils).
func FromPaise(paise int64) Money { return Money(paise) }

// MarshalJSON emits Money as a plain JSON integer (minor units).
// Explicit so the contract is documented and paired with UnmarshalJSON.
func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(m), 10)), nil
}

// UnmarshalJSON accepts three input shapes for backwards compatibility
// with publishers that haven't migrated to integer minor units:
//
//   - integer (canonical): treated as minor units as-is. e.g. 16949 → 16949 paise.
//   - float (legacy): treated as INR major units and converted via FromMajor.
//     e.g. 169.49 → 16949 paise. Assumes INR because legacy publishers in
//     this codebase only ever sent INR amounts; non-INR float payloads
//     will be wrong and should migrate to the integer form.
//   - string: parsed via ParseString (INR) — covers older string payloads
//     used by some payment-gateway adapters.
//   - null: treated as zero Money.
//
// Once all publishers emit integer minor units, the float branch becomes
// dead code but is kept as a defensive fallback.
func (m *Money) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s == "" || s == "null" {
		*m = 0
		return nil
	}
	if s[0] == '"' {
		var raw string
		if err := json.Unmarshal(data, &raw); err != nil {
			return fmt.Errorf("money: invalid string: %w", err)
		}
		parsed, err := ParseString(raw)
		if err != nil {
			return fmt.Errorf("money: %w", err)
		}
		*m = parsed
		return nil
	}
	if !strings.ContainsAny(s, ".eE") {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("money: invalid integer: %w", err)
		}
		*m = Money(n)
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("money: invalid number: %w", err)
	}
	*m = FromMajor(CurrencyINR, f)
	return nil
}

// MarshalBSONValue emits Money as a BSON Int64 (minor units), mirroring
// the JSON contract.
func (m Money) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(int64(m))
}

// UnmarshalBSONValue accepts BSON int32, int64, double, or string for
// the same backwards-compat reasons as UnmarshalJSON. Documents written
// by legacy publishers that stored basePrice as a BSON double (e.g.
// 169.49) decode as INR major units; canonical int64 documents pass
// through unchanged.
func (m *Money) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	rv := bson.RawValue{Type: t, Value: data}
	switch t {
	case bsontype.Null, bsontype.Undefined:
		*m = 0
		return nil
	case bsontype.Int64:
		*m = Money(rv.Int64())
		return nil
	case bsontype.Int32:
		*m = Money(int64(rv.Int32()))
		return nil
	case bsontype.Double:
		*m = FromMajor(CurrencyINR, rv.Double())
		return nil
	case bsontype.String:
		parsed, err := ParseString(rv.StringValue())
		if err != nil {
			return fmt.Errorf("money: %w", err)
		}
		*m = parsed
		return nil
	default:
		return fmt.Errorf("money: unsupported BSON type %v", t)
	}
}

// FromMajor converts a major-unit value (rupees, dollars, dinars) to
// Money using the currency's minor-unit factor. Banker's rounding to
// keep aggregate sums consistent with downstream financial reporting.
func FromMajor(c Currency, amount float64) Money {
	factor := float64(c.MinorUnitsPerMajor())
	return Money(int64(math.RoundToEven(amount * factor)))
}

// FromRupees converts a rupee value to Money. Kept as a wrapper for
// backwards compatibility; new code should prefer FromMajor.
func FromRupees(rupees float64) Money { return FromMajor(CurrencyINR, rupees) }

// MustFromString parses a decimal string into Money (assumes INR) and
// panics on bad input. Use in seeds or tests where the input is
// hardcoded; production callers should use ParseFor or ParseString.
func MustFromString(s string) Money {
	m, err := ParseString(s)
	if err != nil {
		panic(err)
	}
	return m
}

// ParseFor parses a decimal string into Money using the currency's
// minor-unit precision. Accepts "499.99" / "499" / "0.5" / signed
// forms. Rejects strings with more fractional digits than the
// currency allows — paise/cents/fils is the precision floor; anything
// finer is operator error and silent rounding would lose audit trail.
func ParseFor(c Currency, s string) (Money, error) {
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
	factor := c.MinorUnitsPerMajor()
	maxFrac := decimalsForFactor(factor)
	var frac int64
	if len(parts) == 2 {
		fracStr := parts[1]
		if maxFrac == 0 {
			return 0, fmt.Errorf("money: %s does not allow fractional digits, got %q", c, s)
		}
		if len(fracStr) > maxFrac {
			return 0, fmt.Errorf("money: more than %d decimal places in %q for %s", maxFrac, s, c)
		}
		for len(fracStr) < maxFrac {
			fracStr += "0"
		}
		frac, err = strconv.ParseInt(fracStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("money: invalid fractional part %q: %w", fracStr, err)
		}
	}
	total := whole*factor + frac
	if negative {
		total = -total
	}
	return Money(total), nil
}

// ParseString parses an INR decimal string. Backwards-compatible
// wrapper around ParseFor; new code should pass an explicit currency.
func ParseString(s string) (Money, error) { return ParseFor(CurrencyINR, s) }

// decimalsForFactor returns the count of trailing zeros in factor,
// i.e. how many decimal digits the currency permits. 100 → 2,
// 1 → 0, 1000 → 3. Defined here rather than as a Currency method
// to keep the precision derivation in one place.
func decimalsForFactor(factor int64) int {
	d := 0
	for f := factor; f > 1; f /= 10 {
		d++
	}
	return d
}

// Paise returns the underlying minor-unit count. Named "Paise" for
// historical reasons; the value is whatever minor unit the carrying
// model declares (paise/cents/yen/fils).
func (m Money) Paise() int64 { return int64(m) }

// Major returns the major-unit float for the given currency. Display
// only — round-tripping through float reintroduces the very precision
// loss this type exists to prevent.
func (m Money) Major(c Currency) float64 {
	return float64(m) / float64(c.MinorUnitsPerMajor())
}

// Rupees returns the rupee equivalent. Wrapper around Major(CurrencyINR)
// kept for backwards compatibility; new code should use Major.
func (m Money) Rupees() float64 { return m.Major(CurrencyINR) }

// IsZero reports whether the value is exactly zero minor units.
func (m Money) IsZero() bool { return m == 0 }

// IsNegative reports whether the value is below zero. Used by schema
// validators that need to reject negative values without dropping
// down to int64.
func (m Money) IsNegative() bool { return m < 0 }

// Add returns m + other. Caller is responsible for ensuring both
// values are in the same currency — in practice the parent model
// (Invoice.Currency, Payment.Currency) pins it. For boundaries where
// currency comes from two independent sources, use AddSafe.
func (m Money) Add(other Money) Money { return m + other }

// Sub returns m - other. Same currency assumption as Add.
func (m Money) Sub(other Money) Money { return m - other }

// AddSafe returns m+other only if mc==oc. Use when currency comes
// from two independent sources (merging a refund into an invoice,
// reconciling a webhook payment against an internal record). Returns
// an error rather than producing a nonsense paise+yen sum.
func (m Money) AddSafe(mc Currency, other Money, oc Currency) (Money, error) {
	if mc != oc {
		return 0, fmt.Errorf("money: cannot add %s and %s", mc, oc)
	}
	return m + other, nil
}

// SubSafe returns m-other only if mc==oc. Same rationale as AddSafe.
func (m Money) SubSafe(mc Currency, other Money, oc Currency) (Money, error) {
	if mc != oc {
		return 0, fmt.Errorf("money: cannot subtract %s from %s", oc, mc)
	}
	return m - other, nil
}

// MulInt scales the amount by an integer multiplier (e.g. quantity).
// Use this rather than (m * Money(qty)) to make the unit-less scalar
// explicit at the call site. No overflow check — see MulIntChecked
// when the multiplier comes from untrusted input.
func (m Money) MulInt(n int64) Money { return Money(int64(m) * n) }

// MulIntChecked scales m by n and returns an error if the product
// would overflow int64. Use when n is caller-supplied (API quantity,
// CSV import) rather than internally derived.
func (m Money) MulIntChecked(n int64) (Money, error) {
	if m == 0 || n == 0 {
		return 0, nil
	}
	a := int64(m)
	prod := a * n
	if prod/n != a {
		return 0, fmt.Errorf("money: int64 overflow multiplying %d by %d", a, n)
	}
	return Money(prod), nil
}

// FormatFor renders "<code> <major>.<minor>" using the currency's
// precision: "INR 499.99", "JPY 1500", "KWD 1.234". For UI formatting
// (locale-aware ₹ symbol, thousands separators) defer to the
// dashboard's display layer — this method exists for log lines and
// error messages, not customer invoices.
func (m Money) FormatFor(c Currency) string {
	decimals := decimalsForFactor(c.MinorUnitsPerMajor())
	return fmt.Sprintf("%s %.*f", c, decimals, m.Major(c))
}

// String renders as "INR <rupees>.<paise>" — backwards-compatible
// formatting for log lines that pre-date currency-aware code. New
// code should prefer FormatFor with an explicit currency.
func (m Money) String() string { return m.FormatFor(CurrencyINR) }

// Equal reports minor-unit equality. Defined explicitly so a future
// switch to a struct-backed Money (e.g. carrying a Currency tag)
// doesn't silently change comparison semantics.
func (m Money) Equal(other Money) bool { return m == other }

// Converter exchanges Money between currencies. Implementations apply
// the rate at a defined point in time and round to the target
// currency's minor-unit precision. The package ships only the
// FixedRateConverter (deterministic, for tests and operator-locked
// rates); production services wire their own provider (RBI feed,
// FX vendor, etc.) behind this interface.
type Converter interface {
	Convert(m Money, from, to Currency) (Money, error)
}

// FixedRateConverter is a deterministic Converter for tests and
// configs where the rate is operator-supplied (e.g. a tenant-locked
// rate for a quarter). Not for production market quotes.
type FixedRateConverter struct {
	// Rates[from][to] = how many "to" major units per one "from" major unit.
	Rates map[Currency]map[Currency]float64
}

// Convert exchanges m from→to using the configured rate. Returns an
// error if no rate is configured. Same-currency conversion is a no-op.
func (f FixedRateConverter) Convert(m Money, from, to Currency) (Money, error) {
	if from == to {
		return m, nil
	}
	row, ok := f.Rates[from]
	if !ok {
		return 0, fmt.Errorf("money: no rates configured for %s", from)
	}
	rate, ok := row[to]
	if !ok {
		return 0, fmt.Errorf("money: no rate %s→%s configured", from, to)
	}
	major := m.Major(from) * rate
	return FromMajor(to, major), nil
}
