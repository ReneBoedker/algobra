package bivariate

import (
	"fmt"
	"strings"

	"algobra/errors"
)

// Ideal is the implementation of a polynomial ideal.
type Ideal struct {
	*ring
	generators []*Polynomial
	isGroebner int8 // 0=undecided, 1=true, -1=false
	isMinimal  int8 // 0=undecided, 1=true, -1=false
	isReduced  int8 // 0=undecided, 1=true, -1=false
}

// String returns the string representation of id.
func (id *Ideal) String() string {
	var sb strings.Builder
	for _, g := range id.generators {
		fmt.Fprint(&sb, g)
	}
	return fmt.Sprintf("Ideal <%s> over %v", sb.String(), id.ring)
}

// NewIdeal returns a new polynomial ideal over the given ring. If the
// generators are not defined over the given ring, the function panics.
// Internally, this function computes a reduced Gröbner basis.
func (r *QuotientRing) NewIdeal(generators ...*Polynomial) (*Ideal, error) {
	const op = "Defining ideal"
	id := &Ideal{
		ring:       r.ring,
		generators: make([]*Polynomial, 0, len(generators)),
		isGroebner: 0,
		isMinimal:  0,
		isReduced:  0,
	}

	for _, g := range generators {
		if g.baseRing != r {
			return nil, errors.New(
				op, errors.InputIncompatible,
				"Generators defined over different rings",
			)
		}
		if g.IsZero() {
			// Skip zero polynomials
			continue
		}
		id.generators = append(id.generators, g.Copy())
	}

	if len(id.generators) == 0 {
		return nil, errors.New(
			op, errors.InputValue,
			"Generators %v define an empty ideal", generators,
		)
	}

	return id, nil
}

// Copy creates a copy of id.
func (id *Ideal) Copy() *Ideal {
	generators := make([]*Polynomial, len(id.generators), len(id.generators))
	for i, g := range id.generators {
		generators[i] = g.Copy()
	}

	return &Ideal{
		ring:       id.ring,
		generators: generators,
		isGroebner: id.isGroebner,
		isMinimal:  id.isMinimal,
		isReduced:  id.isReduced,
	}
}

// Reduce sets f to f modulo id
func (id *Ideal) Reduce(f *Polynomial) error {
	// TODO: Ought id to be a Gröbner basis here?
	const op = "Reducing polynomial"

	_, r, err := f.QuoRem(id.generators...)
	if err != nil {
		return errors.Wrap(op, errors.Inherit, err)
	}

	*f = *r
	return nil
}
