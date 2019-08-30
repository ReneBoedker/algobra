package bivariate

import (
	"fmt"
	"regexp"
	"strconv"
)

func byteToDegreeAndCoef(string) (degs [2]uint, coef uint, err error) {

	return
}

func polynomialStringToMap(s string) (map[[2]uint]uint, error) {
	matches := regexp.MustCompile(
		`\s*(?P<sign>\+|-)\s*`+
			`(?P<coef>[0-9]*)\s*\*?\s*`+
			`(?P<var1>(?i:x|y))?\^?(?P<deg1>[0-9]*)\s*\*?\s*`+
			`(?P<var2>(?i:x|y))?\^?(?P<deg2>[0-9]*)\s*`).FindAllStringSubmatch(s, -1)
	// Check that total match length is the full input string
	matchLen := 0
	for _, m := range matches {
		matchLen += len(m[0])
	}
	if matchLen != len(s) {
		return nil, fmt.Errorf("Failed to parse string %s as polynomial", s)
	}
	out := make(map[[2]uint]uint)
	for _, m := range monomials {
		deg, coef, err := byteToDegreeAndCoef(m)
		if err != nil {
			return nil, err
		}
		if _, ok := out[deg]; !ok {
			out[deg] = coef
		}
	}
	return out, nil
}
