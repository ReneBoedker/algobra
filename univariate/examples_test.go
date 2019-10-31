package univariate_test

import (
	"fmt"
	"log"

	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/univariate"
)

func ExamplePolynomial_SetCoef() {
	gf4, _ := finitefield.Define(4)
	ring := univariate.DefRing(gf4)

	f := ring.One().SetCoef(25, gf4.One())
	fmt.Println(f)
	// Output:
	// X^25 + 1
}

func ExamplePolynomial_Lc() {
	gf9, _ := finitefield.Define(9)
	ring := univariate.DefRing(gf9)
	f, _ := ring.PolynomialFromString("(a+2)X^4 + aX^2 + (2a+2)X - 1")

	fmt.Println(f.Lc())
	// Output:
	// a + 2
}

func ExamplePolynomial_Ld() {
	gf9, _ := finitefield.Define(9)
	ring := univariate.DefRing(gf9)
	f, _ := ring.PolynomialFromString("(a+2)X^4 + aX^2 + (2a+2)X - 1")

	fmt.Println(f.Ld())
	// Output:
	// 4
}
func ExamplePolynomial_Lt() {
	gf9, _ := finitefield.Define(9)
	ring := univariate.DefRing(gf9)
	f, _ := ring.PolynomialFromString("(a+2)X^4 + aX^2 + (2a+2)X - 1")

	fmt.Println(f.Lt())
	// Output:
	// (a + 2)X^4
}

func ExampleIdeal_ShortString() {
	gf7, _ := finitefield.Define(7)
	ring := univariate.DefRing(gf7)

	f, _ := ring.PolynomialFromString("X^7-X")
	id, _ := ring.NewIdeal(f)

	fmt.Println(id.ShortString())
	// Output:
	// <X^7 + 6X>
}

func ExampleQuotientRing_PolynomialFromString() {
	gf9, _ := finitefield.Define(9)
	ring := univariate.DefRing(gf9)

	f, err := ring.PolynomialFromString("(a+2)X^4 + aX^2 + 2a+2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("f(X) = %v", f)
	// Output:
	// f(X) = (a + 2)X^4 + aX^2 + (2a + 2)
}

func ExampleQuotientRing_NewIdeal() {
	gf7, _ := finitefield.Define(7)
	ring := univariate.DefRing(gf7)

	f, _ := ring.PolynomialFromString("X^7-X")
	id, _ := ring.NewIdeal(f)

	fmt.Println(id)
	// Output:
	// Ideal <X^7 + 6X> of Univariate polynomial ring in X over Finite field of 7 elements
}

func ExampleQuotientRing_Quotient() {
	gf7, _ := finitefield.Define(7)
	ring := univariate.DefRing(gf7)

	f, _ := ring.PolynomialFromString("X^7-X")
	id, _ := ring.NewIdeal(f)

	ring, _ = ring.Quotient(id)
	fmt.Println(ring)
	// Output:
	// Quotient ring of univariate polynomials in X over Finite field of 7 elements modulo <X^7 + 6X>
}

func ExampleQuotientRing_SetVarName() {
	gf4, _ := finitefield.Define(4)
	ring := univariate.DefRing(gf4)

	f, _ := ring.PolynomialFromString("X^4 + (a+1)X^3 + a")

	// Change the variable names
	err := ring.SetVarName("Y")
	if err != nil {
		fmt.Printf("Could not set variable names: %q", err)
	}

	g, _ := ring.PolynomialFromString("aY^3 + (a+1)Y^2 + Y")

	// Both f and g are affected by the change
	fmt.Printf("f(Y) = %v\ng(Y) = %v", f, g)
	// Output:
	// f(Y) = Y^4 + (a + 1)Y^3 + a
	// g(Y) = aY^3 + (a + 1)Y^2 + Y
}
