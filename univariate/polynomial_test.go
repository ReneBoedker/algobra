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

func TestArithmetic(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field)

	f := ring.PolynomialFromUnsigned([]uint{2, 0, 3})
	g := ring.PolynomialFromUnsigned([]uint{1, 1, 1, 1})
	h := ring.PolynomialFromSigned([]int{-1, 0, 0, 3})

	res := []*Polynomial{
		f.Plus(g),
		f.Plus(h),
		f.Minus(g),
		f.Minus(h),
		f.Mult(g),
	}
	desc := []string{
		"Plus", "Plus",
		"Minus", "Minus",
		"Multiplication",
	}
	expected := []*Polynomial{
		ring.PolynomialFromUnsigned([]uint{3, 1, 4, 1}),
		ring.PolynomialFromUnsigned([]uint{1, 0, 3, 3}),
		ring.PolynomialFromUnsigned([]uint{1, 6, 2, 6}),
		ring.PolynomialFromSigned([]int{3, 0, 3, -3}),
		ring.PolynomialFromUnsigned([]uint{2, 2, 5, 5, 3, 3}),
	}

	for i := range res {
		if !res[i].Equal(expected[i]) {
			t.Errorf(
				"Test %d (%s) failed. Got %v, but expected %v",
				i, desc[i], res[i], expected[i],
			)
		}
	}
}
