package univariate

import (
	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

func (f *Polynomial) SetError(err error) {
	f.err = err
}

func HasErr(op errors.Op, f *Polynomial, g ...*Polynomial) *Polynomial {
	return hasErr(op, f, g...)
}

func CheckCompatible(op errors.Op, f *Polynomial, g ...*Polynomial) *Polynomial {
	return checkCompatible(op, f, g...)
}

func AllDistinct(points []ff.Element) bool {
	return allDistinct(points)
}

func (r *QuotientRing) LagrangeBasis(points []ff.Element, ignore ff.Element) *Polynomial {
	return r.lagrangeBasis(points, ignore)
}
