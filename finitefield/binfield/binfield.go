package binfield

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ReneBoedker/algobra/auxmath"
	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/conway"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

func init() {
	// Set a new seed for the pseudo-random generator.
	// Note that this is not cryptographically safe.
	rand.Seed(time.Now().UTC().UnixNano())
}

// Field is the implementation of a finite field.
type Field struct {
	extDeg     uint
	conwayPoly uint
	varName    string
}

// Define creates a new finite field with given cardinality.
//
// If card is not a power of two, the package returns an InputValue-error. If
// card implies that multiplication will overflow uint, the function returns an
// InputTooLarge-error.
func Define(card uint) (*Field, error) {
	const op = "Defining prime field"

	if card == 0 {
		return nil, errors.New(
			op, errors.InputValue,
			"Field characteristic cannot be zero",
		)
	}

	char, extDeg, err := auxmath.FactorizePrimePower(card)
	if err != nil {
		return nil, errors.Wrap(op, errors.Inherit, err)
	}

	if char != 2 {
		return nil, errors.New(
			op, errors.InputValue,
			"The cardinality of a binary field must be a power of 2",
		)
	}

	conwayCoefs, err := conway.Lookup(2, extDeg)
	if err != nil {
		return nil, errors.Wrap(op, errors.Inherit, err)
	}

	conwayPoly := uint(0)
	for i, c := range conwayCoefs {
		// Fill the i'th entry into the i'th bit
		conwayPoly += c << i
	}

	return &Field{
		extDeg:     extDeg,
		conwayPoly: conwayPoly,
		varName:    "a",
	}, nil
}

// String returns the string representation of f.
func (f *Field) String() string {
	return fmt.Sprintf("Finite field of %d elements", f.Card())
}

// Char returns the characteristic of f.
func (f *Field) Char() uint {
	return 2
}

// Card returns the cardinality of f.
func (f *Field) Card() uint {
	return 1 << f.extDeg
}

// MultGenerator returns an element that generates the units of f.
func (f *Field) MultGenerator() ff.Element {
	// The field is defined from a Conway polynomial, so alpha is a generator
	return f.ElementFromBits(2)
}

// checkErrAndCompatible is a wrapper for the two functions hasErr and
// checkCompatible. It is used in arithmetic functions to check that the inputs
// are 'good' to use.
func checkErrAndCompatible(op errors.Op, a, b *Element) *Element {
	if tmp := hasErr(op, a, b); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, a, b); tmp != nil {
		return tmp
	}

	return nil
}

// hasErr is an internal method for checking if a or b has a non-nil error
// field.
//
// It returns the first element with non-nil error status after wrapping the
// error. The new error inherits the kind from the old.
func hasErr(op errors.Op, a, b *Element) *Element {
	switch {
	case a.err != nil:
		a.err = errors.Wrap(
			op, errors.Inherit,
			a.err,
		)
		return a
	case b.err != nil:
		b.err = errors.Wrap(
			op, errors.Inherit,
			b.err,
		)
		return b
	}
	return nil
}

// checkCompatible is an internal method for checking if a and b are compatible;
// that is, if they are defined over the same field.
//
// If not, the return value is an element with error status set to
// ArithmeticIncompat.
func checkCompatible(op errors.Op, a, b *Element) *Element {
	if a.field != b.field {
		o := a.field.Zero()
		out := o.(*Element)
		out.err = errors.New(
			op, errors.ArithmeticIncompat,
			"%v and %v defined over different fields", a, b,
		)
		return out
	}
	return nil
}
