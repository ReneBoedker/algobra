package auxmath_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/auxmath"
)

func ExampleCombinIter() {
	for ci := auxmath.NewCombinIter(6, 3); ci.Active(); ci.Next() {
		fmt.Println(ci.Current())
	}
	// Output:
	// [0 1 2]
	// [0 1 3]
	// [0 1 4]
	// [0 1 5]
	// [0 2 3]
	// [0 2 4]
	// [0 2 5]
	// [0 3 4]
	// [0 3 5]
	// [0 4 5]
	// [1 2 3]
	// [1 2 4]
	// [1 2 5]
	// [1 3 4]
	// [1 3 5]
	// [1 4 5]
	// [2 3 4]
	// [2 3 5]
	// [2 4 5]
	// [3 4 5]
}
