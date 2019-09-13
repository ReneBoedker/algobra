package bivariate

import (
	"algobra/finitefield"
	"testing"
)

// Pseudo-random generator prg defined in polynomial_test.go

func TestLagrangeBasis(t *testing.T) {
	field := defineField(17, t)
	ring := DefRing(field, Lex(true))

	for rep := 0; rep < 100; rep++ {
		const nPoints = 2
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

		for j := range points {
			// f evaluates to 1 in points[j] and 0 in other components of 0
			f := ring.lagrangeBasis(points, j)

			if f.Zero() {
				t.Errorf("Lagrange basis is zero with points %v and index %d",
					points, j)
			} else if ld := f.Ld(); ld[0]+ld[1] > 2*(nPoints-1) {
				t.Errorf("Lagrange basis has too large total degree (%d) with point %v",
					ld[0]+ld[1], points[j],
				)
			}

			for k, p := range points {
				ev := f.Eval(p)
				switch {
				case j == k && !ev.One():
					t.Errorf("f(%v)=%v instead of 1 with f = %v", p, ev, f)
				case j != k && ev.Nonzero():
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
