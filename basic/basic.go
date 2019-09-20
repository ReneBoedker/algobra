package basic

import (
	"algobra/errors"
	"math/bits"
)

// CeilLog returns the ceiling function applied to the base 2-logarithm of n.
func CeilLog(n uint) uint {
	if n == 0 {
		return 0
	}
	b := uint(bits.Len(n))
	if bits.OnesCount(n) == 1 {
		return b - 1
	}
	return b
}

// Pow returns a to the power of n
func Pow(a, n uint) uint {
	res := uint(1)
	for ; n > 0; n-- {
		res *= a
	}
	return res
}

// Gcd computes the greatest common divisor between a and b.
func Gcd(a, b uint) uint {
	for a > 0 {
		q := b / a
		b, a = a, b-q*a
	}
	return b
}

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
	maxP := CeilLog(q)
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
			exponents = append(exponents, p)

			rFact, rExp, _ := Factorize(n)
			factors = append(factors, rFact...)
			exponents = append(exponents, rExp...)

			return
		}
	}

	maxFactor := CeilLog(n)
	for k := uint(6); k-1 <= maxFactor; k += 6 {
		for _, p := range []uint{k - 1, k + 1} {
			exp := uint(0)
			for n%p == 0 {
				exp++
				n /= p
			}
			if exp > 0 {
				factors = append(factors, p)
				exponents = append(exponents, p)

				rFact, rExp, _ := Factorize(n)
				factors = append(factors, rFact...)
				exponents = append(exponents, rExp...)

				return
			}
		}
	}

	// If no divisor was found so far, p is prime
	factors = append(factors, n)
	return
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
