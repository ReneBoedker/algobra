package binfield

import (
	"math/bits"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

// reduce computes the reduction of a modulo the Conway polynomial of the field,
// and sets a to this value.
func (a *Element) reduce() *Element {
	conwayLen := uint(bits.Len(a.field.conwayPoly))
	for l := uint(bits.Len(a.val)); l > a.field.extDeg; l = uint(bits.Len(a.val)) {
		a.val ^= (a.field.conwayPoly << (l - conwayLen))
	}
	return a
}

// Add sets a to the sum of a and b. It then returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Add(b ff.Element) ff.Element {
	const op = "Adding elements"

	bb, ok := b.(*Element)
	if !ok {
		a.err = errors.New(
			op, errors.InputIncompatible,
			"Cannot add %v (%[1]T) and %v (%[2]T)", a, b,
		)
	}

	if tmp := checkErrAndCompatible(op, a, bb); tmp != nil {
		a = tmp
		return a
	}

	a.val ^= bb.val

	return a
}

// Plus returns the sum of elements a and b.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Plus(b ff.Element) ff.Element {
	return a.Copy().Add(b)
}

// Sub sets a to the difference of elements a and b. It then returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Sub(b ff.Element) ff.Element {
	return a.Add(b)
}

// Minus returns the difference of elements a and b.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Minus(b ff.Element) ff.Element {
	return a.Copy().Sub(b)
}

// Prod sets a to the product of b and c. It then returns a.
//
// The function returns an ArithmeticIncompat-error if b, and c are not defined
// over the same field.
//
// When b or c has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Prod(b, c ff.Element) ff.Element {
	const op = "Multiplying elements"

	bb, okB := b.(*Element)
	cc, okC := c.(*Element)
	if !okB || !okC {
		a.err = errors.New(
			op, errors.InputIncompatible,
			"Cannot set type %T to product of %v (%[1]T) and %v (%[2]T)", a, b, c,
		)
	}

	if tmp := checkErrAndCompatible(op, bb, cc); tmp != nil {
		a = tmp
		return a
	}

	// Set the correct field of a
	a.field = bb.field

	// TODO: Check that no overflow will occur
	res := uint(0)
	for tmp, d := cc.val, 0; tmp > 0; tmp, d = tmp>>1, d+1 {
		res ^= ((tmp & 1) * bb.val) << d
	}

	a.val = res
	a.reduce()

	return a
}

// Times returns the product of elements a and b.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Times(b ff.Element) ff.Element {
	return a.Copy().Mult(b)
}

// Mult sets a to the product of elements a and b. It then returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Mult(b ff.Element) ff.Element {
	return a.Prod(a, b)
}

// Neg returns a scaled by negative one (modulo the characteristic).
func (a *Element) Neg() ff.Element {
	return a.Copy()
}

// SetNeg sets a to a scaled by negative one (modulo the characteristic). It
// then returns a.
func (a *Element) SetNeg() ff.Element {
	return a
}

// Pow returns a raised to the power of n.
func (a *Element) Pow(n uint) ff.Element {
	if a.IsZero() {
		if n == 0 {
			return a.field.One()
		}
		return a.field.Zero()
	}

	if n >= a.field.Card() {
		// Use that a^(q-1)=1 for units
		n = n % (a.field.Card() - 1)
	}

	out := a.field.One()
	b := a.Copy()
	for n > 0 {
		if n%2 == 1 {
			out.Mult(b)
		}
		n /= 2
		b.Mult(b)
	}
	return out
}

// bitQuoRem computes the quotient and remainder of a divided by b when viewed
// as binary polynomials.
func bitQuoRem(a, b uint) (quo, rem uint) {
	l := bits.Len(b)
	for a > 0 {
		if tmp := bits.Len(a); tmp >= l {
			quo ^= 1 << (tmp - l)
			a ^= b << (tmp - l)
		} else {
			break
		}
	}
	return quo, a
}

// bitProd computes the product of a and b when viewed as binary polynomials.
func bitProd(a, b uint) (out uint) {
	for ; b > 0; b >>= 1 {
		out ^= a * (b & 1)
		a <<= 1
	}
	return out
}

// Inv returns the inverse of a.
//
// If a is the zero element, the return value is an element with
// InputValue-error as error status.
func (a *Element) Inv() ff.Element {
	const op = "Inverting element"

	if a.IsZero() {
		o := a.field.Zero()
		out := o.(*Element)
		out.err = errors.New(
			op, errors.InputValue,
			"Cannot invert zero element",
		)
		return out
	}

	if a.IsOne() {
		return a.Copy()
	}

	// Implemented using the extended euclidean algorithm (see for instance
	// [GG13; Algorithm 3.14])
	r0 := a.field.conwayPoly
	r1 := a.val

	i0 := uint(0)
	i1 := uint(1)
	for r1 > 0 {
		quo, rem := bitQuoRem(r0, r1)

		r0, r1 = r1, rem
		i0, i1 = i1, i0^bitProd(i1, quo)
	}

	return (&Element{
		field: a.field,
		val:   i0,
	}).reduce()
}
