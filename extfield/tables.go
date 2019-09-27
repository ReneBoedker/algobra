package extfield

import (
	"algobra/univariate"
	"math/bits"
)

const defaultMaxMem uint = 1 << 19 // Maximal memory allowed per table in KiB (default: 512 MiB)

type table struct {
	t [][]*univariate.Polynomial
}

// func newTable(f *Field, op func(i, j uint) uint, maxMem ...uint) (*table, error) {
// 	if len(maxMem) == 0 {
// 		maxMem = append(maxMem, defaultMaxMem)
// 	}
// 	if m := estimateMemory(f); m > maxMem[0] {
// 		return nil, errors.New(
// 			"Creating arithmetic table", errors.InputTooLarge,
// 			"Requires %d KiB, which exceeds maxMem (%d KiB)", m, maxMem,
// 		)
// 	}
// 	t := make([][]*Element, f.Card(), f.Card())
// 	for i := uint(0); i < f.Card(); i++ {
// 		t[i] = make([]*Element, f.Card()-i, f.Card()-i)
// 		for j := i; j < f.Card(); j++ {
// 			t[i][j-i] = op(i, j)
// 		}
// 	}
// 	return &table{t: t}, nil
// }

func (t *table) lookup(i, j *univariate.Polynomial) *univariate.Polynomial {
	if j.Ld() < i.Ld() {
		return t.lookup(j, i)
	}
	return nil
	// return t.t[i][j-i]
}

// estimateMemory gives a lower bound on the memory required to store a table.
// This estimate ignores overhead from the slices. Return value is in KiB
func estimateMemory(f *Field) uint {
	b := f.Char() * (f.Char() + 1) * (bits.UintSize / 16)
	return b >> 10
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