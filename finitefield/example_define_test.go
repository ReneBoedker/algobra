package finitefield_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/finitefield"
)

func ExampleDefine() {
	for _, card := range []uint{5, 9, 256, 10} {
		field, err := finitefield.Define(card)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Printf("%v (type %[1]T)\n", field)
	}
	// Output:
	// Finite field of 5 elements (type *primefield.Field)
	// Finite field of 9 elements (type *extfield.Field)
	// Finite field of 256 elements (type *binfield.Field)
	// Error: Defining finite field: Factorizing prime power: 10 does not seem to be a prime power.
}
