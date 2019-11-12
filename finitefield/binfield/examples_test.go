package binfield_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/finitefield/binfield"
)

// Set up the finitefield of 4 elements for examples where the cardinality does
// not matter
func getGf4() *binfield.Field {
	out, _ := binfield.Define(4)
	return out
}

var gf4 *binfield.Field = getGf4()

func ExampleDefine() {
	field, _ := binfield.Define(256)
	fmt.Println(field)

	_, err := binfield.Define(49)
	fmt.Printf("Error: %v", err)
	// Output:
	// Finite field of 256 elements
	// Error: Defining binary field: The cardinality of a binary field must be a power of 2
}

func ExampleElement_Err() {
	a := gf4.Zero().Inv()
	if a.Err() != nil {
		fmt.Println(a.Err())
	}
	// Output:
	// Inverting element: Cannot invert zero element
}

func ExampleField_Element() {
	field, _ := binfield.Define(16)

	a, _ := field.Element(uint(7))
	b, _ := field.Element(-2)
	c, _ := field.Element("a + a^2 + 1")

	fmt.Printf("%v, %v, %v", a, b, c)
	// Output:
	// 1, 0, a^2 + a + 1
}

func ExampleField_ElementFromBits() {
	field, _ := binfield.Define(16)

	a := field.ElementFromBits(0b1101)
	b := field.ElementFromBits(0b11)

	fmt.Printf("%v, %v", a, b)
	// Output:
	// a^3 + a^2 + 1, a + 1
}

func ExampleField_Elements() {
	field, _ := binfield.Define(8)
	for _, v := range field.Elements() {
		fmt.Println(v)
	}
	// Unordered output:
	// 0
	// 1
	// a
	// a + 1
	// a^2
	// a^2 + 1
	// a^2 + a
	// a^2 + a + 1
}

func ExampleElement_NTerms() {
	a := gf4.ElementFromBits(0b11) // a + 1
	b := gf4.One()

	fmt.Printf("%d, %d", a.NTerms(), b.NTerms())
	// Output:
	// 2, 1
}
