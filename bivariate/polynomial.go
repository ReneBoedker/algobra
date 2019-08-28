package bivariate

import (
	"algobra/primefield"
	"fmt"
	"sort"
	"strings"
)

func addDegs(deg1, deg2 [2]uint) [2]uint {
	return [2]uint{deg1[0] + deg2[0], deg1[1] + deg2[1]}
}

func subtractDegs(deg1, deg2 [2]uint) (deg [2]uint, ok bool) {
	if deg1[0] >= deg2[0] && deg1[1] >= deg2[1] {
		return [2]uint{deg1[0] - deg2[0], deg1[1] - deg2[1]}, true
	}
	return deg, false
}

type Polynomial struct {
	baseRing *QuotientRing
	degrees  map[[2]uint]*primefield.Element
}

func (f *Polynomial) baseField() *primefield.Field {
	return f.baseRing.baseField
}

// Coef returns the coefficient of the monomial with degree specified by the
// input. The return value is a finite field element.
func (f *Polynomial) Coef(deg [2]uint) *primefield.Element {
	if c, ok := f.degrees[deg]; ok {
		return c
	}
	return f.baseField().Element(0)
}

// Copy returns a new polynomial object over the same ring and with the same
// coefficients as f.
func (f *Polynomial) Copy() *Polynomial {
	h := f.baseRing.Zero()
	for deg, c := range f.degrees {
		h.degrees[deg] = c
	}
	return h
}

// Plus returns the sum of the two polynomials f and g.
func (f *Polynomial) Plus(g *Polynomial) *Polynomial {
	if f.baseRing != g.baseRing {
		fNew, gNew, err := embedInCommonRing(f, g)
		if err != nil {
			panic(fmt.Sprintf("Polynomial.Plus: Polynomials incompatible\n'%#v'\n'%#v')",
				f.baseRing, g.baseRing))
		}
		return fNew.Plus(gNew)
	}
	h := f.Copy()
	for deg, c := range g.degrees {
		if _, ok := h.degrees[deg]; !ok {
			h.degrees[deg] = c
			continue
		}
		tmp := h.Coef(deg).Plus(c)
		if tmp.Nonzero() {
			h.degrees[deg] = tmp
		} else {
			delete(h.degrees, deg)
		}
	}
	return h
}

// Neg returns the polynomial obtained by scaling f by -1 (modulo the
// characteristic).
func (f *Polynomial) Neg() *Polynomial {
	g := f.baseRing.Zero()
	for deg, c := range f.degrees {
		g.degrees[deg] = c.Neg()
	}
	return g
}

// Equal determines whether two polynomials are equal. That is, whether they are
// defined over the same ring, and have the same coefficients.
func (f *Polynomial) Equal(g *Polynomial) bool {
	if f.baseRing != g.baseRing {
		return false
	}
	if len(f.degrees) != len(g.degrees) {
		return false
	}
	for d, cf := range f.degrees {
		if cg, ok := g.degrees[d]; !ok || !cg.Equal(cf) {
			return false
		}
	}
	return true
}

// Minus returns polynomial difference f-g.
func (f *Polynomial) Minus(g *Polynomial) *Polynomial {
	return f.Plus(g.Neg())
}

// Internal method. Multiplies the two polynomials f and g, but does not reduce
// the result according to the specified ring.
func (f *Polynomial) multNoReduce(g *Polynomial) *Polynomial {
	h := f.baseRing.Zero()
	for degf, cf := range f.degrees {
		for degg, cg := range g.degrees {
			tmp := cf.Mult(cg)
			if tmp.Nonzero() {
				if c, ok := h.degrees[addDegs(degf, degg)]; ok {
					if c.Plus(tmp).Nonzero() {
						h.degrees[addDegs(degf, degg)] = c.Plus(tmp)
					} else {
						delete(h.degrees, addDegs(degf, degg))
					}
				} else {
					h.degrees[addDegs(degf, degg)] = tmp
				}
			}
		}
	}
	return h
}

// Mult returns the product of the polynomials f and g
func (f *Polynomial) Mult(g *Polynomial) *Polynomial {
	h := f.multNoReduce(g)
	h.reduce()
	return h
}

// Normalize creates a new polynomial obtained by normalizing f. That is,
// f.Normalize() multiplied by f.Lc() is f. If f is the zero polynomial, the
// zero polynomial is returned.
func (f *Polynomial) Normalize() *Polynomial {
	if f.Zero() {
		return f
	}
	return f.Scale(f.Lc().Inv())
}

// Scale scales all coefficients of f by the given field element and returns the
// result as a new polynomial.
func (f *Polynomial) Scale(c *primefield.Element) *Polynomial {
	g := f.Copy()
	for d := range g.degrees {
		g.degrees[d] = g.degrees[d].Mult(c)
	}
	return g
}

// Pow raises f to the power of n, and return the result in a new polynomial.
func (f *Polynomial) Pow(n uint) *Polynomial {
	out := f.baseRing.New(map[[2]uint]uint{
		{0, 0}: 1,
	})
	g := f.Copy()
	for n > 0 {
		if n%2 == 1 {
			out = out.Mult(g)
		}
		n /= 2
		g = g.Mult(g)
	}
	return out
}

// SortedDegrees returns a list containing the degrees is the support of f. The
// list is sorted according to the ring order with higher orders preceding
// lower orders in the list.
func (f *Polynomial) SortedDegrees() [][2]uint {
	degs := make([][2]uint, 0, len(f.degrees))
	for deg := range f.degrees {
		degs = append(degs, deg)
	}
	sort.Slice(degs, func(i, j int) bool {
		return (f.baseRing.ord(degs[i], degs[j]) >= 0)
	})
	return degs
}

// Ld returns the leading degree of f.
func (f *Polynomial) Ld() [2]uint {
	return f.SortedDegrees()[0]
}

// Lc returns the leading coefficient of f.
func (f *Polynomial) Lc() *primefield.Element {
	return f.Coef(f.Ld())
}

// Lt returns the leading term of f.
func (f *Polynomial) Lt() *Polynomial {
	h := f.baseRing.Zero()
	ld := f.Ld()
	h.degrees[ld] = f.Coef(ld)
	return h
}

// Zero determines whether f is the zero polynomial.
func (f *Polynomial) Zero() bool {
	if len(f.degrees) == 0 {
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
	if len(f.degrees) == 1 {
		return true
	}
	return false
}

// Reduces f in-place
func (f *Polynomial) reduce() {
	if f.baseRing.id != nil {
		f.baseRing.id.reduce(f)
	}
}

// Embed f in another ring
func embedInCommonRing(f, g *Polynomial) (fOut, gOut *Polynomial, err error) {
	fOut = f.Copy()
	gOut = g.Copy()
	if f.baseRing.ring != g.baseRing.ring {
		err = fmt.Errorf("embedInCommonRing: Rings '%v' and '%v' are not compatible",
			f.baseRing.ring, g.baseRing.ring,
		)
		return
	}
	switch {
	case f.baseRing.id == nil && g.baseRing.id == nil:
		err = nil
	case f.baseRing.id == nil && g.baseRing.id != nil:
		fOut.baseRing = g.baseRing
		err = nil
	case f.baseRing != nil && g.baseRing.id == nil:
		gOut.baseRing = f.baseRing
	case f.baseRing != nil && g.baseRing.id != nil:
		err = fmt.Errorf("embedInCommonRing: Polynomials defined over different quotient rings.")
	}
	return
}

// String returns the string representation of f. Variables are named 'X' and
// 'Y'.
func (f *Polynomial) String() string {
	degs := f.SortedDegrees()
	if len(degs) == 0 {
		return "0"
	}
	var b strings.Builder
	for i, d := range degs {
		if i > 0 {
			fmt.Fprint(&b, " + ")
		}
		if tmp := f.Coef(d); !tmp.One() || (d[0] == 0 && d[1] == 0) {
			fmt.Fprintf(&b, "%v", tmp)
		}
		if d[0] == 1 {
			fmt.Fprint(&b, "X")
		}
		if d[0] > 1 {
			fmt.Fprintf(&b, "X^%d", d[0])
		}
		if d[1] == 1 {
			fmt.Fprint(&b, "Y")
		}
		if d[1] > 1 {
			fmt.Fprintf(&b, "Y^%d", d[1])
		}
	}
	return b.String()
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
