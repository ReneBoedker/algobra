![coverage-badge](https://img.shields.io/badge/coverage-96.3%25-brightgreen?cacheSeconds=86400&style=flat)
# Algobra: Prime Fields
This package implements arithmetic in finite fields of prime cardinality.

## Basic usage
To use the package, simply define the field -- or fields -- that you want to work over. Then construct the needed elements.
```go
// import "github.com/ReneBoedker/algobra/primefield"
ff,err:=primefield.Define(7)
if err!=nil {
    // Define returns an error if the characteristic is not a prime (or too large)
}

a:=ff.Element(3)
b:=ff.Element(6)
c:=a.Plus(b)    // c = 2
```

### Arithmetic operations
The Element objects have methods `Add`, `Sub`, `Mult`, and `Inv` for the basic field operations. Of these four, only `Inv` allocates a new object. The other methods store the result in the receiving element object. For instance, `a.Add(b)` would evaluate the sum `a+b`, and then set `a` to this value. If the result is to be stored in a new object, the package provides the methods `Plus`, `Minus`, and `Times`, which evaluate the arithmetic operation and returns the result in a new object.

Additional functions such as `Neg` and `Pow` are also defined. For details, please refer to the documentation.

### Equality testing
To test if two elements `a` and `b` are equal, `a == b` will not work since this compares pointers rather than the underlying data. Instead, use `a.Equal(b)`. In addition, the expressions `a.IsZero()`, `a.IsNonzero()`, and `a.IsOne()` provide shorthands for common comparisons.

### Error handling
In order to allow method chaining for arithmetic operations -- such as `a.Add(b).Mult(c.Inv())` -- the methods themselves do not return errors. Instead, potential errors are tied to the resulting field element, and the error can be retrieved with the `Err`-method. For instance, you might do something like this:
``` go
a:=field.Element(0).Inv()
if a.Err()!=nil {
    // Handle error
}
```

### Using table lookups
By default, the package will perform the computation each time two elements are added or multiplied. If requested by the user, the package will instead precompute the addition and/or multiplication tables for a field, and use a table lookup for the corresponding field operations.
```go
err:=ff.ComputeTables(true,false)   // Precompute the addition table
if err!=nil {
    // Table exceeds maximal memory usage
}
```
