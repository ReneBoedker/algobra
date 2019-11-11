[![Go Report Card](https://goreportcard.com/badge/github.com/ReneBoedker/algobra)](https://goreportcard.com/report/github.com/ReneBoedker/algobra)
![coverage-badge](https://img.shields.io/badge/coverage-87.8%25-green?cacheSeconds=86400&style=flat)
[![GoDoc](https://godoc.org/github.com/ReneBoedker/algobra/finitefield?status.svg)](https://godoc.org/github.com/ReneBoedker/algobra/finitefield)
# Algobra: Finite Fields
This package and its subpackages implement arithmetic in finite fields.

In itself, this package only provides a convenient method for defining finite fields. Based on the cardinality, it will automatically choose the appropriate underlying implementation. The return value is the interface `ff.Field`.

## Basic usage
To use the package, simply define the field &ndash; or fields &ndash; that you want to work:
```go
gf7, err := finitefield.Define(7)
if err != nil {
    // Define returns an error if the characteristic is not a prime power (or too large)
}
```
Elements in the field can be constructed in several different ways. The two methods `ElementFromSigned` and `ElementFromUnsigned` return the element corresponding to a signed or an unsigned integer, respectively. These are guaranteed to never return an error. The same holds true for the convenient `Zero` and `One` methods, which return the additive and multiplicative identities. `ElementFromString` parses a string and returns the corresponding element if the parsing succeeds &ndash; otherwise an error is returned.

Each field also has a general method `Element` which will call the appropriate constructor based on the input type. The accepted types depends on the type of field:

* **Prime fields:** uint, int, string
* **Extension fields:** uint, []uint, int, []int, string

If any other type is given as input, the method returns an error.

For extension fields, `Elements` provides access to constructors that are otherwise not callable from this package. If desired, the user can import `finitefield/extfield` directly and define the fields from there. Then all implemented constructors are available.

The following examples illustrate how elements are constructed.
```go
a := gf7.ElementFromUnsigned(3)	// a = 3
b := gf7.ElementFromSigned(-2)	// b = 5

gf9, _ := finitefield.Define(9)
c, err := gf9.Element([]int{1,2})   // c = 1 + α
d := gf9.Zero()                     // d = 0
```

### Arithmetic operations
The elements have methods `Add`, `Sub`, `Mult`, and `Inv` for the basic field operations. Of these four, only `Inv` allocates a new object. The other methods store the result in the receiving element object. For instance, `a.Add(b)` would evaluate the sum `a + b`, and then set `a` to this value. If the result is to be stored in a new object, the package provides the methods `Plus`, `Minus`, and `Times`, which evaluate the arithmetic operation and returns the result in a new object. As a mnemonic, methods named after the operation are destructive, whereas those named after the mathematical sign are non-destructive.

```go
e := a.Plus(b)  // Returns a new element e with value a + b = 1
e.Add(b)        // Alters e to have value e + b = 6

e := a.Times(b).Minus(gf7.One())
// e now has value (a * b) - 1 = 0
// The values of a and b are unchanged
```

Additional functions such as `Neg`, `SetNeg`, and `Pow` are also defined. For details, please refer to the documentation.

### Equality testing
To test if two elements `a` and `b` are equal, `a == b` will not work since this compares pointers rather than the underlying data. Instead, use `a.Equal(b)`. In addition, the expressions `a.IsZero()`, `a.IsNonzero()`, and `a.IsOne()` provide shorthands for common comparisons.

### Error handling
In order to allow method chaining for arithmetic operations &ndash; such as `a.Add(b).Mult(c.Inv())` &ndash; the methods themselves do not return errors. Instead, potential errors are tied to the resulting field element, and the error can be retrieved with the `Err`-method. For instance, you might do something like this:
``` go
a:=gf9.Element(0).Inv()
if a.Err()!=nil {
    // Handle error
}
```
