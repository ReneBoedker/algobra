package basic

import (
	"testing"
)

func TestFactorize(t *testing.T) {
	bigPrime := uint(92683)
	p, n, err := FactorizePrimePower(bigPrime)
	if err != nil {
		t.Errorf("FactorizePrimePower(%d) returned error", bigPrime)
		return
	}
	if p != bigPrime || n != 1 {
		t.Errorf("FactorizePrimePower(%d)=(%d,%d,nil), but expected (%[1]d,1, nil)",
			bigPrime, p, n)
	}
}

func TestCeilLog(t *testing.T) {
	testPairs := [][2]uint{
		{0, 0},
		{1, 0},
		{2, 1},
		{431, 9},
		{999999, 20},
		{1<<40 - 1, 40},
	}
	for _, pair := range testPairs {
		if tmp := CeilLog(pair[0]); tmp != pair[1] {
			t.Errorf("CeilLog(%d)=%d (Expected %d)", pair[0], tmp, pair[1])
		}
	}
}
