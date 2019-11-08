package binfield

import (
	"testing"

	"github.com/ReneBoedker/algobra/finitefield/ff"
)

func hardcodedTableTest(
	f *Field,
	elems []ff.Element,
	sumTable, prodTable [][]ff.Element,
	//invList []*Element,
	t *testing.T,
) {
	test := func(f *Field) {
		for i := range elems {
			for j := range elems {
				if t1, t2 := elems[i].Plus(elems[j]), sumTable[i][j]; !t1.Equal(t2) {
					t.Errorf(
						"GF(%d) failed: (%v) + (%v) = %v (Expected %v)",
						f.Card(), elems[i], elems[j], t1, t2,
					)
				}
				if t1, t2 := elems[i].Times(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
					t.Errorf(
						"GF(%d) failed: (%v) * (%v) = %v (Expected %v)",
						f.Card(), elems[i], elems[j], t1, t2,
					)
				}
			}
		}
		// for i := 1; i < len(elems); i++ {
		// 	if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
		// 		t.Errorf(
		// 			"GF(%d) failed: inv(%v) = %v (Expected %v)",
		// 			f.Card(), elems[i], t1, invList[i-1],
		// 		)
		// 	}
		// }
	}
	test(f)
}

func TestGf4(t *testing.T) {
	field, _ := Define(4)
	elems := []ff.Element{
		field.ElementFromBits(0),
		field.ElementFromBits(1),
		field.ElementFromBits(2),
		field.ElementFromBits(3),
	}
	sumTable := [][]ff.Element{
		{elems[0], elems[1], elems[2], elems[3]},
		{elems[1], elems[0], elems[3], elems[2]},
		{elems[2], elems[3], elems[0], elems[1]},
		{elems[3], elems[2], elems[1], elems[0]},
	}
	prodTable := [][]ff.Element{
		{elems[0], elems[0], elems[0], elems[0]},
		{elems[0], elems[1], elems[2], elems[3]},
		{elems[0], elems[2], elems[3], elems[1]},
		{elems[0], elems[3], elems[1], elems[2]},
	}
	//invList := []*Element{elems[1], elems[3], elems[2]}

	hardcodedTableTest(
		field,
		elems,
		sumTable,
		prodTable,
		//invList,
		t,
	)
}

func TestGf8(t *testing.T) {
	field, _ := Define(8)
	elems := []ff.Element{
		field.ElementFromBits(0), // 0
		field.ElementFromBits(1), // 1
		field.ElementFromBits(2), // a
		field.ElementFromBits(3), // a+1
		field.ElementFromBits(4), // a^2
		field.ElementFromBits(5), // a^2+1
		field.ElementFromBits(6), // a^2+a
		field.ElementFromBits(7), // a^2+a+1
	}
	sumTable := [][]ff.Element{
		{elems[0], elems[1], elems[2], elems[3], elems[4], elems[5], elems[6], elems[7]},
		{elems[1], elems[0], elems[3], elems[2], elems[5], elems[4], elems[7], elems[6]},
		{elems[2], elems[3], elems[0], elems[1], elems[6], elems[7], elems[4], elems[5]},
		{elems[3], elems[2], elems[1], elems[0], elems[7], elems[6], elems[5], elems[4]},
		{elems[4], elems[5], elems[6], elems[7], elems[0], elems[1], elems[2], elems[3]},
		{elems[5], elems[4], elems[7], elems[6], elems[1], elems[0], elems[3], elems[2]},
		{elems[6], elems[7], elems[4], elems[5], elems[2], elems[3], elems[0], elems[1]},
		{elems[7], elems[6], elems[5], elems[4], elems[3], elems[2], elems[1], elems[0]},
	}
	prodTable := [][]ff.Element{
		{elems[0], elems[0], elems[0], elems[0], elems[0], elems[0], elems[0], elems[0]},
		{elems[0], elems[1], elems[2], elems[3], elems[4], elems[5], elems[6], elems[7]},
		{elems[0], elems[2], elems[4], elems[6], elems[3], elems[1], elems[7], elems[5]},
		{elems[0], elems[3], elems[6], elems[5], elems[7], elems[4], elems[1], elems[2]},
		{elems[0], elems[4], elems[3], elems[7], elems[6], elems[2], elems[5], elems[1]},
		{elems[0], elems[5], elems[1], elems[4], elems[2], elems[7], elems[3], elems[6]},
		{elems[0], elems[6], elems[7], elems[1], elems[5], elems[3], elems[2], elems[4]},
		{elems[0], elems[7], elems[5], elems[2], elems[1], elems[6], elems[4], elems[3]},
	}
	//invList := []*Element{elems[1], elems[3], elems[2]}

	hardcodedTableTest(
		field,
		elems,
		sumTable,
		prodTable,
		//invList,
		t,
	)
}
