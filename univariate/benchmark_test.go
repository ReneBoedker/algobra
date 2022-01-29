package univariate_test

import (
	"testing"

	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/univariate"
)

func benchString(field ff.Field, b *testing.B) {
	ring := univariate.DefRing(field)

	var str string
	for rep := 0; rep < b.N; rep++ {
		// Create random polynomial with up to 50 different terms
		nDegs := (uint(prg.Uint32()) % 50) + 1
		coefs := make([]ff.Element, nDegs, nDegs)
		for i := uint(0); i < nDegs; i++ {
			coefs[i] = field.RandElement()
		}
		f := ring.Polynomial(coefs)

		str = f.String()
	}
	if str == "" {
		b.Fail()
	}
}

func BenchmarkStringBinfield(b *testing.B) {
	field, _ := finitefield.Define(256)
	benchString(field, b)
}

func BenchmarkStringPrime(b *testing.B) {
	field, _ := finitefield.Define(113)
	benchString(field, b)
}

func BenchmarkStringExtension(b *testing.B) {
	field, _ := finitefield.Define(49)
	benchString(field, b)
}

func benchEval(field ff.Field, b *testing.B) {
	ring := univariate.DefRing(field)

	for rep := 0; rep < b.N; rep++ {
		for i := 0; i < 100; i++ {
			// Create random polynomial with up to 50 different terms
			nDegs := (uint(prg.Uint32()) % 50) + 1
			coefs := make([]ff.Element, nDegs, nDegs)
			for i := uint(0); i < nDegs; i++ {
				coefs[i] = field.RandElement()
			}
			f := ring.Polynomial(coefs)

			f.Eval(field.RandElement())
		}
	}
}

func BenchmarkEvalBinfield(b *testing.B) {
	field, _ := finitefield.Define(256)
	benchEval(field, b)
}

func BenchmarkEvalPrime(b *testing.B) {
	field, _ := finitefield.Define(113)
	benchEval(field, b)
}
