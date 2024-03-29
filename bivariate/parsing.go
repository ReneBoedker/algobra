package bivariate

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

type monomialMatch struct {
	qr   *QuotientRing
	sign string
	coef string
	vars [2]string
	degs [2]string
}

func newMonomialMatch(match []string, op errors.Op, qr *QuotientRing) (*monomialMatch, error) {
	if len(match) != 8 {
		return nil, errors.New(
			op, errors.Parsing,
			"Regexp-match has unexpected form (%q)", match,
		)
	}
	out := &monomialMatch{
		qr:   qr,
		sign: match[1],
		coef: match[2] + match[7],
		vars: [2]string{
			strings.ToLower(match[3]),
			strings.ToLower(match[5]),
		},
		degs: [2]string{
			match[4],
			match[6],
		},
	}

	// Check that the match is not only a sign
	if out.coef == "" && out.vars == [2]string{"", ""} && out.degs == [2]string{"", ""} {
		return nil, errors.New(
			op, errors.Parsing,
			"Found regexp-match containing only a sign (full match %q)", match[0],
		)
	}

	err := out.ensureVariableOrder(op)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (m *monomialMatch) ensureVariableOrder(op errors.Op) error {
	varLower := [2]string{
		strings.ToLower(m.qr.varNames[0]),
		strings.ToLower(m.qr.varNames[1]),
	}

	switch {
	case m.vars[0] == varLower[0] && (m.vars[1] == varLower[1] || m.vars[1] == ""):
		// Correct; do nothing
	case m.vars[0] == varLower[1] && (m.vars[1] == varLower[0] || m.vars[1] == ""):
		m.vars[0], m.vars[1] = m.vars[1], m.vars[0]
		m.degs[0], m.degs[1] = m.degs[1], m.degs[0]
	case m.vars[0] == "" && m.vars[1] == "":
		// Correct; do nothing
	default:
		return errors.New(
			op, errors.Parsing,
			"Unexpected variable names in match %v", m,
		)
	}
	return nil
}

func (m *monomialMatch) degreesAndCoef(op errors.Op) (deg [2]uint, coef ff.Element, err error) {
	if m.coef == "" {
		coef = m.qr.baseField.One()
	} else {
		coef, err = m.qr.baseField.ElementFromString(strings.Trim(m.coef, "()"))
		if err != nil {
			return deg, coef, errors.Wrap(op, errors.Conversion, err)
		}
	}
	if m.sign == "-" {
		coef.SetNeg()
	}
	for i := 0; i < 2; i++ {
		if m.vars[i] != "" {
			deg[i], err = parseExponent(m.degs[i])
			if err != nil {
				err = errors.Wrap(op, errors.Conversion, err)
				return
			}
		}
	}
	return
}

func parseExponent(s string) (uint, error) {
	if s == "" {
		return 1, nil
	}
	tmp, err := strconv.ParseUint(s, 10, 0)
	return uint(tmp), err
}

func polynomialStringToMap(s string, varNames *[2]string, qr *QuotientRing) (map[[2]uint]ff.Element, error) {
	const op = "Parsing polynomial from string"

	xOrY := regexp.QuoteMeta((*varNames)[0]) + `|` + regexp.QuoteMeta((*varNames)[1])

	// Construct a regular expression that captures each term of a bivariate
	// polynomial
	pattern, err := regexp.Compile(
		`(?P<sign>^|\+|-)\s*` + // A sign (or start of string)
			`(?:` + // Begin group of alternatives

			// First alternative: Term contains a variable
			`(?P<coef>` + qr.baseField.RegexElement(true) + `)?` + // A coefficient
			`\s*\*?\s*` + // Optional multiplication sign
			`(?P<var1>(?i:` + xOrY + `))\^?(?P<deg1>[0-9]*)` + // Variable and optional degree
			`(?:` + // Begin optional second variable
			`\s*\*?\s*` + // Optional multiplication sign
			`(?P<var2>(?i:` + xOrY + `))?\^?(?P<deg2>[0-9]*)` +
			`)?\s*` + // End optional second variable

			`|` +

			// Second alternative: Term is only a coefficient
			`(?P<coefOnly>` + qr.baseField.RegexElement(true) + `)\s*` +

			`)`, // End group of alternatives
	)
	if err != nil {
		return nil, errors.New(
			op, errors.Internal,
			"Failed to compile regular expression using variable names %q and %q",
			varNames[0], varNames[1],
		)
	}

	matches := pattern.FindAllStringSubmatch(s, -1)
	// Check that total match length is the full input string
	matchLen := 0
	for _, m := range matches {
		matchLen += len(m[0])
	}
	if matchLen != len(s) {
		return nil, errors.New(
			op, errors.Parsing,
			"Cannot parse %s; lengths do not match (%d ≠ %d) %v",
			s, matchLen, len(s), matches,
		)
	}

	out := make(map[[2]uint]ff.Element)
	for _, m := range matches {
		tmp, err := newMonomialMatch(m, op, qr)
		if err != nil {
			return nil, err
		}
		deg, coef, err2 := tmp.degreesAndCoef(op)
		if err2 != nil {
			return nil, err2
		}
		if _, ok := out[deg]; !ok {
			out[deg] = coef
		} else {
			out[deg].Add(coef)
		}
	}

	return out, nil
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
