package univariate_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/univariate"
)

func TestParsingWellFormed(t *testing.T) {
	field := defineField(7)
	ring := univariate.DefRing(field)
	testStrings := []string{
		"4X^4 + X^3 + 3X - 5",
		"4*x4 + X^3  + 3*x -5",
		"4 * X^4 + X^3 + 3 * X - 5",
		"3X4 + X4 +X3+3X-5",
	}
	testPolys := make([]*univariate.Polynomial, len(testStrings), len(testStrings)+1)
	testErrs := make([]error, len(testStrings))
	for i, s := range testStrings {
		testPolys[i], testErrs[i] = ring.PolynomialFromString(s)
	}
	for i, err := range testErrs {
		if err != nil {
			t.Errorf(
				"Failed to parse polynomial %s. Received error %v",
				testStrings[i],
				err,
			)
		}
	}
	testPolys = append(testPolys, ring.PolynomialFromSigned([]int{-5, 3, 0, 1, 4}))
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
	field := defineField(7)
	ring := univariate.DefRing(field)

	testStrings := []string{
		"2^2 X",
		"X^2-+X^3",
		"X^^4+X^2",
		"y^3+y^4",
	}

	testPolys := make([]*univariate.Polynomial, len(testStrings), len(testStrings))
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
	field := defineField(13)
	ring := univariate.DefRing(field)

	testStrings := []string{
		fmt.Sprintf("%d0X", math.MaxInt64),
		fmt.Sprintf("X^%d0", uint(math.MaxUint64)),
	}

	testPolys := make([]*univariate.Polynomial, len(testStrings), len(testStrings))
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
		ring := univariate.DefRing(field)
		for _, varName := range []string{"X", "Y", "α", "\\beta", "e^i", "", "\t"} {
			ring.SetVarName(varName)
			for rep := 0; rep < 100; rep++ {
				// Create random polynomial with up to 50 different terms
				nDegs := (uint(prg.Uint32()) % 50) + 1
				coefs := make([]ff.Element, nDegs, nDegs)
				for i := uint(0); i < nDegs; i++ {
					coefs[i] = field.RandElement()
				}
				f := ring.Polynomial(coefs)

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

func TestSetVarName(t *testing.T) {
	char := uint(31)
	field := defineField(char)
	ring := univariate.DefRing(field)
	f := ring.PolynomialFromUnsigned([]uint{0, 1})
	// Check that the printed variable is correct
	for _, varName := range []string{"X", "y", "(ω^2)", "", "c\td"} {
		if err := ring.SetVarName(varName); err != nil {
			t.Errorf("An error was returned for varName = %q", varName)
		}
		if f.String() != varName {
			t.Errorf("Variable name was not printed as expected")
		}
	}

	// Check that errors are returned if appropriate
	for _, varName := range []string{"\t", "  \n\t", "	", ""} {
		if err := ring.SetVarName(varName); err == nil {
			t.Errorf("No error was returned for varName = %q", varName)
		} else if !errors.Is(errors.InputValue, err) {
			t.Errorf(
				"An error was returned for varName = %q, but it was of "+
					"unexpected type.", varName,
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
