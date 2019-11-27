package extfield_test

import (
	"testing"

	"github.com/ReneBoedker/algobra/finitefield/extfield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

func benchSum(f ff.Field, b *testing.B) {
	res := f.Zero()
	l := make([]ff.Element, b.N, b.N)
	for i := 0; i < b.N; i++ {
		l[i] = f.RandElement()
	}

	b.ResetTimer() // Ignore the cost of generating (and storing) random elements
	for _, v := range l {
		res.Add(v)
	}
}

func benchProd(f ff.Field, b *testing.B) {
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

func benchInv(f ff.Field, b *testing.B) {
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
		res = v.Inv()
	}
	if res.Err() != nil {
		b.FailNow()
	}
}

func BenchmarkProd343(b *testing.B) {
	field, _ := extfield.Define(343)
	benchProd(field, b)
}

func BenchmarkSum343(b *testing.B) {
	field, _ := extfield.Define(343)
	benchSum(field, b)
}

func BenchmarkInv343(b *testing.B) {
	field, _ := extfield.Define(343)
	benchInv(field, b)
}

func BenchmarkProd390625(b *testing.B) {
	field, _ := extfield.Define(390625)
	benchProd(field, b)
}

func BenchmarkSum390625(b *testing.B) {
	field, _ := extfield.Define(390625)
	benchSum(field, b)
}

func BenchmarkInv390625(b *testing.B) {
	field, _ := extfield.Define(390625)
	benchInv(field, b)
}
