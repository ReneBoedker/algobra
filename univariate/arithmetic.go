package univariate

import (
	"algobra/errors"
	"algobra/finitefield/ff"
)

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
		f = tmp
		return f
	}

	for deg, c := range g.coefs {
		f.IncrementCoef(deg, c)
	}
	return f
}

// Plus returns the sum of the two polynomials f and g.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Plus(g *Polynomial) *Polynomial {
	return f.Copy().Add(g)
}

// Neg returns the polynomial obtained by scaling f by -1 (modulo the
// characteristic).
func (f *Polynomial) Neg() *Polynomial {
	return f.Copy().SetNeg()
}

// SetNeg returns the polynomial obtained by scaling f by -1 (modulo the
// characteristic).
func (f *Polynomial) SetNeg() *Polynomial {
	for _, c := range f.coefs {
		c.SetNeg()
	}
	return f
}

// Equal determines whether two polynomials are equal. That is, whether they are
// defined over the same ring, and have the same coefficients.
func (f *Polynomial) Equal(g *Polynomial) bool {
	if f.baseRing != g.baseRing {
		return false
	}
	if len(f.coefs) != len(g.coefs) {
		return false
	}
	for deg, cf := range f.coefs {
		if !g.Coef(deg).Equal(cf) {
			return false
		}
	}
	return true
}

// Sub sets f to the polynomial difference f-g and returns f.
//
// If f and g are defined over different rings, the error status of f is set to
// ArithmeticIncompat.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Sub(g *Polynomial) *Polynomial {
	const op = "Subtracting polynomials"

	if tmp := checkErrAndCompatible(op, f, g); tmp != nil {
		f = tmp
		return f
	}

	return f.Add(g.Neg())
}

// Minus returns polynomial difference f-g.
//
// If f and g are defined over different rings, a new polynomial is returned
// with an ArithmeticIncompat-error as error status.
//
// When f or g has a non-nil error status, its error is wrapped and the same
// polynomial is returned.
func (f *Polynomial) Minus(g *Polynomial) *Polynomial {
	return f.Copy().Sub(g)
}

// Internal method. Multiplies the two polynomials f and g, but does not reduce
// the result according to the specified ring.
func (f *Polynomial) multNoReduce(g *Polynomial) *Polynomial {
	const op = "Multiplying polynomials"

	if tmp := checkErrAndCompatible(op, f, g); tmp != nil {
		return tmp
	}

	if f.IsZero() || g.IsZero() {
		return f.baseRing.Zero()
	}
	h := f.baseRing.zeroWithCap(f.Ld() + g.Ld() + 1)
	for degf, cf := range f.coefs {
		for degg, cg := range g.coefs {
			degSum := degf + degg
			// Check if overflow
			if degSum < degf {
				h = f.baseRing.Zero()
				h.err = errors.New(
					op, errors.Overflow,
					"Degrees %v + %v overflow uint", degf, degg,
				)
				return h
			}
			h.IncrementCoef(degSum, cf.Times(cg))
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

// Pow raises f to the power of n.
//
// If the computation causes the degree of f to overflow, the returned
// polynomial has an Overflow-error as error status.
func (f *Polynomial) Pow(n uint) *Polynomial {
	const op = "Computing polynomial power"

	out := f.baseRing.Polynomial([]ff.Element{
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
