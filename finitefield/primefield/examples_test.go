package primefield_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/finitefield/primefield"
)

// Set up the finitefield of 3 elements for examples where the cardinality does
// not matter
func getGf3() *primefield.Field {
	out, _ := primefield.Define(3)
	return out
}

var gf3 *primefield.Field = getGf3()

func ExampleDefine() {
	field, _ := primefield.Define(7)
	fmt.Println(field)

	_, err := primefield.Define(4)
	fmt.Printf("Error: %v", err)
	// Output:
	// Finite field of 7 elements
	// Error: Defining prime field: 4 is not a prime
}

func ExampleElement_Err() {
	a := gf3.Zero().Inv()
	fmt.Println(a.Err())
	// Output:
	// Inverting element: Cannot invert zero element
}

func ExampleField_Element() {
	field, _ := primefield.Define(13)

	a, _ := field.Element(14)
	b, _ := field.Element(-5)
	c, _ := field.Element("6")

	fmt.Printf("%v, %v, %v", a, b, c)
	// Output:
	// 1, 8, 6
}

func ExampleField_Elements() {
	field, _ := primefield.Define(7)
	for _, v := range field.Elements() {
		fmt.Println(v)
	}
	// Unordered output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
}
