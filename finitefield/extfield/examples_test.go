package extfield_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/finitefield/extfield"
)

// Set up the finitefield of 4 elements for examples where the cardinality does
// not matter
func getGf4() *extfield.Field {
	out, _ := extfield.Define(4)
	return out
}

var gf4 *extfield.Field = getGf4()

func ExampleDefine() {
	field, _ := extfield.Define(9)
	fmt.Println(field)

	_, err := extfield.Define(10)
	fmt.Printf("Error: %v", err)
	// Output:
	// Finite field of 9 elements
	// Error: Defining prime field: Factorizing prime power: 10 does not seem to be a prime power.
}

func ExampleElement_Err() {
	a := gf4.Zero().Inv()
	if a.Err() != nil {
		fmt.Println(a.Err())
	}
	// Output:
	// Inverting element: Cannot invert zero element
}

func ExampleElement_AsSlice() {
	gf27, _ := extfield.Define(27)

	a := gf27.ElementFromUnsignedSlice([]uint{1, 0, 2})
	fmt.Println(a)

	// a is of interface type ff.Element. Type assertion is necessary to access
	// the methods of the extfield.Element type.
	aAssert, _ := a.(*extfield.Element)
	fmt.Println(aAssert.AsSlice())
	// Output:
	// 2a^2 + 1
	// [1 0 2]
}

func ExampleElement_NTerms() {
	a := gf4.ElementFromUnsignedSlice([]uint{1, 1}) // a + 1
	b := gf4.One()

	fmt.Printf("%d, %d", a.NTerms(), b.NTerms())
	// Output:
	// 2, 1
}

func ExampleField_Element() {
	field, _ := extfield.Define(25)

	a, _ := field.Element(6)
	b, _ := field.Element([]int{3, -1})
	c, _ := field.Element("2a+4")

	fmt.Printf("%v, %v, %v", a, b, c)
	// Output:
	// 1, 4a + 3, 2a + 4
}

func ExampleField_Elements() {
	field, _ := extfield.Define(9)
	for _, v := range field.Elements() {
		fmt.Println(v)
	}
	// Unordered output:
	// 0
	// 1
	// 2
	// a
	// a + 1
	// a + 2
	// 2a
	// 2a + 1
	// 2a + 2
}
