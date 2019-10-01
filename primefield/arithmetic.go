package primefield

import (
	"algobra/errors"
)

// Add sets a to the sum of a and b and returns a
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Add(b *Element) *Element {
	const op = "Adding elements"

	if tmp := checkErrAndCompatible(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	if a.field.addTable != nil {
		a.val = a.field.addTable.lookup(a.val, b.val)
	} else {
		a.val = (a.val + b.val) % a.field.Char()
	}

	return a
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

// Sub sets a to the difference of elements a and b and returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Sub(b *Element) *Element {
	const op = "Subtracting elements"

	if tmp := checkErrAndCompatible(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	a.val = (a.field.Char() - a.val) % a.field.Char()
	a.Add(b)
	a.val = (a.field.Char() - a.val) % a.field.Char()
	return a
}

// Minus returns the difference of elements a and b
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Minus(b *Element) *Element {
	return a.Copy().Sub(b)
}

// Prod sets a to the product of b and c, and returns a.
//
// The function returns an ArithmeticIncompat-error if b, and c are not defined
// over the same field.
func (a *Element) Prod(b, c *Element) *Element {
	const op = "Multiplying elements"

	if tmp := checkErrAndCompatible(op, b, c); tmp != nil {
		a = tmp
	}

	// Set the correct field of a
	a.field = b.field

	if a.field.multTable != nil {
		a.val = a.field.multTable.lookup(b.val, c.val)
	} else {
		a.val = (b.val * c.val) % a.field.Char()
	}
	return a
}

// Times returns the product of elements a and b
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Times(b *Element) *Element {
	return a.Copy().Mult(b)
}

// Mult sets a to the product of elements a and b and returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Mult(b *Element) *Element {
	return a.Prod(a, b)
}

// Neg returns -a (modulo the characteristic)
func (a *Element) Neg() *Element {
	return a.Copy().SetNeg()
}

// SetNeg sets a to -a (modulo the characteristic), and returns a
func (a *Element) SetNeg() *Element {
	a.val = a.field.char - a.val
	return a
}

// Pow returns a raised to the power of n.
func (a *Element) Pow(n uint) *Element {
	if a.IsZero() {
		if n == 0 {
			return a.field.Element(1)
		}
		return a.field.Element(0)
	}

	if n >= a.field.Card() {
		// Use that a^(p-1)=1 for units
		n = n % (a.field.Card() - 1)
	}

	out := a.field.Element(1)
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

// Inv returns the inverse of a
//
// If a is the zero element, the return value is an element with
// InputValue-error as error status.
func (a *Element) Inv() *Element {
	const op = "Inverting element"

	if a.val == 0 {
		out := a.field.Element(0)
		out.err = errors.New(
			op, errors.InputValue,
			"Cannot invert zero element",
		)
		return out
	}

	// Implemented using the extended euclidean algorithm (see for instance
	// [GG13])
	r0 := a.field.char
	r1 := a.val
	i0, i1 := 0, 1
	for r1 > 0 {
		q := r0 / r1
		r0, r1 = r1, r0-q*r1
		i0, i1 = i1, i0-int(q)*i1
	}
	for i0 < 0 {
		i0 += int(a.field.char)
	}
	return a.field.ElementFromSigned(i0)
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
