package algobra_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

// The so-called places of the Hermitian function field over the finite field of
// q^2 elements can be represented by pairs (α, β) satisfying α^(q+1)=β^q+β.
// It is well-known that there are q^3 such pairs, see for instance
// [Stichtenoth, 2009].
//
// This example computes the pairs (α, β) for q=3.
func Example_hermitianPlaces() {
	// Define the norm map given by N(a)=a^(q+1) in the Hermitian case
	norm := func(field ff.Field, a ff.Element) ff.Element {
		return a.Pow(field.Char() + 1)
	}

	// Define the trace map given by Tr(a)=a^q+a in the Hermitian case
	trace := func(field ff.Field, a ff.Element) ff.Element {
		return a.Pow(field.Char()).Add(a)
	}

	// trInv computes the preimage of a with respect to the trace map.
	trInv := func(field ff.Field, a ff.Element) ([]ff.Element, error) {
		if !a.Equal(a.Pow(field.Char())) {
			// The input must be in the image of the trace.
			// That is, it must be in F_q
			return nil, fmt.Errorf("%v is not in the image of the trace", a)
		}

		out := make([]ff.Element, 0, field.Char())
		for _, b := range field.Elements() {
			// Search for those elements that map to a under the trace map
			if trace(field, b).Equal(a) {
				out = append(out, b)
			}
			// Stop searching once all q solutions have been found
			if uint(len(out)) == field.Char() {
				break
			}
		}

		return out, nil
	}

	// Define the field and the list of places.
	field, _ := finitefield.Define(9)
	places := make([][2]ff.Element, 0, 27)

	for _, a := range field.Elements() {
		// Find the preimage of a
		tmp, err := trInv(field, norm(field, a))
		if err != nil {
			fmt.Println(err)
			return
		}

		// Append the found pairs
		for _, b := range tmp {
			places = append(places, [2]ff.Element{a, b})
		}
	}

	for _, p := range places {
		fmt.Printf("α = %v, β = %v\n", p[0], p[1])
	}

	//Unordered output:
	// α = 0, β = 0
	// α = 0, β = a + 1
	// α = 0, β = 2a + 2
	// α = 1, β = a
	// α = 1, β = 2a + 1
	// α = 1, β = 2
	// α = a, β = 1
	// α = a, β = 2a
	// α = a, β = a + 2
	// α = a + 1, β = a
	// α = a + 1, β = 2a + 1
	// α = a + 1, β = 2
	// α = 2a + 1, β = 1
	// α = 2a + 1, β = 2a
	// α = 2a + 1, β = a + 2
	// α = 2, β = a
	// α = 2, β = 2a + 1
	// α = 2, β = 2
	// α = 2a, β = 1
	// α = 2a, β = 2a
	// α = 2a, β = a + 2
	// α = 2a + 2, β = a
	// α = 2a + 2, β = 2a + 1
	// α = 2a + 2, β = 2
	// α = a + 2, β = 1
	// α = a + 2, β = 2a
	// α = a + 2, β = a + 2
}
