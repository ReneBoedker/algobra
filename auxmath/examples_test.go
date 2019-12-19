package auxmath_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/auxmath"
)

func ExampleFactorize() {
	var n uint = 2 * 2 * 2 * 3 * 3 * 13 * 13 * 13
	fact, exps := auxmath.Factorize(n)
	fmt.Println(fact, exps)
	// Output:
	// [2 3 13] [3 2 3]
}

func ExampleFactorizePrimePower() {
	q, err := auxmath.Pow(7, 12)
	if err != nil {
		fmt.Println(err)
		return
	}

	p, n, _ := auxmath.FactorizePrimePower(q)
	fmt.Printf("%d = %d^%d\n", q, p, n)

	_, _, err = auxmath.FactorizePrimePower(2 * 5)
	fmt.Println(err)
	// Output:
	// 13841287201 = 7^12
	// Factorizing prime power: 10 does not seem to be a prime power.
}

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

func ExamplePow_overflow() {
	_, err := auxmath.Pow(5, 28)
	fmt.Println(err)
	// Output:
	// Computing power of unsigned integer: 5^28 is likely to overflow uint
}

func ExampleCombinIter_Next() {
	ci := auxmath.NewCombinIter(3, 2)
	fmt.Println(ci.Current())

	// Iterate to the last combination
	for ci.Active() {
		ci.Next()
	}
	fmt.Println(ci.Current())

	// Attempt incrementing, but already at end
	ci.Next()
	fmt.Println(ci.Current())
	// Output:
	// [0 1]
	// [1 2]
	// [1 2]
}
