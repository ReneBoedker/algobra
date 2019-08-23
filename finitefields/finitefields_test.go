package finitefields

import (
	"testing"
)

func elem(v, m uint) *Element {
	return &Element{val: v % m, mod: m}
}

func TestGf2(t *testing.T) {
	elems := []*Element{elem(0, 2), elem(1, 2)}
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
				t.Errorf("GF(2) failed: %d+%d=%d (Expected %d)",
					elems[i].val, elems[j].val, t1, t2)
			}
			if t1, t2 := elems[i].Mult(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
				t.Errorf("GF(2) failed: %d*%d=%d (Expected %d)",
					elems[i].val, elems[j].val, t1.val, t2.val)
			}
		}
	}
	if t1 := elems[1].Inv(); !t1.Equal(elems[1]) {
		t.Errorf("GF(2) failed: inv(1)=%d (Expected 1)", t1.val)
	}
}

func TestGf3(t *testing.T) {
	elems := []*Element{elem(0, 3), elem(1, 3), elem(2, 3)}
	sumTable := [][]*Element{
		{elems[0], elems[1], elems[2]},
		{elems[1], elems[2], elems[0]},
		{elems[2], elems[0], elems[1]},
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
			if t1, t2 := elems[i].Mult(elems[j]).val, prodTable[i][j].val; t1 != t2 {
				t.Errorf("GF(3) failed: %d*%d=%d (Expected %d)",
					elems[i].val, elems[j].val, t1, t2)
			}
		}
	}
	invList := []*Element{elem(1, 3), elem(2, 3)}
	for i := 1; i < len(elems); i++ {
		if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
			t.Errorf("GF(3) failed: inv(%d)=%d (Expected %d)", elems[i].val, t1.val, invList[i].val)
		}
	}
}

func TestGf7(t *testing.T) {
	elems := []*Element{
		elem(0, 7), elem(1, 7), elem(2, 7), elem(3, 7),
		elem(4, 7), elem(5, 7), elem(6, 7),
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
			if t1, t2 := elems[i].Mult(elems[j]).val, prodTable[i][j].val; t1 != t2 {
				t.Errorf("GF(7) failed: %d*%d=%d (Expected %d)",
					elems[i].val, elems[j].val, t1, t2)
			}
		}
	}
	invList := []*Element{
		elem(1, 7), elem(4, 7), elem(5, 7), elem(2, 7), elem(3, 7), elem(6, 7),
	}
	for i := 1; i < len(elems); i++ {
		if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
			t.Errorf("GF(7) failed: inv(%d)=%d (Expected %d)", elems[i].val, t1.val, invList[i].val)
		}
	}
}
