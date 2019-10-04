package univariate

import (
	"algobra/errors"
	"algobra/finitefield/ff"
)

// Interpolate computes the Lagrange interpolation polynomial evaluating to
// values in the specified points. The resulting polynomial has degree at most
// len(points).
//
// It returns an InputValue-error if the number of points and values differ, or
// if points are not distinct.
func (r *QuotientRing) Interpolate(
	points []ff.Element,
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
			"Interpolation points must be distinct",
		)
	}

	f := r.zeroWithCap(len(points))
	for i, p := range points {
		if values[i].IsZero() {
			// No need to compute the Lagrange basis polynomial since we will
			// scale it by zero anyway
			continue
		}
		f.Add(r.lagrangeBasis(points, p).Scale(values[i]))
	}
	return f, f.Err()
}

// allDistinct checks if given points are all distinct
func allDistinct(points []ff.Element) bool {
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
// func (r *QuotientRing) lagrangeBasis(points []ff.Element, index int) *Polynomial {
// 	f := r.PolynomialFromUnsigned([]uint{1})

// 	for i, p := range points {
// 		if i == index {
// 			continue
// 		}

// 		f.Mult(r.Polynomial([]ff.Element{
// 			p.Neg(),
// 			r.baseField.One(),
// 		})).Scale(points[index].Minus(p).Inv())
// 	}

// 	return f
// }

func (r *QuotientRing) lagrangeBasis(
	points []ff.Element,
	ignore ff.Element,
) *Polynomial {
	f := r.zeroWithCap(len(points))
	denom := r.baseField.One()

	// Find the index of ignore-element
	ignoreIndex := 0
	for i, p := range points {
		if p.Equal(ignore) {
			ignoreIndex = i
		}
	}

	// Compute the coefficients directly
	for k := 0; k < len(points); k++ {
		f.SetCoef(
			k,
			r.coefK(points, ignoreIndex, len(points)-1-k),
		)
	}

	// Compute the denominator
	for i, p := range points {
		if i == ignoreIndex {
			continue
		}
		denom.Mult(ignore.Minus(p))
	}

	f = f.Scale(denom.Inv())
	return f
}

// combinIter is an iterator for combinations.
//
// It will iterate over all possible ways to choose a given number of elements
// from n elements. For instance, it will generate the sequence [0,1,2],
// [0,1,3],..., [3,4,5] if defined with n=5 and k=3.
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

// coefK computes the coefficient of X^k or Y^k in the numerator of a Lagrange
// basis polynomial. Such polynomials have the form (X-p_1)(X-p_2)...(X-p_n),
// where we skip the p_i corresponding to ignore.
func (r *QuotientRing) coefK(points []ff.Element, ignore, k int) ff.Element {
	out := r.baseField.Zero()
	tmp := r.baseField.Zero()

outer:
	for ci := newCombinIter(len(points)-1, k); ci.active(); ci.next() {
		tmp.SetUnsigned(1)

		for _, i := range ci.current() {
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
