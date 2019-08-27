package bivariate

import (
	"algobra/primefield"
	"fmt"
)

type Ring struct {
	baseField  *primefield.Field
	ord        order
	reduceFunc func(*Polynomial)
}

// DefRing defines a new polynomial ring with the given characteristic, using
// the order function ord. It returns a new ring-object
func DefRing(field *primefield.Field, ord order) *Ring {
	return &Ring{
		baseField:  field,
		ord:        ord,
		reduceFunc: func(f *Polynomial) { return },
	}
}

// Zero returns a zero polynomial over the specified ring.
func (r *Ring) Zero() *Polynomial {
	return &Polynomial{
		baseRing: r,
		degrees:  map[[2]uint]*primefield.Element{},
	}
}

// New defines a new polynomial with the given coefficients
func (r *Ring) New(coefs map[[2]uint]uint) *Polynomial {
	m := make(map[[2]uint]*primefield.Element)
	for d, c := range coefs {
		e := r.baseField.Element(c)
		if e.Nonzero() {
			m[d] = e
		}
	}
	out := &Polynomial{baseRing: r, degrees: m}
	r.reduceFunc(out)
	return out
}

// func reduce(f *Polynomial) *Polynomial {
// 	return f
// }

// Quotient defines the quotient of the given ring modulo the input ideal.
// The return type is a new ring-object
func (r *Ring) Quotient(id ideal) (*Ring, error) {
	if r != id[0].baseRing {
		return r, fmt.Errorf("ring.Quotient: Input argument not ideal of r")
	}
	for _, f := range id {
		if f.baseRing != r {
			return nil, fmt.Errorf(
				"ring.Quotient: Ideal member %v not in ring",
				f,
			)
		}
	}
	return &Ring{
		baseField:  r.baseField,
		ord:        r.ord,
		reduceFunc: id.reduce,
	}, nil
}
