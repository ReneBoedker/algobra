package extfield

import (
	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

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

	a.val.Add(bb.val)

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
	const op = "Subtracting elements"

	bb, ok := b.(*Element)
	if !ok {
		a.err = errors.New(
			op, errors.InputIncompatible,
			"Cannot subtract %v (%[1]T) from %v (%[2]T)", b, a,
		)
	}

	if tmp := checkErrAndCompatible(op, a, bb); tmp != nil {
		a = tmp
		return a
	}

	a.val.Sub(bb.val)
	return a
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

	if bb.IsZero() || cc.IsZero() {
		a.val.SetZero()
		return a
	}

	if a.field.logTable != nil {
		s := a.field.logTable.lookup(bb)
		t := a.field.logTable.lookup(cc)

		a = a.field.logTable.lookupReverse((s + t) % (a.field.Card() - 1))
	} else {
		a.val = (bb.val.Times(cc.val))
	}
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
	return a.Copy().SetNeg()
}

// SetNeg sets a to a scaled by negative one (modulo the characteristic). It
// then returns a.
func (a *Element) SetNeg() ff.Element {
	a.val.SetNeg()
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

	if a.field.logTable != nil {
		s := a.field.logTable.lookup(a)
		return a.field.logTable.lookupReverse(a.field.Card() - 1 - s)
	}

	// Implemented using the extended euclidean algorithm (see for instance
	// [GG13; Algorithm 3.14])
	r0 := a.field.conwayPoly.Normalize()
	r0.EmbedIn(a.field.polyRing, false)
	r1 := a.val.Normalize()

	i0 := a.field.polyRing.Zero()
	i1 := a.field.polyRing.Polynomial([]ff.Element{a.val.Lc().Inv()})
	for true {
		quo, rem, err := r0.QuoRem(r1)
		if err != nil {
			o := a.field.Zero()
			out := o.(*Element)
			out.err = err
			return out
		}

		if rem.IsZero() {
			break
		}

		lcInv := rem.Lc().Inv()
		r0, r1 = r1, rem.SetScale(lcInv)
		i0, i1 = i1, i0.Sub(quo[0].Mult(i1)).SetScale(lcInv)
	}
	return &Element{
		field: a.field,
		val:   i1,
	}
}

// Trace computes the field trace of a.
//
// The resulting object is considered as an element of the extension field (i.e.
// the same field as a).
func (a *Element) Trace() ff.Element {
	out := a.Copy()

	for i := uint(1); i < a.field.extDeg; i++ {
		out = out.Pow(a.field.Char())
		out.Add(a)
	}

	return out
}

/* Copyright 2019 René Bødker Christensen
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
