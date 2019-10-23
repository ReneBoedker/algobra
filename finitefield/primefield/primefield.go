package primefield

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"

	"algobra/auxmath"
	"algobra/errors"
	"algobra/finitefield/ff"
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

	_, n, err := auxmath.FactorizePrimePower(card)
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

// RegexElement returns a string containing a regular expression describing an
// element of f.
//
// The input argument requireParens indicates whether parentheses are required
// around elements containing several terms. This has no effect for prime fields.
func (f *Field) RegexElement(requireParens bool) string {
	const pattern = `[0-9]*`

	return pattern
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
	factors, _, _ := auxmath.Factorize(f.Card() - 1)

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
