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

	for r1.Nonzero() {
		_, rem := r0.QuoRem(r1)
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
	_, r := f.QuoRem(id.generator)
	*f = *r // For some reason using pointers alone is not enough
}
