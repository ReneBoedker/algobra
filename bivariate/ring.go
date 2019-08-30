package bivariate

import (
	"algobra/primefield"
	"fmt"
)

type ring struct {
	baseField *primefield.Field
	ord       order
}

type QuotientRing struct {
	*ring
	id *Ideal
}

// DefRing defines a new polynomial ring with the given characteristic, using
// the order function ord. It returns a new ring-object
func DefRing(field *primefield.Field, ord order) *QuotientRing {
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
		degrees:  map[[2]uint]*primefield.Element{},
	}
}

// New defines a new polynomial with the given coefficients
func (r *QuotientRing) New(coefs map[[2]uint]uint) *Polynomial {
	m := make(map[[2]uint]*primefield.Element)
	for d, c := range coefs {
		e := r.baseField.Element(c)
		if e.Nonzero() {
			m[d] = e
		}
	}
	out := &Polynomial{baseRing: r, degrees: m}
	out.reduce()
	return out
}

// New defines a new polynomial with the given coefficients
func (r *QuotientRing) NewFromString(s string) (*Polynomial, error) {
	m, err := polynomialStringToMap(s)
	if err != nil {
		return r.Zero(), err
	}
	return r.New(m), nil
}

// Quotient defines the quotient of the given ring modulo the input ideal.
// The return type is a new ring-object
func (r *QuotientRing) Quotient(id *Ideal) (*QuotientRing, error) {
	if r.id != nil {
		return r, fmt.Errorf("Quotient: Given ring is already reduced modulo an ideal")
	}
	if id.isGroebner != 1 {
		id = id.GroebnerBasis()
		_ = id.ReduceBasis()
	}
	if r.ring != id.ring {
		return r, fmt.Errorf("ring.Quotient: Input argument not ideal of r")
	}
	for _, f := range id.generators {
		if f.baseRing != r {
			return nil, fmt.Errorf(
				"ring.Quotient: Ideal member %v not in ring",
				f,
			)
		}
	}
	return &QuotientRing{
		ring: r.ring,
		id:   id,
	}, nil
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
