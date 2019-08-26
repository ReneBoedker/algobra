package basic

import (
	"fmt"
)

func Pow(a, n uint) uint {
	res := uint(1)
	for ; n > 0; n-- {
		res *= a
	}
	return res
}

func FactorizePrimePower(q uint) (p uint, n uint, err error) {
	if q%2 == 0 {
		p = 2
	} else if q%3 == 0 {
		p = 3
	}
	for k := uint(6); p == 0; k += 6 {
		if q%(k-1) == 0 {
			p = k - 1
		} else if q%(k+1) == 0 {
			p = k + 1
		}
	}
	for q > 1 {
		if q%p != 0 {
			return 0, 0, fmt.Errorf("factorizePrimePower: %d does not seem to be a prime power.", q*Pow(p, n))
		}
		q /= p
		n++
	}
	return p, n, nil
}
