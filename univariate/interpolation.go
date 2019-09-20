package univariate

import (
	"algobra/errors"
	"algobra/finitefield"
)

// Interpolate computes the Lagrange interpolation polynomial evaluating to
// values in the specified points. The resulting polynomial has degree at most
// len(points).
//
// It returns an InputValue-error if the number of points and values differ, or
// if points are not distinct.
func (r *QuotientRing) Interpolate(
	points []*finitefield.Element,
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
	return f, f.Err()
}

// allDistinct checks if given points are all distinct
func allDistinct(points []*finitefield.Element) bool {
	unique := make(map[string]struct{})
	for _, p := range points {
		if _, ok := unique[p.String()]; ok {
			return false
		}
		unique[p.String()] = struct{}{}
	}
	return true
}

// lagrangeBasis computes a "lagrange-type" basis element. That is, it computes
// a polynomial that evaluates to 1 in point at index and to 0 in any other
// point of points.
func (r *QuotientRing) lagrangeBasis(points []*finitefield.Element, index int) *Polynomial {
	f := r.PolynomialFromUnsigned([]uint{1})

	for i, p := range points {
		if i == index {
			continue
		}

		f = f.Mult(r.Polynomial([]*finitefield.Element{
			p.Neg(),
			r.baseField.One(),
		})).Scale(points[index].Minus(p).Inv())
	}

	return f
}
