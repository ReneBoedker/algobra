# Algobra: Prime Fields
This package implements arithmetic in finite fields of prime characteristic.

## Using table lookups
By default, the package will do the computation each time two elements are added or multiplied. If requested by the user, the package will instead precompute the addition and/or multiplication tables for a field, and use a table lookup for the corresponding field operations.
```go
f,_:=Define(7)						// Ignoring errors
err:=f.ComputeTables(true,false)	// Precompute the addition table
if err!=nil {
   // Table exceeds maximal memory usage
}
```
