package primefield

import (
	"algobra/errors"
)

var maxMem uint = 1 << 19 // Maximal memory allowed per table in KiB (default: 512 MiB)

type table struct {
	t [][]uint
}

func newTable(f *Field, op func(i, j uint) uint) (*table, error) {
	if m := estimateMemory(f); m > maxMem {
		return nil, errors.New(
			"Creating arithmetic table", errors.InputTooLarge,
			"Requires %d KiB, which exceeds maxMem (%d)", m, maxMem,
		)
	}
	t := make([][]uint, f.char, f.char)
	for i := uint(0); i < f.char; i++ {
		t[i] = make([]uint, f.char-i, f.char-i)
		for j := i; j < f.char; j++ {
			t[i][j-i] = op(i, j)
		}
	}
	return &table{t: t}, nil
}

func (t *table) lookup(i, j uint) uint {
	if j < i {
		return t.lookup(j, i)
	}
	return t.t[i][j-i]
}

// estimateMemory gives a lower bound on the memory required to store a table.
// This estimate ignores overhead from the slices. Return value is in KiB
func estimateMemory(f *Field) uint {
	b := f.char * (f.char + 1) * (uintBitSize / 16)
	return b >> 10
}
