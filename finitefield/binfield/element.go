package binfield

import (
	"math/bits"
	"math/rand"
	"regexp"
	"strings"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"strconv"
)

// Ensure that elements in binary fields satisfy the ff.Element interface
var _ ff.Element = &Element{}

// Element is the implementation of an element in a finite field.
type Element struct {
	field *Field
	val   uint
	err   error
}

// Zero returns the additive identity in f.
func (f *Field) Zero() ff.Element {
	return &Element{
		field: f,
		val:   0,
	}
}

// One returns the multiplicative identity in f.
func (f *Field) One() ff.Element {
	return &Element{
		field: f,
		val:   1,
	}
}

// RandElement returns a pseudo-random element in f.
//
// This function uses the default source from the math/rand package. The seed is
// set automatically when loading the primefield package, but a new seed can be
// set by calling rand.Seed().
//
// The pseudo-random generator used is not cryptographically safe.
func (f *Field) RandElement() ff.Element {
	prg := func() uint {
		return uint(rand.Uint64())
	}
	if bits.UintSize == 32 {
		prg = func() uint {
			return uint(rand.Uint32())
		}
	}

	return &Element{
		field: f,
		val:   prg() % (1 << f.extDeg), // Discard anything but the first extDeg bits
	}
}

// Element defines a new element over f with value val, which must be either
// uint, int, or string.
//
// If type of val is unsupported, the function returns an Input-error.
func (f *Field) Element(val interface{}) (ff.Element, error) {
	const op = "Defining element"

	switch v := val.(type) {
	case uint:
		return f.ElementFromUnsigned(v), nil
	case int:
		return f.ElementFromSigned(v), nil
	case string:
		return f.ElementFromString(v)
	default:
		return nil, errors.New(
			op, errors.Input,
			"Cannot define element in %v from type %T", f, v,
		)
	}
}

// ElementFromBits defines a new element over f with value specified by the
// bitstring val. That is, the i'th bit in val determines the coefficient of the
// i'th term in the representation of the element.
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) ElementFromBits(val uint) ff.Element {
	a := &Element{
		field: f,
		val:   val,
	}
	a.reduce()
	return a
}

// ElementFromUnsigned defines a new element over f with value specified by val.
//
// The returned element will automatically be reduced modulo the characteristic.
// That is, the additive identity is returned if val is even, and the
// multiplicative identity is returned if val is odd.
func (f *Field) ElementFromUnsigned(val uint) ff.Element {
	return &Element{
		field: f,
		val:   val % 2,
	}
}

// ElementFromSigned defines a new element over f with value val.
//
// The returned element will be reduced modulo the characteristic automatically.
// That is, the additive identity is returned if val is even, and the
// multiplicative identity is returned if val is odd.
func (f *Field) ElementFromSigned(val int) ff.Element {
	val %= 2
	if val < 0 {
		val += 2
	}
	return &Element{
		field: f,
		val:   uint(val),
	}
}

// ElementFromString defines a new element over f from the given string.
//
// A Parsing-error is returned if the string cannot be parsed.
func (f *Field) ElementFromString(s string) (ff.Element, error) {
	const op = "Defining element from string"

	pattern, err := regexp.Compile(
		`\s*(?:^|\+|-)\s*` + // A sign
			`(` + // Consider two options:
			`(?:0|1)` + // Option 1: A constant coefficient
			`|` + // or
			regexp.QuoteMeta(f.varName) + // Option 2: The variable name
			`(?:\^?([0-9]+))?` + // followed by an optional exponent
			`)\s*`,
	)
	if err != nil {
		return nil, errors.New(
			op, errors.InputValue,
			"Cannot construct regular expression with variable name %q. "+
				"Received error %q.",
			f.varName, err,
		)
	}

	matches := pattern.FindAllStringSubmatch(s, -1)

	// Check that total match length is the full input string
	matchLen := 0
	for _, m := range matches {
		matchLen += len(m[0])
	}
	if matchLen != len(s) {
		return nil, errors.New(
			op, errors.Parsing,
			"Cannot parse %s; lengths do not match (%d ≠ %d).",
			s, matchLen, len(s),
		)
	}

	val := uint(0)
	for _, m := range matches {
		switch m[1] {
		case "0":
			// The term is zero
			continue
		case "1":
			// The term contains no variable. Add one to constant term
			val ^= 1
			continue
		}

		var deg uint
		if m[2] == "" {
			// No degree was given. Implicit degree is one
			deg = 1
		} else {
			// Convert the given degree
			tmp, err := strconv.ParseUint(m[2], 10, 0)
			if err != nil {
				return nil, errors.New(
					op, errors.InputValue,
					"Failed to parse exponent %q as unsigned integer",
					m[2],
				)
			}
			deg = uint(tmp)

			// Ensure that the degree is less than the uint size
			if deg >= bits.UintSize {
				return nil, errors.New(
					op, errors.InputTooLarge,
					"Degree %d exceeds the maximal value (%d)",
					deg, bits.UintSize-1,
				)
			}
		}

		// Add one to the term with given degree
		val ^= (1 << deg)
	}

	return f.ElementFromBits(val), nil
}

// Copy returns a copy of a.
func (a *Element) Copy() ff.Element {
	return &Element{
		field: a.field,
		val:   a.val,
		err:   a.err,
	}
}

// Err returns the error status of a.
func (a *Element) Err() error {
	return a.err
}

// SetUnsigned sets the value of a to the element corresponding to val. It then
// returns a.
//
// The value is automatically reduced modulo the characteristic.
func (a *Element) SetUnsigned(val uint) ff.Element {
	a.val = val % 2
	return a
}

// Equal tests equality of elements a and b.
func (a *Element) Equal(b ff.Element) bool {
	bb, ok := b.(*Element)
	if !ok {
		return false
	}

	if a.field == bb.field && a.val == bb.val {
		return true
	}
	return false
}

// IsZero returns a boolean describing whether a is the additive identity.
func (a *Element) IsZero() bool {
	return a.val == 0
}

// IsNonzero returns a boolean describing whether a is a non-zero element.
func (a *Element) IsNonzero() bool {
	return a.val != 0
}

// IsOne returns a boolean describing whether a is the multiplicative identity.
func (a *Element) IsOne() bool {
	return a.val == 1
}

// String returns the string representation of a.
func (a *Element) String() string {
	if a.IsZero() {
		return "0"
	}

	var b strings.Builder

	nPlus := a.NTerms() - 1
	for term, d := uint(1)<<a.field.extDeg, a.field.extDeg; term > 0; term, d = term>>1, d-1 {
		if term&a.val == 0 {
			continue
		}

		if d == 0 {
			b.WriteByte('1') // Always returns nil error
			break
		}

		b.WriteString(a.field.varName) // Always returns nil error
		if d > 1 {
			b.WriteByte('^')
			b.WriteString(strconv.FormatUint(uint64(d), 10))
		}

		// Print a plus if necessary
		if nPlus > 0 {
			b.Write([]byte(" + "))
			nPlus--
		}
	}

	return b.String()
}

// NTerms returns the number of terms in the representation of a.
func (a *Element) NTerms() uint {
	return uint(bits.OnesCount(a.val))
}
