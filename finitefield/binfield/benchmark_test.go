package binfield_test

import (
	"testing"

	"github.com/ReneBoedker/algobra/finitefield/binfield"
	"github.com/ReneBoedker/algobra/finitefield/extfield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

func bench(f ff.Field, b *testing.B) {
	res := f.One()
	l := make([]ff.Element, b.N, b.N)
	for i := 0; i < b.N; i++ {
		tmp := f.RandElement()
		for tmp.IsZero() {
			tmp = f.RandElement()
		}
		l[i] = tmp
	}

	b.ResetTimer() // Ignore the cost of generating (and storing) random elements
	for _, v := range l {
		res.Mult(v)
	}
}

func BenchmarkGF64Binary(b *testing.B) {
	field, _ := binfield.Define(64)
	bench(field, b)
}

func BenchmarkGF64Extension(b *testing.B) {
	field, _ := extfield.Define(64)
	bench(field, b)
}

func BenchmarkGF256Binary(b *testing.B) {
	field, _ := binfield.Define(256)
	bench(field, b)
}

func BenchmarkGF256Extension(b *testing.B) {
	field, _ := extfield.Define(256)
	bench(field, b)
}

func BenchmarkGFHugeExtension(b *testing.B) {
	field, _ := extfield.Define(282429536481)
	bench(field, b)
}
