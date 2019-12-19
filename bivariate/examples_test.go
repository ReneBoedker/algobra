package bivariate_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/bivariate"
	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

func ExampleIdeal_IsMinimal() {
	gf9, _ := finitefield.Define(9)
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
		ring.Polynomial(map[[2]uint]ff.Element{
			{0, 9}: gf9.One(),
			{0, 1}: gf9.One().SetNeg(),
		}),
	)

	fmt.Printf("%s is a Gröbner basis: %t\n", id.ShortString(), id.IsGroebner())
	fmt.Printf("%s is a minimal Gröbner basis: %t\n", id.ShortString(), id.IsMinimal())

	id.MinimizeBasis()
	fmt.Printf("%s is a minimal Gröbner basis: %t\n", id.ShortString(), id.IsMinimal())
	// Output:
	// <Y^3 + 2X^4 + Y, X^9 + 2X, Y^9 + 2Y> is a Gröbner basis: true
	// <Y^3 + 2X^4 + Y, X^9 + 2X, Y^9 + 2Y> is a minimal Gröbner basis: false
	// <Y^3 + 2X^4 + Y, X^9 + 2X> is a minimal Gröbner basis: true
}

func ExampleIdeal_ShortString() {
	field, _ := finitefield.Define(7)
	ring := bivariate.DefRing(field, bivariate.Lex(true))

	id, _ := ring.NewIdeal(
		ring.Polynomial(map[[2]uint]ff.Element{
			{2, 0}: field.One(),
			{0, 1}: field.One(),
		}),
		ring.Polynomial(map[[2]uint]ff.Element{
			{2, 1}: field.One(),
			{0, 0}: field.One(),
		}),
	)

	fmt.Println(id.ShortString())
	// Output:
	// <X^2 + Y, X^2Y + 1>
}

func ExamplePolynomial_Lc() {
	gf9, _ := finitefield.Define(9)
	ring := bivariate.DefRing(gf9, bivariate.DegLex(true))
	f, _ := ring.PolynomialFromString("(a+2)X^4Y + aX^2Y^3 + Y^4 - 1")

	fmt.Println(f.Lc())
	// Output:
	// a + 2
}

func ExamplePolynomial_Ld() {
	gf9, _ := finitefield.Define(9)
	ring := bivariate.DefRing(gf9, bivariate.DegLex(true))
	f, _ := ring.PolynomialFromString("(a+2)X^4Y + aX^2Y^3 + Y^4 - 1")

	fmt.Println(f.Ld())
	// Output:
	// [4 1]
}
func ExamplePolynomial_Lt() {
	gf9, _ := finitefield.Define(9)
	ring := bivariate.DefRing(gf9, bivariate.DegLex(true))
	f, _ := ring.PolynomialFromString("(a+2)X^4Y + aX^2Y^3 + Y^4 - 1")

	fmt.Println(f.Lt())
	// Output:
	// (a + 2)X^4Y
}

func ExamplePolynomial_SetCoef() {
	gf5, _ := finitefield.Define(5)
	ring := bivariate.DefRing(gf5, bivariate.WDegLex(2, 1, false))
	f := ring.PolynomialFromUnsigned(map[[2]uint]uint{
		{3, 1}: 3,
		{1, 4}: 1,
		{0, 0}: 3,
	})
	fmt.Println(f)

	f.SetCoef([2]uint{1, 4}, gf5.ElementFromUnsigned(4))
	f.SetCoef([2]uint{1, 1}, gf5.One())
	fmt.Println(f)
	// Output:
	// 3X^3Y + XY^4 + 3
	// 3X^3Y + 4XY^4 + XY + 3
}

func ExamplePolynomial_SortedDegrees() {
	gf9, _ := finitefield.Define(9)
	ring := bivariate.DefRing(gf9, bivariate.DegLex(true))
	f, _ := ring.PolynomialFromString("(a+2)X^4Y + aX^2Y^3 + Y^4 - 1")

	fmt.Println(f.SortedDegrees())
	// Output:
	// [[4 1] [2 3] [0 4] [0 0]]
}

func ExampleQuotientRing_NewIdeal() {
	gf9, _ := finitefield.Define(9)
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
	// Ideal <Y^3 + 2X^4 + Y, X^9 + 2X> of Bivariate polynomial ring over Finite field of 9 elements
}

func ExampleQuotientRing_Quotient() {
	gf9, _ := finitefield.Define(9)
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

	ring, _ = ring.Quotient(id)

	fmt.Println(ring)
	// Output:
	// Quotient ring of bivariate polynomials over Finite field of 9 elements modulo <Y^3 + 2X^4 + Y, X^9 + 2X>
}

func ExampleQuotientRing_SetVarNames() {
	gf4, _ := finitefield.Define(4)
	ring := bivariate.DefRing(gf4, bivariate.Lex(false))

	f, _ := ring.PolynomialFromString("Y^4 + (a+1)Y^2X^3 + a")

	// Change the variable names
	err := ring.SetVarNames([2]string{"S", "T"})
	if err != nil {
		fmt.Printf("Could not set variable names: %q", err)
	}

	g, _ := ring.PolynomialFromString("ST^3 + T^2 + S")

	// Both f and g are affected by the change
	fmt.Printf("f(S, T) = %v\ng(S, T) = %v", f, g)
	// Output:
	// f(S, T) = T^4 + (a + 1)S^3T^2 + a
	// g(S, T) = ST^3 + T^2 + S
}
