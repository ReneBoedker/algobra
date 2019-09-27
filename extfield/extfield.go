package extfield

import (
	"algobra/basic"
	"algobra/errors"
	"algobra/extfield/conway"
	"algobra/finitefield"
	"algobra/univariate"
	"fmt"
)

// Field is the implementation of a finite field
type Field struct {
	baseField *finitefield.Field
	extDeg    uint
	polyRing  *univariate.QuotientRing
	generator *Element
	addTable  *table
	multTable *table
}

// Define creates a new finite field with given cardinality.
//
// If card is not a prime power, the package returns an InputValue-error. If
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

	char, extDeg, err := basic.FactorizePrimePower(card)
	if err != nil {
		return nil, errors.Wrap(op, errors.Inherit, err)
	}

	baseField, err := finitefield.Define(char)
	if err != nil {
		return nil, err
	}

	polyRing := univariate.DefRing(baseField)
	conwayCoefs, err := conway.Lookup(char, extDeg)
	if err != nil {
		errors.Wrap(op, errors.Inherit, err)
	}

	id, err := polyRing.NewIdeal(polyRing.PolynomialFromUnsigned(conwayCoefs))
	if err != nil {
		return nil, err
	}
	polyRing, err = polyRing.Quotient(id)
	if err != nil {
		return nil, err
	}

	return &Field{
		baseField: baseField,
		extDeg:    extDeg,
		addTable:  nil,
		multTable: nil,
	}, nil
}

// String returns the string representation of f
func (f *Field) String() string {
	return fmt.Sprintf("Finite field of %d elements", f.Card())
}

// Char returns the characteristic of f
func (f *Field) Char() uint {
	return f.baseField.Char()
}

// Card returns the cardinality of f
func (f *Field) Card() uint {
	return basic.Pow(f.Char(), f.extDeg)
}

// MultGenerator returns an element that generates the units of f.
func (f *Field) MultGenerator() *Element {
	return f.generator.Copy()
}

// Elements returns a slice containing all elements of f.
func (f *Field) Elements() []*Element {
	out := make([]*Element, f.Card(), f.Card())
	out[0] = f.Zero()

	gen := f.MultGenerator()
	for i, e := uint(1), gen.Copy(); i < f.Char(); i, e = i+1, e.Mult(gen) {
		out[i] = e.Copy()
	}
	return out
}

// Element is the implementation of an element in a finite field.
type Element struct {
	field *Field
	val   *univariate.Polynomial
	err   error
}

// Zero defines a new zero element over f.
func (f *Field) Zero() *Element {
	return &Element{
		field: f,
		val:   f.polyRing.Zero(),
	}
}

// Element defines a new element over f with value val
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) Element(val []uint) *Element {
	return &Element{
		field: f,
		val:   f.polyRing.PolynomialFromUnsigned(val),
	}
}

// ElementFromSigned defines a new element over f with value val
//
// The returned element will be reduced modulo the characteristic automatically.
// Negative values are reduced to a positive remainder (as opposed to the
// %-operator in Go).
func (f *Field) ElementFromSigned(val []int) *Element {
	return &Element{
		field: f,
		val:   f.polyRing.PolynomialFromSigned(val),
	}
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

// Equal tests equality of elements a and b.
func (a *Element) Equal(b *Element) bool {
	if a.field == b.field && a.val.Equal(b.val) {
		return true
	}
	return false
}

// Zero returns a boolean describing whether a is the zero element
func (a *Element) Zero() bool {
	return a.val.Zero()
}

// Nonzero returns a boolean describing whether a is a non-zero element
func (a *Element) Nonzero() bool {
	return a.val.Nonzero()
}

// One returns a boolean describing whether a is one
func (a *Element) One() bool {
	return a.val.One()
}

// String returns the string representation of a.
func (a *Element) String() string {
	return a.val.String()
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
		out := a.field.Zero()
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
