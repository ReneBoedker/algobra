package univariate

import (
	"algobra/errors"
	"fmt"
)

// Ideal is the implementation of a polynomial ideal.
type Ideal struct {
	*ring
	generator *Polynomial
}

// String returns the string representation of id.
func (id *Ideal) String() string {
	return fmt.Sprintf("Ideal <%v> over %v", id.generator, id.ring)
}

// NewIdeal returns a new polynomial ideal over the given ring. If the
// generators are not defined over the given ring, the function panics.
// Internally, this function computes a reduced Gröbner basis.
func (r *QuotientRing) NewIdeal(generators ...*Polynomial) (*Ideal, error) {
	const op = "Defining ideal"

	gcd, err := Gcd(generators[0], generators[1:]...)

	if err != nil {
		return nil, errors.Wrap(op, errors.Inherit, err)
	}

	id := &Ideal{
		ring:      r.ring,
		generator: gcd,
	}
	return id, nil
}

// Gcd returns the greatest common divisor of the given polynomials.
//
// An InputIncompatible-error is returned if the polynomials are not defined
// over the same ring.
func Gcd(f *Polynomial, g ...*Polynomial) (*Polynomial, error) {
	const op = "Computing polynomial GCD"

	for _, h := range g {
		if h.baseRing != f.baseRing {
			return nil, errors.New(
				op, errors.InputIncompatible,
				"Generators defined over different rings",
			)
		}
	}

	switch len(g) {
	case 0:
		return f, nil
	case 1:
		// Base case. Do computations below
	default:
		rec, err := Gcd(f, g[0])
		if err != nil {
			return rec, err
		}
		return Gcd(rec, g[1:]...)
	}

	r0 := f.Copy()
	r1 := g[0].Copy()

	for r1.IsNonzero() {
		_, rem, err := r0.QuoRem(r1)
		if err != nil {
			return nil, errors.Wrap(op, errors.Inherit, err)
		}
		r0, r1 = r1, rem
	}

	return r0, nil
}

// Copy creates a copy of id.
func (id *Ideal) Copy() *Ideal {
	return &Ideal{
		ring:      id.ring,
		generator: id.generator.Copy(),
	}
}

// Reduce sets f to f modulo id
func (id *Ideal) Reduce(f *Polynomial) {
	// TODO: Ought id to be a Gröbner basis here?
	_, r, err := f.QuoRem(id.generator)
	if err != nil {
		panic(err)
	}
	*f = *r // For some reason using pointers alone is not enough
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
