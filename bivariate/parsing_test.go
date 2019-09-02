package bivariate

import (
	"testing"
)

func TestParsing(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field, Lex(true))
	testStrings := []string{
		"2X^3Y^2+X^3-2Y+2",
		"2 x3y2 + x^3 - 2Y+ 2",
		"2 X^3 Y^2 + x3 - 2 Y + 2",
	}
	testPolys := make([]*Polynomial, len(testStrings), len(testStrings)+1)
	testErrs := make([]error, len(testStrings))
	for i, s := range testStrings {
		testPolys[i], testErrs[i] = ring.NewFromString(s)
	}
	for i, err := range testErrs {
		if err != nil {
			t.Errorf("Failed to parse polynomial %s. Received error %v",
				testStrings[i],
				err)
		}
	}
	testPolys = append(testPolys, ring.New(map[[2]uint]uint{
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
					i, f, j, testPolys[j])
			}
		}
	}
}
