package auxmath

import (
	"algobra/errors"
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
			return 0, 0, errors.New(
				op, errors.InputValue,
				"%d does not seem to be a prime power.", q*Pow(p, n),
			)
		}
		q /= p
		n++
	}
	return p, n, nil
}

// Factorize computes the prime factorization of n.
func Factorize(n uint) (factors, exponents []uint, err error) {
	const op = "Factorizing integer"
	factors = make([]uint, 0)
	exponents = make([]uint, 0)

	switch n {
	case 0:
		return nil, nil, errors.New(
			op, errors.InputValue,
			"Cannot factorize 0",
		)
	case 1:
		// Factorization is the empty sum
		return
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

			rFact, rExp, _ := Factorize(n)
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

				rFact, rExp, _ := Factorize(n)
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
