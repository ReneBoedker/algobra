package univariate

import (
	"fmt"
	"strings"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

type ring struct {
	baseField ff.Field
	varName   string
}

// QuotientRing denotes a polynomial quotient ring. The quotient may be trivial,
// in which case, the object acts as a ring.
type QuotientRing struct {
	*ring
	id *Ideal
}

// String returns the strings representation of r.
func (r *ring) String() string {
	return fmt.Sprintf(
		"Univariate polynomial ring in %s over %v", r.varName, r.baseField)
}

// String returns the string representation of r.
func (r *QuotientRing) String() string {
	if r.id == nil {
		return r.ring.String()
	}
	return fmt.Sprintf(
		"Quotient ring of univariate polynomials in %s over %v modulo %v",
		r.varName, r.baseField, r.id,
	)
}

// DefRing defines a new polynomial ring over the given field. It returns a new
// ring-object
func DefRing(field ff.Field) *QuotientRing {
	return &QuotientRing{
		ring: &ring{
			baseField: field,
			varName:   "X",
		},
		id: nil,
	}
}

// SetVarName sets the variable name to be used in the given quotient ring.
//
// Leading and trailing whitespace characters are removed before setting the
// variable name. If the string consists solely of whitespace characters, an
// InputValue-error is returned.
func (r *QuotientRing) SetVarName(varName string) error {
	// TODO: Do more strings have to be disallowed (eg. +, -)?
	const op = "Setting variable name"

	varName = strings.TrimSpace(varName)
	if len(varName) == 0 {
		return errors.New(
			op, errors.InputValue,
			"Cannot use whitespace characters as variable name",
		)
	}
	r.varName = varName
	return nil
}

// VarName returns the string used to represent the variable of r.
func (r *QuotientRing) VarName() string {
	return r.varName
}

// zeroWithCap returns a zero polynomial over the specified ring, where the
// underlying representation has given capacity. If cap is negative, the default
// map capacity is used.
func (r *QuotientRing) zeroWithCap(cap int) *Polynomial {
	var coefs []ff.Element
	if cap < 0 {
		coefs = make([]ff.Element, 1)
	} else {
		coefs = make([]ff.Element, 1, cap)
	}
	coefs[0] = r.baseField.Zero()
	return &Polynomial{
		baseRing: r,
		coefs:    coefs,
	}
}

// Zero returns a zero polynomial over the specified ring.
func (r *QuotientRing) Zero() *Polynomial {
	return r.zeroWithCap(-1)
}

// One returns the degree zero polynomial over the specified ring with
// coefficient 1.
func (r *QuotientRing) One() *Polynomial {
	coefs := make([]ff.Element, 1)
	coefs[0] = r.baseField.One()
	return &Polynomial{
		baseRing: r,
		coefs:    coefs,
	}
}

// Polynomial defines a new polynomial with the given coefficients
func (r *QuotientRing) Polynomial(coefs []ff.Element) *Polynomial {
	out := r.Zero()
	for d, e := range coefs {
		if e.IsNonzero() {
			out.SetCoef(d, e.Copy())
		}
	}
	out.reduce()
	return out
}

// PolynomialFromUnsigned defines a new polynomial with the given coefficients
func (r *QuotientRing) PolynomialFromUnsigned(coefs []uint) *Polynomial {
	out := r.Zero()
	for d, c := range coefs {
		e := r.baseField.ElementFromUnsigned(c)
		if e.IsNonzero() {
			out.SetCoef(d, e)
		}
	}
	out.reduce()
	return out
}

// PolynomialFromSigned defines a new polynomial with the given coefficients
func (r *QuotientRing) PolynomialFromSigned(coefs []int) *Polynomial {
	out := r.Zero()
	for d, e := range coefs {
		if e != 0 {
			out.SetCoef(d, r.baseField.ElementFromSigned(e))
		}
	}
	out.reduce()
	return out
}

// PolynomialFromString defines a polynomial by parsing s.
//
// The string s must use 'X' of 'x' as variable names. Multiplication symbol '*'
// is allowed, but not necessary. Additionally, Singular-style exponents are
// allowed, meaning that "X2" is interpreted as "X^2".
//
// If the string cannot be parsed, the function returns the zero polynomial and
// a Parsing-error.
func (r *QuotientRing) PolynomialFromString(s string) (*Polynomial, error) {
	const op = "Defining polynomial from string"

	m, err := polynomialStringToMap(s, &r.varName, r.baseField)
	f := r.Zero()
	if err != nil {
		return f, errors.Wrap(op, errors.Inherit, err)
	}
	for deg, coef := range m {
		if deg < 0 {
			return r.Zero(), errors.New(
				op, errors.InputValue,
				"Input %q contains negative degree %d", s, deg,
			)
		}
		f.SetCoef(deg, f.Coef(deg).Plus(coef))
	}
	return f, nil
}

// Quotient defines the quotient of the given ring modulo the input ideal.
//
// The return type is a new ring-object
func (r *QuotientRing) Quotient(id *Ideal) (*QuotientRing, error) {
	const op = "Define quotient ring"
	if r.id != nil {
		return r, errors.New(
			op, errors.InputValue,
			"Given ring is already reduced modulo an ideal",
		)
	}
	if r.ring != id.ring {
		return r, errors.New(
			op, errors.InputIncompatible,
			"Input argument not ideal of ring '%v'", r,
		)
	}

	qr := &QuotientRing{
		ring: r.ring,
		id:   nil,
	}
	idConv := id.Copy()
	// Mark the generators as 'belonging' to the new ring
	idConv.generator.baseRing = qr
	qr.id = idConv
	return qr, nil
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
