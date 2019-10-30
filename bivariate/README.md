[![Go Report Card](https://goreportcard.com/badge/github.com/ReneBoedker/algobra)](https://goreportcard.com/report/github.com/ReneBoedker/algobra)
![coverage-badge](https://img.shields.io/badge/coverage-90.1%25-brightgreen?cacheSeconds=86400&style=flat)
[![GoDoc](https://godoc.org/github.com/ReneBoedker/algobra/bivariate?status.svg)](https://godoc.org/github.com/ReneBoedker/algobra/bivariate)
# Algobra: Bivariate Polynomials
This package implements bivariate polynomials over prime fields.

## Basic usage
To perform computations on bivariate polynomials, first define the finite field and the polynomial ring. Polynomials can then be constructed in several different ways.
```go
field, _ := finitefield.Define(7)	// Ignoring errors in example
ring := bivariate.DefRing(field, bivariate.Lex(true))	// Use lex ordering with X>Y

// Define f = X^2Y + 6Y 
f := ring.PolynomialFromUnsigned(map[[2]uint]uint{
	{2,1}:	1,
	{0,3}:	6,
})

// The same polynomial can be defined from 
g := ring.PolynomialFromSigned(map[[2]uint]int{
	{2,1}:	1,
	{0,3}:	-1,
})

fmt.Println(f.Equal(g))	// Prints 'true'
```

In addition to the polynomial definition from maps as above, it is also possible to define polynomials from strings in a natural way by using `PolynomialFromString`. When doing so, each monomial can contain at most one of each variable. The order of the variables does not matter, and capitalization is ignored. Using `*` to indicate multiplication is optional. In addition, the parser supports _Singular-style_ exponents, meaning that `5X2Y3` is interpreted as `5X^2Y^3`.

### Ideals
The package provides support for computations modulo an ideal.

``` go
// Let ring be defined as above
id, _ := ring.NewIdeal(
	ring.PolynomialFromString("X^49-X"),
	ring.PolynomialFromString("Y^7-X^8+Y"),
)
qRing, _ := ring.Quotient(id)
```
Once the quotient ring has been defined, polynomials are defined as before. For instance, `h := qRing.PolynomialFromString("X^50")` sets `h` to `X^2` since the polynomial is automatically reduced modulo the ideal.

Internally, this is achieved by transforming the ideal such that its generators form a reduced Gröbner basis. Hence, calling `Generators` at a later point will not necessarily return the polynomials that were used to define the ideal.

### Monomial orderings
The following monomial orderings are defined by default.
* Lexicographical
* Degree lexicographical
* Degree reverse lexicographical
* Weighted degree lexicographical
* Weighted degree reverse lexicographical

Additional orderings can be defined by writing a function with signature `func(deg1, deg2 [2]uint) int`. For more information, see the documentation for the `Order` type.

### Error handling
In order to allow method chaining for arithmetic operations &ndash; such as `f.Plus(g).Mult(h.Inv())` &ndash; the methods themselves do not return errors. Instead, potential errors are tied to the resulting polynomial, and the error can be retrieved with the `Err`-method.

