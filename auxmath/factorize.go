package auxmath

import (
	"github.com/ReneBoedker/algobra/errors"
)

// FactorizePrimePower computes p and n such that q=p^n.
//
// If q is not a prime power, the function returns an InputValue-error.
func FactorizePrimePower(q uint) (p uint, n uint, err error) {
	const op = "Factorizing prime power"

	if q == 0 || q == 1 {
		return 0, 0, errors.New(
			op, errors.InputValue,
			"%d is not a prime power.", q,
		)
	}

	if q%2 == 0 {
		p = 2
	} else if q%3 == 0 {
		p = 3
	}
	maxP := BoundSqrt(q)
	for k := uint(6); p == 0 && k-1 <= maxP; k += 6 {
		if q%(k-1) == 0 {
			p = k - 1
		} else if q%(k+1) == 0 {
			p = k + 1
		}
	}
	if p == 0 {
		// If no divisor was found so far, p is prime
		return q, 1, nil
	}
	for q > 1 {
		if q%p != 0 {
			// Restore the original value of q for the error message
			tmp, _ := Pow(p, n) // Ignore errors since p^n smaller than input q
			return 0, 0, errors.New(
				op, errors.InputValue,
				"%d does not seem to be a prime power.", q*tmp,
			)
		}
		q /= p
		n++
	}
	return p, n, nil
}

// Factorize computes the prime factorization of n.
//
// The output slice factors contains each distinct prime factor, and exponents
// contains the corresponding exponents. When n is one, both will be empty,
// representing the empty product.
//
// As an exception, the function will return 0 as the only factor and
// 1 as the exponent when given zero as input.
func Factorize(n uint) (factors, exponents []uint) {
	const op = "Factorizing integer"

	factors = make([]uint, 0)
	exponents = make([]uint, 0)

	switch n {
	case 0:
		return []uint{0}, []uint{1}
	case 1:
		// Factorization is the empty product
		return []uint{}, []uint{}
	}

	for p := uint(2); p <= 3; p++ {
		exp := uint(0)
		for n%p == 0 {
			exp++
			n /= p
		}
		if exp > 0 {
			factors = append(factors, p)
			exponents = append(exponents, exp)

			rFact, rExp := Factorize(n)
			factors = append(factors, rFact...)
			exponents = append(exponents, rExp...)

			return
		}
	}

	maxFactor := BoundSqrt(n)
	for k := uint(6); k-1 <= maxFactor; k += 6 {
		for _, p := range []uint{k - 1, k + 1} {
			exp := uint(0)
			for n%p == 0 {
				exp++
				n /= p
			}
			if exp > 0 {
				factors = append(factors, p)
				exponents = append(exponents, exp)

				rFact, rExp := Factorize(n)
				factors = append(factors, rFact...)
				exponents = append(exponents, rExp...)

				return
			}
		}
	}

	// If no divisor was found so far, p is prime
	factors = append(factors, n)
	exponents = append(exponents, 1)
	return
}
