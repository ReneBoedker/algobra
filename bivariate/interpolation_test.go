package bivariate

import (
	"algobra/finitefield"
	"testing"
)

// Pseudo-random generator prg defined in polynomial_test.go

func TestLagrangeBasis(t *testing.T) {
	field := defineField(17, t)
	ring := DefRing(field, Lex(true))

	for i := 0; i < 50; i++ {
		point := [2]*finitefield.Element{
			field.ElementFromUnsigned(uint(prg.Uint32())),
			field.ElementFromUnsigned(uint(prg.Uint32())),
		}

		f := ring.lagrangeBasis(point)

		for _, a := range field.Elements() {
			for _, b := range field.Elements() {
				ev := f.Eval([2]*finitefield.Element{a, b})
				isPoint := a.Equal(point[0]) && b.Equal(point[1])
				switch {
				case isPoint && !ev.One():
					t.Errorf("f(%v)=%v instead of 1", point, ev)
				case !isPoint && ev.Nonzero():
					t.Errorf("f([%v, %v])=%v instead of 0", a, b, ev)
				}
			}
		}
	}
}

func TestInterpolation(t *testing.T) {
	field := defineField(7, t)
	ring := DefRing(field, Lex(false))
	for i := 0; i < 10; i++ {
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

		values := make([]*finitefield.Element, nPoints, nPoints)
		for j := 0; j < nPoints; j++ {
			values[j] = field.ElementFromUnsigned(uint(prg.Uint32()))
		}

		f, err := ring.Interpolate(points, values)

		if err != nil {
			t.Errorf("Interpolation returned error: %q", err)
		} else {
			// Test that all evaluations are correct
			var testVals [nPoints]*finitefield.Element
			overAllSuccess := true
			for j, p := range points {
				testVals[j] = f.Eval(p)
				overAllSuccess = overAllSuccess && testVals[j].Equal(values[j])
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
