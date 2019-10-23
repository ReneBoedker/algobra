package bivariate

import (
	"algobra/errors"
	"algobra/finitefield/ff"
	"fmt"
	"math"
	"testing"
)

func TestParsingWellFormed(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field, Lex(true))
	testStrings := []string{
		"2X^3Y^2+X^3-2Y+2",
		"2 x3y2 + x^3 - 2Y+ 2",
		"2 X^3 Y^2 + x3 - 2 Y + 2",
		"X^3Y^2 + x^3*y^2 +X3    - 2  y  +2",
	}
	testPolys := make([]*Polynomial, len(testStrings), len(testStrings)+1)
	testErrs := make([]error, len(testStrings))
	for i, s := range testStrings {
		testPolys[i], testErrs[i] = ring.PolynomialFromString(s)
	}
	for i, err := range testErrs {
		if err != nil {
			t.Errorf(
				"Failed to parse polynomial %s. Received error %v",
				testStrings[i], err,
			)
		}
	}
	testPolys = append(testPolys, ring.PolynomialFromUnsigned(map[[2]uint]uint{
		{3, 2}: 2,
		{3, 0}: 1,
		{0, 1}: 5,
		{0, 0}: 2,
	}))
	for i, f := range testPolys {
		for j := i + 1; j < len(testPolys); j++ {
			if !f.Equal(testPolys[j]) {
				t.Errorf(
					"The two polynomials f_%d=%v and f_%d=%v are not equal (but they should be)",
					i, f, j, testPolys[j],
				)
			}
		}
	}
}

func TestParsingIllFormed(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field, Lex(true))

	testStrings := []string{
		"2^2 X",
		"X^2-+Y^3",
		"X^^4Y^5",
		"a^3y^4",
	}

	testPolys := make([]*Polynomial, len(testStrings), len(testStrings))
	testErrs := make([]error, len(testStrings))

	for i, s := range testStrings {
		testPolys[i], testErrs[i] = ring.PolynomialFromString(s)
	}

	for i, err := range testErrs {
		if err == nil {
			t.Errorf(
				"Parsing %q returned polynomial %v instead of an error",
				testStrings[i], testPolys[i],
			)
		} else if !errors.Is(errors.Parsing, err) {
			t.Errorf(
				"Expected errors.Parsing while parsing %q, but received error %q",
				testStrings[i], err.Error(),
			)
		}
	}
}

func TestConversionErrors(t *testing.T) {
	field := defineField(13, t)
	ring := DefRing(field, Lex(true))

	testStrings := []string{
		fmt.Sprintf("%d0X", math.MaxInt64),
		fmt.Sprintf("Y^%d0", uint(math.MaxUint64)),
	}

	testPolys := make([]*Polynomial, len(testStrings), len(testStrings))
	testErrs := make([]error, len(testStrings))

	for i, s := range testStrings {
		testPolys[i], testErrs[i] = ring.PolynomialFromString(s)
	}

	for i, err := range testErrs {
		if err == nil {
			t.Errorf(
				"Parsing %q returned polynomial %v instead of an error",
				testStrings[i], testPolys[i],
			)
		} else if !errors.Is(errors.Conversion, err) {
			t.Errorf(
				"Expected errors.Conversion while parsing %q, but received error %q",
				testStrings[i], err.Error(),
			)
		}
	}
}

func TestParseOutput(t *testing.T) {
	do := func(field ff.Field) {
		ring := DefRing(field, Lex(false))
		for _, varNames := range [][2]string{
			{"Y", "X"},
			{"\\alpha", "\\beta"},
			{"e^i", "x^i"},
			{"", "\t"},
		} {
			ring.SetVarNames(varNames)
			for rep := 0; rep < 100; rep++ {
				// Create random polynomial with up to 50 different terms
				nDegs := (uint(prg.Uint32()) % 4) + 1
				coefMap := make(map[[2]uint]ff.Element)
				coefMap[[2]uint{1, 1}] = field.One() // Cover printing cases with degrees 1
				for i := uint(0); i < nDegs; i++ {
					deg := [2]uint{
						uint(prg.Uint32()),
						uint(prg.Uint32()),
					}
					coefMap[deg] = field.RandElement()
				}
				f := ring.Polynomial(coefMap)

				if g, err := ring.PolynomialFromString(f.String()); err != nil {
					t.Errorf("Parsing formatted output of %v returns error %q", f, err)
				} else if !f.Equal(g) {
					t.Errorf("Formatted output of %v is parsed as %v", f, g)
				}
			}
		}
	}
	fieldLoop(do)
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
