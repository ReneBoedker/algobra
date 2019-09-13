package finitefield

import (
	"algobra/errors"
	"algobra/primefield"
)

type Field struct {
	pf *primefield.Field
	//Non-prime field type at some point
}

type Element struct {
	pf *primefield.Element
	// Non-prime field type at some point
	err error
}

type kind uint8

const (
	primeKind = iota
	nonPrimeKind
)

func (f *Field) kind() kind {
	switch {
	case f.pf != nil:
		return primeKind
	default:
		panic("Error")
	}
}

func (a *Element) kind() kind {
	switch {
	case a.pf != nil:
		return primeKind
	default:
		panic("Error")
	}
}

func Define(char uint) (*Field, error) {
	// Must be changed when non-prime is implemented
	pf, err := primefield.Define(char)
	if err != nil {
		return nil, err
	}
	return &Field{
		pf: pf,
	}, nil
}

func (f *Field) ElementFromUnsigned(val uint) *Element {
	switch {
	case f.pf != nil:
		return &Element{
			pf: f.pf.Element(val),
		}
	default:
		panic("Error")
	}
}

func (f *Field) ElementFromSigned(val int) *Element {
	switch {
	case f.pf != nil:
		return &Element{
			pf: f.pf.ElementFromSigned(val),
		}
	default:
		panic("Error")
	}
}

func (f *Field) Element(val interface{}) (*Element, error) {
	const op = "Defining field element"

	switch v := val.(type) {
	case uint:
		return f.ElementFromUnsigned(v), nil
	case int:
		return f.ElementFromSigned(v), nil
	default:
		return nil, errors.New(
			op, errors.Input,
			"Cannot create element from type %T", v,
		)
	}
}

func (f *Field) Elements() []*Element {
	out := make([]*Element, f.Card(), f.Card())
	switch f.kind() {
	case primeKind:
		for i, e := range f.pf.Elements() {
			out[i] = &Element{
				pf: e,
			}
		}
	default:
		panic("Error")
	}
	return out
}

func (f *Field) Zero() *Element {
	switch {
	case f.pf != nil:
		return &Element{
			pf: f.pf.Element(0),
		}
	default:
		panic("Error")
	}
}

func (f *Field) One() *Element {
	switch {
	case f.pf != nil:
		return &Element{
			pf: f.pf.Element(1),
		}
	default:
		panic("Error")
	}
}

func (f *Field) Card() uint {
	switch f.kind() {
	case primeKind:
		return f.pf.Card()
	default:
		panic("Error")
	}
}

func (a *Element) Copy() *Element {
	switch a.kind() {
	case primeKind:
		return &Element{
			pf:  a.pf.Copy(),
			err: a.err,
		}
	default:
		panic("Error")
	}
}

func (a *Element) Plus(b *Element) *Element {
	const op = "Adding elements"

	if a.kind() != b.kind() {
		return &Element{
			err: errors.New(
				op, errors.ArithmeticIncompat,
				"Cannot add elements from different fields",
			),
		}
	}

	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Plus(b.pf),
		}
	default:
		panic("Error")
	}
}

func (a *Element) Mult(b *Element) *Element {
	const op = "Multiplying elements"

	if a.kind() != b.kind() {
		return &Element{
			err: errors.New(
				op, errors.ArithmeticIncompat,
				"Cannot multiply elements from different fields",
			),
		}
	}

	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Mult(b.pf),
		}
	default:
		panic("Error")
	}
}

func (a *Element) Minus(b *Element) *Element {
	const op = "Subtracting elements"

	if a.kind() != b.kind() {
		return &Element{
			err: errors.New(
				op, errors.ArithmeticIncompat,
				"Cannot subtract elements from different fields",
			),
		}
	}

	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Minus(b.pf),
		}
	default:
		panic("Error")
	}
}

func (a *Element) Equal(b *Element) bool {
	if a.kind() != b.kind() {
		return false
	}

	switch a.kind() {
	case primeKind:
		return a.pf.Equal(b.pf)
	default:
		panic("Error")
	}
}

func (a *Element) Inv() *Element {
	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Inv(),
		}
	default:
		panic("Error")
	}
}

func (a *Element) Neg() *Element {
	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Neg(),
		}
	default:
		panic("Error")
	}
}

func (a *Element) Pow(n uint) *Element {
	switch a.kind() {
	case primeKind:
		return &Element{
			pf: a.pf.Pow(n),
		}
	default:
		panic("Error")
	}
}

func (a *Element) Zero() bool {
	switch a.kind() {
	case primeKind:
		return a.pf.Zero()
	default:
		panic("Error")
	}
}

func (a *Element) Nonzero() bool {
	switch a.kind() {
	case primeKind:
		return a.pf.Nonzero()
	default:
		panic("Error")
	}
}

func (a *Element) One() bool {
	switch a.kind() {
	case primeKind:
		return a.pf.One()
	default:
		panic("Error")
	}
}

func (a *Element) Err() error {
	if a.err != nil {
		return a.err
	}
	switch a.kind() {
	case primeKind:
		return a.pf.Err()
	default:
		panic("Error")
	}
}

func (a *Element) String() string {
	switch a.kind() {
	case primeKind:
		return a.pf.String()
	default:
		panic("Error")
	}
}
