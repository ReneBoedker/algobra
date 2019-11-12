[![Go Report Card](https://goreportcard.com/badge/github.com/ReneBoedker/algobra)](https://goreportcard.com/report/github.com/ReneBoedker/algobra)
![coverage-badge](https://img.shields.io/badge/coverage-84.4%25-green?cacheSeconds=86400&style=flat)
[![GoDoc](https://godoc.org/github.com/ReneBoedker/algobra/finitefield/binfield?status.svg)](https://godoc.org/github.com/ReneBoedker/algobra/finitefield/binfield)
# Algobra: Binary Fields
This package implements arithmetic in finite fields of characteristic two. These can also be obtained from the general implementation in [algobra/extfield](https://github.com/ReneBoedker/algobra/tree/master/extfield), but the implementation in the binfield-package is more efficient.

## Basic usage
### Arithmetic operations
The Element objects have methods `Add`, `Sub`, `Mult`, and `Inv` for the basic field operations. Of these four, only `Inv` allocates a new object. The other methods store the result in the receiving element object. For instance, `a.Add(b)` would evaluate the sum `a+b`, and then set `a` to this value. If the result is to be stored in a new object, the package provides the methods `Plus`, `Minus`, and `Times`, which evaluate the arithmetic operation and returns the result in a new object.

Additional functions such as `Neg` and `Pow` are also defined. For details, please refer to the documentation.

### Equality testing
To test if two elements `a` and `b` are equal, `a == b` will not work since this compares pointers rather than the underlying data. Instead, use `a.Equal(b)`. In addition, the expressions `a.IsZero()`, `a.IsNonzero()`, and `a.IsOne()` provide shorthands for common comparisons.

### Error handling
In order to allow method chaining for arithmetic operations &ndash; such as `a.Add(b).Mult(c.Inv())` &ndash; the methods themselves do not return errors. Instead, potential errors are tied to the resulting field element, and the error can be retrieved with the `Err`-method. For instance, you might do something like this:
``` go
a:=field.Element(0).Inv()
if a.Err()!=nil {
    // Handle error
}
```

## References
* Frank Lübeck: [Conway polynomials for finite fields](http://www.math.rwth-aachen.de/~Frank.Luebeck/data/ConwayPol/index.html)
