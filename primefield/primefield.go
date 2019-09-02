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
	char uint
}

func Define(char uint) (*Field, error) {
	const op = "Defining prime field"
	if char-1 >= 1<<(uintBitSize/2) {
		return nil, errors.New(
			op, errors.InputTooLarge,
			"%d exceeds maximal field size (2^%d)", char, uintBitSize/2,
		)
	}
	_, _, err := basic.FactorizePrimePower(char)
	if err != nil {
		return nil, errors.Wrap(op, errors.Other, err)
	}
	return &Field{char: char}, nil
}

func (f *Field) String() string {
	return fmt.Sprintf("Finite field of %d elements", f.char)
}

type Element struct {
	field *Field
	val   uint
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

func (a *Element) Equal(b *Element) bool {
	if a.field == b.field && a.val == b.val {
		return true
	}
	return false
}

func (a *Element) Plus(b *Element) *Element {
	if a.field != b.field {
		panic("Element.Plus: Elements are from different fields.")
	}
	return a.field.Element(a.val + b.val)
}

func (a *Element) Neg() *Element {
	return a.field.Element(a.field.char - a.val)
}

func (a *Element) Minus(b *Element) *Element {
	if a.field != b.field {
		panic("Element.Minus: Elements are from different fields.")
	}
	return a.Plus(b.Neg())
}

func (a *Element) Mult(b *Element) *Element {
	if a.field != b.field {
		panic("Element.Mult: Elements are from different fields.")
	}
	return a.field.Element(a.val * b.val)
}

func (a *Element) Inv() *Element {
	if a.val == 0 {
		panic("Element.Inv: Cannot invert zero element")
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

func (a *Element) Nonzero() bool {
	return (a.val != 0)
}

func (a *Element) One() bool {
	return (a.val == 1)
}

func (a *Element) String() string {
	return strconv.FormatUint(uint64(a.val), 10)
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
