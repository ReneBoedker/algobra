package bivariate

import (
	"algobra/errors"
	"algobra/primefield"
)

// Interpolate computes an interpolation polynomial evaluating to values in the
// specified points.
//
// It returns an InputValue-error if the number of points and values differ, or
// if points are not distinct.
func (r *QuotientRing) Interpolate(
	points [][2]*primefield.Element,
	values []*primefield.Element,
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
func allDistinct(points [][2]*primefield.Element) bool {
	unique := make(map[[2]uint]struct{})
	for _, p := range points {
		asUints := [2]uint{p[0].Uint(), p[1].Uint()}
		if _, ok := unique[asUints]; ok {
			return false
		}
		unique[asUints] = struct{}{}
	}
	return true
}

// lagrangeBasis computes a "lagrange-type" basis element. That is, it computes
// a polynomial that evaluates to 1 in point and to 0 in any other point.
func (r *QuotientRing) lagrangeBasis(point [2]*primefield.Element) *Polynomial {
	f := r.New(map[[2]uint]uint{
		{0, 0}: 1,
	})

	for i := 0; i < 2; i++ {
		for j := uint(0); j < r.baseField.Char(); j++ {
			if j == point[i].Uint() {
				continue
			}
			ld := [2]uint{0, 0}
			ld[i] = 1

			f = f.Mult(r.New(map[[2]uint]uint{
				ld:     1,
				{0, 0}: r.baseField.Char() - j,
			})).Scale(point[i].Minus(r.baseField.Element(j)).Inv())
		}
	}

	return f
}
