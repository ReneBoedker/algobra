package finitefields

import (
	"strconv"
)

type Element struct {
	val uint
	mod uint
}

func New(val, mod uint) *Element {
	return &Element{val: val % mod, mod: mod}
}

func (a *Element) Equal(b *Element) bool {
	if a.mod == b.mod && a.val == b.val {
		return true
	}
	return false
}

func (a *Element) Plus(b *Element) *Element {
	if a.mod != b.mod {
		panic("Element.Plus: Elements are from different fields.")
	}
	return &Element{val: (a.val + b.val) % a.mod, mod: a.mod}
}

func (a *Element) Neg() *Element {
	return &Element{val: a.mod - a.val, mod: a.mod}
}

func (a *Element) Minus(b *Element) *Element {
	if a.mod != b.mod {
		panic("Element.Minus: Elements are from different fields.")
	}
	return a.Plus(b.Neg())
}

func (a *Element) Mult(b *Element) *Element {
	if a.mod != b.mod {
		panic("Element.Mult: Elements are from different fields.")
	}
	return &Element{val: (a.val * b.val) % a.mod, mod: a.mod}
}

func (a *Element) Inv() *Element {
	if a.val == 0 {
		panic("Element.Inv: Cannot invert zero element")
	}
	// Implemented using the extended euclidean algorithm (see for instance
	// [GG13])
	r0 := a.mod
	r1 := a.val
	i0, i1 := 0, 1
	for r1 > 0 {
		q := r0 / r1
		r0, r1 = r1, r0-q*r1
		i0, i1 = i1, i0-int(q)*i1
	}
	for i0 < 0 {
		i0 += int(a.mod)
	}
	return &Element{val: uint(i0), mod: a.mod}
}

func (a *Element) Nonzero() bool {
	return (a.val != 0)
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
