package bivariate

import (
	"algobra/errors"
	"algobra/finitefield/ff"
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

// monomialLcm returns the least common multiple of two monomials.
//
// If either of the inputs is not a monomial, or they are defined over different
// rings, nil is returned
func monomialLcm(f, g *Polynomial) (lcm *Polynomial) {
	if !f.IsMonomial() || !g.IsMonomial() || f.baseRing != g.baseRing {
		return nil
	}
	ldf, ldg := f.Ld(), g.Ld()
	lcm = f.baseRing.Polynomial(map[[2]uint]ff.Element{
		{max(ldf[0], ldg[0]), max(ldf[1], ldg[1])}: f.BaseField().One(),
	})
	return lcm
}

// SPolynomial computes the S-polynomial of f and g.
//
// It returns an ArithmeticIncompat-error if f and g are defined over different
// rings.
func SPolynomial(f, g *Polynomial) (*Polynomial, error) {
	const op = "Computing S-polynomial"

	if tmp := checkErrAndCompatible(op, f, g); tmp != nil {
		return nil, tmp.Err()
	}

	lcm := monomialLcm(f.Lt(), g.Lt()) // Ignore error since Lt() is monomial
	q1, _, _ := lcm.QuoRem(f.Lt())     // Remainder always zero, err has been checked
	q2, _, _ := lcm.QuoRem(g.Lt())     // (as above)
	return q1[0].Mult(f).Minus(q2[0].Mult(g)), nil
}

// GroebnerBasis computes a Gröbner basis for id. The result is returned as a
// new ideal object.
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
				id.Reduce(r)
				if r.IsNonzero() {
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

// MinimizeBasis transforms the generators of id into a minimal Gröbner basis.
//
// If the generators of id do not form a Gröbner basis, the function returns an
// InputValue-error.
func (id *Ideal) MinimizeBasis() error {
	const op = "Minimizing Gröbner basis"

	if id.isGroebner != 1 {
		return errors.New(
			op, errors.InputValue,
			"Given ideal is not a Gröbner basis.",
		)
	}

	lts := make([]*Polynomial, len(id.generators))
	for i := range id.generators {
		id.generators[i] = id.generators[i].Normalize()
		lts[i] = id.generators[i].Lt()
	}

	for i := 0; i < len(id.generators); {
		if _, r, _ := lts[i].quoRemWithIgnore(i, lts...); r.IsZero() {
			id.generators = append(id.generators[:i], id.generators[i+1:]...)
			lts = append(lts[:i], lts[i+1:]...)
		} else {
			i++
		}
	}

	id.isMinimal = 1
	return nil
}

// ReduceBasis transforms the generators of id into a reduced Gröbner basis.
//
// If the generators of id do not form a Gröbner basis, the function returns an
// InputValue-error.
func (id *Ideal) ReduceBasis() error {
	const op = "Reducing Gröbner basis"

	if id.isGroebner != 1 {
		return errors.New(
			op, errors.InputValue,
			"Given ideal is not a Gröbner basis.",
		)
	}
	if id.isMinimal != 1 {
		_ = id.MinimizeBasis()
	}

	for i := range id.generators {
		_, id.generators[i], _ = id.generators[i].quoRemWithIgnore(i, id.generators...)
	}

	id.isReduced = 1
	return nil
}

// Write f=qg if possible; otherwise set ok=false
func (f *Polynomial) monomialDivideBy(g *Polynomial) (q *Polynomial, ok bool, err error) {
	const op = "Dividing monomials"

	if !f.IsMonomial() {
		return nil, false, errors.New(
			op, errors.InputValue,
			"Object %v is not a monomial", f,
		)
	}
	if !g.IsMonomial() {
		return nil, false, errors.New(
			op, errors.InputValue,
			"Input %v is not a monomial", g,
		)
	}

	ldf, ldg := f.Ld(), g.Ld()
	if d, ok := subtractDegs(ldf, ldg); ok {
		h := f.baseRing.Zero()
		h.coefs[d] = f.Coef(ldf).Times(g.Coef(ldg).Inv())
		return h, true, nil
	}
	return nil, false, nil
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
