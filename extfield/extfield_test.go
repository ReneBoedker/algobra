package extfield

import (
	"testing"
)

func DefineField(card uint) *Field {
	f, err := Define(card)
	if err != nil {
		// Testing code is wrong, so panic
		panic(err)
	}
	return f
}

func TestGf4(t *testing.T) {
	field := DefineField(4)
	test := func(field *Field) {
		elems := []*Element{
			field.Element([]uint{0, 0}),
			field.Element([]uint{1, 0}),
			field.Element([]uint{0, 1}),
			field.Element([]uint{1, 1}),
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
		for i := range elems {
			for j := range elems {
				if t1, t2 := elems[i].Plus(elems[j]), sumTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(4) failed: (%v) + (%v) = %v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
				if t1, t2 := elems[i].Minus(elems[j]), sumTable[i][j]; !t1.Equal(t2) { // Note that Plus=Minus
					t.Errorf("GF(4) failed: (%v) - (%v) = %v (Expected %v)",
						elems[i].val, elems[j].val, t1, t2)
				}
				if t1, t2 := elems[i].Times(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(4) failed: %v*%v=%v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
			}
		}
		// if t1 := elems[1].Inv(); !t1.Equal(elems[1]) {
		// 	t.Errorf("GF(2) failed: inv(1)=%d (Expected 1)", t1.val)
		// }
	}
	// Without tables
	test(field)
	// With tables
	//field.ComputeTables(true, true)
	//test(field)
}

func TestGf9(t *testing.T) {
	field := DefineField(9)
	test := func(field *Field) {
		elems := [9]*Element{
			field.Element([]uint{0, 0}), //0
			field.Element([]uint{1, 0}), //1
			field.Element([]uint{2, 0}), //2
			field.Element([]uint{0, 1}), //3
			field.Element([]uint{0, 2}), //4
			field.Element([]uint{1, 1}), //5
			field.Element([]uint{2, 1}), //6
			field.Element([]uint{1, 2}), //7
			field.Element([]uint{2, 2}), //8
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
		for i := range elems {
			for j := range elems {
				if t1, t2 := elems[i].Plus(elems[j]), sumTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(4) failed: (%v) + (%v) = %v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
				// if t1, t2 := elems[i].Minus(elems[j]), sumTable[i][j]; !t1.Equal(t2) { // Note that Plus=Minus
				// 	t.Errorf("GF(4) failed: (%v) - (%v) = %v (Expected %v)",
				// 		elems[i].val, elems[j].val, t1, t2)
				// }
				if t1, t2 := elems[i].Times(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(4) failed: %v*%v=%v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
			}
		}
		// if t1 := elems[1].Inv(); !t1.Equal(elems[1]) {
		// 	t.Errorf("GF(2) failed: inv(1)=%d (Expected 1)", t1.val)
		// }
	}
	// Without tables
	test(field)
	// With tables
	//field.ComputeTables(true, true)
	//test(field)
}
