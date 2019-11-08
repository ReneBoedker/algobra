// Package conway contains a database of Conway polynomials. The list of
// polynomials was compiled by Frank Lübeck:
// http://www.math.rwth-aachen.de/~Frank.Luebeck/data/ConwayPol/index.html?LANG=en
package conway

import (
	"fmt"
	"github.com/ReneBoedker/algobra/errors"
	"regexp"
	"strconv"
	"strings"
)

// The list of conway polynomials is defined as a constant in cpimport.go

// Lookup returns the coefficients for a Conway polynomial for the finite field
// of characteristic char and extension degree extDeg. The element at position i
// is the coefficient of X^i.
//
// If no such polynomial is in the database, an InputValue-error is returned.
func Lookup(char, extDeg uint) (coefs []uint, err error) {
	return lookupInternal(char, extDeg, cpimport)
}

// lookupInternal has an additional input parameter which allows testing
func lookupInternal(char, extDeg uint, conwayList string) (coefs []uint, err error) {
	const op = "Searching for Conway polynomial"

	pattern, err := regexp.Compile(fmt.Sprintf(`\[%d,%d,\[([^]]*)\]\]`, char, extDeg))
	if err != nil {
		return nil, errors.New(
			op, errors.InputValue,
			"Failed to construct search pattern for characteristic %d and "+
				"extension degree %d", char, extDeg,
		)
	}

	match := pattern.FindStringSubmatch(conwayList)
	if match == nil {
		return nil, errors.New(
			op, errors.InputValue,
			"No polynomial was found for characteristic %d and extension degree %d",
			char, extDeg,
		)
	}

	coefStr := strings.Split(match[1], ",")
	if tmp := uint(len(coefStr)); tmp != extDeg+1 {
		return nil, errors.New(
			op, errors.Internal,
			"Polynomial in database has degree %d rather than %d", tmp, extDeg+1,
		)
	}

	coefs = make([]uint, len(coefStr), len(coefStr))
	for i, c := range coefStr {
		tmp, err := strconv.ParseUint(c, 10, 0)
		if err != nil {
			return nil, errors.New(
				op, errors.Internal,
				"Could not parse coefficient %s of database polynomial for "+
					"characteristic %d and extension degree %d", c, char, extDeg,
			)
		}
		coefs[i] = uint(tmp)
	}
	return coefs, nil
}
