package primefield_test

import (
	"testing"

	"github.com/ReneBoedker/algobra/finitefield/ff"
	"github.com/ReneBoedker/algobra/finitefield/primefield"
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

func BenchmarkProd101(b *testing.B) {
	field, _ := primefield.Define(101)
	benchProd(field, b)
}

func BenchmarkSum101(b *testing.B) {
	field, _ := primefield.Define(101)
	benchSum(field, b)
}

func BenchmarkInv101(b *testing.B) {
	field, _ := primefield.Define(101)
	benchInv(field, b)
}

func BenchmarkProd10007(b *testing.B) {
	field, _ := primefield.Define(10007)
	benchProd(field, b)
}

func BenchmarkSum10007(b *testing.B) {
	field, _ := primefield.Define(10007)
	benchSum(field, b)
}

func BenchmarkInv10007(b *testing.B) {
	field, _ := primefield.Define(10007)
	benchInv(field, b)
}
