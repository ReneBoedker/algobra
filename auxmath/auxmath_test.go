package auxmath

import (
	"algobra/errors"
	"fmt"
	"math/big"
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

	// Check that non prime powers return an error
	testsErr := []uint{
		0,
		1,
		14,
		100,
		40320,
	}
	for _, q := range testsErr {
		if _, _, err := FactorizePrimePower(q); err == nil {
			t.Errorf(
				"FactorizePrimePower(%d) returned no error, but %[1]d is not a prime power",
				q,
			)
		} else if !errors.Is(errors.InputValue, err) {
			t.Errorf(
				"FactorizePrimePower(%d) returned error, but not of kind InputValue",
				q,
			)
		}
	}
}

func TestFactorize(t *testing.T) {
	tests := []uint{
		0,
		75,
		139,
		444,
		3670,
		952875,
	}
	expected := [][2][]uint{
		{{}, {}},
		{{3, 5}, {1, 2}},
		{{139}, {1}},
		{{2, 3, 37}, {2, 1, 1}},
		{{2, 5, 367}, {1, 1, 1}},
		{{3, 5, 7, 11}, {2, 3, 1, 2}},
	}

	for i, q := range tests {
		p, n, err := Factorize(q)
		if err != nil {
			t.Errorf("Factorize(%d) returned error", q)
			return
		}
		if len(expected[i][0]) != len(p) {
			t.Errorf(
				"Factorize(%d) gave factors %v, but expected %v",
				q, p, expected[i][0],
			)
		}
		for j, f := range p {
			if f != expected[i][0][j] || n[j] != expected[i][1][j] {
				t.Errorf("Factorize(%d) gave factor %d^%d but expected %d^%d",
					q, f, n[j], expected[i][0][j], expected[i][1][j])
			}
		}
	}
}

func TestPow(t *testing.T) {
	tests := [][3]uint{
		{0, 0, 1},
		{0, 4, 0},
		{2, 5, 32},
		{5, 20, 95367431640625},
		{10, 7, 10000000},
		{17, 6, 24137569},
	}
	for _, n := range tests {
		tmp, err := Pow(n[0], n[1])
		if tmp != n[2] {
			t.Errorf("Pow(%d, %d) = %d, but expected %d", n[0], n[1], tmp, n[2])
		}
		if err != nil {
			t.Errorf("Pow(%d, %d) returned error %q", n[0], n[1], err)
		}
	}

	// Check that overflow detection works
	testsErr := [][2]uint{
		{2, 64},
		{5, 28},
		{10, 20},
		{17, 16},
	}
	for _, n := range testsErr {
		if _, err := Pow(n[0], n[1]); err == nil {
			t.Errorf(
				"Pow(%d, %d) returned no error, but overflow was expected",
				n[0], n[1],
			)
		} else if !errors.Is(errors.Overflow, err) {
			t.Errorf(
				"Pow(%d, %d) returned an error, but not of kind Overflow",
				n[0], n[1],
			)
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

func TestCombinIter(t *testing.T) {
	one := big.NewInt(1)
	expected := big.NewInt(0)

	for n := 5; n < 20; n++ {
		for k := 0; k <= n; k++ {
			expected.Binomial(int64(n), int64(k))
			count := big.NewInt(0)

			unique := make(map[string]struct{})
			for ci := NewCombinIter(n, k); ci.Active(); ci.Next() {
				tmp := fmt.Sprint(ci.Current())
				if _, exists := unique[tmp]; exists {
					t.Errorf("Found combination %s twice for n=%d and k=%d", tmp, n, k)
				} else {
					unique[tmp] = struct{}{}
				}
				count.Add(count, one)
			}

			if count.Cmp(expected) != 0 {
				t.Errorf(
					"Found %v combinations for n=%d and k=%d, but expected %d",
					count, n, k, expected,
				)
			}
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
