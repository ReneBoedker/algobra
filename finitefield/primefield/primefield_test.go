package primefield

import (
	"algobra/errors"
	"math/bits"
	"math/rand"
	"testing"
	"time"
)

var prg = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func DefineField(char uint) *Field {
	f, err := Define(char)
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

func TestNonPrimeInput(t *testing.T) {
	testCases := []uint{0, 1, 8, 10, 77}
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
	field := DefineField(23)
	if !field.element(20).Equal(field.ElementFromSigned(-3)) {
		t.Errorf("Reported 20!=20 (mod 23)")
	}
	field2 := DefineField(13)
	if field.element(7).Equal(field2.ElementFromSigned(7)) {
		t.Errorf("Reported equality for elements from different fields")
	}
}

func TestTableMemory(t *testing.T) {
	var bigPrime uint
	if bits.UintSize == 32 {
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

func TestArithmeticErrors(t *testing.T) {
	fieldA := DefineField(11)
	fieldB := DefineField(17)

	a := fieldA.element(0)
	b := fieldB.element(10)
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

func TestBools(t *testing.T) {
	field := DefineField(47)
	if field.element(0).IsNonzero() {
		t.Errorf("Element(0) element considered non-zero")
	}
	if !field.element(0).IsZero() {
		t.Errorf("Element(0) element not considered zero")
	}
	if !field.element(1).IsOne() {
		t.Errorf("Element(1) not considered as one")
	}
	if !field.element(1).IsNonzero() {
		t.Errorf("Element(1) not considered non-zero")
	}
}

func TestConversion(t *testing.T) {
	field := DefineField(71)

	for i := 0; i < 1000; i++ {
		val := uint(prg.Uint32())

		a := field.element(val)
		if a.Uint() != val%71 {
			t.Errorf("Conversion to uint failed for a = %d. Received %d",
				val, a.Uint(),
			)
		}

		b := field.element(0)
		b.SetUnsigned(val)
		if b.Uint() != val%71 {
			t.Errorf("Conversion to uint failed for b = %d. Received %d",
				val, b.Uint(),
			)
		}
	}
}

func TestPow(t *testing.T) {
	field := DefineField(13)
	elems := []uint{0, 1, 2, 3, 4}
	expectedPows := [][][2]uint{
		{{0, 1}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {5, 1}, {14, 1}},
		{{0, 1}, {2, 4}, {3, 8}, {4, 3}, {13, 2}, {16, 3}},
		{{0, 1}, {1, 3}, {2, 9}, {3, 1}, {4, 3}, {14, 9}},
		{{0, 1}, {1, 4}, {2, 3}, {3, 12}},
	}
	for i, val := range elems {
		elem := field.element(val)
		for _, j := range expectedPows[i] {
			res := elem.Pow(j[0])
			if !res.Equal(field.element(j[1])) {
				t.Errorf("Pow failed: %v^%d = %v (Expected %v)", elem, j[0], res, j[1])
			}
		}
	}
}

func TestGenerator(t *testing.T) {
	for _, p := range []uint{2, 3, 5, 7, 11} {
		unique := make(map[uint]struct{})

		field, err := Define(p)
		if err != nil {
			panic(err)
		}

		g := field.MultGenerator()
		for i, e := uint(0), g.Copy(); i < p-1; i, e = i+1, e.Times(g) {
			ee := e.(*Element)
			if _, ok := unique[ee.Uint()]; ok {
				t.Errorf("Found element %v twice for p=%v (generator = %v)", e, p, g)
			} else {
				unique[ee.Uint()] = struct{}{}
			}
		}
	}
}

func TestElements(t *testing.T) {
	for _, p := range []uint{2, 3, 5, 7, 11} {
		unique := make(map[uint]struct{})

		field, err := Define(p)
		if err != nil {
			panic(err)
		}

		for _, e := range field.Elements() {
			ee := e.(*Element)
			if _, ok := unique[ee.Uint()]; ok {
				t.Errorf("Found element %v twice for p=%v", e, p)
			} else {
				unique[ee.Uint()] = struct{}{}
			}
		}
	}
}

func TestNeg(t *testing.T) {
	for _, card := range []uint{3, 7, 13, 31} {
		field := DefineField(card)
		for _, e := range field.Elements() {
			if tmp := e.Plus(e.Neg()); tmp.IsNonzero() {
				t.Errorf("%[1]v + (-%[1]v) returned %[2]v rather than 0", e, tmp)
			}
		}
	}
}

func TestProd(t *testing.T) {
	field := DefineField(11)

	test := func(field *Field) {
		a := field.element(0)

		prods := [][2]uint{
			{2, 2},
			{2, 5},
			{2, 7},
			{7, 2},
			{3, 3},
			{6, 7},
			{10, 10},
		}
		expected := []uint{
			4,
			10,
			3,
			3,
			9,
			9,
			1,
		}
		for i, p := range prods {
			if a.Prod(field.element(p[0]), field.element(p[1])); !a.Equal(field.element(expected[i])) {
				t.Errorf("Prod failed: a was set to %v for %v * %v (Expected %v)",
					a, p[0], p[1], expected[i])
			}
		}

	}

	test(field)
	// With tables
	field.ComputeTables(false, true)
	test(field)
}

func TestGf2(t *testing.T) {
	field := DefineField(2)
	test := func(field *Field) {
		elems := []*Element{field.element(0), field.element(1)}
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
				if t1, t2 := elems[i].Minus(elems[j]), sumTable[i][j]; !t1.Equal(t2) { // Note that Plus=Minus
					t.Errorf("GF(2) failed: %v - %v = %v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
				if t1, t2 := elems[i].Times(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(2) failed: %v*%v=%v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
			}
		}
		if t1 := elems[1].Inv(); !t1.Equal(elems[1]) {
			t.Errorf("GF(2) failed: inv(1)=%v (Expected 1)", t1)
		}
	}
	// Without tables
	test(field)
	// With tables
	field.ComputeTables(true, true)
	test(field)
}

func TestGf3(t *testing.T) {
	field := DefineField(3)
	test := func(field *Field) {
		elems := []*Element{field.element(0), field.element(1), field.element(2)}
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
				if t1, t2 := elems[i].Plus(elems[j]), sumTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(3) failed: %v + %v = %v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
				if t1, t2 := elems[i].Minus(elems[j]), diffTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(3) failed: %v - %v =%v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
				if t1, t2 := elems[i].Times(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(3) failed: %v * %v = %v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
			}
		}
		invList := []*Element{field.element(1), field.element(2)}
		for i := 1; i < len(elems); i++ {
			if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
				t.Errorf("GF(3) failed: inv(%v) = %v (Expected %v)", elems[i], t1, invList[i])
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
	field := DefineField(7)
	test := func(field *Field) {
		elems := []*Element{
			field.element(0), field.element(1), field.element(2), field.element(3),
			field.element(4), field.element(5), field.element(6),
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
				if t1, t2 := elems[i].Plus(elems[j]), sumTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(7) failed: %v + %v = %v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
				if t1, t2 := elems[i].Minus(elems[j]), diffTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(7) failed: %v - %v = %v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
				if t1, t2 := elems[i].Times(elems[j]), prodTable[i][j]; !t1.Equal(t2) {
					t.Errorf("GF(7) failed: %v * %v = %v (Expected %v)",
						elems[i], elems[j], t1, t2)
				}
			}
		}
		invList := []*Element{
			field.element(1), field.element(4), field.element(5), field.element(2), field.element(3), field.element(6),
		}
		for i := 1; i < len(elems); i++ {
			if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
				t.Errorf("GF(7) failed: inv(%v) = %v (Expected %v)",
					elems[i], t1, invList[i-1],
				)
			}
		}
	}
	// Without tables
	test(field)
	// With tables
	field.ComputeTables(true, true)
	test(field)
}

/* Copyright 2019 René Bødker Christensen
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
