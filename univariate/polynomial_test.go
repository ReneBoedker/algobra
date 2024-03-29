package univariate_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/univariate"
)

var prg = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func defineField(card uint) ff.Field {
	field, err := finitefield.Define(card)
	if err != nil {
		// Error is in tests, so panic is OK
		panic(err)
	}
	return field
}

func fieldLoop(do func(field ff.Field), minCard ...uint) {
	for _, card := range [...]uint{2, 3, 4, 5, 9, 16, 25, 49, 64, 125} {
		if len(minCard) > 0 && card < minCard[0] {
			continue
		}
		field, err := finitefield.Define(card)
		if err != nil {
			// Error is in tests, so panic is OK
			panic(err)
		}

		do(field)
	}
	return
}

func TestAssignment(t *testing.T) {
	field := defineField(13)
	ring := univariate.DefRing(field)

	// Define 2X^4-X^2+3X-3 in three separate ways
	polynomials := []*univariate.Polynomial{
		ring.Polynomial([]ff.Element{
			field.ElementFromSigned(-3),
			field.ElementFromSigned(3),
			field.ElementFromSigned(-1),
			field.Zero(),
			field.ElementFromSigned(2),
		}),
		ring.PolynomialFromUnsigned([]uint{10, 3, 12, 0, 2}),
		ring.PolynomialFromSigned([]int{-3, 3, -1, 0, 2}),
	}

	for i, f := range polynomials {
		// Check that the leading term, coefficient, and degree are correct
		if !f.Lt().Equal(ring.PolynomialFromUnsigned([]uint{0, 0, 0, 0, 2})) {
			t.Errorf(
				"Assignment %d has wrong leading term (%v; expected 2X^4)",
				i+1, f.Lt(),
			)
		}

		if !f.Lc().Equal(field.ElementFromUnsigned(2)) {
			t.Errorf(
				"Assignment %d has wrong leading coefficient (%v; expected 2)",
				i+1, f.Lc(),
			)
		}

		if f.Ld() != 4 {
			t.Errorf(
				"Assignment %d has wrong leading degree (%v; expected 4)",
				i+1, f.Ld(),
			)
		}

		// Check that degrees are correct
		degs := f.Degrees()
		if len(degs) != 4 {
			t.Errorf(
				"Assignment %d has wrong degrees (%v; expected [4,2,1,0])",
				i+1, degs,
			)
		} else {
			var degArray [4]int
			copy(degArray[:], degs)
			if degArray != [4]int{4, 2, 1, 0} {
				t.Errorf(
					"Assignment %d has wrong degrees (%v; expected [4,2,1,0])",
					i+1, degs,
				)
			}
		}

		// Check that all polynomials are equal
		for j, g := range polynomials {
			if !f.Equal(g) {
				t.Errorf(
					"Assignments %d and %d gave different polynomials %v and %v",
					i+1, j+1, f, g,
				)
			}
		}
	}
}

func TestCoefs(t *testing.T) {
	do := func(field ff.Field) {
		ring := univariate.DefRing(field)

		for rep := 0; rep < 200; rep++ {
			const nDegs = 20
			coefs := make([]uint, nDegs, nDegs)

			for j := range coefs {
				coefs[j] = uint(prg.Uint32())
			}

			f := ring.PolynomialFromUnsigned(coefs)
			g := ring.Polynomial(f.Coefs())

			if !f.Equal(g) {
				t.Errorf("%v did not test equal to %v", f, g)
			}
		}
	}

	fieldLoop(do)
}

func TestPlus(t *testing.T) {
	field := defineField(7)
	ring := univariate.DefRing(field)

	f := ring.PolynomialFromUnsigned([]uint{2, 0, 3})
	g := ring.PolynomialFromUnsigned([]uint{1, 1, 1, 1})
	h := ring.PolynomialFromSigned([]int{-1, 0, 0, 3})

	tests := [][3]*univariate.Polynomial{ // f, g, expected sum
		{f, f, ring.PolynomialFromUnsigned([]uint{4, 0, 6})},
		{f, g, ring.PolynomialFromUnsigned([]uint{3, 1, 4, 1})},
		{g, f, ring.PolynomialFromUnsigned([]uint{3, 1, 4, 1})},
		{f, h, ring.PolynomialFromUnsigned([]uint{1, 0, 3, 3})},
		{h, f, ring.PolynomialFromUnsigned([]uint{1, 0, 3, 3})},
		{g, g, ring.PolynomialFromUnsigned([]uint{2, 2, 2, 2})},
		{g, h, ring.PolynomialFromUnsigned([]uint{0, 1, 1, 4})},
		{h, g, ring.PolynomialFromUnsigned([]uint{0, 1, 1, 4})},
		{h, h, ring.PolynomialFromSigned([]int{-2, 0, 0, 6})},
	}

	for _, test := range tests {
		res := test[0].Plus(test[1])
		if !res.Equal(test[2]) {
			t.Errorf(
				"(%v) + (%v) = %v, but expected %v",
				test[0], test[1], res, test[2],
			)
		}
	}
}

func TestMinus(t *testing.T) {
	field := defineField(7)
	ring := univariate.DefRing(field)

	f := ring.PolynomialFromUnsigned([]uint{2, 0, 3})
	g := ring.PolynomialFromUnsigned([]uint{1, 1, 1, 1})
	h := ring.PolynomialFromSigned([]int{-1, 0, 0, 3})

	tests := [][3]*univariate.Polynomial{ // f, g, expected difference
		{f, f, ring.Zero()},
		{f, g, ring.PolynomialFromSigned([]int{1, -1, 2, -1})},
		{g, f, ring.PolynomialFromSigned([]int{-1, 1, -2, 1})},
		{f, h, ring.PolynomialFromSigned([]int{3, 0, 3, -3})},
		{h, f, ring.PolynomialFromSigned([]int{-3, 0, -3, 3})},
		{g, g, ring.Zero()},
		{g, h, ring.PolynomialFromSigned([]int{2, 1, 1, -2})},
		{h, g, ring.PolynomialFromSigned([]int{-2, -1, -1, 2})},
		{h, h, ring.Zero()},
	}

	for _, test := range tests {
		res := test[0].Minus(test[1])
		if !res.Equal(test[2]) {
			t.Errorf(
				"(%v) - (%v) = %v, but expected %v",
				test[0], test[1], res, test[2],
			)
		}
	}
}

func TestTimes(t *testing.T) {
	field := defineField(7)
	ring := univariate.DefRing(field)

	f := ring.PolynomialFromUnsigned([]uint{2, 0, 3})
	g := ring.PolynomialFromUnsigned([]uint{1, 1, 1, 1})
	h := ring.PolynomialFromSigned([]int{-1, 0, 0, 3})

	tests := [][3]*univariate.Polynomial{ // f, g, expected product
		{f, f, ring.PolynomialFromSigned([]int{4, 0, 12, 0, 9})},
		{f, g, ring.PolynomialFromSigned([]int{2, 2, 5, 5, 3, 3})},
		{g, f, ring.PolynomialFromSigned([]int{2, 2, 5, 5, 3, 3})},
		{f, h, ring.PolynomialFromSigned([]int{-2, 0, -3, 6, 0, 9})},
		{h, f, ring.PolynomialFromSigned([]int{-2, 0, -3, 6, 0, 9})},
		{g, g, ring.PolynomialFromSigned([]int{1, 2, 3, 4, 3, 2, 1})},
		{g, h, ring.PolynomialFromSigned([]int{-1, -1, -1, 2, 3, 3, 3})},
		{h, g, ring.PolynomialFromSigned([]int{-1, -1, -1, 2, 3, 3, 3})},
		{h, h, ring.PolynomialFromSigned([]int{1, 0, 0, -6, 0, 0, 9})},
	}

	for _, test := range tests {
		res := test[0].Times(test[1])
		if !res.Equal(test[2]) {
			t.Errorf(
				"(%v) * (%v) = %v, but expected %v",
				test[0], test[1], res, test[2],
			)
		}
	}

	zero := ring.Zero()
	one := ring.One()
	for _, test := range []*univariate.Polynomial{f, g, h} {
		if res := test.Times(zero); !res.Equal(zero) {
			t.Errorf(
				"(%v) * 0 = %v, but expected 0",
				test, res,
			)
		}
		if res := zero.Times(test); !res.Equal(zero) {
			t.Errorf(
				"0 * (%v) = %v, but expected 0",
				test, res,
			)
		}

		if res := test.Times(one); !res.Equal(test) {
			t.Errorf(
				"(%v) * 1 = %v, but expected %[1]v",
				test, res,
			)
		}
		if res := one.Times(test); !res.Equal(test) {
			t.Errorf(
				"1 * (%v) = %v, but expected %[1]v",
				test, res,
			)
		}
	}
}

func TestEquality(t *testing.T) {
	field := defineField(5)
	ring1 := univariate.DefRing(field)
	ring2 := univariate.DefRing(field)

	f := ring1.PolynomialFromUnsigned([]uint{1, 2, 3})
	g := ring1.PolynomialFromUnsigned([]uint{1, 2, 2})
	h := ring1.Zero()
	k := ring2.Zero()

	tests := [][2]*univariate.Polynomial{
		{f, f},
		{f, g},
		{f, h},
		{f, k},
		{g, g},
		{g, h},
		{g, k},
		{h, h},
		{h, k},
		{k, k},
	}
	expected := []bool{
		true,
		false,
		false,
		false,
		true,
		false,
		false,
		true,
		false,
		true,
	}

	for i, test := range tests {
		e1 := test[0].Equal(test[1])
		e2 := test[1].Equal(test[0])
		if e1 != e2 {
			t.Errorf(
				"f.Equal(g) is different from g.Equal(f) for f=%v and g=%v",
				f, g,
			)
		} else if e1 != expected[i] {
			t.Errorf(
				"(%v).Equal(%v) gives %t even though %t is expected",
				f, g, e1, expected[i],
			)
		}
	}
}

func TestNormalize(t *testing.T) {
	do := func(field ff.Field) {
		ring := univariate.DefRing(field)
		for rep := 0; rep <= 100; rep++ {
			const nDegs = 5
			degs := make([]uint, nDegs, nDegs)

			for j := range degs {
				degs[j] = uint(prg.Uint32())
			}

			f := ring.PolynomialFromUnsigned(degs)
			if f.IsZero() {
				continue
			}
			g := f.Normalize()
			if g.Err() != nil {
				t.Errorf("Normalize gave error status %q", g.Err())
				continue
			}

			if !g.Lc().IsOne() || !f.Equal(g.Scale(f.Lc())) {
				t.Errorf("%v was normalized as %v (f.coefs = %v)", f, g, f.Coefs())
			}
		}

		// Check that normalizing zero gives zero
		f := ring.Zero().Normalize()
		if !f.Equal(ring.Zero()) {
			t.Errorf("Normalizing zero polynomial gave %v rather than 0", f)
		}
	}

	fieldLoop(do)
}

func TestQuotient(t *testing.T) {
	field := defineField(5)
	ring := univariate.DefRing(field)

	id, err := ring.NewIdeal(
		ring.PolynomialFromUnsigned([]uint{0, 2, 0, 2}),
		ring.PolynomialFromUnsigned([]uint{1, 0, 1, 1, 0, 1}),
	)

	if err != nil {
		panic(err)
	}

	// Check that the ideal generator is automatically reduced to the gcd
	if !id.Generator().Equal(ring.PolynomialFromUnsigned([]uint{1, 0, 1})) {
		t.Fatalf("Ideal had generator %v, but expected X^2 + 1", id.Generator())
	}

	qr, err := ring.Quotient(id)
	if err != nil {
		panic(err)
	}

	// Check that simple calculations in the ring succeed
	if qr.PolynomialFromUnsigned([]uint{1, 0, 2, 0, 1}).IsNonzero() {
		t.Errorf(
			"Polynomial (X^2 + 1)^2 reduced to %v rather than 0",
			qr.PolynomialFromUnsigned([]uint{1, 0, 2, 0, 1}),
		)
	}
	f := qr.PolynomialFromUnsigned([]uint{0, 3})
	g := qr.PolynomialFromUnsigned([]uint{4, 0, 0, 2})

	if !f.Equal(qr.PolynomialFromUnsigned([]uint{0, 3})) {
		t.Errorf("Polynomial 3X reduced to %v rather than 3X", f)
	}
	if !g.Equal(qr.PolynomialFromSigned([]int{4, -2})) {
		t.Errorf("Polynomial 2X^3+4 reduced to %v rather than -2X+4", f)
	}
}

func TestGcd(t *testing.T) {
	field := defineField(3)
	ring := univariate.DefRing(field)

	gcd, err := univariate.Gcd(
		ring.PolynomialFromUnsigned([]uint{2, 1, 2, 1, 2, 0, 2}),    // (2X^4+X+2)(X^2+1)
		ring.PolynomialFromUnsigned([]uint{0, 2, 1, 0, 0, 2}),       // (2X^4+X+2)X
		ring.PolynomialFromUnsigned([]uint{0, 0, 0, 1, 2, 0, 0, 1}), // (2X^4+X+2)(2X^3)
	)
	if err != nil {
		panic(err)
	}

	if !gcd.Equal(ring.PolynomialFromUnsigned([]uint{2, 1, 0, 0, 2})) {
		t.Errorf("Gcd returned %v rather than 2X^4 + X + 2", gcd)
	}

	// Test that an error is returned when polynomials are defined over
	// different rings.
	field2 := defineField(5)
	ring2 := univariate.DefRing(field2)
	gcd, err = univariate.Gcd(
		ring.PolynomialFromUnsigned([]uint{2, 1, 2, 1, 2, 0, 2}),
		ring.PolynomialFromUnsigned([]uint{0, 2, 1, 0, 0, 2}),
		ring2.PolynomialFromUnsigned([]uint{0, 0, 0, 1, 2, 0, 0, 1}),
	)

	if err == nil {
		t.Errorf("Gcd returned no error even though polynomials are defined " +
			"over different rings")
	} else if !errors.Is(errors.InputIncompatible, err) {
		t.Errorf(
			"Gcd returned an error (%s) but different kind than expected\n"+
				"err = %v", err.Error(), err,
		)
	}
}

func TestPow(t *testing.T) {
	field := defineField(11)
	ring := univariate.DefRing(field)

	f := ring.PolynomialFromUnsigned([]uint{5, 4, 3, 2, 1})
	expected := []*univariate.Polynomial{
		ring.PolynomialFromUnsigned([]uint{1}),
		f,
		ring.PolynomialFromUnsigned([]uint{25, 40, 46, 44, 35, 20, 10, 4, 1}),
		ring.PolynomialFromUnsigned([]uint{125, 300, 465, 574, 594, 504, 369, 234, 126, 56, 21, 6, 1}),
	}

	for i, g := range expected {
		fPow := f.Pow(uint(i))
		if !fPow.Equal(g) {
			t.Errorf("(%v)^%d = %v, but expected %v", f, i, fPow, g)
		}
	}
}
