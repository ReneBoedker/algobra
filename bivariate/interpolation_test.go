package bivariate

import (
	"algobra/errors"
	"algobra/finitefield"
	"testing"
)

// Pseudo-random generator prg defined in polynomial_test.go

func TestLagrangeBasis(t *testing.T) {
	field := defineField(17, t)
	ring := DefRing(field, Lex(true))

	for rep := 0; rep < 100; rep++ {
		const nPoints = 4
		points := make([][2]*finitefield.Element, nPoints, nPoints)

		nRuns := 0
		for nRuns == 0 || !allDistinct(points) {
			for j := 0; j < nPoints; j++ {
				points[j][0] = field.ElementFromUnsigned(uint(prg.Uint32()))
				points[j][1] = field.ElementFromUnsigned(uint(prg.Uint32()))
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

			if f.Zero() {
				t.Errorf("Lagrange basis is zero with points %v and index %d",
					points, i)
			} else if ld := f.Ld(); ld[0]+ld[1] > 2*(nPoints-1) {
				t.Errorf("Lagrange basis has too large total degree (%d) with point %v",
					ld[0]+ld[1], points[i],
				)
			}

			for k, p := range points {
				ev := f.Eval(p)
				switch {
				case i == k && !ev.One():
					t.Errorf("f(%v)=%v instead of 1 with f = %v", p, ev, f)
				case i != k && ev.Nonzero():
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
		points := make([][2]*finitefield.Element, nPoints, nPoints)

		nRuns := 0
		for nRuns == 0 || !allDistinct(points) {
			for j := 0; j < nPoints; j++ {
				points[j][0] = field.ElementFromUnsigned(uint(prg.Uint32()))
				points[j][1] = field.ElementFromUnsigned(uint(prg.Uint32()))
			}
			nRuns++
			if nRuns >= 100 {
				t.Log("Skipping generation after 100 attempts")
				break
			}
		}

		values := make([]*finitefield.Element, nPoints, nPoints)
		for i := 0; i < nPoints; i++ {
			values[i] = field.ElementFromUnsigned(uint(prg.Uint32()))
		}

		f, err := ring.Interpolate(points, values)

		if err != nil {
			t.Errorf("Interpolation returned error: %q", err)
		} else {
			// Test that all evaluations are correct
			var testVals [nPoints]*finitefield.Element
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

	a := [2]*finitefield.Element{field.Zero(), field.Zero()}
	b := [2]*finitefield.Element{field.Zero(), field.ElementFromUnsigned(5)}

	_, err := ring.Interpolate(
		[][2]*finitefield.Element{a, b},
		[]*finitefield.Element{field.Zero()})
	if err == nil {
		t.Errorf("Interpolation did not return an error even though there " +
			"were more points than values")
	} else if !errors.Is(errors.InputValue, err) {
		t.Errorf("Interpolation returned an error on different length inputs, "+
			"but it was of unexpected kind (err = %v)", err)
	}

	_, err = ring.Interpolate(
		[][2]*finitefield.Element{a, a},
		[]*finitefield.Element{field.Zero(), field.Zero()})
	if err == nil {
		t.Errorf("Interpolation did not return an error even though input " +
			"contains duplicate points")
	} else if !errors.Is(errors.InputValue, err) {
		t.Errorf("Interpolation returned an error when input contains duplicate"+
			" points, but it was of unexpected kind (err = %v)", err)
	}
}

// func TestIter(t *testing.T) {
// 	for ci := newCombinIter(5, 3); ci.Active(); ci.Next() {
// 		fmt.Println(ci.slice)
// 	}
// }
