package univariate

import (
	"algobra/errors"
	"algobra/finitefield"
	"fmt"
)

type ring struct {
	baseField *finitefield.Field
}

type QuotientRing struct {
	*ring
	id *Ideal
}

// String returns the strings representation of r.
func (r *ring) String() string {
	return fmt.Sprintf("Univariate polynomial ring over %v", r.baseField)
}

// String returns the string representation of r.
func (r *QuotientRing) String() string {
	if r.id == nil {
		return r.ring.String()
	}
	return fmt.Sprintf(
		"Quotient ring of univariate polynomials over %v modulo %v",
		r.ring, r.id,
	)
}

// DefRing defines a new polynomial ring over the given field. It returns a new
// ring-object
func DefRing(field *finitefield.Field) *QuotientRing {
	return &QuotientRing{
		ring: &ring{
			baseField: field,
		},
		id: nil,
	}
}

// Zero returns a zero polynomial over the specified ring.
func (r *QuotientRing) Zero() *Polynomial {
	coefs := make([]*finitefield.Element, 1)
	coefs[0] = r.baseField.Zero()
	return &Polynomial{
		baseRing: r,
		coefs:    coefs,
	}
}

// Polynomial defines a new polynomial with the given coefficients
func (r *QuotientRing) Polynomial(coefs []*finitefield.Element) *Polynomial {
	out := r.Zero()
	for d, e := range coefs {
		if e.Nonzero() {
			out.SetCoef(d, e.Copy())
		}
	}
	out.reduce()
	return out
}

// PolynomialFromUnsigned defines a new polynomial with the given coefficients
func (r *QuotientRing) PolynomialFromUnsigned(coefs []uint) *Polynomial {
	out := r.Zero()
	for d, c := range coefs {
		e := r.baseField.ElementFromUnsigned(c)
		if e.Nonzero() {
			out.SetCoef(d, e)
		}
	}
	out.reduce()
	return out
}

// PolynomialFromSigned defines a new polynomial with the given coefficients
func (r *QuotientRing) PolynomialFromSigned(coefs []int) *Polynomial {
	out := r.Zero()
	for d, e := range coefs {
		if e != 0 {
			out.SetCoef(d, r.baseField.ElementFromSigned(e))
		}
	}
	out.reduce()
	return out
}

func (r *QuotientRing) Quotient(id *Ideal) (*QuotientRing, error) {
	const op = "Define quotient ring"
	if r.id != nil {
		return r, errors.New(
			op, errors.InputValue,
			"Given ring is already reduced modulo an ideal",
		)
	}
	if r.ring != id.ring {
		return r, errors.New(
			op, errors.InputIncompatible,
			"Input argument not ideal of ring '%v'", r,
		)
	}

	qr := &QuotientRing{
		ring: r.ring,
		id:   nil,
	}
	idConv := id.Copy()
	// Mark the generators as 'belonging' to the new ring
	idConv.generator.baseRing = qr
	qr.id = idConv
	return qr, nil
}

// QuoRem return the polynomial quotient and remainder under division by the
// given list of polynomials.
//
// Loosely based on [GG; Algorithm 2.5].
func (f *Polynomial) QuoRem(list ...*Polynomial) (q []*Polynomial, r *Polynomial) {
	r = f.baseRing.Zero()
	p := f.Copy()

	q = make([]*Polynomial, len(list), len(list))
	for i := range list {
		q[i] = f.baseRing.Zero()
	}
outer:
	for p.Nonzero() {
		for i, g := range list {
			if p.Ld() >= g.Ld() {
				tmp := f.baseRing.Zero()
				tmp.SetCoef(
					p.Ld()-g.Ld(),
					p.Lc().Mult(g.Lc().Inv()),
				)
				q[i] = q[i].Plus(tmp)
				p = p.Minus(tmp.multNoReduce(g))
				continue outer
			}
		}
		// No polynomials divide the leading term of f
		r = r.Plus(p.Lt())
		p = p.Minus(p.Lt())
	}
	return q, r
}
