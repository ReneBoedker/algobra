// Package auxmath contains a number of auxiliary mathematical functions.
package auxmath

import (
	"math/bits"

	"github.com/ReneBoedker/algobra/errors"
)

// BoundSqrt returns an upper bound on the square root of n.
func BoundSqrt(a uint) uint {
	if a == 0 {
		return 0
	}
	b := uint(bits.Len(a))
	if b%2 == 0 {
		b = b >> 1
	} else {
		b = (b >> 1) + 1
	}
	return 1 << b
}

// boundLog2 returns an upper bound on the logarithm of a in base 2.
//
//
func boundLog2(a uint) uint {
	if a == 0 {
		return 0
	}

	if bits.OnesCount(a) == 1 {
		// In this case, we can compute Log2 exactly
		return uint(bits.Len(a)) - 1
	}
	return uint(bits.Len(a))
}

// Pow returns a to the power of n
func Pow(a, n uint) (uint, error) {
	const op = "Computing power of unsigned integer"

	// Ensure that this will not overflow
	if a > 0 && boundLog2(a)*n >= bits.UintSize {
		return 0, errors.New(
			op, errors.Overflow,
			"%d^%d is likely to overflow uint", a, n,
		)
	}

	res := uint(1)
	for ; n > 0; n-- {
		res *= a
	}
	return res, nil
}

// Gcd computes the greatest common divisor between a and b.
func Gcd(a, b uint) uint {
	for a > 0 {
		q := b / a
		b, a = a, b-q*a
	}
	return b
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
