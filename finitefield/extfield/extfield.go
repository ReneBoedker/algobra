package extfield

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ReneBoedker/algobra/auxmath"
	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/conway"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/finitefield/primefield"
	"github.com/ReneBoedker/algobra/univariate"
)

func init() {
	// Set a new seed for the pseudo-random generator.
	// Note that this is not cryptographically safe.
	rand.Seed(time.Now().UTC().UnixNano())
}

// Field is the implementation of a finite field.
type Field struct {
	baseField  *primefield.Field
	extDeg     uint
	conwayPoly *univariate.Polynomial
	polyRing   *univariate.QuotientRing
	logTable   *table
}

// Define creates a new finite field with given cardinality.
//
// If card is not a prime power, the package returns an InputValue-error.
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

	baseField, err := primefield.Define(char)
	if err != nil {
		return nil, err
	}

	polyRing := univariate.DefRing(baseField)
	polyRing.SetVarName("a") // Ignoring error since this always succeeds

	conwayCoefs, err := conway.Lookup(char, extDeg)
	if err != nil {
		return nil, errors.Wrap(op, errors.Inherit, err)
	}

	conwayPoly := polyRing.PolynomialFromUnsigned(conwayCoefs)
	id, err := polyRing.NewIdeal(conwayPoly)
	if err != nil {
		return nil, err
	}

	polyRing, err = polyRing.Quotient(id)
	if err != nil {
		return nil, err
	}

	return &Field{
		baseField:  baseField,
		extDeg:     extDeg,
		conwayPoly: conwayPoly,
		polyRing:   polyRing,
		logTable:   nil,
	}, nil
}

// String returns the string representation of f.
func (f *Field) String() string {
	return fmt.Sprintf("Finite field of %d elements", f.Card())
}

// Char returns the characteristic of f.
func (f *Field) Char() uint {
	return f.baseField.Char()
}

// Card returns the cardinality of f.
func (f *Field) Card() uint {
	// Error can be ignored since cardinality was given as uint when
	// constructing the field
	tmp, _ := auxmath.Pow(f.Char(), f.extDeg)
	return tmp
}

// ComputeMultTable will precompute the table of discrete logarithms for the
// field f.
//
// The optional argument maxMem specifies the maximal table size in KiB. If no
// value is given, a default value is used. If more than one value is given,
// only the first is used.
//
// Returns an InputTooLarge-error if the estimated memory usage exceeds the
// maximal value specified by maxMem.
func (f *Field) ComputeMultTable(maxMem ...uint) (err error) {
	if f.logTable == nil {
		f.logTable, err = newLogTable(f, maxMem...)
	}

	if err != nil {
		return err
	}

	return nil
}

// MultGenerator returns an element that generates the units of f.
func (f *Field) MultGenerator() ff.Element {
	// The field is defined from a Conway polynomial, so alpha is a generator
	return f.element([]uint{0, 1})
}

// Elements returns a slice containing all elements of f.
func (f *Field) Elements() []ff.Element {
	out := make([]ff.Element, f.Card(), f.Card())
	out[0] = f.Zero()

	gen := f.MultGenerator()
	for i, e := uint(1), f.One(); i < f.Card(); i, e = i+1, e.Mult(gen) {
		out[i] = e.Copy()
	}
	return out
}

// RegexElement returns a string containing a regular expression describing an
// element of f.
//
// The input argument requireParens indicates whether parentheses should be
// required around elements containing several terms.
func (f *Field) RegexElement(requireParens bool) string {
	termPattern := `(?:[0-9]*(?:` + f.polyRing.VarName() + `(?:\^?[0-9]+)?)|[0-9]+)`
	moreTerms := `(?:` + // Optional group of additional terms consisting of
		`\s*(?:\+|-)\s*` + // a sign
		termPattern + // and a term
		`)*`

	var pattern string

	if requireParens {
		pattern = `(?:\(\s*` + termPattern + moreTerms + `\s*\)|` + // several
			// terms in parentheses
			termPattern + `)` // Or single term

	} else {
		pattern = termPattern + moreTerms
	}

	return pattern
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
