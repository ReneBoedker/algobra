// Package bivariate implements bivariate polynomials over finite fields.
//
// # Basic usage
//
// To perform computations on bivariate polynomials, first define the finite
// field and the polynomial ring. Polynomials can then be constructed in several
// different ways.
//
//	field, _ := finitefield.Define(7)    // Ignoring errors in example
//	ring := bivariate.DefRing(field, bivariate.Lex(true))   // Use lex ordering with X>Y
//
//	// Define f = X^2Y + 6Y
//	f := ring.PolynomialFromUnsigned(map[[2]uint]uint{
//	    {2,1}:  1,
//	    {0,3}:  6,
//	})
//
//	// The same polynomial can be defined from
//	g := ring.PolynomialFromSigned(map[[2]uint]int{
//	    {2,1}:  1,
//	    {0,3}:  -1,
//	})
//
//	fmt.Println(f.Equal(g)) // Prints 'true'
//
// In addition to the polynomial definition from maps as above, it is also
// possible to define polynomials from strings in a natural way by using
// PolynomialFromString. When doing so, each monomial can contain at most one of
// each variable. The order of the variables does not matter, and capitalization
// is ignored. Using * to indicate multiplication is optional. In addition, the
// parser supports _Singular-style_ exponents, meaning that '5X2Y3' is
// interpreted as '5X^2Y^3'.
//
// By default, the variable names 'X' and 'Y' are used, but this can be changed
// via the method SetVarNames. The current variable names can be obtained from
// the VarNames method.
//
// # Ideals
//
// The package provides support for computations modulo an ideal.
//
//	// Let ring be defined as above
//	id, _ := ring.NewIdeal(
//	    ring.PolynomialFromString("X^49-X"),
//	    ring.PolynomialFromString("Y^7-X^8+Y"),
//	)
//	qRing, _ := ring.Quotient(id)
//
// Once the quotient ring has been defined, polynomials are defined as before.
// For instance, h := qRing.PolynomialFromString("X^50") sets h to X^2 since the
// polynomial is automatically reduced modulo the ideal.
//
// Internally, this is achieved by transforming the ideal such that its
// generators form a reduced Gröbner basis. Hence, calling Generators at a
// later point will not necessarily return the polynomials that were used to
// define the ideal.
//
// # Monomial orderings
//
// The following monomial orderings are defined by default.
//   - Lexicographical
//   - Degree lexicographical
//   - Degree reverse lexicographical
//   - Weighted degree lexicographical
//   - Weighted degree reverse lexicographical
//
// Additional orderings can be defined by writing a function with signature
// func(deg1, deg2 [2]uint) int. For more information, see the documentation
// for the Order type.
//
// # Error handling
//
// In order to allow method chaining for arithmetic operations -- such as
// f.Plus(g).Mult(h.Inv()) -- the methods themselves do not return errors.
// Instead, potential errors are tied to the resulting polynomial, and the error
// can be retrieved with the Err-method.
package bivariate
