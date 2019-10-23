package bivariate

import (
	"algobra/errors"
	"testing"
)

func TestHasErr(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field, Lex(false))

	f := ring.Zero()
	f.err = errors.New("Testing", errors.Internal, "Message")

	tmp := hasErr("", f, ring.Zero())
	if tmp.Err() == nil {
		t.Errorf("hasErr returned polynomial with nil error status")
	} else if !errors.Is(errors.Internal, tmp.Err()) {
		t.Errorf("hasErr returned polynomial with error status of unexpected "+
			"kind(err = %v)", tmp.Err())
	}

	tmp = hasErr("", ring.Zero(), f)
	if tmp.Err() == nil {
		t.Errorf("hasErr returned polynomial with nil error status")
	} else if !errors.Is(errors.Internal, tmp.Err()) {
		t.Errorf("hasErr returned polynomial with error status of unexpected "+
			"kind(err = %v)", tmp.Err())
	}
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
		func(f, g *Polynomial) error {
			_, err := ring1.NewIdeal(f, g)
			return err
		},
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
		if err == nil {
			t.Errorf("Function %d did not return an error", i+1)
		} else if !errors.Is(errors.ArithmeticIncompat, err) && !errors.Is(errors.InputIncompatible, err) {
			t.Errorf("Function %d returned an error, but unexpected kind "+
				"(err = %v)", i+1, err)
		}
	}
}

func TestEmptyIdeal(t *testing.T) {
	field := defineField(3, t)
	ring := DefRing(field, Lex(false))

	_, err := ring.NewIdeal(ring.Zero())

	if err == nil {
		t.Errorf("Ideal defined successfully even though all generators are zero")
	} else if !errors.Is(errors.InputValue, err) {
		t.Errorf(
			"Ideal definition returned an error, but not of the expected "+
				"kind (err = %v)", err)
	}
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

	funcs := [](func(*Polynomial, *Polynomial) error){
		func(f, g *Polynomial) error {
			_, _, err := f.monomialDivideBy(g)
			return err
		},
		func(f, g *Polynomial) error {
			_, _, err := g.monomialDivideBy(f)
			return err
		},
	}

	for i, fun := range funcs {
		err := fun(f, g)
		if err == nil {
			t.Errorf("Function %d did not return an error", i+1)
		} else if !errors.Is(errors.InputValue, err) {
			t.Errorf("Function %d returned an error, but unexpected kind "+
				"(err = %v)", i+1, err)
		}
	}
}
