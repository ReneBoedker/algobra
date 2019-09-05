package primefield

import (
	"algobra/errors"
	"testing"
)

func TestInit(t *testing.T) {
	if uintBitSize == 0 {
		t.Error("init() failed to detect size of uint")
	}
}

func TestOverflowDetection(t *testing.T) {
	var bigPrime uint
	if uintBitSize == 32 {
		bigPrime = uint(65537)
	} else {
		bigPrime = uint(4294967311)
	}
	_, err := Define(bigPrime)
	if err == nil {
		t.Errorf("Define succeeded even though p=%d", bigPrime)
	} else if !errors.Is(errors.InputTooLarge, err) {
		t.Errorf("Define failed, but the error kind was unexpected")
	}
}

func TestGf2(t *testing.T) {
	field, _ := Define(2)
	test := func(field *Field) {
		elems := []*Element{field.Element(0), field.Element(1)}
		sumTable := [][]*Element{
			{elems[0], elems[1]},
			{elems[1], elems[0]},
		}
		prodTable := [][]*Element{
			{elems[0], elems[0]},
			{elems[0], elems[1]},
		}
		for i := range elems {
			for j := range elems {
				if t1, t2 := elems[i].Plus(elems[j]), sumTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(2) failed: %v+%v=%v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
				if t1, t2 := elems[i].Minus(elems[j]).val, sumTable[i][j].val; t1 != t2 { // Note that Plus=Minus
					t.Errorf("GF(7) failed: %d-%d=%d (Expected %d)",
						elems[i].val, elems[j].val, t1, t2)
				}
				if t1, t2 := elems[i].Mult(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(2) failed: %v*%v=%v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
			}
		}
		if t1 := elems[1].Inv(); !t1.Equal(elems[1]) {
			t.Errorf("GF(2) failed: inv(1)=%d (Expected 1)", t1.val)
		}
	}
	// Without tables
	test(field)
	// With tables
	field.ComputeTables(true, true)
	test(field)
}

func TestGf3(t *testing.T) {
	field, _ := Define(3)
	test := func(field *Field) {
		elems := []*Element{field.Element(0), field.Element(1), field.Element(2)}
		sumTable := [][]*Element{
			{elems[0], elems[1], elems[2]},
			{elems[1], elems[2], elems[0]},
			{elems[2], elems[0], elems[1]},
		}
		diffTable := [][]*Element{
			{elems[0], elems[2], elems[1]},
			{elems[1], elems[0], elems[2]},
			{elems[2], elems[1], elems[0]},
		}
		prodTable := [][]*Element{
			{elems[0], elems[0], elems[0]},
			{elems[0], elems[1], elems[2]},
			{elems[0], elems[2], elems[1]},
		}
		for i := range elems {
			for j := range elems {
				if t1, t2 := elems[i].Plus(elems[j]).val, sumTable[i][j].val; t1 != t2 {
					t.Errorf("GF(3) failed: %d+%d=%d (Expected %d)",
						elems[i].val, elems[j].val, t1, t2)
				}
				if t1, t2 := elems[i].Minus(elems[j]).val, diffTable[i][j].val; t1 != t2 {
					t.Errorf("GF(7) failed: %d-%d=%d (Expected %d)",
						elems[i].val, elems[j].val, t1, t2)
				}
				if t1, t2 := elems[i].Mult(elems[j]).val, prodTable[i][j].val; t1 != t2 {
					t.Errorf("GF(3) failed: %d*%d=%d (Expected %d)",
						elems[i].val, elems[j].val, t1, t2)
				}
			}
		}
		invList := []*Element{field.Element(1), field.Element(2)}
		for i := 1; i < len(elems); i++ {
			if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
				t.Errorf("GF(3) failed: inv(%d)=%d (Expected %d)", elems[i].val, t1.val, invList[i].val)
			}
		}
	}
	// Without tables
	test(field)
	// With tables
	field.ComputeTables(true, true)
	test(field)
}

func TestGf7(t *testing.T) {
	field, _ := Define(7)
	test := func(field *Field) {
		elems := []*Element{
			field.Element(0), field.Element(1), field.Element(2), field.Element(3),
			field.Element(4), field.Element(5), field.Element(6),
		}
		sumTable := [][]*Element{
			{elems[0], elems[1], elems[2], elems[3], elems[4], elems[5], elems[6]},
			{elems[1], elems[2], elems[3], elems[4], elems[5], elems[6], elems[0]},
			{elems[2], elems[3], elems[4], elems[5], elems[6], elems[0], elems[1]},
			{elems[3], elems[4], elems[5], elems[6], elems[0], elems[1], elems[2]},
			{elems[4], elems[5], elems[6], elems[0], elems[1], elems[2], elems[3]},
			{elems[5], elems[6], elems[0], elems[1], elems[2], elems[3], elems[4]},
			{elems[6], elems[0], elems[1], elems[2], elems[3], elems[4], elems[5]},
		}
		diffTable := [][]*Element{
			{elems[0], elems[6], elems[5], elems[4], elems[3], elems[2], elems[1]},
			{elems[1], elems[0], elems[6], elems[5], elems[4], elems[3], elems[2]},
			{elems[2], elems[1], elems[0], elems[6], elems[5], elems[4], elems[3]},
			{elems[3], elems[2], elems[1], elems[0], elems[6], elems[5], elems[4]},
			{elems[4], elems[3], elems[2], elems[1], elems[0], elems[6], elems[5]},
			{elems[5], elems[4], elems[3], elems[2], elems[1], elems[0], elems[6]},
			{elems[6], elems[5], elems[4], elems[3], elems[2], elems[1], elems[0]},
		}
		prodTable := [][]*Element{
			{elems[0], elems[0], elems[0], elems[0], elems[0], elems[0], elems[0]},
			{elems[0], elems[1], elems[2], elems[3], elems[4], elems[5], elems[6]},
			{elems[0], elems[2], elems[4], elems[6], elems[1], elems[3], elems[5]},
			{elems[0], elems[3], elems[6], elems[2], elems[5], elems[1], elems[4]},
			{elems[0], elems[4], elems[1], elems[5], elems[2], elems[6], elems[3]},
			{elems[0], elems[5], elems[3], elems[1], elems[6], elems[4], elems[2]},
			{elems[0], elems[6], elems[5], elems[4], elems[3], elems[2], elems[1]},
		}
		for i := range elems {
			for j := range elems {
				if t1, t2 := elems[i].Plus(elems[j]).val, sumTable[i][j].val; t1 != t2 {
					t.Errorf("GF(7) failed: %d+%d=%d (Expected %d)",
						elems[i].val, elems[j].val, t1, t2)
				}
				if t1, t2 := elems[i].Minus(elems[j]).val, diffTable[i][j].val; t1 != t2 {
					t.Errorf("GF(7) failed: %d-%d=%d (Expected %d)",
						elems[i].val, elems[j].val, t1, t2)
				}
				if t1, t2 := elems[i].Mult(elems[j]).val, prodTable[i][j].val; t1 != t2 {
					t.Errorf("GF(7) failed: %d*%d=%d (Expected %d)",
						elems[i].val, elems[j].val, t1, t2)
				}
			}
		}
		invList := []*Element{
			field.Element(1), field.Element(4), field.Element(5), field.Element(2), field.Element(3), field.Element(6),
		}
		for i := 1; i < len(elems); i++ {
			if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
				t.Errorf("GF(7) failed: inv(%d)=%d (Expected %d)", elems[i].val, t1.val, invList[i].val)
			}
		}
	}
	// Without tables
	test(field)
	// With tables
	field.ComputeTables(true, true)
	test(field)
}

func TestEqual(t *testing.T) {
	field, _ := Define(23)
	if !field.Element(20).Equal(field.ElementFromSigned(-3)) {
		t.Errorf("Reported 20!=20 (mod 23)")
	}
	field2, _ := Define(13)
	if field.Element(7).Equal(field2.ElementFromSigned(7)) {
		t.Errorf("Reported equality for elements from different fields")
	}
}

func TestTableMemory(t *testing.T) {
	var bigPrime uint
	if uintBitSize == 32 {
		bigPrime = uint(16411)
	} else {
		bigPrime = uint(11587)
	}
	f, _ := Define(bigPrime)
	if err := f.ComputeTables(true, false); err == nil {
		t.Errorf("No error returned")
	} else if !errors.Is(errors.InputTooLarge, err) {
		t.Errorf("Error returned has wrong kind. Expected errors.InputTooLarge, "+
			"but received error %q", err.Error())
	}
}
