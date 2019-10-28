package univariate_test

import (
	"testing"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/univariate"
)

func TestHasErr(t *testing.T) {
	field := defineField(7)
	ring := univariate.DefRing(field)

	f := ring.Zero()
	f.SetError(errors.New("Testing", errors.Internal, "Message"))

	tmp := univariate.HasErr("", f, ring.Zero())
	if tmp.Err() == nil {
		t.Errorf("hasErr returned polynomial with nil error status")
	} else if !errors.Is(errors.Internal, tmp.Err()) {
		t.Errorf(
			"hasErr returned polynomial with error status of unexpected "+
				"kind(err = %v)", tmp.Err(),
		)
	}

	tmp = univariate.HasErr("", ring.Zero(), f)
	if tmp.Err() == nil {
		t.Errorf("hasErr returned polynomial with nil error status")
	} else if !errors.Is(errors.Internal, tmp.Err()) {
		t.Errorf(
			"hasErr returned polynomial with error status of unexpected "+
				"kind(err = %v)", tmp.Err(),
		)
	}
}

func TestDiffRings(t *testing.T) {
	field1 := defineField(11)
	field2 := defineField(7)
	ring1 := univariate.DefRing(field1)
	ring2 := univariate.DefRing(field2)

	f := ring1.Zero()
	g := ring2.Zero()

	funcs := [](func(*univariate.Polynomial, *univariate.Polynomial) error){
		func(f, g *univariate.Polynomial) error {
			return univariate.CheckCompatible("", f, g).Err()
		},
		func(f, g *univariate.Polynomial) error {
			_, err := ring1.NewIdeal(f, g)
			return err
		},
		func(f, g *univariate.Polynomial) error {
			h := f.Plus(g)
			return h.Err()
		},
		func(f, g *univariate.Polynomial) error {
			h := f.Minus(g)
			return h.Err()
		},
		func(f, g *univariate.Polynomial) error {
			h := f.Mult(g)
			return h.Err()
		},
	}

	for i, fun := range funcs {
		err := fun(f, g)
		if err == nil {
			t.Errorf("Function %d did not return an error", i+1)
		} else if !errors.Is(errors.ArithmeticIncompat, err) && !errors.Is(errors.InputIncompatible, err) {
			t.Errorf(
				"Function %d returned an error, but unexpected kind (err = %v)",
				i+1, err,
			)
		}
	}
}
