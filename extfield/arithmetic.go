package extfield

import (
//"fmt"
)

// Add sets a to the sum of a and b and returns a
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Add(b *Element) *Element {
	const op = "Adding elements"

	if tmp := checkErrAndCompatible(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	if a.field.addTable != nil {
		a.val = a.field.addTable.lookup(a.val, b.val)
	} else {
		a.val.Add(b.val)
	}

	return a
}

// Plus returns the sum of elements a and b
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Plus(b *Element) *Element {
	return a.Copy().Add(b)
}

// Sub sets a to the difference of elements a and b and returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Sub(b *Element) *Element {
	const op = "Subtracting elements"

	if tmp := checkErrAndCompatible(op, a, b); tmp != nil {
		a = tmp
		return a
	}

	a.val.Sub(b.val)
	return a
}

// Minus returns the difference of elements a and b
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Minus(b *Element) *Element {
	return a.Copy().Sub(b)
}

// Prod sets a to the product of b and c, and returns a.
//
// The function returns an ArithmeticIncompat-error if b, and c are not defined
// over the same field.
func (a *Element) Prod(b, c *Element) *Element {
	const op = "Multiplying elements"

	if tmp := checkErrAndCompatible(op, b, c); tmp != nil {
		a = tmp
	}

	// Set the correct field of a
	a.field = b.field

	if a.field.multTable != nil {
		a.val = a.field.multTable.lookup(b.val, c.val)
	} else {
		a.val = (b.val.Times(c.val))
	}
	return a
}

// Times returns the product of elements a and b
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Times(b *Element) *Element {
	return a.Copy().Mult(b)
}

// Mult sets a to the product of elements a and b and returns a.
//
// If a and b are defined over different fields, a new element is returned with
// an ArithmeticIncompat-error as error status.
//
// When a or b has a non-nil error status, its error is wrapped and the same
// element is returned.
func (a *Element) Mult(b *Element) *Element {
	return a.Prod(a, b)
}
