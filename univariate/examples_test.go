package univariate_test

import (
	"fmt"
	"log"

	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/univariate"
)

// Set up the finitefield of 4 elements for examples where the cardinality does
// not matter
func getGf9() ff.Field {
	out, _ := finitefield.Define(9)
	return out
}

var gf9 ff.Field = getGf9()

func ExamplePolynomial_SetCoef() {
	ring := univariate.DefRing(gf9)

	f := ring.One().SetCoef(25, gf9.One())
	fmt.Println(f)
	// Output:
	// X^25 + 1
}

func ExampleQuotientRing_PolynomialFromString() {
	ring := univariate.DefRing(gf9)

	f, err := ring.PolynomialFromString("(a+2)X^4 + aX^2 + 2a+2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("f(X) = %v", f)
	// Output:
	// f(X) = (a + 2)X^4 + aX^2 + (2a + 2)
}
