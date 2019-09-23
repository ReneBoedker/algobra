package bivariate

import (
	"algobra/errors"
	"algobra/finitefield"
	"fmt"
	"sort"
	"strings"
)

// addDegs computes the component-wise sum of deg1 and deg2
//
// If either component overflows the size of the uint type, the function returns
// an Overflow-error
func addDegs(deg1, deg2 [2]uint) (sum [2]uint, err error) {
	const op = "Adding degrees"
	sum = [2]uint{deg1[0] + deg2[0], deg1[1] + deg2[1]}
	if sum[0] < deg1[0] || sum[1] < deg1[1] {
		err = errors.New(
			op, errors.Overflow,
			"%v + %v overflows uint", deg1, deg2,
		)
	}
	return
}

// subtractDegs computes the component-wise difference deg1 and deg2
//
// The return value ok indicates whether each component of deg1 is at least as
// large as the corresponding component of deg2.
func subtractDegs(deg1, deg2 [2]uint) (deg [2]uint, ok bool) {
	if deg1[0] >= deg2[0] && deg1[1] >= deg2[1] {
		return [2]uint{deg1[0] - deg2[0], deg1[1] - deg2[1]}, true
	}
	return deg, false
}

// Polynomial denotes a bivariate polynomial.
type Polynomial struct {
	baseRing *QuotientRing
	coefs    map[[2]uint]*finitefield.Element
	err      error
}

// BaseField returns the field over which the coefficients of f are defined.
func (f *Polynomial) BaseField() *finitefield.Field {
	return f.baseRing.baseField
}

// Err returns the error status of f.
func (f *Polynomial) Err() error {
	return f.err
}

// Coef returns the coefficient of the monomial with degree specified by the
// input. The return value is a finite field element.
func (f *Polynomial) Coef(deg [2]uint) *finitefield.Element {
	if c, ok := f.coefs[deg]; ok {
		return c
	}
	return f.BaseField().Zero()
}

// SetCoef sets the coefficient of the monomial with degree deg in f to val by
// copying. See also SetCoefPtr.
func (f *Polynomial) SetCoef(deg [2]uint, val *finitefield.Element) {
	if val.Zero() {
		delete(f.coefs, deg)
	} else {
		f.coefs[deg] = val.Copy()
	}
}

// SetCoefPtr sets the coefficient of the monomial with degree deg in f to ptr
// as a pointer. To set coefficient as a value, use SetCoef instead.
func (f *Polynomial) SetCoefPtr(deg [2]uint, ptr *finitefield.Element) {
	if ptr.Zero() {
		delete(f.coefs, deg)
	} else {
		f.coefs[deg] = ptr
	}
}

// IncrementCoef sets the coefficient of the monomial with degree deg in f to val.
func (f *Polynomial) IncrementCoef(deg [2]uint, val *finitefield.Element) {
	if val.Zero() {
		return
	}
	if c, ok := f.coefs[deg]; ok {
		c.Add(val)
		if c.Zero() {
			delete(f.coefs, deg)
		}
	} else {
		f.coefs[deg] = val.Copy()
	}
}

// SetCoef sets the coefficient of the monomial with degree deg in f to ptr.
// func (f *Polynomial) IncrementCoefPtr(deg [2]uint, val *finitefield.Element) {
// 	if val.Zero() {
// 		return
// 	}
// 	if c, ok := f.coefs[deg]; ok {
// 		c.Add(val)
// 		if c.Zero() {
// 			delete(f.coefs, deg)
// 		}
// 	} else {
// 		f.coefs[deg] = val
// 	}
// }

// Copy returns a new polynomial object over the same ring and with the same
// coefficients as f.
func (f *Polynomial) Copy() *Polynomial {
	h := f.baseRing.zeroWithCap(len(f.coefs))
	for deg, c := range f.coefs {
		h.coefs[deg] = c.Copy()
	}
	return h
}

func (f *Polynomial) clean() {
	for d, c := range f.coefs {
		if c.Zero() {
			delete(f.coefs, d)
		}
	}
}

// Eval evaluates f at the given point.
func (f *Polynomial) Eval(point [2]*finitefield.Element) *finitefield.Element {
	out := f.BaseField().Zero()
	for deg, coef := range f.coefs {
		out = out.Plus(coef.Times(point[0].Pow(deg[0])).Times(point[1].Pow(deg[1])))
	}
	return out
}

// Plus returns the sum of the two polynomials f and g.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Plus(g *Polynomial) *Polynomial {
	const op = "Adding polynomials"

	if tmp := hasErr(op, f, g); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g); tmp != nil {
		return tmp
	}

	h := f.Copy()
	for deg, c := range g.coefs {
		h.IncrementCoef(deg, c)
	}
	return h
}

// Add sets f to the sum of the two polynomials f and g and returns f.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Add(g *Polynomial) *Polynomial {
	const op = "Adding polynomials"

	if tmp := hasErr(op, f, g); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g); tmp != nil {
		return tmp
	}

	for deg, c := range g.coefs {
		f.IncrementCoef(deg, c)
	}
	return f
}

// Neg returns the polynomial obtained by scaling f by -1 (modulo the
// characteristic).
func (f *Polynomial) Neg() *Polynomial {
	g := f.baseRing.zeroWithCap(len(f.coefs))
	for deg, c := range f.coefs {
		g.coefs[deg] = c.Neg()
	}
	return g
}

// Equal determines whether two polynomials are equal. That is, whether they are
// defined over the same ring, and have the same coefficients.
func (f *Polynomial) Equal(g *Polynomial) bool {
	if f.baseRing != g.baseRing {
		return false
	}
	if len(f.coefs) != len(g.coefs) {
		return false
	}
	for d, cf := range f.coefs {
		if cg, ok := g.coefs[d]; !ok || !cg.Equal(cf) {
			return false
		}
	}
	return true
}

// Minus returns polynomial difference f-g.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Minus(g *Polynomial) *Polynomial {
	const op = "Subtracting polynomials"

	if tmp := hasErr(op, f, g); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g); tmp != nil {
		return tmp
	}

	return f.Plus(g.Neg())
}

// Internal method. Multiplies the two polynomials f and g, but does not reduce
// the result according to the specified ring.
func (f *Polynomial) multNoReduce(g *Polynomial) *Polynomial {
	const op = "Multiplying polynomials"

	if tmp := hasErr(op, f, g); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g); tmp != nil {
		return tmp
	}

	h := f.baseRing.zeroWithCap(len(f.coefs) * len(g.coefs))
	tmp := f.BaseField().One()
	for degf, cf := range f.coefs {
		for degg, cg := range g.coefs {
			degSum, err := addDegs(degf, degg)
			if err != nil {
				h = f.baseRing.Zero()
				h.err = errors.Wrap(op, errors.Inherit, err)
				return h
			}
			tmp.Prod(cf, cg)
			h.IncrementCoef(degSum, tmp)
		}
	}
	return h
}

// Times returns the product of the polynomials f and g
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Times(g *Polynomial) *Polynomial {
	h := f.multNoReduce(g)
	if h.Err() != nil {
		return h
	}
	h.reduce()
	return h
}

// Mult sets f to the product of the polynomials f and g and returns f.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Mult(g *Polynomial) *Polynomial {
	*f = *f.multNoReduce(g)
	if f.Err() != nil {
		return f
	}
	f.reduce()
	return f
}

// Normalize creates a new polynomial obtained by normalizing f. That is,
// f.Normalize() multiplied by f.Lc() is f.
//
// If f is the zero polynomial, a copy of f is returned.
func (f *Polynomial) Normalize() *Polynomial {
	if f.Zero() {
		return f.Copy()
	}
	return f.Scale(f.Lc().Inv())
}

// Scale scales all coefficients of f by the field element c and returns the
// result as a new polynomial.
func (f *Polynomial) Scale(c *finitefield.Element) *Polynomial {
	g := f.Copy()
	for d := range g.coefs {
		g.coefs[d] = g.coefs[d].Times(c)
	}
	return g
}

// ScaleInPlace scales all coefficients of f by the field element c and returns
// f.
func (f *Polynomial) ScaleInPlace(c *finitefield.Element) *Polynomial {
	for d := range f.coefs {
		f.coefs[d].Mult(c)
	}
	return f
}

// Pow raises f to the power of n.
//
// If the computation causes the degree of f to overflow, the returned
// polynomial has an Overflow-error as error status.
func (f *Polynomial) Pow(n uint) *Polynomial {
	const op = "Computing polynomial power"

	out := f.baseRing.Polynomial(map[[2]uint]*finitefield.Element{
		{0, 0}: f.BaseField().One(),
	})
	g := f.Copy()

	for n > 0 {
		if n%2 == 1 {
			out.Mult(g)
			if out.Err() != nil {
				out = f.baseRing.Zero()
				out.err = errors.Wrap(op, errors.Inherit, out.Err())
				return out
			}
		}
		n /= 2
		g.Mult(g)
	}
	return out
}

// SortedDegrees returns a list containing the degrees is the support of f.
//
// The list is sorted according to the ring order with higher orders preceding
// lower orders in the list.
func (f *Polynomial) SortedDegrees() [][2]uint {
	degs := make([][2]uint, 0, len(f.coefs))
	for deg := range f.coefs {
		degs = append(degs, deg)
	}
	sort.Slice(degs, func(i, j int) bool {
		return (f.baseRing.ord(degs[i], degs[j]) >= 0)
	})
	return degs
}

// Ld returns the leading degree of f.
func (f *Polynomial) Ld() [2]uint {
	return f.SortedDegrees()[0]
}

// Lc returns the leading coefficient of f.
func (f *Polynomial) Lc() *finitefield.Element {
	return f.Coef(f.Ld())
}

// Lt returns the leading term of f.
func (f *Polynomial) Lt() *Polynomial {
	h := f.baseRing.Zero()
	ld := f.Ld()
	h.coefs[ld] = f.Coef(ld)
	return h
}

// Zero determines whether f is the zero polynomial.
func (f *Polynomial) Zero() bool {
	if len(f.coefs) == 0 {
		return true
	}
	return false
}

// Nonzero determines whether f contains some monomial with nonzero coefficient.
func (f *Polynomial) Nonzero() bool {
	return !f.Zero()
}

// Monomial returns a bool describing whether f consists of a single monomial.
func (f *Polynomial) Monomial() bool {
	if len(f.coefs) == 1 {
		return true
	}
	return false
}

// Reduces f in-place
func (f *Polynomial) reduce() {
	if f.baseRing.id != nil {
		f.baseRing.id.Reduce(f)
	}
}

// String returns the string representation of f. Variables are named 'X' and
// 'Y'.
func (f *Polynomial) String() string {
	degs := f.SortedDegrees()
	if len(degs) == 0 {
		return "0"
	}
	var b strings.Builder
	for i, d := range degs {
		if i > 0 {
			fmt.Fprint(&b, " + ")
		}
		if tmp := f.Coef(d); !tmp.One() || (d[0] == 0 && d[1] == 0) {
			fmt.Fprintf(&b, "%v", tmp)
		}
		if d[0] == 1 {
			fmt.Fprint(&b, "X")
		}
		if d[0] > 1 {
			fmt.Fprintf(&b, "X^%d", d[0])
		}
		if d[1] == 1 {
			fmt.Fprint(&b, "Y")
		}
		if d[1] > 1 {
			fmt.Fprintf(&b, "Y^%d", d[1])
		}
	}
	return b.String()
}

// hasErr is an internal method for checking if f or g has a non-nil error
// field.
//
// It returns the first polynomial with non-nil error status after wrapping the
// error. The new error inherits the kind from the old.
func hasErr(op errors.Op, f, g *Polynomial) *Polynomial {
	switch {
	case f.err != nil:
		f.err = errors.Wrap(
			op, errors.Inherit,
			f.err,
		)
		return f
	case g.err != nil:
		g.err = errors.Wrap(
			op, errors.Inherit,
			g.err,
		)
		return g
	}
	return nil
}

// checkCompatible is an internal method for checking if f and g are compatible;
// that is, if they are defined over the same ring.
//
// If not, the return value is an element with error status set to
// ArithmeticIncompat.
func checkCompatible(op errors.Op, f, g *Polynomial) *Polynomial {
	if f.baseRing != g.baseRing {
		out := f.baseRing.Zero()
		out.err = errors.New(
			op, errors.ArithmeticIncompat,
			"%v and %v defined over different rings", f, g,
		)
		return out
	}
	return nil
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
