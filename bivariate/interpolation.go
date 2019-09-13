package bivariate

import (
	"algobra/errors"
	"algobra/finitefield"
)

// Interpolate computes an interpolation polynomial evaluating to values in the
// specified points. The resulting polynomial has degree at most 2*len(points)
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
	for i := range points {
		f = f.Plus(r.lagrangeBasis(points, i).Scale(values[i]))
	}
	if f.Err() != nil {
		panic(f.Err())
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
// a polynomial that evaluates to 1 in point at index and to 0 in any other
// point of points.
func (r *QuotientRing) lagrangeBasis(points [][2]*finitefield.Element, index int) *Polynomial {
	f := r.PolynomialFromUnsigned(map[[2]uint]uint{
		{0, 0}: 1,
	})

	for i := 0; i < 2; i++ {
		ld := [2]uint{0, 0}
		ld[i] = 1
		for _, p := range points {
			if p[i].Equal(points[index][i]) {
				continue
			}

			f = f.Mult(r.Polynomial(map[[2]uint]*finitefield.Element{
				ld:     r.baseField.One(),
				{0, 0}: p[i].Neg(),
			})).Scale(points[index][i].Minus(p[i]).Inv())
		}
	}

	return f
}
