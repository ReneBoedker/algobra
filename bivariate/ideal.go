package bivariate

import (
	"algobra/errors"
	"fmt"
	"strings"
)

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
		generators: make([]*Polynomial, len(generators)),
		isGroebner: 0,
		isMinimal:  0,
		isReduced:  0,
	}
	for i, g := range generators {
		if g.baseRing != r {
			return nil, errors.New(
				op, errors.InputIncompatible,
				"Generators defined over different rings",
			)
		}
		id.generators[i] = g.Copy()
	}
	return id, nil
}

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
func (id *Ideal) Reduce(f *Polynomial) {
	// TODO: Ought id to be a Gröbner basis here?
	_, r := f.QuoRem(id.generators...)
	*f = *r // For some reason using pointers alone is not enough
}

// QuoRem return the polynomial quotient and remainder under division by the
// given list of polynomials.
func (f *Polynomial) QuoRem(list ...*Polynomial) (q []*Polynomial, r *Polynomial) {
	return f.quoRemWithIgnore(-1, list...)
}

func (f *Polynomial) quoRemWithIgnore(ignoreIndex int, list ...*Polynomial) (q []*Polynomial, r *Polynomial) {
	r = f.baseRing.Zero()
	p := f.Copy()

	q = make([]*Polynomial, len(list), len(list))
	for i, _ := range list {
		q[i] = f.baseRing.Zero()
	}
outer:
	for p.Nonzero() {
		for i, g := range list {
			if i == ignoreIndex {
				continue
			}
			if mquo, ok, err := p.Lt().monomialDivideBy(g.Lt()); err != nil {
				// Should not occur
				panic(err)
			} else if ok {
				// Lt(g) divides p.Lt()
				q[i] = q[i].Plus(mquo)
				p = p.Minus(g.multNoReduce(mquo))
				continue outer
			}
		}
		// No generators divide
		tmp := p.Lt()
		r = r.Plus(tmp)
		p = p.Minus(tmp)
	}
	return q, r
}