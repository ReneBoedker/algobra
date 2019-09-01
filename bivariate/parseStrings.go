package bivariate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type monomialMatch struct {
	sign string
	coef string
	var1 string
	deg1 string
	var2 string
	deg2 string
}

func newMonomialMatch(match []string) (*monomialMatch, error) {
	if len(match) != 7 {
		return nil, fmt.Errorf(
			"newMonomialmatch: Input match has unexpected form (%v)",
			match)
	}
	out := &monomialMatch{
		sign: match[1],
		coef: match[2],
		var1: strings.ToLower(match[3]),
		deg1: match[4],
		var2: strings.ToLower(match[5]),
		deg2: match[6],
	}
	err := out.ensureVariableOrder()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (m *monomialMatch) ensureVariableOrder() error {
	switch {
	case m.var1 == "x" && (m.var2 == "y" || m.var2 == ""):
		// Correct; do nothing
	case m.var1 == "y" && (m.var2 == "x" || m.var2 == ""):
		m.var1, m.var2 = m.var2, m.var1
		m.deg1, m.deg2 = m.deg2, m.deg1
	case m.var1 == "" && m.var2 == "":
		// Correct; do nothing
	default:
		return fmt.Errorf("monomialMatch: Cannot parse regex-match %v", m)
	}
	return nil
}

func (m *monomialMatch) degreesAndCoef() (deg [2]uint, coef int, err error) {
	if m.coef == "" {
		coef = 1
	} else {
		tmp, err := strconv.ParseInt(m.coef, 10, 0)
		if err != nil {
			return deg, coef, err
		} else {
			coef = int(tmp)
		}
	}
	if m.sign == "-" {
		coef *= -1
	}
	if m.var1 != "" {
		deg[0], err = parseDegree(m.deg1)
		if err != nil {
			return
		}
	}
	if m.var2 != "" {
		deg[1], err = parseDegree(m.deg2)
		if err != nil {
			return
		}
	}
	return
}

func parseDegree(s string) (uint, error) {
	if s == "" {
		return 1, nil
	}
	tmp, err := strconv.ParseUint(s, 10, 0)
	return uint(tmp), err
}

func polynomialStringToSignedMap(s string) (map[[2]uint]int, error) {
	matches := regexp.MustCompile(
		`\s*(?P<sign>^|\+|-)\s*`+
			`(?P<coef>[0-9]*)\s*\*?\s*`+
			`(?P<var1>(?i:x|y))?\^?(?P<deg1>[0-9]*)\s*\*?\s*`+
			`(?P<var2>(?i:x|y))?\^?(?P<deg2>[0-9]*)\s*`).FindAllStringSubmatch(s, -1)
	// Check that total match length is the full input string
	matchLen := 0
	for _, m := range matches {
		matchLen += len(m[0])
	}
	if matchLen != len(s) {
		return nil, fmt.Errorf(
			"Failed to parse string %s as polynomial. Lengths do not match (%d ≠ %d)",
			s, matchLen, len(s))
	}
	out := make(map[[2]uint]int)
	for _, m := range matches {
		tmp, err := newMonomialMatch(m)
		if err != nil {
			return nil, err
		}
		fmt.Println(tmp)
		deg, coef, err2 := tmp.degreesAndCoef()
		fmt.Printf("deg: %v, coef: %v\n\n", deg, coef)
		if err2 != nil {
			return nil, err2
		}
		if _, ok := out[deg]; !ok {
			out[deg] = coef
		}
	}
	return out, nil
}
