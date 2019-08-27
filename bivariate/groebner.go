package bivariate

import (
	"fmt"
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

type ideal []*Polynomial

// NewIdeal returns a new polynomial ideal over the given ring. If the
// generators are not defined over the given ring, the function panics.
// Internally, this function computes a reduced Gröbner basis.
func (r *Ring) NewIdeal(generators ...*Polynomial) ideal {
	targetRing := generators[0].baseRing
	for _, g := range generators {
		if g.baseRing != targetRing {
			panic("ring.NewIdeal: Generators defined over different rings")
		}
	}
	return groebner(generators)
}

func (id ideal) reduce(f *Polynomial) {
	_, r := f.QuoRem(id...)
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
	if f.baseRing != g.baseRing {
		return nil, fmt.Errorf("sPoly: Inputs are defined over different rings")
	}
	lcm, _ := monomialLcm(f.Lt(), g.Lt())
	q1, _ := lcm.QuoRem(f.Lt())
	q2, _ := lcm.QuoRem(g.Lt())
	return q1[0].Mult(f).Minus(q2[0].Mult(g)), nil
}

func groebner(id ideal) ideal {
	gb := make(ideal, len(id))
	for i, g := range id {
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
	// ==== Convert gb to a minimal Gröbner basis ====
	lts := make([]*Polynomial, len(gb))
	for i := range gb {
		gb[i] = gb[i].Normalize() //.Scale(gb[i].Lc().Inv())
		lts[i] = gb[i].Lt()
	}
	for i := 0; i < len(gb); {
		if _, r := lts[i].quoRemWithIgnore(i, id...); r.Zero() {
			gb = append(gb[:i], gb[i+1:]...)
			lts = append(lts[:i], lts[i+1:]...)
		} else {
			i++
		}
	}
	// ==== Convert to a reduced Gröbner basis ====
	for i := range gb {
		_, gb[i] = gb[i].quoRemWithIgnore(i, gb...)
	}
	return gb
}

// Write f=qg if possible; otherwise set ok=false
func (f *Polynomial) monomialDivideBy(g *Polynomial) (q *Polynomial, ok bool, err error) {
	if !f.Monomial() || !g.Monomial() {
		return nil, false, fmt.Errorf("Polynomial.monomialDividesBy: Input is not a monomial")
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
