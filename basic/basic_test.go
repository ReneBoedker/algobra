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
