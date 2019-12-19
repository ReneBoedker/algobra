package bivariate

import (
	"fmt"
	"strings"

	"github.com/ReneBoedker/algobra/errors"
)

// Ideal is the implementation of a polynomial ideal.
type Ideal struct {
	*ring
	generators []*Polynomial
	isGroebner int8 // 0=undecided, 1=true, -1=false
	isMinimal  int8 // 0=undecided, 1=true, -1=false
	isReduced  int8 // 0=undecided, 1=true, -1=false
}

// ShortString returns a short string description of id. More precisely, it
// returns the string representation of the generators.
func (id *Ideal) ShortString() string {
	var sb strings.Builder
	fmt.Fprint(&sb, "<")

	for i, g := range id.generators {
		if i > 0 {
			fmt.Fprint(&sb, ", ")
		}
		fmt.Fprint(&sb, g)
	}

	fmt.Fprint(&sb, ">")
	return sb.String()
}

// String returns the string representation of id. See also ShortString.
func (id *Ideal) String() string {
	return fmt.Sprintf("Ideal %s of %v", id.ShortString(), id.ring)
}

// NewIdeal returns a new polynomial ideal over the given ring. If the
// generators are not defined over the given ring, the function returns an
// InputIncompatible-error.
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

// Generators returns a copy of the generators of id.
func (id *Ideal) Generators() []*Polynomial {
	gens := make([]*Polynomial, len(id.generators), len(id.generators))
	for i, g := range id.generators {
		gens[i] = g.Copy()
	}
	return gens
}

// Reduce sets f to f modulo id.
//
// Note that when the generators of id do not form a Gröbner basis, such a basis
// will be computed. This alters the representation of id.
func (id *Ideal) Reduce(f *Polynomial) error {
	const op = "Reducing polynomial"

	if !id.IsGroebner() {
		id = id.GroebnerBasis()
	}

	r, err := f.Rem(id.generators...)
	if err != nil {
		return errors.Wrap(op, errors.Inherit, err)
	}

	*f = *r
	return nil
}
