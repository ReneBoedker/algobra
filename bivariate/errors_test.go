package bivariate

import (
	"testing"

	"github.com/ReneBoedker/algobra/errors"
)

func TestHasErr(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field, Lex(false))

	f := ring.Zero()
	f.err = errors.New("Testing", errors.Internal, "Message")

	// Check that error in first position is found
	tmp := hasErr("", f, ring.Zero())
	assertError(t, tmp.Err(), errors.Internal, "hasErr")

	// Check that error in second position is found
	tmp = hasErr("", ring.Zero(), f)
	assertError(t, tmp.Err(), errors.Internal, "hasErr")
}

func TestDiffRings(t *testing.T) {
	field := defineField(11, t)
	ring1 := DefRing(field, DegLex(false))
	ring2 := DefRing(field, WDegLex(11, 12, true))

	f := ring1.Zero()
	g := ring2.Zero()

	funcs := [](func(*Polynomial, *Polynomial) error){
		func(f, g *Polynomial) error {
			return checkCompatible("", f, g).Err()
		},
		func(f, g *Polynomial) error {
			_, err := SPolynomial(f, g)
			return err
		},
		// func(f, g *Polynomial) error {
		// 	_, err := ring1.NewIdeal(f, g)
		// 	return err
		// },
		func(f, g *Polynomial) error {
			h := f.Plus(g)
			return h.Err()
		},
		func(f, g *Polynomial) error {
			h := f.Minus(g)
			return h.Err()
		},
		func(f, g *Polynomial) error {
			h := f.Mult(g)
			return h.Err()
		},
	}

	for i, fun := range funcs {
		err := fun(f, g)
		assertError(t, err, errors.ArithmeticIncompat, "Function %d", i+1)
	}
}

func TestEmptyIdeal(t *testing.T) {
	field := defineField(3, t)
	ring := DefRing(field, Lex(false))

	_, err := ring.NewIdeal(ring.Zero())
	assertError(t, err, errors.InputValue, "Defining empty ideal")
}

func TestNotMonomials(t *testing.T) {
	field := defineField(11, t)
	ring := DefRing(field, WDegLex(11, 12, false))

	f := ring.PolynomialFromUnsigned(map[[2]uint]uint{{2, 0}: 1})
	g := ring.PolynomialFromUnsigned(map[[2]uint]uint{{1, 2}: 1, {0, 0}: 4})

	if tmp := monomialLcm(f, g); tmp != nil {
		t.Errorf(
			"monomialLcm returned nil, even though one input was not a monomial.",
		)
	}

	if !f.IsMonomial() {
		t.Errorf("(%v).IsMonomial() returned false", f)
	}

	if g.IsMonomial() {
		t.Errorf("(%v).IsMonomial() returned true", f)
	}
}
