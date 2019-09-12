package bivariate

import (
	"algobra/errors"
	"algobra/finitefield"
)

// Interpolate computes an interpolation polynomial evaluating to values in the
// specified points.
//
// It returns an InputValue-error if the number of points and values differ, or
// if points are not distinct.
func (r *QuotientRing) Interpolate(
	points [][2]*finitefield.Element,
	values []*finitefield.Element,
) (*Polynomial, error) {
	const op = "Computing interpolation"

	if len(points) != len(values) {
		return nil, errors.New(
			op, errors.InputValue,
			"Different number of interpolation points and values (%d and %d)",
			len(points), len(values),
		)
	}

	if !allDistinct(points) {
		return nil, errors.New(
			op, errors.InputValue,
			"Interpolation points must me distinct",
		)
	}

	f := r.Zero()
	for i, p := range points {
		f = f.Plus(r.lagrangeBasis(p).Scale(values[i]))
	}

	return f, nil
}

// allDistinct checks if given points are all distinct
func allDistinct(points [][2]*finitefield.Element) bool {
	unique := make(map[[2]string]struct{})
	for _, p := range points {
		asStrings := [2]string{p[0].String(), p[1].String()}
		if _, ok := unique[asStrings]; ok {
			return false
		}
		unique[asStrings] = struct{}{}
	}
	return true
}

// lagrangeBasis computes a "lagrange-type" basis element. That is, it computes
// a polynomial that evaluates to 1 in point and to 0 in any other point.
func (r *QuotientRing) lagrangeBasis(point [2]*finitefield.Element) *Polynomial {
	f := r.Polynomial(map[[2]uint]*finitefield.Element{
		{0, 0}: r.baseField.One(),
	})

	for i := 0; i < 2; i++ {
		for j := uint(0); j < r.baseField.Char(); j++ {
			if j == point[i].Uint() {
				continue
			}
			ld := [2]uint{0, 0}
			ld[i] = 1

			jElem, err := r.baseField.Element(j)
			if err != nil {
				panic(err)
			}

			f = f.Mult(r.Polynomial(map[[2]uint]uint{
				ld:     1,
				{0, 0}: r.baseField.Char() - j,
			})).Scale(point[i].Minus(jElem).Inv())
		}
	}

	return f
}
