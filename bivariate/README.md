# Algobra: Bivariate polynomials
This package implements bivariate polynomials over prime fields.

## Basic usage
To perform computations on bivariate polynomials, first define the finite field and the polynomial ring. Polynomials can then be constructed in several different ways.
```go
field, _ := primefield.Define(7)	// Ignoring errors in example
ring := bivariate.DefRing(field, bivariate.Lex(true))	// Use lex ordering with X>Y

// Define f = X^2Y + 6Y 
f := ring.New(map[[2]uint]uint{
	{2,1}:	1,
	{0,3}:	6,
})

// The same polynomial can be defined from 
g := ring.New(map[[2]uint]int{
	{2,1}:	1,
	{0,3}:	-1,
})

fmt.Println(f.Equal(g))	// Prints 'true'
```

In addition to the polynomial definition from maps as above, it is also possible to define polynomials from strings in a natural way by using `NewFromString`. When doing so, the variable names must be 'X' and 'Y' -- but not necessarily capitalized -- and each monomial can contain at most one of each variable. The order of the variables does not matter. Using `*` to indicate multiplication is optional. In addition, the parser supports _Singular-style_ exponents, meaning that `5X2Y3` is interpreted as `5X^2Y^3`.

## Ideals

### Error handling
In order to allow method chaining for arithmetic operations -- such as `f.Plus(g).Mult(h.Inv())` -- the methods themselves do not return errors. Instead, potential errors are tied to the resulting polynomial, and the error can be retrieved with the `Err`-method.

