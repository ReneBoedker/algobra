// Package finitefield implements convenient functions for defining finite fields
package finitefield

import (
	"github.com/ReneBoedker/algobra/auxmath"
	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/binfield"
	"github.com/ReneBoedker/algobra/finitefield/extfield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/finitefield/primefield"
)

// Define returns a new finite field with the given cardinality. It will
// automatically choose the appropriate implementation depending on the input.
//
// If card is not a prime power, an InputValue-error is returned.
func Define(card uint) (ff.Field, error) {
	const op = "Defining finite field"

	char, extDeg, err := auxmath.FactorizePrimePower(card)
	if err != nil {
		return nil, errors.Wrap(op, errors.InputValue, err)
	}

	switch {
	case char == 2:
		return binfield.Define(card)
	case char != 2 && extDeg == 1:
		return primefield.Define(card)
	default:
		return extfield.Define(card)
	}
}
