package basic

import (
	"testing"
)

func TestFactorizePower(t *testing.T) {
	tests := []uint{
		2,
		4,
		8,
		9,
		25,
		24389,
		92683,
	}
	expected := [][2]uint{
		{2, 1},
		{2, 2},
		{2, 3},
		{3, 2},
		{5, 2},
		{29, 3},
		{92683, 1},
	}

	for i, q := range tests {
		p, n, err := FactorizePrimePower(q)
		if err != nil {
			t.Errorf("FactorizePrimePower(%d) returned error", q)
			return
		}
		if p != expected[i][0] || n != expected[i][1] {
			t.Errorf("FactorizePrimePower(%d)=(%d, %d, nil), but expected (%d, %d, nil)",
				q, p, n, expected[i][0], expected[i][1])
		}
	}
}

func TestBoundSqrt(t *testing.T) {
	tests := []uint{
		0,
		1,
		2,
		25,
		431,
		999999,
		1<<40 - 1,
	}
	for _, n := range tests {
		if tmp := BoundSqrt(n); tmp*tmp < n {
			t.Errorf("BoundSqrt(%d) = %d, but %[2]d^2 = %d", n, tmp, tmp*tmp)
		}
	}
}

func TestGcd(t *testing.T) {
	testTriples := [][3]uint{
		{1, 1, 1},
		{1, 2, 1},
		{4, 2, 2},
		{2, 4, 2},
		{24, 18, 6},
	}
	for _, triple := range testTriples {
		if tmp := Gcd(triple[0], triple[1]); tmp != triple[2] {
			t.Errorf("Gcd(%d, %d)=%d (Expected %d)", triple[0], triple[1], tmp, triple[2])
		}
	}
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
