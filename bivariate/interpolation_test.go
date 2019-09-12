package bivariate

import (
	"algobra/primefield"
	"testing"
)

// Pseudo-random generator prg defined in polynomial_test.go

func TestInterpolation(t *testing.T) {
	field := defineField(23, t)
	ring := DefRing(field, Lex(false))
	for i := 0; i < 10; i++ {
		const nPoints = 5
		points := make([][2]*primefield.Element, nPoints, nPoints)

		nRuns := 0
		for nRuns == 0 || !allDistinct(points) {
			for j := 0; j < nPoints; j++ {
				points[j][0] = field.Element(uint(prg.Uint32()))
				points[j][1] = field.Element(uint(prg.Uint32()))
			}
			nRuns++
			if nRuns >= 100 {
				t.Log("Skipping generation after 100 attempts")
				break
			}
		}

		values := make([]*primefield.Element, nPoints, nPoints)
		for j := 0; j < nPoints; j++ {
			values[j] = field.Element(uint(prg.Uint32()))
		}

		f, err := ring.Interpolate(points, values)

		if err != nil {
			t.Errorf("Interpolation returned error: %q", err)
		} else {
			// Test that all evaluations are correct
			var testSuccess [nPoints]bool
			overAllSuccess := true
			for j, p := range points {
				testSuccess[j] = f.Eval(p).Equal(values[j])
				overAllSuccess = overAllSuccess && testSuccess[j]
			}
			if !overAllSuccess {
				t.Errorf(
					"Interpolation failed for points = %v and values = %v",
					points, values,
				)
			}
		}
	}
}
