package univariate

import (
	"algobra/errors"
	"testing"
)

func TestHasErr(t *testing.T) {
	field := defineField(7)
	ring := DefRing(field)

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
	field1 := defineField(11)
	field2 := defineField(7)
	ring1 := DefRing(field1)
	ring2 := DefRing(field2)

	f := ring1.Zero()
	g := ring2.Zero()

	funcs := [](func(*Polynomial, *Polynomial) error){
		func(f, g *Polynomial) error {
			return checkCompatible("", f, g).Err()
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
