package algobra_test

import (
	"fmt"

	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/univariate"
	"math/rand"
)

// This example contains a simple illustration of Shamir's secret sharing scheme.
func Example_secretSharing() {
	rand.Seed(314159265) // Use a fixed seed for the example

	// We want 5 participants and any 3 should be able to reconstruct
	n := 5
	recon := 3

	// Define the finite field of 9 elements (a smaller field could have been used)
	gf9, err := finitefield.Define(9)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := univariate.DefRing(gf9)

	// Let the secret be 2a+1
	s, err := gf9.ElementFromString("2a+1")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Define random polynomial f of degree less than recon conditioned on f(0)=s
	f := r.Polynomial([]ff.Element{s})
	for i := 1; i < recon; i++ {
		f.SetCoef(i, gf9.RandElement())
	}
	fmt.Printf("The sharing polynomal is f(X) = %v\n\n", f)

	// Evaluate f in n distinct points, giving the shares
	points := make([]ff.Element, 0, n)
	shares := make([]ff.Element, 0, n)
	for _, p := range gf9.Elements() {
		if p.IsZero() {
			// Avoid zero since this evaluates to the secret
			continue
		}

		points = append(points, p)
		shares = append(shares, f.Eval(p))
		if len(shares) == n {
			break
		}
	}
	fmt.Printf("Points: %v\nShares: %v\n\n", points, shares)

	// Assume that the first three shares are known
	knownPoints := points[:3]
	knownShares := shares[:3]

	// Use interpolation to find the secret
	g, err := r.Interpolate(knownPoints, knownShares)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The reconstructed polynomial is g(X) = %v\n", g)
	fmt.Printf("g(0) = %v\n", g.Eval(gf9.Zero()))
	// Output:
	// The sharing polynomal is f(X) = (a + 2)X^2 + 2aX + (2a + 1)
	//
	// Points: [1 a a + 1 2a + 1 2]
	// Shares: [2a 2a 2a + 1 a + 2 a]
	//
	// The reconstructed polynomial is g(X) = (a + 2)X^2 + 2aX + (2a + 1)
	// g(0) = 2a + 1
}
