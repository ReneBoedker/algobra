// Package extfield implements finite fields as extensions of prime fields.
//
// If you need finite fields of extension degree one -- that is, prime fields --
// algobra/primefield provides a more efficient implementation. If both prime
// fields and extension fields are needed, consider using algobra/finitefields
// instead.
//
// # Arithmetic operations
//
// The Element objects have methods Add, Sub, Mult, and Inv for the
// basic field operations. Of these four, only Inv allocates a new object. The
// other methods store the result in the receiving element object. For instance,
// a.Add(b) would evaluate the sum a+b, and then set a to this value. If
// the result is to be stored in a new object, the package provides the methods
// Plus, Minus, and Times, which evaluate the arithmetic operation and
// returns the result in a new object.
//
// Additional functions such as Neg and Pow are also defined.
//
// # Equality testing
//
// To test if two elements a and b are equal, a == b will not work since this
// compares pointers rather than the underlying data. Instead, use a.Equal(b).
// In addition, the expressions a.IsZero(), a.IsNonzero(), and a.IsOne() provide
// shorthands for common comparisons.
//
// # Error handling
//
// In order to allow method chaining for arithmetic operations -- such as
// a.Add(b).Mult(c.Inv()) -- the methods themselves do not return errors.
// Instead, potential errors are tied to the resulting field element, and the
// error can be retrieved with the Err-method. For instance, you might do
// something like this:
//
//	a:=field.Element(0).Inv()
//	if a.Err()!=nil {
//		// Handle error
//	}
package extfield
