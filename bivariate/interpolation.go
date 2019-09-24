package bivariate

import (
	"algobra/errors"
	"algobra/finitefield"
	"fmt"
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
	distinct := distinct(points)

	f := r.zeroWithCap(2 * len(points))
	for i, p := range points {
		tmp := r.zeroWithCap(2 * len(points))
		tmp.SetCoefPtr([2]uint{0, 0}, r.baseField.One())
		for j := 0; j < 2; j++ {
			tmp.Mult(r.lagrangeBasis2(distinct, p[j], j))
			if tmp.Zero() {
				panic(fmt.Sprintf("Distinct: %v\np[j]: %v\nj: %v", distinct, p[j], j))
			}
		}
		// fmt.Printf("\n\nIndex %v,\tf = %v\n", i, tmp)
		// for _, j := range points {
		// 	fmt.Printf("f(%v) = %v\n", j, tmp.Eval(j))
		// }
		f.Add(tmp.ScaleInPlace(values[i]))

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

// distinct returns the distinct X- and Y-values
func distinct(points [][2]*finitefield.Element) (out [2][]*finitefield.Element) {
	for i := 0; i < 2; i++ {
		unique := make(map[string]*finitefield.Element, len(points))
		for j := range points {
			if _, ok := unique[points[j][i].String()]; ok {
				continue
			}
			unique[points[j][i].String()] = points[j][i]
		}
		// Transfer the keys to the output
		out[i] = make([]*finitefield.Element, 0, len(unique))
		for _, e := range unique {
			out[i] = append(out[i], e)
		}
	}
	return out
}

func (r *QuotientRing) lagrangeBasis2(
	points [2][]*finitefield.Element,
	ignore *finitefield.Element,
	variable int,
) *Polynomial {
	var deg func(int, int) [2]uint
	if variable == 0 {
		deg = func(variable, i int) [2]uint {
			return [2]uint{uint(i), 0}
		}
	} else {
		deg = func(variable, i int) [2]uint {
			return [2]uint{0, uint(i)}
		}
	}
	f := r.zeroWithCap(len(points))
	denom := r.baseField.One()
	ignoreIndex := 0
	for i, p := range points[variable] {
		if p.Equal(ignore) {
			ignoreIndex = i
		}
	}
	for k := 0; k < len(points[variable]); k++ {
		f.SetCoefPtr(deg(variable, k), r.CoefK(points[variable], ignoreIndex, len(points[variable])-1-k))
	}
	// Compute the denominator
	for i, p := range points[variable] {
		if i == ignoreIndex {
			continue
		}
		denom.Mult(ignore.Minus(p))
	}

	if denom.Zero() {
		panic(fmt.Sprint(points))
	}

	f.ScaleInPlace(denom.Inv())
	return f
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

// combinIter is an iterator for combinations
type combinIter struct {
	n     int
	slice []int
	atEnd bool
}

func newCombinIter(n, k int) *combinIter {
	s := make([]int, k, k)
	for i := range s {
		s[i] = i
	}
	return &combinIter{
		n:     n,
		slice: s,
	}
}

func (ci *combinIter) current() []int {
	return ci.slice
}

func (ci *combinIter) active() bool {
	return !ci.atEnd
}

func (ci *combinIter) next() {
	for i := range ci.slice {
		j := len(ci.slice) - 1 - i
		if ci.slice[j] < (ci.n - i) {
			ci.slice[j]++
			for l := 1; l <= i; l++ {
				ci.slice[j+l] = ci.slice[j] + l
			}
			return
		}
	}
	ci.atEnd = true
}

func (r *QuotientRing) CoefK(points []*finitefield.Element, ignore, k int) *finitefield.Element {
	out := r.baseField.Zero()

outer:
	for ci := newCombinIter(len(points)-1, k); ci.active(); ci.next() {
		tmp := r.baseField.One()

		for _, i := range ci.current() {
			if i == ignore {
				continue outer
			}
			tmp.Mult(points[i])
		}

		out.Add(tmp)
	}
	if k%2 != 0 {
		return out.Neg()
	}
	return out
}
