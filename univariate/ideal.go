package univariate

import (
	"fmt"

	"github.com/ReneBoedker/algobra/errors"
)

// Ideal is the implementation of a polynomial ideal.
type Ideal struct {
	*ring
	generator *Polynomial
}

// Generator returns a copy of the generator of id.
func (id *Ideal) Generator() *Polynomial {
	return id.generator.Copy()
}

// ShortString returns a short string description of id. More precisely, it
// returns the string representation of the generators.
func (id *Ideal) ShortString() string {
	return "<" + id.generator.String() + ">"
}

// String returns the string representation of id. See also ShortString.
func (id *Ideal) String() string {
	return fmt.Sprintf("Ideal %s of %v", id.ShortString(), id.ring)
}

// NewIdeal returns a new polynomial ideal over the given ring. If the
// generators are not defined over the given ring, the function returns an
// InputIncompatible-error.
//
// Internally, this computes the greatest common divisor of the generators to
// find a single generator.
func (r *QuotientRing) NewIdeal(generators ...*Polynomial) (*Ideal, error) {
	const op = "Defining ideal"

	gcd, err := Gcd(generators[0], generators[1:]...)

	if err != nil {
		return nil, errors.Wrap(op, errors.Inherit, err)
	}

	id := &Ideal{
		ring:      r.ring,
		generator: gcd.Normalize(),
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

// Reduce sets f to f modulo id.
//
// Note that when a Gröbner basis has not been computed for id, the reduction is
// not necessarily unique.
func (id *Ideal) Reduce(f *Polynomial) error {
	const op = "Reducing polynomial"

	if tmp := checkErrAndCompatible(op, f, id.generator); tmp != nil {
		return tmp.Err()
	}

	for d := f.Ld(); d >= id.generator.Ld(); d = f.Ld() {
		f.subWithShiftAndScale(
			id.generator,
			d-id.generator.Ld(),
			f.coefs[f.Ld()], // Note that id.generator is normalized
		)
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
