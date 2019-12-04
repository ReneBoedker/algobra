package extfield

import (
	"math/bits"
	"math/rand"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/univariate"
)

// Ensure that Element implements the ff.Element interface
var _ ff.Element = &Element{}

// Element is the implementation of an element in a finite field.
type Element struct {
	field *Field
	val   *univariate.Polynomial
	err   error
}

// Zero returns the additive identity in f.
func (f *Field) Zero() ff.Element {
	return &Element{
		field: f,
		val:   f.polyRing.Zero(),
	}
}

// One returns the multiplicative identity in f.
func (f *Field) One() ff.Element {
	return &Element{
		field: f,
		val:   f.polyRing.One(),
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

	coefs := make([]uint, f.extDeg, f.extDeg)
	for i := range coefs {
		coefs[i] = prg()
	}

	return f.ElementFromUnsignedSlice(coefs)
}

// Element defines a new element over f with value val, which must be either
// uint, int, []uint, []int, or string.
//
// If type of val is unsupported, the function returns an Input-error.
func (f *Field) Element(val interface{}) (ff.Element, error) {
	const op = "Defining element"

	switch v := val.(type) {
	case uint:
		return f.ElementFromUnsigned(v), nil
	case int:
		return f.ElementFromSigned(v), nil
	case []uint:
		return f.ElementFromUnsignedSlice(v), nil
	case []int:
		return f.ElementFromSignedSlice(v), nil
	case string:
		return f.ElementFromString(v)
	default:
		return nil, errors.New(
			op, errors.Input,
			"Cannot define element in %v from type %T", f, v,
		)
	}
}

// element defines a new element over f with value specified by val.
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) element(val []uint) *Element {
	return &Element{
		field: f,
		val:   f.polyRing.PolynomialFromUnsigned(val),
	}
}

// ElementFromUnsigned defines a new element over f with value specified by val.
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) ElementFromUnsigned(val uint) ff.Element {
	return &Element{
		field: f,
		val:   f.polyRing.PolynomialFromUnsigned([]uint{val}),
	}
}

// ElementFromUnsignedSlice defines a new element over f with value specified by val.
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) ElementFromUnsignedSlice(val []uint) ff.Element {
	return &Element{
		field: f,
		val:   f.polyRing.PolynomialFromUnsigned(val),
	}
}

// ElementFromSigned defines a new element over f with values specified by val.
//
// The returned element will be reduced modulo the characteristic automatically.
// Negative values are reduced to a positive remainder (as opposed to the
// %-operator in Go).
func (f *Field) ElementFromSigned(val int) ff.Element {
	return &Element{
		field: f,
		val:   f.polyRing.PolynomialFromSigned([]int{val}),
	}
}

// ElementFromSignedSlice defines a new element over f with values specified by val.
//
// The returned element will be reduced modulo the characteristic automatically.
// Negative values are reduced to a positive remainder (as opposed to the
// %-operator in Go).
func (f *Field) ElementFromSignedSlice(val []int) ff.Element {
	return &Element{
		field: f,
		val:   f.polyRing.PolynomialFromSigned(val),
	}
}

// ElementFromString defines a new element over f from the given string.
//
// A Parsing-error is returned if the string cannot be parsed.
func (f *Field) ElementFromString(val string) (ff.Element, error) {
	const op = "Defining element from string"

	v, err := f.polyRing.PolynomialFromString(val)
	if err != nil {
		return nil, errors.Wrap(op, errors.Parsing, err)
	}

	return &Element{
		field: f,
		val:   v,
	}, nil
}

// Copy returns a copy of a.
func (a *Element) Copy() ff.Element {
	return &Element{
		field: a.field,
		val:   a.val.Copy(),
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
	a.val = a.field.polyRing.Polynomial(
		[]ff.Element{a.field.baseField.ElementFromUnsigned(val)},
	)
	return a
}

// Equal tests equality of elements a and b.
func (a *Element) Equal(b ff.Element) bool {
	bb, ok := b.(*Element)
	if !ok {
		return false
	}

	if a.field == bb.field && a.val.Equal(bb.val) {
		return true
	}
	return false
}

// IsZero returns a boolean describing whether a is the additive identity.
func (a *Element) IsZero() bool {
	return a.val.IsZero()
}

// IsNonzero returns a boolean describing whether a is a non-zero element.
func (a *Element) IsNonzero() bool {
	return a.val.IsNonzero()
}

// IsOne returns a boolean describing whether a is the multiplicative identity.
func (a *Element) IsOne() bool {
	return a.val.IsOne()
}

// String returns the string representation of a.
func (a *Element) String() string {
	return a.val.String()
}

// NTerms returns the number of terms in the representation of a.
func (a *Element) NTerms() uint {
	return a.val.NTerms()
}
