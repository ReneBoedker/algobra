package primefield

import (
	"algobra/basic"
	"algobra/errors"
	"algobra/finitefield/ff"
	"fmt"
	"math/bits"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	// Set a new seed for the pseudo-random generator.
	// Note that this is not cryptographically safe.
	rand.Seed(time.Now().UTC().UnixNano())
}

// Field is the implementation of a finite field
type Field struct {
	char      uint
	addTable  *table
	multTable *table
}

// Define creates a new finite field with prime cardinality.
//
// If card is not a prime, the package returns an InputValue-error. If card
// implies that multiplication will overflow uint, the function returns an
// InputTooLarge-error.
func Define(card uint) (*Field, error) {
	const op = "Defining prime field"

	if card == 0 {
		return nil, errors.New(
			op, errors.InputValue,
			"Field characteristic cannot be zero",
		)
	}

	if card-1 >= 1<<(bits.UintSize/2) {
		return nil, errors.New(
			op, errors.InputTooLarge,
			"%d exceeds maximal field size (2^%d)", card, bits.UintSize/2,
		)
	}

	_, n, err := basic.FactorizePrimePower(card)
	if err != nil {
		return nil, errors.Wrap(op, errors.Inherit, err)
	} else if n != 1 {
		return nil, errors.New(
			op, errors.InputValue,
			"%d is not a prime", card,
		)
	}
	return &Field{char: card, addTable: nil, multTable: nil}, nil
}

// String returns the string representation of f.
func (f *Field) String() string {
	return fmt.Sprintf("Finite field of %d elements", f.char)
}

// Char returns the characteristic of f.
func (f *Field) Char() uint {
	return f.char
}

// Card returns the cardinality of f.
func (f *Field) Card() uint {
	return f.char
}

// ComputeTables will precompute either the addition or multiplication tables
// (or both) for the field f.
//
// The optional argument maxMem specifies the maximal table size in KiB. If no
// value is given, a DefaultMaxMem is used. If more than one value is given,
// only the first is used.
//
// Returns an InputTooLarge-error if the estimated memory usage exceeds the
// maximal value specified by maxMem.
func (f *Field) ComputeTables(add, mult bool, maxMem ...uint) (err error) {
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
func (f *Field) MultGenerator() ff.Element {
	if f.Card() == 2 {
		return f.One()
	}

	// The possible orders of elements divide p-1 (we can ignore errors since
	// input is non-zero)
	factors, _, _ := basic.Factorize(f.Card() - 1)

	var e *Element
outer:
	for i := uint(2); true; i++ {
		e = f.element(i)
		for _, p := range factors {
			// We need to check if p is a non-trivial factor
			if p != f.Card()-1 && e.Pow(p).IsOne() {
				// Not a generator
				continue outer
			}
		}
		break
	}
	return e
}

// Elements returns a slice containing all elements of f.
func (f *Field) Elements() []ff.Element {
	out := make([]ff.Element, f.Card(), f.Card())
	for i := uint(0); i < f.Card(); i++ {
		out[i] = f.element(i)
	}
	return out
}

// Element is the implementation of a finite field element.
type Element struct {
	field *Field
	val   uint
	err   error
}

// Zero returns the additive identity in f.
func (f *Field) Zero() ff.Element {
	return &Element{field: f, val: 0}
}

// One returns the multiplicative identity in f.
func (f *Field) One() ff.Element {
	return &Element{field: f, val: 1}
}

// RandElement returns a pseudo-random element in f.
//
// The pseudo-random generator used is not cryptographically safe.
func (f *Field) RandElement() ff.Element {
	if bits.UintSize == 32 {
		return f.ElementFromUnsigned(uint(rand.Uint32()))
	}
	return f.ElementFromUnsigned(uint(rand.Uint64()))
}

// Element defines a new element over f with value val, which must be either
// uint or int.
//
// If type of val is unsupported, the function returns an Input-error.
func (f *Field) Element(val interface{}) (ff.Element, error) {
	const op = "Defining element"

	switch v := val.(type) {
	case uint:
		return f.element(v), nil
	case int:
		return f.ElementFromSigned(v), nil
	default:
		return nil, errors.New(
			op, errors.Input,
			"Cannot define element in %v from type %T", f, v,
		)
	}
}

// element defines a new element over f with value val.
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) element(val uint) *Element {
	return &Element{field: f, val: val % f.char}
}

// ElementFromUnsigned defines a new element over f with value val.
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) ElementFromUnsigned(val uint) ff.Element {
	return f.element(val)
}

// ElementFromSigned defines a new element over f with value val.
//
// The returned element will be reduced modulo the characteristic automatically.
// Negative values are reduced to a positive remainder (as opposed to the
// %-operator in Go).
func (f *Field) ElementFromSigned(val int) ff.Element {
	val %= int(f.char)
	if val < 0 {
		val += int(f.char)
	}
	return f.element(uint(val))
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

// Uint returns the value of a represented as an unsigned integer.
func (a *Element) Uint() uint {
	return a.val
}

// SetUnsigned sets the value of a to the element corresponding to val.
//
// The value is automatically reduced modulo the characteristic.
func (a *Element) SetUnsigned(val uint) {
	a.val = val % a.field.Char()
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
	return (a.val == 0)
}

// IsNonzero returns a boolean describing whether a is a nonzero element.
func (a *Element) IsNonzero() bool {
	return (a.val != 0)
}

// IsOne returns a boolean describing whether a is the multiplicative identity.
func (a *Element) IsOne() bool {
	return (a.val == 1)
}

// String returns the string representation of a.
func (a *Element) String() string {
	return strconv.FormatUint(uint64(a.val), 10)
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
