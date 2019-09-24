package primefield

import (
	"algobra/basic"
	"algobra/errors"
	"fmt"
	"strconv"
)

// Number of bits in a uint type
var uintBitSize uint

func init() {
	// Determine bit-size of uint type
	i := ^uint(0)
	for i > 0 {
		uintBitSize++
		i >>= 1
	}
}

// Field is the implementation of a finite field
type Field struct {
	char      uint
	addTable  *table
	multTable *table
}

// Define creates a new finite field with prime characteristic.
//
// If char is not a prime, the package returns an InputValue-error. If char
// implies that multiplication will overflow uint, the function returns an
// InputTooLarge-error.
func Define(char uint) (*Field, error) {
	const op = "Defining prime field"

	if char == 0 {
		return nil, errors.New(
			op, errors.InputValue,
			"Field characteristic cannot be zero",
		)
	}

	if char-1 >= 1<<(uintBitSize/2) {
		return nil, errors.New(
			op, errors.InputTooLarge,
			"%d exceeds maximal field size (2^%d)", char, uintBitSize/2,
		)
	}

	_, n, err := basic.FactorizePrimePower(char)
	if err != nil {
		return nil, errors.Wrap(op, errors.Inherit, err)
	} else if n != 1 {
		return nil, errors.New(
			op, errors.InputValue,
			"%d is not a prime", char,
		)
	}
	return &Field{char: char, addTable: nil, multTable: nil}, nil
}

// String returns the string representation of f
func (f *Field) String() string {
	return fmt.Sprintf("Finite field of %d elements", f.char)
}

// Char returns the characteristic of f
func (f *Field) Char() uint {
	return f.char
}

// Card returns the cardinality of f
func (f *Field) Card() uint {
	return f.char
}

// ComputeTables will precompute either the addition or multiplication tables
// (or both) for the field f.
//
// Returns an InputTooLarge-error if the estimated memory usage exceeds the
// maximal value.
func (f *Field) ComputeTables(add, mult bool) (err error) {
	if add && f.addTable == nil {
		f.addTable, err = newTable(f, func(i, j uint) uint {
			return (i + j) % f.char
		})
	}

	if mult && f.multTable == nil {
		f.multTable, err = newTable(f, func(i, j uint) uint {
			return (i * j) % f.char
		})
	}

	if err != nil {
		return err
	}

	return nil
}

// MultGenerator returns an element that generates the units of f.
func (f *Field) MultGenerator() *Element {
	if f.Card() == 2 {
		return f.Element(1)
	}

	// The possible orders of elements divide p-1 (we can ignore errors since
	// input is non-zero)
	factors, _, _ := basic.Factorize(f.Card() - 1)

	var e *Element
outer:
	for i := uint(2); true; i++ {
		e = f.Element(i)
		for _, p := range factors {
			// We need to check if p is a non-trivial factor
			if p != f.Card()-1 && e.Pow(p).One() {
				// Not a generator
				continue outer
			}
		}
		break
	}
	return e
}

// Elements returns a slice containing all elements of f.
// TODO: This can be done much easier since card is prime
func (f *Field) Elements() []*Element {
	out := make([]*Element, f.Card(), f.Card())
	out[0] = f.Element(0)

	gen := f.MultGenerator()
	for i, e := uint(1), gen.Copy(); i < f.Char(); i, e = i+1, e.Mult(gen) {
		out[i] = e.Copy()
	}
	return out
}

// Element is the implementation of an element in a finite field.
type Element struct {
	field *Field
	val   uint
	err   error
}

// Element defines a new element over f with value val
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) Element(val uint) *Element {
	return &Element{field: f, val: val % f.char}
}

// ElementFromSigned defines a new element over f with value val
//
// The returned element will be reduced modulo the characteristic automatically.
// Negative values are reduced to a positive remainder (as opposed to the
// %-operator in Go).
func (f *Field) ElementFromSigned(val int) *Element {
	val %= int(f.char)
	if val < 0 {
		val += int(f.char)
	}
	return f.Element(uint(val))
}

// Copy returns a copy of a
func (a *Element) Copy() *Element {
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

// Uint returns the value of a represented as a uint.
func (a *Element) Uint() uint {
	return a.val
}

// SetUnsigned sets the value of a to the element corresponding to val.
func (a *Element) SetUnsigned(val uint) {
	a.val = val % a.field.Char()
}

// Equal tests equality of elements a and b.
func (a *Element) Equal(b *Element) bool {
	if a.field == b.field && a.val == b.val {
		return true
	}
	return false
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

// Add sets a to the sum of a and b and returns a
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Add(b *Element) *Element {
	const op = "Adding elements"

	if tmp := hasErr(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	if tmp := checkCompatible(op, a, b); tmp != nil {
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

// Neg returns -a (modulo the characteristic)
func (a *Element) Neg() *Element {
	return a.field.Element(a.field.char - a.val)
}

// NegInPlace sets a to -a (modulo the characteristic), and returns a
func (a *Element) NegInPlace() *Element {
	a.val = a.field.char - a.val
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

// Sub sets a to the difference of elements a and b and returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Sub(b *Element) *Element {
	const op = "Subtracting elements"

	if tmp := hasErr(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	if tmp := checkCompatible(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	a.val = (a.field.Char() - a.val) % a.field.Char()
	a.Add(b)
	a.val = (a.field.Char() - a.val) % a.field.Char()
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
	const op = "Multiplying elements"

	if tmp := hasErr(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	if tmp := checkCompatible(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	if a.field.multTable != nil {
		a.val = a.field.multTable.lookup(a.val, b.val)
	} else {
		a.val = (a.val * b.val) % a.field.Char()
	}

	return a
}

// Prod sets a to the product of b and c, and returns a.
//
// The function returns an ArithmeticIncompat-error if b, and c are not defined
// over the same field.
func (a *Element) Prod(b, c *Element) {
	const op = "Multiplying elements"

	if tmp := hasErr(op, b, c); tmp != nil {
		a = tmp
	}

	if tmp := checkCompatible(op, b, c); tmp != nil {
		a = tmp
	}

	// Set the correct field of a
	a.field = b.field

	if a.field.multTable != nil {
		a.val = a.field.multTable.lookup(b.val, c.val)
	} else {
		a.val = (b.val * c.val) % a.field.Char()
	}
}

// Pow returns a raised to the power of n
func (a *Element) Pow(n uint) *Element {
	if a.Zero() {
		if n == 0 {
			return a.field.Element(1)
		}
		return a.field.Element(0)
	}

	if n >= a.field.Char() {
		// Use that a^(p-1)=1 for units
		n = n % (a.field.Char() - 1)
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

// Zero returns a boolean describing whether a is the zero element
func (a *Element) Zero() bool {
	return (a.val == 0)
}

// Nonzero returns a boolean describing whether a is a non-zero element
func (a *Element) Nonzero() bool {
	return (a.val != 0)
}

// One returns a boolean describing whether a is one
func (a *Element) One() bool {
	return (a.val == 1)
}

// String returns the string representation of a.
func (a *Element) String() string {
	return strconv.FormatUint(uint64(a.val), 10)
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
		out := a.field.Element(0)
		out.err = errors.New(
			op, errors.ArithmeticIncompat,
			"%v and %v defined over different fields", a, b,
		)
		return out
	}
	return nil
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
