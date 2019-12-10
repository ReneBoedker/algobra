package bivariate

import (
	"sort"
	"strconv"
	"strings"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

// Polynomial implements bivariate polynomials.
type Polynomial struct {
	baseRing *QuotientRing
	coefs    map[[2]uint]ff.Element
	err      error
}

// BaseField returns the field over which the coefficients of f are defined.
func (f *Polynomial) BaseField() ff.Field {
	return f.baseRing.baseField
}

// Err returns the error status of f.
func (f *Polynomial) Err() error {
	return f.err
}

// Zero returns a zero polynomial over the specified ring.
func (r *QuotientRing) Zero() *Polynomial {
	return &Polynomial{
		baseRing: r,
		coefs:    make(map[[2]uint]ff.Element),
	}
}

// zeroWithCap returns a zero polynomial over the specified ring where the
// underlying map has given capacity.
func (r *QuotientRing) zeroWithCap(cap int) *Polynomial {
	return &Polynomial{
		baseRing: r,
		coefs:    make(map[[2]uint]ff.Element, cap),
	}
}

// Polynomial defines a new polynomial with the given coefficients.
func (r *QuotientRing) Polynomial(coefs map[[2]uint]ff.Element) *Polynomial {
	m := make(map[[2]uint]ff.Element, len(coefs))
	for d, e := range coefs {
		if e.IsNonzero() {
			m[d] = e
		}
	}
	out := &Polynomial{baseRing: r, coefs: m}
	out.reduce()
	return out
}

// PolynomialFromUnsigned defines a new polynomial with the given coefficients.
func (r *QuotientRing) PolynomialFromUnsigned(coefs map[[2]uint]uint) *Polynomial {
	m := make(map[[2]uint]ff.Element, len(coefs))
	for d, c := range coefs {
		e := r.baseField.ElementFromUnsigned(c)
		if e.IsNonzero() {
			m[d] = e
		}
	}
	out := &Polynomial{baseRing: r, coefs: m}
	out.reduce()
	return out
}

// PolynomialFromSigned defines a new polynomial with the given coefficients.
func (r *QuotientRing) PolynomialFromSigned(coefs map[[2]uint]int) *Polynomial {
	m := make(map[[2]uint]ff.Element, len(coefs))
	for d, c := range coefs {
		e := r.baseField.ElementFromSigned(c)
		if e.IsNonzero() {
			m[d] = e
		}
	}
	out := &Polynomial{baseRing: r, coefs: m}
	out.reduce()
	return out
}

// PolynomialFromString defines a polynomial by parsing s.
//
// The string s must use the variable names specified by r, but capitalization
// is ignored. Multiplication symbol '*' is allowed, but not necessary.
// Additionally, Singular-style exponents are allowed, meaning that "X2Y3" is
// interpreted as "X^2Y^3".
//
// If the string cannot be parsed, the function returns the zero polynomial and
// a Parsing-error.
func (r *QuotientRing) PolynomialFromString(s string) (*Polynomial, error) {
	m, err := polynomialStringToMap(s, &r.varNames, r)
	if err != nil {
		return r.Zero(), err
	}
	return r.Polynomial(m), nil
}

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

// Coef returns the coefficient of the monomial with degree specified by the
// input. The return value is a finite field element.
func (f *Polynomial) Coef(deg [2]uint) ff.Element {
	if c, ok := f.coefs[deg]; ok {
		return c.Copy()
	}
	return f.BaseField().Zero()
}

// coefPtr returns a pointer to the coefficient of the monomial with degree
// specified by the input. The return value is nil if the coefficient does not
// exist.
func (f *Polynomial) coefPtr(deg [2]uint) ff.Element {
	if c, ok := f.coefs[deg]; ok {
		return c
	}
	return nil
}

// SetCoef sets the coefficient of the monomial with degree deg in f to val by
// copying. See also SetCoefPtr.
func (f *Polynomial) SetCoef(deg [2]uint, val ff.Element) {
	if val.IsZero() {
		delete(f.coefs, deg)
	} else {
		f.coefs[deg] = val.Copy()
	}
}

// EmbedIn embeds f in the ring r if possible. The input reduce determines if f
// is reduced in the new ring.
//
// An InputIncompatible-error is returned if r and the polynomial ring of f are
// not compatible.
func (f *Polynomial) EmbedIn(r *QuotientRing, reduce bool) error {
	const op = "Embedding polynomial in ring"

	if f.baseRing.ring != r.ring {
		return errors.New(
			op, errors.InputIncompatible,
			"Cannot embed polynomial over %v in %v", f.baseRing, r,
		)
	}

	f.baseRing = r
	if reduce {
		f.reduce()
	}
	return nil
}

// SetCoefPtr sets the coefficient of the monomial with degree deg in f to ptr
// as a pointer. To set coefficient as a value, use SetCoef instead.
func (f *Polynomial) SetCoefPtr(deg [2]uint, ptr ff.Element) {
	if ptr.IsZero() {
		delete(f.coefs, deg)
	} else {
		f.coefs[deg] = ptr
	}
}

// IncrementCoef increments the coefficient of the monomial with degree deg in f
// by val.
func (f *Polynomial) IncrementCoef(deg [2]uint, val ff.Element) {
	if val.IsZero() {
		return
	}

	if c, ok := f.coefs[deg]; ok {
		c.Add(val)
		if c.IsZero() {
			delete(f.coefs, deg)
		}
	} else {
		f.coefs[deg] = val.Copy()
	}
}

// DecrementCoef decrements the coefficient of the monomial with degree deg in f
// by val.
func (f *Polynomial) DecrementCoef(deg [2]uint, val ff.Element) {
	if val.IsZero() {
		return
	}

	if c, ok := f.coefs[deg]; ok {
		c.Sub(val)
		if c.IsZero() {
			delete(f.coefs, deg)
		}
	} else {
		f.coefs[deg] = val.Neg()
	}
}

// Copy returns a new polynomial object over the same ring and with the same
// coefficients as f.
func (f *Polynomial) Copy() *Polynomial {
	h := f.baseRing.zeroWithCap(len(f.coefs))
	for deg, c := range f.coefs {
		h.coefs[deg] = c.Copy()
	}
	return h
}

// clean removes any zero coefficients from the underlying map of f.
// func (f *Polynomial) clean() {
// 	for d, c := range f.coefs {
// 		if c.IsZero() {
// 			delete(f.coefs, d)
// 		}
// 	}
// }

// Eval evaluates f at the given point.
func (f *Polynomial) Eval(point [2]ff.Element) ff.Element {
	out := f.BaseField().Zero()
	for deg, coef := range f.coefs {
		out = out.Plus(coef.Times(point[0].Pow(deg[0])).Times(point[1].Pow(deg[1])))
	}
	return out
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

// SortedDegrees returns a list containing the degrees is the support of f.
//
// The list is sorted according to the ring order with higher orders preceding
// lower orders in the list.
func (f *Polynomial) SortedDegrees() [][2]uint {
	degs := make([][2]uint, 0, len(f.coefs))
	for deg := range f.coefs {
		degs = append(degs, deg)
	}

	if len(degs) > 1 {
		sort.Slice(degs, func(i, j int) bool {
			return (f.baseRing.ord(degs[i], degs[j]) >= 0)
		})
	}
	return degs
}

// Ld returns the leading degree of f.
func (f *Polynomial) Ld() [2]uint {
	ld := [2]uint{0, 0}
	for deg := range f.coefs {
		if f.baseRing.ord(deg, ld) == 1 {
			ld = deg
		}
	}
	return ld
}

// Lc returns the leading coefficient of f.
func (f *Polynomial) Lc() ff.Element {
	return f.Coef(f.Ld())
}

// lcPtr returns a pointer to the leading coefficient of f.
func (f *Polynomial) lcPtr() ff.Element {
	return f.coefs[f.Ld()]
}

// Lt returns the leading term of f.
func (f *Polynomial) Lt() *Polynomial {
	h := f.baseRing.Zero()
	ld := f.Ld()
	h.coefs[ld] = f.Coef(ld)
	return h
}

// IsZero determines whether f is the zero polynomial.
func (f *Polynomial) IsZero() bool {
	if len(f.coefs) == 0 {
		return true
	}
	return false
}

// IsNonzero determines whether f contains some monomial with nonzero coefficient.
func (f *Polynomial) IsNonzero() bool {
	return !f.IsZero()
}

// IsMonomial returns a bool describing whether f consists of a single monomial.
func (f *Polynomial) IsMonomial() bool {
	if len(f.coefs) == 1 {
		return true
	}
	return false
}

// Reduces f in-place and sets its error state if needed.
func (f *Polynomial) reduce() {
	const op = "Reducing polynomial"

	if tmp := hasErr(op, f); tmp != nil {
		return
	}

	if f.baseRing.id == nil {
		return
	}

	err := f.baseRing.id.Reduce(f)
	if err != nil {
		f.err = err
	}
}

// String returns the string representation of f. Variables are 'X' and 'Y' by
// default. To change this, see the SetVarNames method.
func (f *Polynomial) String() string {
	if f.IsZero() {
		return "0"
	}

	degs := f.SortedDegrees()

	var b strings.Builder
	for i, d := range degs {
		if i > 0 {
			b.Write([]byte(" + "))
		}

		// Append the coefficient
		if tmp := f.coefPtr(d); !tmp.IsOne() || (d[0] == 0 && d[1] == 0) {
			if tmp.NTerms() > 1 {
				b.WriteByte('(')
				b.WriteString(tmp.String())
				b.WriteByte(')')
			} else {
				b.WriteString(tmp.String())
			}
		}

		// Append the first variable name and its degree
		if d[0] >= 1 {
			b.WriteString(f.baseRing.VarNames()[0])
		}
		if d[0] > 1 {
			b.WriteByte('^')
			b.WriteString(strconv.FormatUint(uint64(d[0]), 10))
		}

		// Append the second variable name and its degree
		if d[1] >= 1 {
			b.WriteString(f.baseRing.VarNames()[1])
		}
		if d[1] > 1 {
			b.WriteByte('^')
			b.WriteString(strconv.FormatUint(uint64(d[1]), 10))
		}
	}
	return b.String()
}

// checkErrAndCompatible is a wrapper for the two functions hasErr and
// checkCompatible. It is used in arithmetic functions to check that the inputs
// are 'good' to use.
func checkErrAndCompatible(op errors.Op, f *Polynomial, g ...*Polynomial) *Polynomial {
	if tmp := hasErr(op, f); tmp != nil {
		return tmp
	}

	if tmp := hasErr(op, g...); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g...); tmp != nil {
		return tmp
	}

	return nil
}

// hasErr is an internal method for checking if one of the given polynomials has
// a non-nil error field.
//
// It returns the first polynomial with non-nil error status after wrapping the
// error. The new error inherits the kind from the old.
func hasErr(op errors.Op, f ...*Polynomial) *Polynomial {
	for _, g := range f {
		if g.err != nil {
			g.err = errors.Wrap(
				op, errors.Inherit,
				g.err,
			)
			return g
		}
	}
	return nil
}

// checkCompatible is an internal method for checking if f and g are compatible;
// that is, if they are defined over the same ring.
//
// If not, the return value is an element with error status set to
// ArithmeticIncompat.
func checkCompatible(op errors.Op, f *Polynomial, g ...*Polynomial) *Polynomial {
	for _, h := range g {
		if f.baseRing != h.baseRing {
			out := f.baseRing.Zero()
			out.err = errors.New(
				op, errors.ArithmeticIncompat,
				"%v and %v defined over different rings", f, g,
			)
			return out
		}
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
