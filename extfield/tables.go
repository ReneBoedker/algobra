package extfield

import (
	"algobra/errors"
	"math/bits"
)

const defaultMaxMem uint = 1 << 19 // Maximal memory allowed per table in KiB (default: 512 MiB)

type table struct {
	log    []*Element
	invLog map[string]int
	card   int
}

func newMultTable(f *Field, maxMem ...uint) (*table, error) {
	if len(maxMem) == 0 {
		maxMem = append(maxMem, defaultMaxMem)
	}
	if m := estimateMemory(f); m > maxMem[0] {
		return nil, errors.New(
			"Creating arithmetic table", errors.InputTooLarge,
			"Requires %d KiB, which exceeds maxMem (%d KiB)", m, maxMem,
		)
	}

	log := f.Elements()[1:]

	invLog := make(map[string]int, len(log))
	for i, e := range log {
		invLog[e.String()] = i
	}

	return &table{log: log, invLog: invLog, card: int(f.Card())}, nil
}

func (t *table) lookup(i, j *Element) *Element {
	if i.Zero() || j.Zero() {
		return i.field.Zero()
	}
	a, b := t.invLog[i.String()], t.invLog[j.String()]
	return t.log[(a+b)%(t.card-1)].Copy()
}

// estimateMemory gives an estimate on the memory required to store a table.
// This estimate ignores overhead from the map. Return value is in KiB
func estimateMemory(f *Field) uint {
	const keySize = uint(1) // Map keys are uint8 == 1 byte
	elemSize := f.extDeg * bits.UintSize / 8
	b := (f.Card() - 1) * (keySize + elemSize)
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
