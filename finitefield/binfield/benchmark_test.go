package binfield_test

import (
	"testing"

	"github.com/ReneBoedker/algobra/finitefield/binfield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

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

func BenchmarkProd256(b *testing.B) {
	field, _ := binfield.Define(256)
	benchProd(field, b)
}

func BenchmarkSum256(b *testing.B) {
	field, _ := binfield.Define(256)
	benchSum(field, b)
}

func BenchmarkInv256(b *testing.B) {
	field, _ := binfield.Define(256)
	benchInv(field, b)
}

func BenchmarkProd4096(b *testing.B) {
	field, _ := binfield.Define(4096)
	benchProd(field, b)
}

func BenchmarkSum4096(b *testing.B) {
	field, _ := binfield.Define(4096)
	benchSum(field, b)
}

func BenchmarkInv4096(b *testing.B) {
	field, _ := binfield.Define(4096)
	benchInv(field, b)
}
