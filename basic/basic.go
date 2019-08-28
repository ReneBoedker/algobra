package basic

import (
	"fmt"
	"math/bits"
)

func CeilLog(n uint) uint {
	if n == 0 {
		return 0
	}
	b := uint(bits.Len(n))
	if bits.OnesCount(n) == 1 {
		return b - 1
	} else {
		return b
	}
}

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
			return 0, 0, fmt.Errorf("factorizePrimePower: %d does not seem to be a prime power.", q*Pow(p, n))
		}
		q /= p
		n++
	}
	return p, n, nil
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
