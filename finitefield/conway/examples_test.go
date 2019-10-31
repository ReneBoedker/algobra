package conway_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/finitefield/conway"
)

func Example() {
	coefs, err := conway.Lookup(3, 5)
	if err != nil {
		return
	}

	fmt.Println(coefs)
	// Output:
	// [1 2 0 0 0 1]
}
