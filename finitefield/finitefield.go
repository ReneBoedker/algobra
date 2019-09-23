package finitefield

import (
	"algobra/errors"
	"algobra/primefield"
)

// Field is the implementation of a generic finite field.
type Field struct {
	pf *primefield.Field
	//Non-prime field type at some point
}

// Element is the implementation of a finite field element.
type Element struct {
	pf *primefield.Element
	// Non-prime field type at some point
	err error
}

type kind uint8

const (
	primeKind = iota
	nonPrimeKind
)

// Kind is an internal method to easily distinguish between fields of different
// types.
func (f *Field) kind() kind {
	switch {
	case f.pf != nil:
		return primeKind
	default:
		panic("Error")
	}
}

// Kind is an internal method to easily distinguish between field elements of
// different types.
func (a *Element) kind() kind {
	switch {
	case a.pf != nil:
		return primeKind
	default:
		panic("Error")
	}
}

// Define defines a finite field with the given cardinality
//
// Currently, only prime cardinality is supported
func Define(card uint) (*Field, error) {
	// Must be changed when non-prime is implemented
	pf, err := primefield.Define(card)
	if err != nil {
		return nil, err
	}
	return &Field{
		pf: pf,
	}, nil
}

// ElementFromUnsigned returns a new element with value specified by val. That
// is, val will be reduced modulo the characteristic of f.
func (f *Field) ElementFromUnsigned(val uint) *Element {
	switch {
	case f.pf != nil:
		return &Element{
			pf: f.pf.Element(val),
		}
	default:
		panic("Error")
	}
}

// ElementFromSigned returns a new element with value specified by val.
//
// The element is automatically reduced modulo the characteristic. Negative
// values are reduced to a positive remainder (as opposed to the %-operator in
// Go).
func (f *Field) ElementFromSigned(val int) *Element {
	switch {
	case f.pf != nil:
		return &Element{
			pf: f.pf.ElementFromSigned(val),
		}
	default:
		panic("Error")
	}
}

// Element is a general method for defining field element.
//
// Based on the type of val, it calls the appropriate element constructor. It
// returns an Input-error if the type of val is not supported, and an
// InputValue-error if a val is a slice and f is a prime field.
func (f *Field) Element(val interface{}) (*Element, error) {
	const op = "Defining field element"

	switch v := val.(type) {
	case uint:
		return f.ElementFromUnsigned(v), nil
	case int:
		return f.ElementFromSigned(v), nil
	default:
		return nil, errors.New(
			op, errors.Input,
			"Cannot create element from type %T", v,
		)
	}
}

// Elements returns a slice containing all elements of f.
func (f *Field) MultGenerator() *Element {
	switch f.kind() {
	case primeKind:
		return &Element{
			pf: f.pf.MultGenerator(),
		}
	default:
		panic("Error")
	}
}

// Elements returns a slice containing all elements of f.
func (f *Field) Elements() []*Element {
	out := make([]*Element, f.Card(), f.Card())
	switch f.kind() {
	case primeKind:
		for i, e := range f.pf.Elements() {
			out[i] = &Element{
				pf: e,
			}
		}
	default:
		panic("Error")
	}
	return out
}

// Zero returns the additive identity of field f.
func (f *Field) Zero() *Element {
	switch {
	case f.pf != nil:
		return &Element{
			pf: f.pf.Element(0),
		}
	default:
		panic("Error")
	}
}

// One returns the multiplicative identity in field f.
func (f *Field) One() *Element {
	switch {
	case f.pf != nil:
		return &Element{
			pf: f.pf.Element(1),
		}
	default:
		panic("Error")
	}
}

// Card returns the cardinality of f.
func (f *Field) Card() uint {
	switch f.kind() {
	case primeKind:
		return f.pf.Card()
	default:
		panic("Error")
	}
}

// Char returns the characteristic of f.
func (f *Field) Char() uint {
	switch f.kind() {
	case primeKind:
		return f.pf.Char()
	default:
		panic("Error")
	}
}

func (f *Field) ComputeTables(add, mult bool) (err error) {
	switch f.kind() {
	case primeKind:
		return f.pf.ComputeTables(add, mult)
	default:
		panic("Error")
	}
}

// Copy return a new element with the same field and value as a.
func (a *Element) Copy() *Element {
	switch a.kind() {
	case primeKind:
		return &Element{
			pf:  a.pf.Copy(),
			err: a.err,
		}
	default:
		panic("Error")
	}
}

// Plus returns the sum of elements a and b
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Plus(b *Element) *Element {
	return a.Copy().Add(b)
}

// Add set a to the sum of elements a and b and returns a
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Add(b *Element) *Element {
	const op = "Adding elements"

	if a.kind() != b.kind() {
		a.err = errors.New(
			op, errors.ArithmeticIncompat,
			"Cannot add elements from different fields",
		)
	}

	switch a.kind() {
	case primeKind:
		a.pf.Add(b.pf)
	default:
		panic("Error")
	}
	return a
}

// Minus returns the sum of elements a and b
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Minus(b *Element) *Element {
	return a.Copy().Sub(b)
}

// Sub sets a to the sum of elements a and b and returns a
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Sub(b *Element) *Element {
	const op = "Subtracting elements"

	if a.kind() != b.kind() {
		return &Element{
			err: errors.New(
				op, errors.ArithmeticIncompat,
				"Cannot subtract elements from different fields",
			),
		}
	}

	switch a.kind() {
	case primeKind:
		a.pf.Sub(b.pf)
	default:
		panic("Error")
	}
	return a
}

// Times returns the sum of elements a and b
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Times(b *Element) *Element {
	return a.Copy().Mult(b)
}

// Mult set a to the sum of elements a and b and returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Mult(b *Element) *Element {
	const op = "Multiplying elements"

	if a.kind() != b.kind() {
		return &Element{
			err: errors.New(
				op, errors.ArithmeticIncompat,
				"Cannot multiply elements from different fields",
			),
		}
	}

	switch a.kind() {
	case primeKind:
		a.pf.Mult(b.pf)
	default:
		panic("Error")
	}
	return a
}

func (a *Element) Prod(b, c *Element) {
	const op = "Multiplying elements"

	if a.kind() != b.kind() && a.kind() != c.kind() {
		a.err = errors.New(
			op, errors.ArithmeticIncompat,
			"Cannot multiply elements from different fields",
		)
	}

	switch a.kind() {
	case primeKind:
		a.pf.Prod(b.pf, c.pf)
	default:
		panic("Error")
	}
}

// Equal tests equality of elements a and b.
func (a *Element) Equal(b *Element) bool {
	if a.kind() != b.kind() {
		return false
	}

	switch a.kind() {
	case primeKind:
		return a.pf.Equal(b.pf)
	default:
		panic("Error")
	}
}

// Inv returns the inverse of a
//
// If a is the zero element, the return value is an element with
// InputValue-error as error status.
func (a *Element) Inv() *Element {
	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Inv(),
		}
	default:
		panic("Error")
	}
}

// Neg returns -a (modulo the characteristic)
func (a *Element) Neg() *Element {
	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Neg(),
		}
	default:
		panic("Error")
	}
}

// Pow returns a raised to the power of n
func (a *Element) Pow(n uint) *Element {
	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Pow(n),
		}
	default:
		panic("Error")
	}
}

// Zero returns a boolean describing whether a is the additive identity.
func (a *Element) Zero() bool {
	switch a.kind() {
	case primeKind:
		return a.pf.Zero()
	default:
		panic("Error")
	}
}

// Nonzero returns a boolean describing whether a is a non-zero element
func (a *Element) Nonzero() bool {
	switch a.kind() {
	case primeKind:
		return a.pf.Nonzero()
	default:
		panic("Error")
	}
}

// One returns a boolean describing whether a is one
func (a *Element) One() bool {
	switch a.kind() {
	case primeKind:
		return a.pf.One()
	default:
		panic("Error")
	}
}

// Err returns the error status of a.
func (a *Element) Err() error {
	if a.err != nil {
		return a.err
	}
	switch a.kind() {
	case primeKind:
		return a.pf.Err()
	default:
		panic("Error")
	}
}

// String returns the string representation of a.
func (a *Element) String() string {
	switch a.kind() {
	case primeKind:
		return a.pf.String()
	default:
		panic("Error")
	}
}
