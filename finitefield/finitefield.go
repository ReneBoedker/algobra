package finitefield

import (
	"algobra/basic"
	"algobra/errors"
	"algobra/finitefield/extfield"
	"algobra/finitefield/ff"
	"algobra/finitefield/primefield"
)

// Define returns a new finite field with the given cardinality.
//
// If card is not a prime power, an InputValue-error is returned.
func Define(card uint) (ff.Field, error) {
	const op = "Defining finite field"

	_, extDeg, err := basic.FactorizePrimePower(card)
	if err != nil {
		return nil, errors.Wrap(op, errors.InputValue, err)
	}

	if extDeg == 1 {
		return extfield.Define(card)
	}
	return primefield.Define(card)
}
