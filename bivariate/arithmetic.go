package bivariate

import (
	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

// Plus returns the sum of the two polynomials f and g.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Plus(g *Polynomial) *Polynomial {
	const op = "Adding polynomials"

	if tmp := checkErrAndCompatible(op, f, g); tmp != nil {
		return tmp
	}

	h := f.Copy()
	for deg, c := range g.coefs {
		h.IncrementCoef(deg, c)
	}
	return h
}

// Add sets f to the sum of the two polynomials f and g and returns f.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Add(g *Polynomial) *Polynomial {
	const op = "Adding polynomials"

	if tmp := checkErrAndCompatible(op, f, g); tmp != nil {
		return tmp
	}

	for deg, c := range g.coefs {
		f.IncrementCoef(deg, c)
	}
	return f
}

// Neg returns the polynomial obtained by scaling f by -1 (modulo the
// characteristic).
func (f *Polynomial) Neg() *Polynomial {
	g := f.baseRing.zeroWithCap(len(f.coefs))
	for deg, c := range f.coefs {
		g.coefs[deg] = c.Neg()
	}
	return g
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

	if tmp := checkErrAndCompatible(op, f, g); tmp != nil {
		return tmp
	}

	return f.Plus(g.Neg())
}

// Internal method. Multiplies the two polynomials f and g, but does not reduce
// the result according to the specified ring.
func (f *Polynomial) multNoReduce(g *Polynomial) *Polynomial {
	const op = "Multiplying polynomials"

	if tmp := checkErrAndCompatible(op, f, g); tmp != nil {
		return tmp
	}

	h := f.baseRing.zeroWithCap(len(f.coefs) * len(g.coefs))
	tmp := f.BaseField().One()
	for degf, cf := range f.coefs {
		for degg, cg := range g.coefs {
			degSum, err := addDegs(degf, degg)
			if err != nil {
				h = f.baseRing.Zero()
				h.err = errors.Wrap(op, errors.Inherit, err)
				return h
			}
			tmp.Prod(cf, cg)
			h.IncrementCoef(degSum, tmp)
		}
	}
	return h
}

// Times returns the product of the polynomials f and g
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Times(g *Polynomial) *Polynomial {
	h := f.multNoReduce(g)
	if h.Err() != nil {
		return h
	}
	h.reduce()
	return h
}

// Mult sets f to the product of the polynomials f and g and returns f.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Mult(g *Polynomial) *Polynomial {
	*f = *f.multNoReduce(g)
	if f.Err() != nil {
		return f
	}
	f.reduce()
	return f
}

// Normalize creates a new polynomial obtained by normalizing f. That is,
// f.Normalize() multiplied by f.Lc() is f.
//
// If f is the zero polynomial, a copy of f is returned.
func (f *Polynomial) Normalize() *Polynomial {
	if f.IsZero() {
		return f.Copy()
	}
	return f.Scale(f.lcPtr().Inv())
}

// Scale scales all coefficients of f by the field element c and returns the
// result as a new polynomial. See also SetScale.
func (f *Polynomial) Scale(c ff.Element) *Polynomial {
	if c.IsZero() {
		return f.baseRing.Zero()
	}

	g := f.Copy()
	for d := range g.coefs {
		g.coefs[d].Mult(c)
	}
	return g
}

// SetScale scales all coefficients of f by the field element c and returns
// f. See also Scale.
func (f *Polynomial) SetScale(c ff.Element) *Polynomial {
	if c.IsZero() {
		return f.baseRing.Zero()
	}

	for d := range f.coefs {
		f.coefs[d].Mult(c)
	}
	return f
}

// Pow raises f to the power of n.
//
// If the computation causes the degree of f to overflow, the returned
// polynomial has an Overflow-error as error status.
func (f *Polynomial) Pow(n uint) *Polynomial {
	const op = "Computing polynomial power"

	out := f.baseRing.Polynomial(map[[2]uint]ff.Element{
		{0, 0}: f.BaseField().One(),
	})
	g := f.Copy()

	for n > 0 {
		if n%2 == 1 {
			out.Mult(g)
			if out.Err() != nil {
				out = f.baseRing.Zero()
				out.err = errors.Wrap(op, errors.Inherit, out.Err())
				return out
			}
		}
		n /= 2
		g.Mult(g)
	}
	return out
}

// QuoRem returns the polynomial quotient and remainder under division by the
// given list of polynomials.
func (f *Polynomial) QuoRem(list ...*Polynomial) (q []*Polynomial, r *Polynomial, err error) {
	return f.quoRemWithIgnore(-1, list...)
}

func (f *Polynomial) quoRemWithIgnore(
	ignoreIndex int,
	list ...*Polynomial,
) (q []*Polynomial, r *Polynomial, err error) {
	const op = "Computing polynomial quotient and remainder"

	if tmp := checkErrAndCompatible(op, f, list...); tmp != nil {
		err = tmp.Err()
		return
	}

	r = f.baseRing.Zero()
	p := f.Copy()

	q = make([]*Polynomial, len(list), len(list))
	for i := range list {
		q[i] = f.baseRing.Zero()
	}
outer:
	for p.IsNonzero() {
		for i, g := range list {
			if i == ignoreIndex {
				continue
			}
			// Below, err is ignored since both p and g are nonzero (so both
			// leading terms are well defined, and monomialDivideBy will not
			// return an error)
			if mquo, ok, _ := p.Lt().monomialDivideBy(g.Lt()); ok {
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
	return q, r, nil
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
