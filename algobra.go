// Package algobra is a collection of packages for finite field arithmetic.
//
// This package does not provide any functionality itself. Instead, this is
// found in each of the subpackages.
package algobra

import (
	_ "algobra/bivariate"
	_ "algobra/finitefield"
	_ "algobra/univariate"
)
