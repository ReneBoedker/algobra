package univariate

import (
	"algobra/errors"
	"algobra/finitefield"
	"fmt"
	"strings"
)

// Polynomial denotes a bivariate polynomial.
type Polynomial struct {
	baseRing *QuotientRing
	coefs    []*finitefield.Element
	err      error
}

// BaseField returns the field over which the coefficients of f are defined.
func (f *Polynomial) BaseField() *finitefield.Field {
	return f.baseRing.baseField
}

// Err returns the error status of f.
func (f *Polynomial) Err() error {
	return f.err
}

// Coef returns the coefficient of the monomial with degree specified by the
// input. The return value is a finite field element.
func (f *Polynomial) Coef(deg uint) *finitefield.Element {
	if deg < uint(len(f.coefs)) {
		return f.coefs[deg]
	}
	return f.BaseField().Zero()
}

// SetCoef sets the coefficient of the monomial with degree deg in f to val.
func (f *Polynomial) SetCoef(deg uint, val *finitefield.Element) {
	if deg < uint(len(f.coefs)) {
		f.coefs[deg] = val
		return
	}
	// Otherwise, grow the slice to needed length
	grow := make([]*finitefield.Element, deg-f.Ld())
	for i := range grow {
		grow[i] = f.BaseField().Zero()
	}
	f.coefs = append(f.coefs, grow...)
	f.coefs[deg] = val
}

// Copy returns a new polynomial object over the same ring and with the same
// coefficients as f.
func (f *Polynomial) Copy() *Polynomial {
	h := f.baseRing.Zero()
	h.coefs = make([]*finitefield.Element, len(f.coefs))
	for deg, c := range f.coefs {
		h.coefs[deg] = c
	}
	return h
}

// Eval evaluates f at the given point.
func (f *Polynomial) Eval(point *finitefield.Element) *finitefield.Element {
	out := f.BaseField().Zero()
	for deg, c := range f.coefs {
		out = out.Plus(c.Mult(point.Pow(uint(deg))))
	}
	return out
}

// Plus returns the sum of the two polynomials f and g.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Plus(g *Polynomial) *Polynomial {
	const op = "Adding polynomials"

	if tmp := hasErr(op, f, g); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g); tmp != nil {
		return tmp
	}

	h := f.Copy()
	for d, c := range g.coefs {
		deg := uint(d)
		h.SetCoef(deg, h.Coef(deg).Plus(c))
	}
	return h
}

// Neg returns the polynomial obtained by scaling f by -1 (modulo the
// characteristic).
func (f *Polynomial) Neg() *Polynomial {
	g := f.baseRing.Zero()
	for deg, c := range f.coefs {
		g.SetCoef(uint(deg), c.Neg())
	}
	return g
}

// Equal determines whether two polynomials are equal. That is, whether they are
// defined over the same ring, and have the same coefficients.
func (f *Polynomial) Equal(g *Polynomial) bool {
	if f.baseRing != g.baseRing {
		return false
	}
	if len(f.coefs) != len(g.coefs) { // TODO: Relies on bookkeeping
		return false
	}
	for d, cf := range f.coefs {
		if !g.Coef(uint(d)).Equal(cf) {
			return false
		}
	}
	return true
}

// Minus returns polynomial difference f-g.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Minus(g *Polynomial) *Polynomial {
	const op = "Subtracting polynomials"

	if tmp := hasErr(op, f, g); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g); tmp != nil {
		return tmp
	}

	return f.Plus(g.Neg())
}

// Internal method. Multiplies the two polynomials f and g, but does not reduce
// the result according to the specified ring.
func (f *Polynomial) multNoReduce(g *Polynomial) *Polynomial {
	const op = "Multiplying polynomials"

	if tmp := hasErr(op, f, g); tmp != nil {
		return tmp
	}

	if tmp := checkCompatible(op, f, g); tmp != nil {
		return tmp
	}

	h := f.baseRing.Zero()
	for degf, cf := range f.coefs {
		for degg, cg := range g.coefs {
			degSum := uint(degf) + uint(degg)
			// Check if overflow
			if degSum < uint(degf) {
				h = f.baseRing.Zero()
				h.err = errors.New(
					op, errors.Overflow,
					"Degrees %v + %v overflow uint", degf, degg,
				)
				return h
			}
			h.SetCoef(degSum, h.Coef(degSum).Plus(cf.Mult(cg)))
		}
	}
	return h
}

// Mult returns the product of the polynomials f and g
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Mult(g *Polynomial) *Polynomial {
	h := f.multNoReduce(g)
	if h.Err() != nil {
		return h
	}
	h.reduce()
	return h
}

// Normalize creates a new polynomial obtained by normalizing f. That is,
// f.Normalize() multiplied by f.Lc() is f.
//
// If f is the zero polynomial, a copy of f is returned.
func (f *Polynomial) Normalize() *Polynomial {
	if f.Zero() {
		return f.Copy()
	}
	return f.Scale(f.Lc().Inv())
}

// Scale scales all coefficients of f by the field element c and returns the
// result as a new polynomial.
func (f *Polynomial) Scale(c *finitefield.Element) *Polynomial {
	g := f.Copy()
	for d := range g.coefs {
		deg := uint(d)
		g.SetCoef(deg, g.Coef(deg).Mult(c))
	}
	return g
}

// Pow raises f to the power of n.
//
// If the computation causes the degree of f to overflow, the returned
// polynomial has an Overflow-error as error status.
func (f *Polynomial) Pow(n uint) *Polynomial {
	const op = "Computing polynomial power"

	out := f.baseRing.Polynomial([]*finitefield.Element{
		f.BaseField().One(),
	})
	g := f.Copy()

	for n > 0 {
		if n%2 == 1 {
			out = out.Mult(g)
			if out.Err() != nil {
				out = f.baseRing.Zero()
				out.err = errors.Wrap(op, errors.Inherit, out.Err())
				return out
			}
		}
		n /= 2
		g = g.Mult(g)
	}
	return out
}

// Degrees returns a list containing the degrees is the support of f.
//
// The list is sorted according to the ring order with higher orders preceding
// lower orders in the list.
func (f *Polynomial) Degrees() []uint {
	degs := make([]uint, 0, len(f.coefs))
	for deg := len(f.coefs) - 1; deg >= 0; deg-- {
		if f.Coef(uint(deg)).Zero() {
			continue
		}
		degs = append(degs, uint(deg))
	}
	return degs
}

// Ld returns the leading degree of f.
func (f *Polynomial) Ld() uint {
	return uint(len(f.coefs)) - 1 // TODO: requires bookkeeping
}

// Lc returns the leading coefficient of f.
func (f *Polynomial) Lc() *finitefield.Element {
	return f.Coef(f.Ld())
}

// Lt returns the leading term of f.
func (f *Polynomial) Lt() *Polynomial {
	h := f.baseRing.Zero()
	ld := f.Ld()
	h.SetCoef(ld, f.Coef(ld))
	return h
}

// Zero determines whether f is the zero polynomial.
func (f *Polynomial) Zero() bool {
	if len(f.coefs) == 0 {
		return true
	}
	return false
}

// Nonzero determines whether f contains some monomial with nonzero coefficient.
func (f *Polynomial) Nonzero() bool {
	return !f.Zero()
}

// Monomial returns a bool describing whether f consists of a single monomial.
func (f *Polynomial) Monomial() bool {
	if len(f.Degrees()) == 1 {
		return true
	}
	return false
}

// Reduces f in-place
func (f *Polynomial) reduce() {
	if f.baseRing.id != nil {
		f.baseRing.id.Reduce(f)
	}
}

// hasErr is an internal method for checking if f or g has a non-nil error
// field.
//
// It returns the first polynomial with non-nil error status after wrapping the
// error. The new error inherits the kind from the old.
func hasErr(op errors.Op, f, g *Polynomial) *Polynomial {
	switch {
	case f.err != nil:
		f.err = errors.Wrap(
			op, errors.Inherit,
			f.err,
		)
		return f
	case g.err != nil:
		g.err = errors.Wrap(
			op, errors.Inherit,
			g.err,
		)
		return g
	}
	return nil
}

// checkCompatible is an internal method for checking if f and g are compatible;
// that is, if they are defined over the same ring.
//
// If not, the return value is an element with error status set to
// ArithmeticIncompat.
func checkCompatible(op errors.Op, f, g *Polynomial) *Polynomial {
	if f.baseRing != g.baseRing {
		out := f.baseRing.Zero()
		out.err = errors.New(
			op, errors.ArithmeticIncompat,
			"%v and %v defined over different rings", f, g,
		)
		return out
	}
	return nil
}

// String returns the string representation of f. Variables are named 'X' and
// 'Y'.
func (f *Polynomial) String() string {
	var b strings.Builder
	for i, d := range f.Degrees() {
		if i > 0 {
			fmt.Fprint(&b, " + ")
		}
		if tmp := f.Coef(d); !tmp.One() || d == 0 {
			fmt.Fprintf(&b, "%v", tmp)
		}
		if d == 1 {
			fmt.Fprint(&b, "X")
		}
		if d > 1 {
			fmt.Fprintf(&b, "X^%d", d)
		}
	}
	return b.String()
}
