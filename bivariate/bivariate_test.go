package bivariate

import (
	"testing"
)

func TestReduce(t *testing.T) {
	r := DefRing(3, WDegLex(3, 4))
	mod := r.New(map[[2]uint]uint{{9, 0}: 1, {1, 0}: 2})
	id := r.NewIdeal(mod)
	qr, err := r.Quotient(id)
	if err != nil {
		t.Errorf("Failed to construct quotient ring")
	}
	f := qr.New(map[[2]uint]uint{
		{12, 3}: 1,
	})
	// mod.QuoRem(id)
	// fmt.Print("\n\n")
	// id.reduce(mod)
	// fmt.Println(mod)
	//fmt.Printf("q: %v\nr: %v\n", quo, rem)
	if f.Ld() != [2]uint{4, 3} {
		t.Errorf("Reduce failed: Got %v", f.Ld())
	}
}

func TestGroebner1(t *testing.T) {
	r := DefRing(7, Lex(true))
	id := r.NewIdeal(
		r.New(map[[2]uint]uint{
			{1, 2}: 1,
			{0, 3}: 6,
		}),
		r.New(map[[2]uint]uint{
			{0, 3}: 1,
			{0, 2}: 6,
		}),
	)
	expectedGens := []*Polynomial{
		r.New(map[[2]uint]uint{
			{1, 2}: 1,
			{0, 2}: 6,
		}),
		r.New(map[[2]uint]uint{
			{0, 3}: 1,
			{0, 2}: 6,
		}),
	}
	if len(id) != 2 || (!id[0].Equal(expectedGens[0]) && !id[0].Equal(expectedGens[1])) || (!id[1].Equal(expectedGens[0]) && !id[1].Equal(expectedGens[1])) {
		t.Errorf("Gröbner basis has wrong number of elements")
	}
}

func TestPow(t *testing.T) {
	r := DefRing(3, Lex(true))
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

func TestLexOrder(t *testing.T) {
	ord := Lex(true)
	degrees := [][2][2]uint{
		{{1, 0}, {0, 1}}, {{0, 1}, {0, 2}}, {{2, 1}, {2, 1}}, {{3, 4}, {2, 7}},
	}
	expectedOrd := []int{
		1, -1, 0, 1,
	}
	for i, d := range degrees {
		if tmp := ord(d[0], d[1]); tmp != expectedOrd[i] {
			t.Errorf("Lex(%v, %v)=%d (Expected %d)", d[0], d[1], tmp, expectedOrd[i])
		}
	}
}
