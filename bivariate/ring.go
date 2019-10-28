package bivariate

import (
	"fmt"
	"strings"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

type ring struct {
	baseField ff.Field
	varNames  [2]string
	ord       Order
}

// QuotientRing denotes a polynomial quotient ring. The quotient may be trivial,
// in which case, the object acts as a ring.
type QuotientRing struct {
	*ring
	id *Ideal
}

// String returns the string representation of r.
func (r *ring) String() string {
	return fmt.Sprintf("Bivariate polynomial ring over %v", r.baseField)
}

// String returns the string representation of r.
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
// the order function ord. It returns a new ring-object.
func DefRing(field ff.Field, ord Order) *QuotientRing {
	return &QuotientRing{
		ring: &ring{
			baseField: field,
			varNames:  [2]string{"X", "Y"},
			ord:       ord,
		},
		id: nil,
	}
}

// SetVarNames sets the variable names to be used in the given quotient ring.
//
// Leading and trailing whitespace characters are removed before setting the
// variable name. If the one of the strings consists solely of whitespace
// characters, an InputValue-error is returned.
//
// If the two variable names are identical when ignoring leading and trailing
// whitespace and capitalization, an InputValue-error is returned.
func (r *QuotientRing) SetVarNames(varNames [2]string) error {
	// TODO: Do more strings have to be disallowed (eg. +, -)?
	const op = "Setting variable name"

	for i, v := range varNames {
		varName := strings.TrimSpace(v)
		if len(varName) == 0 {
			return errors.New(
				op, errors.InputValue,
				"Cannot use whitespace characters as variable name",
			)
		}
		r.varNames[i] = varName
	}

	if strings.ToLower(r.varNames[0]) == strings.ToLower(r.varNames[1]) {
		return errors.New(
			op, errors.InputValue,
			"Variable names %q and %q are considered identical. Parsing from "+
				"strings is unlikely to work.",
			r.varNames[0], r.varNames[1],
		)
	}

	return nil
}

// VarNames returns the strings used to represent the variables of r.
func (r *QuotientRing) VarNames() [2]string {
	return r.varNames
}

// Quotient defines the quotient of the given ring modulo the input ideal.
//
// The return type is a new QuotientRing-object
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
	// This is tested when the ideal is created
	// for _, f := range id.generators {
	// 	if f.baseRing != r {
	// 		return nil, errors.New(
	// 			op, errors.InputIncompatible,
	// 			"Ideal member %v not in ring", f,
	// 		)
	// 	}
	// }

	qr := &QuotientRing{
		ring: r.ring,
		id:   id,
	}

	// Mark the generators as belonging to the new ring, but do not reduce
	for _, g := range qr.id.generators {
		g.EmbedIn(qr, false)
	}

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
