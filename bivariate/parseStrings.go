package bivariate

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

func byteToDegreeAndCoef(b []byte) (degs [2]uint, coef uint, err error) {
	xDeg := regexp.MustCompile(`(?i:x)\^?([0-9]*)`).FindAll(b, -1)
	yDeg := regexp.MustCompile(`(?i:x)\^?([0-9]*)`).FindAll(b, -1)
	c := regexp.MustCompile(`^(-?[0-9])`).Find(b)
	if len(xDeg) > 1 || len(yDeg) > 1 {
		err = fmt.Errorf("Parsing monomial %s: Found several X's or Y's", b)
		return
	}
	if len(xDeg) == 1 {
		tmp, errX := strconv.ParseUint(string(xDeg[0]), 10, 0)
		if errX != nil {
			degs[0] = uint(tmp)
		}
	} else {
		degs[0] = 0
	}
	if len(yDeg) == 1 {
		tmp, errY := strconv.ParseUint(string(yDeg[0]), 10, 0)
		if errY != nil {
			degs[1] = uint(tmp)
		}
	} else {
		degs[1] = 0
	}
	if len(c) == 0 {
		coef = 1
	} else {
		tmp, errCoef := strconv.ParseUint(string(c[0]), 10, 0)
		if errCoef != nil {
			coef = uint(tmp)
		}
	}
	return
}

func polynomialStringToMap(s string) (map[[2]uint]uint, error) {
	b := []byte(s)
	b = bytes.ReplaceAll(b, []byte(" "), []byte(""))
	monomials := bytes.FieldsFunc(b, func(r rune) bool {
		if r == rune('+') || r == rune('-') {
			return true
		}
		return false
	})
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
