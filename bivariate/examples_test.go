package bivariate_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/bivariate"
	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

// Set up the finitefield of 9 elements for examples where the cardinality does
// not matter
func getGf9() ff.Field {
	out, _ := finitefield.Define(9)
	return out
}

var gf9 ff.Field = getGf9()

func ExampleQuotientRing_NewIdeal() {
	ring := bivariate.DefRing(gf9, bivariate.WDegLex(3, 4, false))
	id, _ := ring.NewIdeal(
		ring.Polynomial(map[[2]uint]ff.Element{
			{0, 3}: gf9.One(),
			{4, 0}: gf9.One().SetNeg(),
			{0, 1}: gf9.One(),
		}),
		ring.Polynomial(map[[2]uint]ff.Element{
			{9, 0}: gf9.One(),
			{1, 0}: gf9.One().SetNeg(),
		}),
	)

	fmt.Println(id)
	// Output:
	// Ideal <Y^3 + 2X^4 + Y, X^9 + 2X> in Bivariate polynomial ring over Finite field of 9 elements
}
