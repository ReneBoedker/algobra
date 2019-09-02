package bivariate

import (
	"algobra/errors"
)

func max(values ...uint) uint {
	m := values[0]
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}

type Ideal struct {
	*ring
	generators []*Polynomial
	isGroebner int8 // 0=undecided, 1=true, -1=false
	isMinimal  int8 // 0=undecided, 1=true, -1=false
	isReduced  int8 // 0=undecided, 1=true, -1=false
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
		id.generators[i] = g
	}
	return id, nil
}

func (id *Ideal) reduce(f *Polynomial) {
	_, r := f.QuoRem(id.generators...)
	*f = *r // For some reason using pointers alone is not enough
}

func monomialLcm(f, g *Polynomial) (lcm *Polynomial, ok bool) {
	if !f.Monomial() || !g.Monomial() || f.baseRing != g.baseRing {
		return nil, false
	}
	ldf, ldg := f.Ld(), g.Ld()
	lcm = f.baseRing.New(map[[2]uint]uint{
		{max(ldf[0], ldg[0]), max(ldf[1], ldg[1])}: 1,
	})
	return lcm, true
}

func SPolynomial(f, g *Polynomial) (*Polynomial, error) {
	const op = "Computing S-polynomial"
	if f.baseRing != g.baseRing {
		return nil, errors.New(
			op, errors.InputIncompatible,
			"Inputs are defined over different rings",
		)
	}
	lcm, _ := monomialLcm(f.Lt(), g.Lt())
	q1, _ := lcm.QuoRem(f.Lt())
	q2, _ := lcm.QuoRem(g.Lt())
	return q1[0].Mult(f).Minus(q2[0].Mult(g)), nil
}

func (id *Ideal) GroebnerBasis() *Ideal {
	gb := make([]*Polynomial, len(id.generators))
	for i, g := range id.generators {
		gb[i] = g.Copy()
	}
	for true {
		newGens := make([]*Polynomial, 0)
		for i, f := range gb {
			for j, g := range gb {
				if j <= i {
					continue
				}
				r, _ := SPolynomial(f, g)
				id.reduce(r)
				if r.Nonzero() {
					newGens = append(newGens, r)
				}
			}
		}
		if len(newGens) == 0 {
			break
		}
		gb = append(gb, newGens...)
	}
	return &Ideal{
		ring:       id.ring,
		generators: gb,
		isGroebner: 1,
		isMinimal:  0,
		isReduced:  0,
	}
}

func (id *Ideal) MinimizeBasis() error {
	const op = "Minimizing Gröbner basis"
	if id.isGroebner != 1 {
		return errors.New(
			op, errors.InputIncompatible,
			"Given ideal is not a Gröbner basis.",
		)
	}
	lts := make([]*Polynomial, len(id.generators))
	for i := range id.generators {
		id.generators[i] = id.generators[i].Normalize()
		lts[i] = id.generators[i].Lt()
	}
	for i := 0; i < len(id.generators); {
		if _, r := lts[i].quoRemWithIgnore(i, lts...); r.Zero() {
			id.generators = append(id.generators[:i], id.generators[i+1:]...)
			lts = append(lts[:i], lts[i+1:]...)
		} else {
			i++
		}
	}
	id.isMinimal = 1
	return nil
}

func (id *Ideal) ReduceBasis() error {
	const op = "Reducing Gröbner basis"
	if id.isGroebner != 1 {
		return errors.New(
			op, errors.InputIncompatible,
			"Given ideal is not a Gröbner basis.",
		)
	}
	if id.isMinimal != 1 {
		_ = id.MinimizeBasis()
	}
	for i := range id.generators {
		_, id.generators[i] = id.generators[i].quoRemWithIgnore(i, id.generators...)
	}
	id.isReduced = 1
	return nil
}

// Write f=qg if possible; otherwise set ok=false
func (f *Polynomial) monomialDivideBy(g *Polynomial) (q *Polynomial, ok bool, err error) {
	const op = "Dividing monomials"
	if !f.Monomial() {
		return nil, false, errors.New(
			op, errors.InputIncompatible,
			"Object %v is not a monomial", f,
		)
	}
	if !g.Monomial() {
		return nil, false, errors.New(
			op, errors.InputIncompatible,
			"Input %v is not a monomial", g,
		)
	}
	ldf, ldg := f.Ld(), g.Ld()
	if d, ok := subtractDegs(ldf, ldg); ok {
		h := f.baseRing.Zero()
		h.degrees[d] = f.Coef(ldf).Mult(g.Coef(ldg).Inv())
		return h, true, nil
	}
	return nil, false, nil
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

// QuoRem return the polynomial quotient and remainder under division by the
// given list of polynomials.
func (f *Polynomial) QuoRem(list ...*Polynomial) (q []*Polynomial, r *Polynomial) {
	return f.quoRemWithIgnore(-1, list...)
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
