package bivariate

import (
	"algobra/auxmath"
	"algobra/errors"
	"algobra/finitefield/ff"
)

// Interpolate computes an interpolation polynomial evaluating to values in the
// specified points. The resulting polynomial has degree at most 2*len(points)
//
// It returns an InputValue-error if the number of points and values differ, or
// if points are not distinct.
func (r *QuotientRing) Interpolate(
	points [][2]ff.Element,
	values []ff.Element,
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
	dist := distinct(points)

	f := r.zeroWithCap(2 * len(points))
	for i, p := range points {
		if values[i].IsZero() {
			// No need to compute the Lagrange basis polynomial since we will
			// scale it by zero anyway
			continue
		}
		tmp := r.zeroWithCap(2 * len(points))
		tmp.SetCoefPtr([2]uint{0, 0}, r.baseField.One())
		for j := 0; j < 2; j++ {
			tmp.Mult(r.lagrangeBasis(dist, p[j], j))
		}
		f.Add(tmp.ScaleInPlace(values[i]))

	}
	if f.Err() != nil {
		return f, errors.Wrap(op, errors.Inherit, f.Err())
	}

	return f, nil
}

// allDistinct checks if given points are all distinct
func allDistinct(points [][2]ff.Element) bool {
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
func distinct(points [][2]ff.Element) (out [2][]ff.Element) {
	for i := 0; i < 2; i++ {
		unique := make(map[string]ff.Element, len(points))
		for j := range points {
			if _, ok := unique[points[j][i].String()]; ok {
				continue
			}
			unique[points[j][i].String()] = points[j][i]
		}
		// Transfer the keys to the output
		out[i] = make([]ff.Element, 0, len(unique))
		for _, e := range unique {
			out[i] = append(out[i], e)
		}
	}
	return out
}

// lagrangeBasis computes a "lagrange-type" basis element in one-variable. That
// is, it computes a polynomial that evaluates to 1 in ignore and to 0 in all
// elements other of points corresponding to given variable.
func (r *QuotientRing) lagrangeBasis(
	points [2][]ff.Element,
	ignore ff.Element,
	variable int,
) *Polynomial {
	// deg gives the monomial with given univariate degree
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

	// Find the index of ignore-element
	ignoreIndex := 0
	for i, p := range points[variable] {
		if p.Equal(ignore) {
			ignoreIndex = i
		}
	}

	// Compute the coefficients directly
	for k := 0; k < len(points[variable]); k++ {
		f.SetCoefPtr(
			deg(variable, k),
			r.coefK(points[variable], ignoreIndex, len(points[variable])-1-k),
		)
	}

	// Compute the denominator
	for i, p := range points[variable] {
		if i == ignoreIndex {
			continue
		}
		denom.Mult(ignore.Minus(p))
	}

	f.ScaleInPlace(denom.Inv())
	return f
}

// coefK computes the coefficient of X^k or Y^k in the numerator of a Lagrange
// basis polynomial. Such polynomials have the form (X-p_1)(X-p_2)...(X-p_n),
// where we skip the p_i corresponding to ignore.
func (r *QuotientRing) coefK(points []ff.Element, ignore, k int) ff.Element {
	out := r.baseField.Zero()
	tmp := r.baseField.Zero()

outer:
	for ci := auxmath.NewCombinIter(len(points), k); ci.Active(); ci.Next() {
		tmp.SetUnsigned(1)

		for _, i := range ci.Current() {
			if i == ignore {
				continue outer
			}
			tmp.Mult(points[i])
		}

		out.Add(tmp)
	}
	if k%2 != 0 {
		// (-1)^k == -1
		out.SetNeg()
	}
	return out
}

/* Copyright 2019 René Bødker Christensen
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
