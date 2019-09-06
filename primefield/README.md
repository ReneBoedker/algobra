# Algobra: Prime Fields
This package implements arithmetic in finite fields of prime characteristic.

## Basic usage
To use the package, simply define the field -- or fields -- that you want to work over. Then construct the needed elements.
```go
// import "github.com/ReneBoedker/algobra/primefield"
ff,err:=Define(7)
if err!=nil {
    // Define returns an error if the characteristic is not a prime (or too large)
}

a:=ff.Element(3)
b:=ff.Element(6)
c:=a.Plus(b)    // c = 2
```

### Error handling
In order to allow method chaining for arithmetic operations -- such as `a.Plus(b).Mult(c.Inv())` -- the methods themselves do not return errors. Instead, potential errors are tied to the resulting field element, and the error can be retrieved with the `Err`-method. For instance, you might do something like this:
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
