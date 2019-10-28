// Package algobra is a collection of packages for finite field arithmetic.
//
// This package does not provide any functionality itself. Instead, this is
// found in each of the subpackages.
package algobra

import (
	_ "github.com/ReneBoedker/algobra/bivariate"   // Bivariate polynomials
	_ "github.com/ReneBoedker/algobra/finitefield" // Finite field elements
	_ "github.com/ReneBoedker/algobra/univariate"  // Univariate polynomials
)
