package extfield

import (
	"algobra/errors"
	"algobra/finitefield/ff"
	"math/bits"
	"math/rand"
	"testing"
	"time"
)

var prg = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func defineField(card uint) *Field {
	f, err := Define(card)
	if err != nil {
		// Testing code is wrong, so panic
		panic(err)
	}
	return f
}

func TestOverflowDetection(t *testing.T) {
	var bigPrime uint
	if bits.UintSize == 32 {
		bigPrime = uint(65537)
	} else {
		bigPrime = uint(4294967311)
	}
	_, err := Define(bigPrime)
	if err == nil {
		t.Errorf(
			"Define succeeded for bit size %d even though p=%d",
			bits.UintSize, bigPrime,
		)
	} else if !errors.Is(errors.InputTooLarge, err) {
		t.Errorf("Define failed, but the error kind was unexpected")
	}
}

func TestNonPrimePowInput(t *testing.T) {
	testCases := []uint{0, 1, 10, 77, 100}
	for _, char := range testCases {
		if _, err := Define(char); err == nil {
			t.Errorf(
				"Defining field with characteristic %d did not return an error",
				char,
			)
		} else if !errors.Is(errors.InputValue, err) {
			t.Errorf(
				"Defining field with characteristic %d returned error, but wrong kind",
				char,
			)
		}
	}
}

func TestEqual(t *testing.T) {
	field := defineField(125)
	if !field.element([]uint{4, 3, 1}).Equal(field.ElementFromSignedSlice([]int{-1, 3, -4})) {
		t.Errorf("Reported 4+3α+α^2 != 4+3α+α^2 in finite field of 125 elements")
	}
	field2 := defineField(25)
	if field.element([]uint{1, 2}).Equal(field2.element([]uint{1, 2})) {
		t.Errorf("Reported equality for elements from different fields")
	}
}

func TestTableMemory(t *testing.T) {
	var bigPrimePow uint
	if bits.UintSize == 32 {
		bigPrimePow = uint(41781923) // 347^3
	} else {
		bigPrimePow = uint(22188041) // 281^3
	}
	f := defineField(bigPrimePow)
	if err := f.ComputeMultTable(); err == nil {
		t.Errorf("No error returned")
	} else if !errors.Is(errors.InputTooLarge, err) {
		t.Errorf("Error returned has wrong kind. Expected errors.InputTooLarge, "+
			"but received error %q", err.Error())
	}
}

func TestConstructors(t *testing.T) {
	field := defineField(125)

	a, _ := field.Element(uint(4))
	b, _ := field.Element([]uint{4})
	c, _ := field.Element(-1)
	d, _ := field.Element([]int{-1})
	e := field.Zero()
	e.SetUnsigned(4)

	elems := []ff.Element{
		a, b, c, d, e,
		field.ElementFromUnsigned(4),
		field.ElementFromSigned(-1),
	}

	for i, a := range elems {
		for j, b := range elems {
			if !a.Equal(b) {
				t.Errorf(
					"Elements %v and %v are not equal (indices %d and %d)",
					a, b, i, j,
				)
			}
		}
	}
}

func TestArithmeticErrors(t *testing.T) {
	fieldA := defineField(8)
	fieldB := defineField(25)

	a := fieldA.element([]uint{0})
	b := fieldB.element([]uint{10})
	// Cannot invert zero
	if e := a.Inv(); e.Err() == nil {
		t.Errorf("Inverting zero did not set error status")
	} else if !errors.Is(errors.InputValue, e.Err()) {
		t.Errorf("Inverting zero set error status, but not InputValue-error")
	}

	// Cannot use elements from different fields
	if e := a.Plus(b); e.Err() == nil {
		t.Errorf("Adding elements from different fields did not set error status")
	}
	if e := a.Times(b); e.Err() == nil {
		t.Errorf("Multiplying elements from different fields did not set error status")
	}
	if e := a.Minus(b); e.Err() == nil {
		t.Errorf("Subtracting elements from different fields did not set error status")
	}

	// Error is passed on to last result
	if e := b.Plus(b.Minus(a.Inv())); e.Err() == nil {
		t.Errorf("Last result in b+b-a^(-1) did not have error status")
	} else if !errors.Is(errors.InputValue, e.Err()) {
		// Inverting gives InputValue-error. This should be the last kind as well
		t.Errorf("Last result did not retain the original error status")
	}
	if e := b.Minus(b).Inv().Times(b); e.Err() == nil {
		t.Errorf("Last result in b-b^(-1)*b did not have error status")
	} else if !errors.Is(errors.InputValue, e.Err()) {
		// Inverting gives InputValue-error. This should be the last kind as well
		t.Errorf("Last result did not retain the original error status")
	}
}

func TestPow(t *testing.T) {
	field := defineField(9)
	for rep := 0; rep < 50; rep++ {
		a0 := uint(prg.Uint32())
		a1 := uint(prg.Uint32())
		a := field.ElementFromUnsignedSlice([]uint{a0, a1})

		n := uint(prg.Uint32()) % 20

		expected := field.One()
		if n > 0 && a.IsZero() {
			expected = field.Zero()
		} else {
			for i := uint(0); i < n; i++ {
				expected.Mult(a)
			}
		}

		if tmp := a.Pow(n); !tmp.Equal(expected) {
			t.Errorf("(%v)^%d = %v, but expected %v", a, n, tmp, expected)
		}
	}
}

func TestBools(t *testing.T) {
	field := defineField(49)
	if field.element([]uint{0}).IsNonzero() {
		t.Errorf("Element(0) element considered non-zero")
	}
	if !field.element([]uint{0}).IsZero() {
		t.Errorf("Element(0) element not considered zero")
	}
	if !field.element([]uint{1}).IsOne() {
		t.Errorf("Element(1) not considered as one")
	}
	if !field.element([]uint{1}).IsNonzero() {
		t.Errorf("Element(1) not considered non-zero")
	}
}

func TestGenerator(t *testing.T) {
	for _, q := range []uint{2, 4, 3, 5, 9, 8, 16, 25} {
		unique := make(map[string]struct{})

		field, err := Define(q)
		if err != nil {
			panic(err)
		}

		g := field.MultGenerator()
		for i, e := uint(0), g.Copy(); i < q-1; i, e = i+1, e.Times(g) {
			if _, ok := unique[e.String()]; ok {
				t.Errorf("Found element %v twice for p=%v (generator = %v)", e, q, g)
			} else {
				unique[e.String()] = struct{}{}
			}
		}
	}
}

func hardcodedTableTest(
	f *Field,
	elems []*Element,
	sumTable, diffTable, prodTable [][]*Element,
	invList []*Element,
	t *testing.T,
) {
	test := func(f *Field) {
		for i := range elems {
			for j := range elems {
				if t1, t2 := elems[i].Plus(elems[j]), sumTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(%d) failed: (%v) + (%v) = %v (Expected %v)",
						f.Card(), elems[i], elems[j], t1, t2)
				}
				if t1, t2 := elems[i].Minus(elems[j]), diffTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(%d) failed: (%v) - (%v) = %v (Expected %v)",
						f.Card(), elems[i].val, elems[j].val, t1, t2)
				}
				if t1, t2 := elems[i].Times(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(%d) failed: (%v) * (%v) = %v (Expected %v)",
						f.Card(), elems[i], elems[j], t1, t2)
				}
			}
		}

		for i := 1; i < len(elems); i++ {
			if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
				t.Errorf("GF(%d) failed: inv(%v) = %v (Expected %v)",
					f.Card(), elems[i], t1, invList[i-1],
				)
			}
		}
	}
	// Without tables
	test(f)
	// With tables
	f.ComputeMultTable()
	test(f)
}

func TestGf4(t *testing.T) {
	field := defineField(4)
	elems := []*Element{
		field.element([]uint{0, 0}),
		field.element([]uint{1, 0}),
		field.element([]uint{0, 1}),
		field.element([]uint{1, 1}),
	}
	sumTable := [][]*Element{
		{elems[0], elems[1], elems[2], elems[3]},
		{elems[1], elems[0], elems[3], elems[2]},
		{elems[2], elems[3], elems[0], elems[1]},
		{elems[3], elems[2], elems[1], elems[0]},
	}
	prodTable := [][]*Element{
		{elems[0], elems[0], elems[0], elems[0]},
		{elems[0], elems[1], elems[2], elems[3]},
		{elems[0], elems[2], elems[3], elems[1]},
		{elems[0], elems[3], elems[1], elems[2]},
	}
	invList := []*Element{elems[1], elems[3], elems[2]}

	hardcodedTableTest(
		field,
		elems,
		sumTable,
		sumTable, // Note that addition and subtraction are equivalent
		prodTable,
		invList,
		t,
	)
}

func TestGf9(t *testing.T) {
	field := defineField(9)
	elems := []*Element{
		field.element([]uint{0, 0}), //0
		field.element([]uint{1, 0}), //1
		field.element([]uint{2, 0}), //2
		field.element([]uint{0, 1}), //3
		field.element([]uint{0, 2}), //4
		field.element([]uint{1, 1}), //5
		field.element([]uint{2, 1}), //6
		field.element([]uint{1, 2}), //7
		field.element([]uint{2, 2}), //8
	}
	sumTable := [][]*Element{
		{elems[0], elems[1], elems[2], elems[3], elems[4], elems[5], elems[6], elems[7], elems[8]},
		{elems[1], elems[2], elems[0], elems[5], elems[7], elems[6], elems[3], elems[8], elems[4]},
		{elems[2], elems[0], elems[1], elems[6], elems[8], elems[3], elems[5], elems[4], elems[7]},
		{elems[3], elems[5], elems[6], elems[4], elems[0], elems[7], elems[8], elems[1], elems[2]},
		{elems[4], elems[7], elems[8], elems[0], elems[3], elems[1], elems[2], elems[5], elems[6]},
		{elems[5], elems[6], elems[3], elems[7], elems[1], elems[8], elems[4], elems[2], elems[0]},
		{elems[6], elems[3], elems[5], elems[8], elems[2], elems[4], elems[7], elems[0], elems[1]},
		{elems[7], elems[8], elems[4], elems[1], elems[5], elems[2], elems[0], elems[6], elems[3]},
		{elems[8], elems[4], elems[7], elems[2], elems[6], elems[0], elems[1], elems[3], elems[5]},
	}
	diffTable := [][]*Element{
		{elems[0], elems[2], elems[1], elems[4], elems[3], elems[8], elems[7], elems[6], elems[5]},
		{elems[1], elems[0], elems[2], elems[7], elems[5], elems[4], elems[8], elems[3], elems[6]},
		{elems[2], elems[1], elems[0], elems[8], elems[6], elems[7], elems[4], elems[5], elems[3]},
		{elems[3], elems[6], elems[5], elems[0], elems[4], elems[2], elems[1], elems[8], elems[7]},
		{elems[4], elems[8], elems[7], elems[3], elems[0], elems[6], elems[5], elems[2], elems[1]},
		{elems[5], elems[3], elems[6], elems[1], elems[7], elems[0], elems[2], elems[4], elems[8]},
		{elems[6], elems[5], elems[3], elems[2], elems[8], elems[1], elems[0], elems[7], elems[4]},
		{elems[7], elems[4], elems[8], elems[5], elems[1], elems[3], elems[6], elems[0], elems[2]},
		{elems[8], elems[7], elems[4], elems[6], elems[2], elems[5], elems[3], elems[1], elems[0]},
	}
	prodTable := [][]*Element{
		{elems[0], elems[0], elems[0], elems[0], elems[0], elems[0], elems[0], elems[0], elems[0]},
		{elems[0], elems[1], elems[2], elems[3], elems[4], elems[5], elems[6], elems[7], elems[8]},
		{elems[0], elems[2], elems[1], elems[4], elems[3], elems[8], elems[7], elems[6], elems[5]},
		{elems[0], elems[3], elems[4], elems[5], elems[8], elems[7], elems[1], elems[2], elems[6]},
		{elems[0], elems[4], elems[3], elems[8], elems[5], elems[6], elems[2], elems[1], elems[7]},
		{elems[0], elems[5], elems[8], elems[7], elems[6], elems[2], elems[3], elems[4], elems[1]},
		{elems[0], elems[6], elems[7], elems[1], elems[2], elems[3], elems[8], elems[5], elems[4]},
		{elems[0], elems[7], elems[6], elems[2], elems[1], elems[4], elems[5], elems[8], elems[3]},
		{elems[0], elems[8], elems[5], elems[6], elems[7], elems[1], elems[4], elems[3], elems[2]},
	}
	invList := []*Element{
		elems[1], elems[2], elems[6], elems[7], elems[8], elems[3], elems[4], elems[5],
	}

	hardcodedTableTest(
		field,
		elems,
		sumTable,
		diffTable,
		prodTable,
		invList,
		t,
	)
}
