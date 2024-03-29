package bivariate

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ReneBoedker/algobra/finitefield/ff"
)

var prg = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

// func defineField(char uint, t *testing.T) *primefield.Field. See def in bivariate_test.go

func TestAddAndSubDegs(t *testing.T) {
	for i := 0; i < 1000; i++ {
		a, b := uint(prg.Uint32()), uint(prg.Uint32())
		c, d := uint(prg.Uint32()), uint(prg.Uint32())
		if tmp, _ := addDegs([2]uint{a, b}, [2]uint{c, d}); tmp != [2]uint{a + c, b + d} {
			t.Errorf("addDegs({%d,%d},{%d,%d})=%v (Expected {%d,%d})", a, b, c, d, tmp, a+c, b+d)
		}
		tmp, ok := subtractDegs([2]uint{a, b}, [2]uint{c, d})
		switch {
		case (a < c || b < d) && ok:
			t.Errorf(
				"subtractDegs({%d,%d},{%d,%d}) signalled no error (Expected ok=false)",
				a, b, c, d,
			)
		case (a >= c && b >= d) && !ok:
			t.Errorf(
				"subtractDegs({%d,%d},{%d,%d}) signalled an error (Expected ok=true)",
				a, b, c, d,
			)
		}
		if tmp != [2]uint{a - c, b - d} && ok {
			t.Errorf(
				"subtractDegs({%d,%d},{%d,%d})=%v, err (Expected {%d,%d})",
				a, b, c, d, tmp, a-c, b-d,
			)
		}
	}
}

func TestPow(t *testing.T) {
	field := defineField(3, t)
	r := DefRing(field, Lex(true))
	inDegs := [][2]uint{{0, 0}, {1, 0}, {1, 1}, {0, 2}}
	expectedPows := [][][2]uint{
		{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
		{{0, 0}, {1, 0}, {1, 1}, {0, 2}},
		{{0, 0}, {2, 0}, {2, 2}, {0, 4}},
		{{0, 0}, {3, 0}, {3, 3}, {0, 6}},
	}
	for i, d1 := range inDegs {
		f := r.PolynomialFromUnsigned(map[[2]uint]uint{d1: 1})
		for n, exp := range expectedPows {
			g := f.Pow(uint(n))
			if g.Ld() != exp[i] {
				t.Errorf("Pow failed: %v^%d = %v (Expected %v)", f, n, g, exp[i])
			}
		}
	}
}

func TestTimes(t *testing.T) {
	field := defineField(5, t)
	r := DefRing(field, Lex(true))
	testPolys := [...]*Polynomial{
		r.PolynomialFromUnsigned(map[[2]uint]uint{
			{2, 1}: 2,
			{1, 6}: 4,
			{0, 1}: 3,
		}),
		r.PolynomialFromUnsigned(map[[2]uint]uint{
			{1, 1}: 3,
			{0, 0}: 1,
		}),
		r.PolynomialFromUnsigned(map[[2]uint]uint{
			{0, 4}: 1,
			{0, 3}: 2,
			{0, 0}: 4,
		}),
	}

	expected := map[[2]int]*Polynomial{
		{0, 0}: r.PolynomialFromUnsigned(map[[2]uint]uint{
			{4, 2}:  4,
			{3, 7}:  16,
			{2, 12}: 16,
			{2, 2}:  12,
			{1, 7}:  24,
			{0, 2}:  9,
		}),
		{0, 1}: r.PolynomialFromUnsigned(map[[2]uint]uint{
			{3, 2}: 6,
			{2, 7}: 12,
			{2, 1}: 2,
			{1, 6}: 4,
			{1, 2}: 9,
			{0, 1}: 3,
		}),
		{0, 2}: r.PolynomialFromUnsigned(map[[2]uint]uint{
			{2, 5}:  2,
			{2, 4}:  4,
			{2, 1}:  8,
			{1, 10}: 4,
			{1, 9}:  8,
			{1, 6}:  16,
			{0, 5}:  3,
			{0, 4}:  6,
			{0, 1}:  12,
		}),
		{1, 1}: r.PolynomialFromUnsigned(map[[2]uint]uint{
			{2, 2}: 9,
			{1, 1}: 6,
			{0, 0}: 1,
		}),
		{1, 2}: r.PolynomialFromUnsigned(map[[2]uint]uint{
			{1, 5}: 3,
			{1, 4}: 6,
			{1, 1}: 12,
			{0, 4}: 1,
			{0, 3}: 2,
			{0, 0}: 4,
		}),
		{2, 2}: r.PolynomialFromUnsigned(map[[2]uint]uint{
			{0, 8}: 1,
			{0, 7}: 4,
			{0, 6}: 4,
			{0, 4}: 8,
			{0, 3}: 16,
			{0, 0}: 16,
		}),
	}

	for i, f := range testPolys {
		for j, g := range testPolys {
			key := [2]int{i, j}
			if j < i {
				key = [2]int{j, i}
			}

			if _, ok := expected[key]; !ok {
				t.Logf("Skipping (i,j)=(%d,%d): No expected result given", i, j)
				continue
			}

			if tmp := f.Times(g); !tmp.Equal(expected[key]) {
				t.Errorf("(%v)*(%v) = %v (Expected %v)", f, g, tmp, expected[key])
			}
		}
	}
}

func TestEval(t *testing.T) {
	field := defineField(13, t)
	r := DefRing(field, WDegLex(13, 14, false))
	f, err := r.PolynomialFromString("X^2-1")
	if err != nil {
		panic(err)
	}
	evPoints := [][2]ff.Element{
		{field.ElementFromUnsigned(0), field.ElementFromUnsigned(0)},
		{field.ElementFromUnsigned(1), field.ElementFromUnsigned(0)},
		{field.ElementFromUnsigned(1), field.ElementFromUnsigned(2)},
		{field.ElementFromUnsigned(1), field.ElementFromUnsigned(10)},
		{field.ElementFromUnsigned(2), field.ElementFromUnsigned(5)},
	}
	expected := []ff.Element{
		field.ElementFromSigned(-1),
		field.ElementFromUnsigned(0),
		field.ElementFromUnsigned(0),
		field.ElementFromUnsigned(0),
		field.ElementFromUnsigned(3),
	}

	for i, p := range evPoints {
		if v := f.Eval(p); !v.Equal(expected[i]) {
			t.Errorf(
				"%v evaluated at %v gave %v rather than %v",
				f, p, v, expected[i],
			)
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
