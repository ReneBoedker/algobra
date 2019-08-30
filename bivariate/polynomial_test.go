package bivariate

import (
	//"algobra/primefield"
	"math/rand"
	"testing"
	"time"
)

var prg = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

// func defineField(char uint, t *testing.T) *primefield.Field. See def in bivariate_test.go

func TestAddAndSubDegs(t *testing.T) {
	for i := 0; i < 1000; i++ {
		a, b := uint(prg.Uint64()), uint(prg.Uint64())
		c, d := uint(prg.Uint64()), uint(prg.Uint64())
		if tmp := addDegs([2]uint{a, b}, [2]uint{c, d}); tmp != [2]uint{a + c, b + d} {
			t.Errorf("addDegs({%d,%d},{%d,%d})=%v (Expected {%d,%d})", a, b, c, d, tmp, a+c, b+d)
		}
		tmp, ok := subtractDegs([2]uint{a, b}, [2]uint{c, d})
		switch {
		case (a < c || b < d) && ok:
			t.Errorf("subtractDegs({%d,%d},{%d,%d}) signalled no error (Expected ok=false)",
				a, b, c, d)
		case (a >= c && b >= d) && !ok:
			t.Errorf("subtractDegs({%d,%d},{%d,%d}) signalled an error (Expected ok=true)",
				a, b, c, d)
		}
		if tmp != [2]uint{a - c, b - d} && ok {
			t.Errorf("subtractDegs({%d,%d},{%d,%d})=%v, err (Expected {%d,%d})",
				a, b, c, d, tmp, a-c, b-d)
		}
	}
}

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
			t.Errorf("Failed to parse polynomial %v", testPolys[i])
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
		f := r.New(map[[2]uint]uint{d1: 1})
		for n, exp := range expectedPows {
			g := f.Pow(uint(n))
			if g.Ld() != exp[i] {
				t.Errorf("Pow failed: %v^%d = %v (Expected %v)", f, n, g, exp[i])
			}
		}
	}
}
