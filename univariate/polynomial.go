package univariate

import (
	"algobra/errors"
	"algobra/finitefield/ff"
	"fmt"
	"strings"
)

// Polynomial denotes a bivariate polynomial.
type Polynomial struct {
	baseRing *QuotientRing
	coefs    []ff.Element
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

// Coef returns the coefficient of the monomial with degree specified by the
// input. The return value is a finite field element.
func (f *Polynomial) Coef(deg int) ff.Element {
	if deg < len(f.coefs) {
		return f.coefs[deg]
	}
	return f.BaseField().Zero()
}

// SetCoef sets the coefficient of the monomial with degree deg in f to val.
func (f *Polynomial) SetCoef(deg int, val ff.Element) {
	if deg <= f.Ld() {
		f.coefs[deg] = val.Copy()
		if val.IsZero() {
			f.reslice()
		}
		return
	}
	// Otherwise, grow the slice to needed length
	grow := make([]ff.Element, deg-f.Ld())
	for i := range grow {
		grow[i] = f.BaseField().Zero()
	}
	f.coefs = append(f.coefs, grow...)
	f.coefs[deg] = val.Copy()
}

// IncrementCoef increments the coefficient of the monomial with degree deg in f
// by val.
func (f *Polynomial) IncrementCoef(deg int, val ff.Element) {
	if val.IsZero() {
		return
	}
	if deg <= f.Ld() {
		f.coefs[deg].Add(val)
		f.reslice()
		return
	}
	// Otherwise, grow the slice to needed length
	grow := make([]ff.Element, deg-f.Ld())
	for i := range grow {
		grow[i] = f.BaseField().Zero()
	}
	f.coefs = append(f.coefs, grow...)
	f.coefs[deg].Add(val)
}

// reslice ensures that the coefficients of f do not contain leading zeros
func (f *Polynomial) reslice() {
	for i := len(f.coefs) - 1; i >= 0; i-- {
		if f.Coef(i).IsNonzero() {
			f.coefs = f.coefs[:i+1]
			return
		}
	}
	// No non-zero entries were found
	f.coefs = f.coefs[:1]
}

// Copy returns a new polynomial object over the same ring and with the same
// coefficients as f.
func (f *Polynomial) Copy() *Polynomial {
	h := f.baseRing.Zero()
	h.coefs = make([]ff.Element, len(f.coefs))
	for deg, c := range f.coefs {
		h.coefs[deg] = c.Copy()
	}
	return h
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

// Eval evaluates f at the given point.
func (f *Polynomial) Eval(point ff.Element) ff.Element {
	out := f.BaseField().Zero()
	for deg, c := range f.coefs {
		out = out.Plus(c.Times(point.Pow(uint(deg))))
	}
	return out
}

// Normalize creates a new polynomial obtained by normalizing f. That is,
// f.Normalize() multiplied by f.Lc() is f.
//
// If f is the zero polynomial, a copy of f is returned.
func (f *Polynomial) Normalize() *Polynomial {
	if f.IsZero() {
		return f.Copy()
	}
	return f.Scale(f.Lc().Inv())
}

// Scale scales all coefficients of f by the field element c and returns the
// result as a new polynomial.
func (f *Polynomial) Scale(c ff.Element) *Polynomial {
	g := f.Copy()
	for deg := range g.coefs {
		g.SetCoef(deg, g.Coef(deg).Times(c))
	}
	return g
}

// Degrees returns a slice containing the degrees in the support of f.
//
// The list is sorted with higher degrees preceding lower ones in the list.
func (f *Polynomial) Degrees() []int {
	degs := make([]int, 0, len(f.coefs))
	for deg := len(f.coefs) - 1; deg >= 0; deg-- {
		if f.Coef(deg).IsZero() {
			continue
		}
		degs = append(degs, deg)
	}
	return degs
}

// Coefs returns a slice containing the coefficients of f.
//
// The i'th element of the resulting slice is the coefficient of degree i.
func (f *Polynomial) Coefs() []ff.Element {
	coefs := make([]ff.Element, len(f.coefs), len(f.coefs))
	for i, c := range f.coefs {
		coefs[i] = c.Copy()
	}
	return coefs
}

// Ld returns the leading degree of f.
func (f *Polynomial) Ld() int {
	return len(f.coefs) - 1
}

// Lc returns the leading coefficient of f.
func (f *Polynomial) Lc() ff.Element {
	return f.Coef(f.Ld())
}

// Lt returns the leading term of f.
func (f *Polynomial) Lt() *Polynomial {
	h := f.baseRing.Zero()
	ld := f.Ld()
	h.SetCoef(ld, f.Coef(ld))
	return h
}

// IsZero determines whether f is the zero polynomial.
func (f *Polynomial) IsZero() bool {
	if len(f.coefs) == 1 && f.Coef(0).IsZero() {
		return true
	}
	return false
}

// IsNonzero determines whether f contains some monomial with nonzero coefficient.
func (f *Polynomial) IsNonzero() bool {
	return !f.IsZero()
}

// IsOne determines whether f is the constant 1.
func (f *Polynomial) IsOne() bool {
	if len(f.coefs) == 1 && f.Coef(0).IsOne() {
		return true
	}
	return false
}

// IsMonomial returns a bool describing whether f consists of a single monomial.
func (f *Polynomial) IsMonomial() bool {
	if len(f.Degrees()) == 1 {
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

// checkErrAndCompatible is a wrapper for the two functions hasErr and
// checkCompatible. It is used in arithmetic functions to check that the inputs
// are 'good' to use.
func checkErrAndCompatible(op errors.Op, f *Polynomial, g ...*Polynomial) *Polynomial {
	if tmp := hasErr(op, f, g...); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g...); tmp != nil {
		return tmp
	}

	return nil
}

// hasErr is an internal method for checking if f or g has a non-nil error
// field.
//
// It returns the first polynomial with non-nil error status after wrapping the
// error. The new error inherits the kind from the old.
func hasErr(op errors.Op, f *Polynomial, g ...*Polynomial) *Polynomial {
	if f.err != nil {
		f.err = errors.Wrap(
			op, errors.Inherit,
			f.err,
		)
		return f
	}
	for _, h := range g {
		if h.err != nil {
			h.err = errors.Wrap(
				op, errors.Inherit,
				h.err,
			)
			return h
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
				"%v and %v defined over different rings", f, h,
			)
			return out
		}
	}
	return nil
}

// String returns the string representation of f. The variable is named
// according to the ring used.
func (f *Polynomial) String() string {
	if f.IsZero() {
		return "0"
	}

	var b strings.Builder
	for i, d := range f.Degrees() {
		if i > 0 {
			fmt.Fprint(&b, " + ")
		}
		if tmp := f.Coef(d); !tmp.IsOne() || d == 0 {
			fmt.Fprintf(&b, "%v", tmp)
		}
		if d == 1 {
			fmt.Fprint(&b, f.baseRing.varName)
		}
		if d > 1 {
			fmt.Fprintf(&b, "%s^%d", f.baseRing.varName, d)
		}
	}
	return b.String()
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
