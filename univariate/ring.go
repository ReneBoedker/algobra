package univariate

import (
	"algobra/errors"
	"algobra/finitefield"
	"fmt"
)

type ring struct {
	baseField *finitefield.Field
}

// QuotientRing denotes a polynomial quotient ring. The quotient may be trivial,
// in which case, the object acts as a ring.
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

// PolynomialFromString defines a polynomial by parsing s.
//
// The string s must use 'X' of 'x' as variable names. Multiplication symbol '*'
// is allowed, but not necessary. Additionally, Singular-style exponents are
// allowed, meaning that "X2" is interpreted as "X^2".
//
// If the string cannot be parsed, the function returns the zero polynomial and
// a Parsing-error.
func (r *QuotientRing) PolynomialFromString(s string) (*Polynomial, error) {
	const op = "Defining polynomial from string"

	m, err := polynomialStringToSignedMap(s)
	f := r.Zero()
	if err != nil {
		return f, errors.Wrap(op, errors.Inherit, err)
	}
	for deg, coef := range m {
		if deg < 0 {
			return r.Zero(), errors.New(
				op, errors.InputValue,
				"Input %q contains negative degree %d", s, deg,
			)
		}
		f.SetCoef(deg, f.Coef(deg).Plus(r.baseField.ElementFromSigned(coef)))
	}
	return f, nil
}

// Quotient defines the quotient of the given ring modulo the input ideal.
//
// The return type is a new ring-object
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