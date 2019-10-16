package bivariate

import (
	"algobra/errors"
	"algobra/finitefield/ff"
	"fmt"
)

type ring struct {
	baseField ff.Field
	ord       Order
}

// QuotientRing denotes a polynomial quotient ring. The quotient may be trivial,
// in which case, the object acts as a ring.
type QuotientRing struct {
	*ring
	id *Ideal
}

// String returns the string representation of r
func (r *ring) String() string {
	return fmt.Sprintf("Bivariate polynomial ring over %v", r.baseField)
}

// String returns the string representation of r
func (r *QuotientRing) String() string {
	if r.id == nil {
		return r.ring.String()
	}
	return fmt.Sprintf(
		"Quotient ring of bivariate polynomials over %v modulo %v",
		r.ring, r.id,
	)
}

// DefRing defines a new polynomial ring over the given field, using
// the order function ord. It returns a new ring-object
func DefRing(field ff.Field, ord Order) *QuotientRing {
	return &QuotientRing{
		ring: &ring{
			baseField: field,
			ord:       ord,
		},
		id: nil,
	}
}

// Zero returns a zero polynomial over the specified ring.
func (r *QuotientRing) Zero() *Polynomial {
	return &Polynomial{
		baseRing: r,
		coefs:    map[[2]uint]ff.Element{},
	}
}

// zeroWithCap returns a zero polynomial over the specified ring where the
// underlying map has given capacity.
func (r *QuotientRing) zeroWithCap(cap int) *Polynomial {
	return &Polynomial{
		baseRing: r,
		coefs:    map[[2]uint]ff.Element{},
	}
}

// Polynomial defines a new polynomial with the given coefficients
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

// PolynomialFromUnsigned defines a new polynomial with the given coefficients
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

// PolynomialFromSigned defines a new polynomial with the given coefficients
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
// The string s must use 'X' and 'Y' as variable names, but lowercase letters are
// accepted. Multiplication symbol '*' is allowed, but not necessary.
// Additionally, Singular-style exponents are allowed, meaning that "X2Y3" is
// interpreted as "X^2Y^3".
//
// If the string cannot be parsed, the function returns the zero polynomial and
// a Parsing-error.
func (r *QuotientRing) PolynomialFromString(s string) (*Polynomial, error) {
	m, err := polynomialStringToMap(s, r.baseField)
	if err != nil {
		return r.Zero(), err
	}
	return r.Polynomial(m), nil
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
	if id.isGroebner != 1 {
		id = id.GroebnerBasis()
		_ = id.ReduceBasis()
	}
	if r.ring != id.ring {
		return r, errors.New(
			op, errors.InputIncompatible,
			"Input argument not ideal of ring '%v'", r,
		)
	}
	for _, f := range id.generators {
		if f.baseRing != r {
			return nil, errors.New(
				op, errors.InputIncompatible,
				"Ideal member %v not in ring", f,
			)
		}
	}

	qr := &QuotientRing{
		ring: r.ring,
		id:   nil,
	}
	idConv := id.Copy()
	// Mark the generators as 'belonging' to the new ring
	for i := range idConv.generators {
		idConv.generators[i].baseRing = qr
	}
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
