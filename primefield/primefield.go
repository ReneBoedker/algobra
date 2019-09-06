package primefield

import (
	"algobra/basic"
	"algobra/errors"
	"fmt"
	"strconv"
)

var uintBitSize uint

func init() {
	// Determine bit-size of uint type
	i := ^uint(0)
	for i > 0 {
		uintBitSize++
		i >>= 1
	}
}

type Field struct {
	char      uint
	addTable  *table
	multTable *table
}

// Define creates a new finite field with prime characteristic. If char is not a
// prime, the package returns
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

func (f *Field) String() string {
	return fmt.Sprintf("Finite field of %d elements", f.char)
}

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

type Element struct {
	field *Field
	val   uint
	err   error
}

func (pf *Field) Element(val uint) *Element {
	return &Element{field: pf, val: val % pf.char}
}

func (pf *Field) ElementFromSigned(val int) *Element {
	val %= int(pf.char)
	if val < 0 {
		val += int(pf.char)
	}
	return pf.Element(uint(val))
}

func (a *Element) Err() error {
	return a.err
}

func (a *Element) Equal(b *Element) bool {
	if a.field == b.field && a.val == b.val {
		return true
	}
	return false
}

func (a *Element) Plus(b *Element) *Element {
	const op = "Adding elements"
	if tmp := hasErr(op, a, b); tmp != nil {
		return tmp
	}
	if tmp := checkCompatible(op, a, b); tmp != nil {
		return tmp
	}
	if a.field.addTable != nil {
		return a.field.Element(a.field.addTable.lookup(a.val, b.val))
	}
	return a.field.Element(a.val + b.val)
}

func (a *Element) Neg() *Element {
	return a.field.Element(a.field.char - a.val)
}

func (a *Element) Minus(b *Element) *Element {
	const op = "Subtracting elements"
	if tmp := hasErr(op, a, b); tmp != nil {
		return tmp
	}
	if tmp := checkCompatible(op, a, b); tmp != nil {
		return tmp
	}
	return a.Plus(b.Neg())
}

func (a *Element) Mult(b *Element) *Element {
	const op = "Multiplying elements"
	if tmp := hasErr(op, a, b); tmp != nil {
		return tmp
	}
	if tmp := checkCompatible(op, a, b); tmp != nil {
		return tmp
	}
	if a.field.multTable != nil {
		return a.field.Element(a.field.multTable.lookup(a.val, b.val))
	}
	return a.field.Element(a.val * b.val)
}

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

func (a *Element) Zero() bool {
	return (a.val == 0)
}

func (a *Element) Nonzero() bool {
	return (a.val != 0)
}

func (a *Element) One() bool {
	return (a.val == 1)
}

func (a *Element) String() string {
	return strconv.FormatUint(uint64(a.val), 10)
}

// hasErr is an internal method for checking if one of two inputs has a non-nil
// error field.
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
