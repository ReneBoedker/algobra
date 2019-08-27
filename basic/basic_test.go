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
