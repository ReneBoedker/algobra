[![Go Report Card](https://goreportcard.com/badge/github.com/ReneBoedker/algobra)](https://goreportcard.com/report/github.com/ReneBoedker/algobra)
![coverage-badge](https://img.shields.io/badge/coverage-94.4%25-brightgreen?cacheSeconds=86400&style=flat)
[![GoDoc](https://godoc.org/github.com/ReneBoedker/algobra/finitefield/conway?status.svg)](https://godoc.org/github.com/ReneBoedker/algobra/finitefield/conway)
# Algobra: Conway polynomials
This package contains the list of Conway polynomials provided on the homepage of [Frank Lübeck](http://www.math.rwth-aachen.de/~Frank.Luebeck/data/ConwayPol/index.html).

If desired, the package can be used on its own by utilizing `Lookup`, which is the only exported function. The return value is the slice uints representing the coefficients of the Conway polynomial.
