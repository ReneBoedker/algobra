[![Go Report Card](https://goreportcard.com/badge/github.com/ReneBoedker/algobra)](https://goreportcard.com/report/github.com/ReneBoedker/algobra)
![coverage-badge](https://img.shields.io/badge/coverage-93.0%25-brightgreen?cacheSeconds=86400&style=flat)
[![GoDoc](https://godoc.org/github.com/ReneBoedker/algobra/univariate?status.svg)](https://godoc.org/github.com/ReneBoedker/algobra/univariate)
# Algobra: Univariate Polynomials
This package implements univariate polynomials over prime fields.

## Basic usage
To perform computations on univariate polynomials, first define the finite field and the polynomial ring. Polynomials can then be constructed in several different ways.
```go
field, _ := finitefield.Define(7)	// Ignoring errors in example
ring := bivariate.DefRing(field)

// Define f = 3X^5+2X^2+6
f := ring.PolynomialFromUnsigned([]uint{6,0,2,0,0,3})

// The same polynomial can be defined from 
g := ring.PolynomialFromSigned([]int{-1,0,2,0,0,-4})

fmt.Println(f.Equal(g))	// Prints 'true'
```

In addition to the polynomial definition from maps as above, it is also possible to define polynomials from strings in a natural way by using `PolynomialFromString`. When doing so, each monomial can contain at most one of each variable. The order of the variables does not matter, and capitalization is ignored. Using `*` to indicate multiplication is optional. In addition, the parser supports _Singular-style_ exponents, meaning that `5X2` is interpreted as `5X^2`.

### Ideals
The package provides support for computations modulo an ideal.

``` go
// Let ring be defined as above
id, _ := ring.NewIdeal(
	ring.PolynomialFromString("X^49-X"),
)
qRing, _ := ring.Quotient(id)
```
Once the quotient ring has been defined, polynomials are defined as before. For instance, `h := qRing.PolynomialFromString("X^50")` sets `h` to `X^2` since the polynomial is automatically reduced modulo the ideal.

Internally, this is achieved by finding a single element that generates the ideal. Hence, calling `Generators` at a later point will not necessarily return the polynomials that were used to define the ideal. Instead, it will return the greatest common divisor of these polynomials.

### Error handling
In order to allow method chaining for arithmetic operations &ndash; such as `f.Plus(g).Mult(h.Inv())` &ndash; the methods themselves do not return errors. Instead, potential errors are tied to the resulting polynomial, and the error can be retrieved with the `Err`-method.

