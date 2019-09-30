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
