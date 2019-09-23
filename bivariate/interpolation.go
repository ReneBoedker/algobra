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

	f := r.zeroWithCap(2 * len(points))
	for i := range points {
		f.Add(r.lagrangeBasis(points, i).ScaleInPlace(values[i]))
	}
	if f.Err() != nil {
		panic(f.Err())
	}

	return f, nil
}

// allDistinct checks if given points are all distinct
func allDistinct(points [][2]*finitefield.Element) bool {
	unique := make(map[[2]string]struct{}, len(points))
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
	f := r.zeroWithCap(2 * len(points))
	f.SetCoefPtr([2]uint{0, 0}, r.baseField.One())

	denominator := r.baseField.ElementFromUnsigned(1)
	for i := 0; i < 2; i++ {
		ld := [2]uint{0, 0}
		ld[i] = 1
		tmp := r.Zero()
		tmp.SetCoefPtr(ld, r.baseField.One())
		for _, p := range points {
			if p[i].Equal(points[index][i]) {
				continue
			}
			tmp.SetCoefPtr([2]uint{0, 0}, p[i].Neg())
			f.Mult(tmp)
			denominator.Mult(points[index][i].Minus(p[i]))
		}
	}

	f.ScaleInPlace(denominator.Inv())

	return f
}
