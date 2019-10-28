package bivariate

import (
	"testing"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

func TestLagrangeBasis(t *testing.T) {
	field := defineField(17, t)
	ring := DefRing(field, Lex(true))

	for rep := 0; rep < 100; rep++ {
		const nPoints = 4
		points := make([][2]ff.Element, nPoints, nPoints)

		nRuns := 0
		for nRuns == 0 || !allDistinct(points) {
			for j := 0; j < nPoints; j++ {
				points[j][0] = field.RandElement()
				points[j][1] = field.RandElement()
			}
			nRuns++
			if nRuns >= 100 {
				t.Log("Skipping generation after 100 attempts")
				break
			}
		}

		dist := distinct(points)

		for i := range points {
			f := ring.PolynomialFromUnsigned(map[[2]uint]uint{{0, 0}: 1})
			for j := 0; j < 2; j++ {
				f.Mult(ring.lagrangeBasis(dist, points[i][j], j))
			}
			// f now evaluates to 1 in points[j] and 0 in other points

			if f.IsZero() {
				t.Errorf(
					"Lagrange basis is zero with points %v and index %d",
					points, i,
				)
			} else if ld := f.Ld(); ld[0]+ld[1] > 2*(nPoints-1) {
				t.Errorf(
					"Lagrange basis has too large total degree (%d) with point %v",
					ld[0]+ld[1], points[i],
				)
			}

			for k, p := range points {
				ev := f.Eval(p)
				switch {
				case i == k && !ev.IsOne():
					t.Errorf("f(%v)=%v instead of 1 with f = %v", p, ev, f)
				case i != k && ev.IsNonzero():
					t.Errorf("f(%v)=%v instead of 0 with f = %v", p, ev, f)
				}
			}
		}
	}
}

func TestInterpolation(t *testing.T) {
	field := defineField(23, t)
	ring := DefRing(field, Lex(false))
	for rep := 0; rep < 100; rep++ {
		const nPoints = 4
		points := make([][2]ff.Element, nPoints, nPoints)

		nRuns := 0
		for nRuns == 0 || !allDistinct(points) {
			for j := 0; j < nPoints; j++ {
				points[j][0] = field.RandElement()
				points[j][1] = field.RandElement()
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
	field := defineField(13, t)
	ring := DefRing(field, DegLex(true))

	a := [2]ff.Element{field.Zero(), field.Zero()}
	b := [2]ff.Element{field.Zero(), field.ElementFromUnsigned(5)}

	_, err := ring.Interpolate(
		[][2]ff.Element{a, b},
		[]ff.Element{field.Zero()},
	)
	assertError(t, err, errors.InputValue, "Interpolation with more points than values")

	_, err = ring.Interpolate(
		[][2]ff.Element{a, a},
		[]ff.Element{field.Zero(), field.Zero()},
	)
	assertError(t, err, errors.InputValue, "Interpolation on duplicate points")
}

// func TestIter(t *testing.T) {
// 	for ci := newCombinIter(5, 3); ci.Active(); ci.Next() {
// 		fmt.Println(ci.slice)
// 	}
// }
