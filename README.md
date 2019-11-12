[![Go Report Card](https://goreportcard.com/badge/github.com/ReneBoedker/algobra)](https://goreportcard.com/report/github.com/ReneBoedker/algobra)
![coverage-badge](https://img.shields.io/badge/coverage-92.2%25-brightgreen?cacheSeconds=86400&style=flat)
[![GoDoc](https://godoc.org/github.com/ReneBoedker/algobra?status.svg)](https://godoc.org/github.com/ReneBoedker/algobra)

# Algobra
Algobra is a collection of packages that implement finite field arithmetic as well as univariate and bivariate polynomials over finite fields.

## Installation
Algobra is a package for the [Go](https://golang.org/) programming language. Therefore, the following assumes that you have already installed Go on your machine. Otherwise, please refer to the [official instructions](https://golang.org/doc/install).

To install the full set of Algobra-packages, you can run
```sh
go get github.com/ReneBoedker/algobra/...
```
Note that this top-level package does not provide any functionality in itself. It simply contains documentation for the package as a whole.

## Example usage
For examples on how to use the package, please refer to the [documentation](https://godoc.org/github.com/ReneBoedker/algobra).

## References
* Gathen, Joachim von zur &amp; Gerhard, Jürgen: _Modern Computer Algebra_ (2013), 3rd edition. Cambridge University Press. ISBN 978-1-107-03903-2.
* Lauritzen, Niels: _Concrete Abstract Algebra_ (2003). Cambridge University Press. ISBN 978-0-521-53410-9
* Lübeck, Frank: [_Conway polynomials for finite fields_](http://www.math.rwth-aachen.de/~Frank.Luebeck/data/ConwayPol/index.html?LANG=en)
* Stichtenoth, Henning: _Algebraic Function Fields and Codes_ (2009), 2nd edition. Springer. ISBN 978-3-540-76877-7.
