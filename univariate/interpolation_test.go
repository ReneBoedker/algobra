package univariate_test

import (
	"testing"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/univariate"
)

func TestLagrangeBasis(t *testing.T) {
	//t.Fatalf("Skipped")
	do := func(field ff.Field) {
		ring := univariate.DefRing(field)

	outer:
		for rep := 0; rep < 25; rep++ {
			const nPoints = 4
			points := make([]ff.Element, nPoints, nPoints)

			nRuns := 0
			for nRuns == 0 || !univariate.AllDistinct(points) {
				for j := 0; j < nPoints; j++ {
					points[j] = field.RandElement()
				}
				nRuns++
				if nRuns >= 100 {
					t.Log("Skipping generation after 100 attempts")
					continue outer
				}
			}

			for j := range points {
				// f evaluates to 1 in points[j] and 0 in other components of 0
				f := ring.LagrangeBasis(points, points[j])

				if f.IsZero() {
					t.Errorf(
						"GF(%d): Lagrange basis is zero with points %v and index %d",
						field.Card(), points, j,
					)
				} else if ld := f.Ld(); ld > (nPoints - 1) {
					t.Errorf(
						"GF(%d): Lagrange basis has too large total degree (%d) with point %v",
						field.Card(), ld, points[j],
					)
				}

				for k, p := range points {
					ev := f.Eval(p)
					switch {
					case j == k && !ev.IsOne():
						t.Errorf("GF(%d): f(%v)=%v instead of 1 with f = %v", field.Card(), p, ev, f)
					case j != k && ev.IsNonzero():
						t.Errorf("GF(%d): f(%v)=%v instead of 0 with f = %v", field.Card(), p, ev, f)
					}
				}
			}
		}
	}

	fieldLoop(do, 4)
}

func TestInterpolation(t *testing.T) {
	field := defineField(23)
	ring := univariate.DefRing(field)
	for rep := 0; rep < 100; rep++ {
		const nPoints = 4
		points := make([]ff.Element, nPoints, nPoints)

		nRuns := 0
		for nRuns == 0 || !univariate.AllDistinct(points) {
			for j := 0; j < nPoints; j++ {
				points[j] = field.RandElement()
			}
			nRuns++
			if nRuns >= 100 {
				t.Log("Skipping generation after 100 attempts")
				break
			}
		}

		values := make([]ff.Element, nPoints, nPoints)
		for i := 0; i < nPoints; i++ {
			values[i] = field.RandElement()
		}

		f, err := ring.Interpolate(points, values)

		if err != nil {
			t.Errorf("Interpolation returned error: %q", err)
		} else {
			// Test that all evaluations are correct
			var testVals [nPoints]ff.Element
			overAllSuccess := true
			for i, p := range points {
				testVals[i] = f.Eval(p)
				overAllSuccess = overAllSuccess && testVals[i].Equal(values[i])
			}
			if !overAllSuccess {
				t.Errorf(
					"Interpolation failed for points = %v and values = %v.\n"+
						"(Returned polynomial %v with values %v)",
					points, values, f, testVals,
				)
			}
		}
	}
}

func TestInterpolationErrors(t *testing.T) {
	field := defineField(13)
	ring := univariate.DefRing(field)

	a := field.Zero()
	b := field.ElementFromUnsigned(5)

	_, err := ring.Interpolate(
		[]ff.Element{a, b},
		[]ff.Element{field.Zero()},
	)
	if err == nil {
		t.Errorf("Interpolation did not return an error even though there " +
			"were more points than values")
	} else if !errors.Is(errors.InputValue, err) {
		t.Errorf(
			"Interpolation returned an error on different length inputs, "+
				"but it was of unexpected kind (err = %v)", err,
		)
	}

	_, err = ring.Interpolate(
		[]ff.Element{a, a},
		[]ff.Element{field.Zero(), field.Zero()},
	)
	if err == nil {
		t.Errorf("Interpolation did not return an error even though input " +
			"contains duplicate points")
	} else if !errors.Is(errors.InputValue, err) {
		t.Errorf(
			"Interpolation returned an error when input contains duplicate"+
				" points, but it was of unexpected kind (err = %v)", err,
		)
	}
}
