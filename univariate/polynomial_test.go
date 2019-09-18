package univariate

import (
	"algobra/finitefield"
	"testing"
)

func defineField(char uint, t *testing.T) *finitefield.Field {
	field, err := finitefield.Define(char)
	if err != nil {
		t.Fatalf("Failed to define finite field of %d elements", char)
	}
	return field
}

func TestAssignment(t *testing.T) {
	field := defineField(13, t)
	ring := DefRing(field)

	// Define 2X^4-X^2+3X-3 in three separate ways
	polynomials := []*Polynomial{
		ring.Polynomial([]*finitefield.Element{
			field.ElementFromSigned(-3),
			field.ElementFromSigned(3),
			field.ElementFromSigned(-1),
			field.Zero(),
			field.ElementFromSigned(2),
		}),
		ring.PolynomialFromUnsigned([]uint{10, 3, 12, 0, 2}),
		ring.PolynomialFromSigned([]int{-3, 3, -1, 0, 2}),
	}

	for _, f := range polynomials {
		for _, g := range polynomials {
			if !f.Equal(g) {
				t.Errorf("Assignment gave different polynomials %v and %v", f, g)
			}
		}
	}
}

func TestPlus(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field)

	f := ring.PolynomialFromUnsigned([]uint{2, 0, 3})
	g := ring.PolynomialFromUnsigned([]uint{1, 1, 1, 1})
	h := ring.PolynomialFromSigned([]int{-1, 0, 0, 3})

	tests := [][3]*Polynomial{ // f, g, expected sum
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
	field := defineField(7, t)
	ring := DefRing(field)

	f := ring.PolynomialFromUnsigned([]uint{2, 0, 3})
	g := ring.PolynomialFromUnsigned([]uint{1, 1, 1, 1})
	h := ring.PolynomialFromSigned([]int{-1, 0, 0, 3})

	tests := [][3]*Polynomial{ // f, g, expected difference
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
				test[0], test[1], res.coefs, test[2].coefs,
			)
		}
	}
}

func TestMult(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field)

	f := ring.PolynomialFromUnsigned([]uint{2, 0, 3})
	g := ring.PolynomialFromUnsigned([]uint{1, 1, 1, 1})
	h := ring.PolynomialFromSigned([]int{-1, 0, 0, 3})

	tests := [][3]*Polynomial{ // f, g, expected product
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
		res := test[0].Mult(test[1])
		if !res.Equal(test[2]) {
			t.Errorf(
				"(%v) * (%v) = %v, but expected %v",
				test[0], test[1], res, test[2],
			)
		}
	}
}
