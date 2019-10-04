package univariate

import (
	"algobra/errors"
	"regexp"
	"strconv"
	"strings"
)

type monomialMatch struct {
	sign string
	coef string
	name string
	deg  string
}

func newMonomialMatch(match []string, op errors.Op) (*monomialMatch, error) {
	if len(match) != 5 {
		return nil, errors.New(
			op, errors.Internal,
			"Regexp-match has unexpected form (%v)", match,
		)
	}
	out := &monomialMatch{
		sign: match[1],
		coef: match[2],
		name: strings.ToLower(match[3]),
		deg:  match[4],
	}

	// Check that the match is not only a sign
	if out.coef == "" && out.name == "" && out.deg == "" {
		return nil, errors.New(
			op, errors.Parsing,
			"Found regexp-match containing only a sign (full match %q)", match[0],
		)
	}

	return out, nil
}

func (m *monomialMatch) degreeAndCoef(op errors.Op) (deg int, coef int, err error) {
	if m.coef == "" {
		coef = 1
	} else {
		tmp, err := strconv.ParseInt(m.coef, 10, 0)
		if err != nil {
			return deg, coef, errors.Wrap(op, errors.Conversion, err)
		}
		coef = int(tmp)
	}
	if m.sign == "-" {
		coef *= -1
	}
	if m.name != "" {
		if m.deg == "" {
			deg = 1

		} else {
			tmp, er := strconv.ParseInt(m.deg, 10, 0)
			if er != nil {
				err = errors.Wrap(op, errors.Conversion, er)
				return
			}
			deg = int(tmp)
		}
	}
	return
}

func polynomialStringToSignedMap(s string, varName *string) (map[int]int, error) {
	const op = "Parsing polynomial from string"

	pattern, err := regexp.Compile(
		`\s*(?P<sign>\+|-)?\s*` +
			`(?P<coef>[0-9]*)\s*\*?\s*` +
			`(?:(?P<name>(?i:` + regexp.QuoteMeta(*varName) + `))\^?(?P<deg1>[0-9]*))?\s*`,
	)
	if err != nil {
		return nil, errors.New(
			op, errors.Internal,
			"Failed to compile regular expression using variable name %q", *varName,
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
			"Cannot parse %s; lengths do not match (%d ≠ %d)",
			s, matchLen, len(s),
		)
	}

	out := make(map[int]int)
	for _, m := range matches {
		tmp, err := newMonomialMatch(m, op)
		if err != nil {
			return nil, err
		}
		deg, coef, err2 := tmp.degreeAndCoef(op)
		if err2 != nil {
			return nil, err2
		}
		if _, ok := out[deg]; !ok {
			out[deg] = coef
		} else {
			out[deg] += coef
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
