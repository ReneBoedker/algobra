package primefield

import (
	"math/bits"
	"math/rand"
	"regexp"
	"strconv"

	"algobra/errors"
	"algobra/finitefield/ff"
)

// Ensure that Element implements the ff.Element interface
var _ ff.Element = &Element{}

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
	case string:
		return f.ElementFromString(v)
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

// ElementFromString defines a new element over f from the given string.
//
// A Parsing-error is returned if the string cannot be parsed.
func (f *Field) ElementFromString(val string) (ff.Element, error) {
	const op = "Defining element from string"

	match := regexp.MustCompile(`(-)?([0-9]+)`).FindStringSubmatch(val)

	// Check that the pattern matches the full string
	if len(match[0]) != len(val) {
		return nil, errors.New(
			op, errors.Parsing,
			"Pattern match %q is not the full input string %q", match[0], val,
		)
	}

	switch {
	case len(match[1]) == 1:
		// The input contains a minus
		tmp, err := strconv.ParseInt(match[0], 10, 0)
		if err != nil {
			return nil, errors.New(
				op, errors.Parsing,
				"Failed to convert input with error %q", err,
			)
		}

		return f.ElementFromSigned(int(tmp)), nil
	default:
		// The input does not contain a minus
		tmp, err := strconv.ParseUint(match[0], 10, 0)
		if err != nil {
			return nil, errors.New(
				op, errors.Parsing,
				"Failed to convert input with error %q", err,
			)
		}

		return f.ElementFromUnsigned(uint(tmp)), nil
	}
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

// NTerms returns the number of terms in the representation of a.
func (a *Element) NTerms() uint {
	return 1
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
