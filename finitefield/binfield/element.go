package binfield

import (
	"fmt"
	"math/bits"
	"math/rand"
	"strings"

	"github.com/ReneBoedker/algobra/finitefield/ff"
)

// Element is the implementation of an element in a finite field.
type Element struct {
	field *Field
	val   uint
	err   error
}

// Zero returns the additive identity in f.
func (f *Field) Zero() ff.Element {
	return &Element{
		field: f,
		val:   0,
	}
}

// One returns the multiplicative identity in f.
func (f *Field) One() ff.Element {
	return &Element{
		field: f,
		val:   1,
	}
}

// RandElement returns a pseudo-random element in f.
//
// This function uses the default source from the math/rand package. The seed is
// set automatically when loading the primefield package, but a new seed can be
// set by calling rand.Seed().
//
// The pseudo-random generator used is not cryptographically safe.
func (f *Field) RandElement() ff.Element {
	prg := func() uint {
		return uint(rand.Uint64())
	}
	if bits.UintSize == 32 {
		prg = func() uint {
			return uint(rand.Uint32())
		}
	}

	return &Element{
		field: f,
		val:   prg() % (2 << f.extDeg), // Discard anything but the first extDeg+1 bits
	}
}

// ElementFromBits defines a new element over f with value specified by the
// bitstring val.
//
// The returned element will automatically be reduced modulo the characteristic.
func (f *Field) ElementFromBits(val uint) ff.Element {
	a := &Element{
		field: f,
		val:   val,
	}
	a.reduce()
	return a
}

// ElementFromUnsigned defines a new element over f with value specified by val.
//
// The returned element will automatically be reduced modulo the characteristic.
// That is, the additive identity is returned if val is even, and the
// multiplicative identity is returned if val is odd.
func (f *Field) ElementFromUnsigned(val uint) ff.Element {
	return &Element{
		field: f,
		val:   val % 2,
	}
}

// Copy returns a copy of a.
func (a *Element) Copy() ff.Element {
	return &Element{
		field: a.field,
		val:   a.val,
		err:   a.err,
	}
}

// Err returns the error status of a.
func (a *Element) Err() error {
	return a.err
}

// SetUnsigned sets the value of a to the element corresponding to val. It then
// returns a.
//
// The value is automatically reduced modulo the characteristic.
func (a *Element) SetUnsigned(val uint) ff.Element {
	a.val = val % 2
	return a
}

// Equal tests equality of elements a and b.
func (a *Element) Equal(b ff.Element) bool {
	bb, ok := b.(*Element)
	if !ok {
		return false
	}

	if a.field == bb.field && a.val == bb.val {
		return true
	}
	return false
}

// IsZero returns a boolean describing whether a is the additive identity.
func (a *Element) IsZero() bool {
	return a.val == 0
}

// IsNonzero returns a boolean describing whether a is a non-zero element.
func (a *Element) IsNonzero() bool {
	return a.val != 0
}

// IsOne returns a boolean describing whether a is the multiplicative identity.
func (a *Element) IsOne() bool {
	return a.val == 1
}

// String returns the string representation of a.
func (a *Element) String() string {
	if a.IsZero() {
		return "0"
	}

	var b strings.Builder
	nPlus := a.NTerms() - 1
	for term, d := uint(1)<<a.field.extDeg, a.field.extDeg; term > 0; term, d = term>>1, d-1 {
		if term&a.val == 0 {
			continue
		}

		if d == 0 {
			fmt.Fprint(&b, 1)
			break
		}

		fmt.Fprint(&b, a.field.varName)
		if d > 1 {
			fmt.Fprintf(&b, "^%d", d)
		}

		// Print a plus if necessary
		if nPlus > 0 {
			fmt.Fprint(&b, " + ")
			nPlus--
		}
	}

	return b.String()
}

// NTerms returns the number of terms in the representation of a.
func (a *Element) NTerms() uint {
	return uint(bits.OnesCount(a.val))
}
