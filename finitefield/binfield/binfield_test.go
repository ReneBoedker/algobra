package binfield

import (
	"math/bits"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

var prg = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

// assertError checks that err is non-nil and that it has kind k.
func assertError(t *testing.T, err error, k errors.Kind, desc string, args ...interface{}) {
	if err == nil {
		t.Errorf(desc+" returned no error", args)
	} else if !errors.Is(k, err) {
		t.Errorf(desc+" returned an error but not of the correct type", args)
	}
}

func TestOverflowDetection(t *testing.T) {
	bigCard := uint(1) << (bits.UintSize/2 + 1)
	_, err := Define(bigCard)
	assertError(t, err, errors.InputTooLarge, "Define(%d)", bigCard)
}

func TestNonBinaryCard(t *testing.T) {
	testCases := []uint{0, 5, 7, 9, 100}
	for _, card := range testCases {
		_, err := Define(card)
		assertError(t, err, errors.InputValue, "Define(%d)", card)
	}
}

func TestGenerator(t *testing.T) {
	for _, p := range []uint{2, 4, 8, 16, 32} {
		unique := make(map[uint]struct{})

		field, err := Define(p)
		if err != nil {
			panic(err)
		}

		g := field.MultGenerator()
		for i, e := uint(0), g.Copy(); i < p-1; i, e = i+1, e.Times(g) {
			ee := e.(*Element)
			if _, ok := unique[ee.val]; ok {
				t.Errorf("Found element %v twice for p=%v (generator = %v)", e, p, g)
			} else {
				unique[ee.val] = struct{}{}
			}
		}
	}
}

func TestElements(t *testing.T) {
	for _, p := range []uint{2, 4, 8, 16, 32} {
		unique := make(map[uint]struct{})

		field, err := Define(p)
		if err != nil {
			panic(err)
		}

		for _, e := range field.Elements() {
			ee := e.(*Element)
			if _, ok := unique[ee.val]; ok {
				t.Errorf("Found element %v twice for p=%v", e, p)
			} else {
				unique[ee.val] = struct{}{}
			}
		}
	}
}

func TestPow(t *testing.T) {
	field, _ := Define(128)
	for rep := 0; rep < 50; rep++ {
		a := field.RandElement()

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

	// Check the zero cases separately
	if tmp := field.Zero().Pow(0); !tmp.IsOne() {
		t.Errorf("0^0 = %v, but expected 1", tmp)
	}
	if tmp := field.Zero().Pow(1); !tmp.IsZero() {
		t.Errorf("0^1 = %v, but expected 0", tmp)
	}
}

func TestBools(t *testing.T) {
	field, _ := Define(256)
	if field.Zero().IsNonzero() {
		t.Errorf("Element(0) element considered non-zero")
	}
	if !field.Zero().IsZero() {
		t.Errorf("Element(0) element not considered zero")
	}
	if !field.One().IsOne() {
		t.Errorf("Element(1) not considered as one")
	}
	if !field.One().IsNonzero() {
		t.Errorf("Element(1) not considered non-zero")
	}
}

func TestRegexElement(t *testing.T) {
	for _, card := range []uint{2, 4, 8, 16, 32, 64} {
		field, _ := Define(card)

		pattern, err := regexp.Compile(field.RegexElement(false))
		if err != nil {
			t.Fatalf("Failed to compile regular expression %q", field.RegexElement(false))
		}

		patternParens, err := regexp.Compile(field.RegexElement(true))
		if err != nil {
			t.Fatalf("Failed to compile regular expression %q", field.RegexElement(true))
		}

		for rep := 0; rep < 50; rep++ {
			a := field.RandElement()

			s := a.String()
			if tmp := pattern.FindString(s); tmp != s {
				// The pattern without parentheses must match the entire string
				t.Errorf(
					"%q was matched as %q without requiring parentheses",
					s, tmp,
				)
			}
			if tmp := pattern.FindString("(" + s + ")"); tmp != s {
				// If the pattern does not require parentheses, these should be
				// ignored when matching
				t.Errorf(
					"(%q) was matched as %q when requiring parentheses",
					s, tmp,
				)
			}

			if tmp := patternParens.FindString(s); a.NTerms() > 1 && tmp == s {
				// Matching without parentheses is only allowed for single terms
				// when parentheses are required
				t.Errorf(
					"%q was matched as %q even though parentheses were required",
					s, tmp,
				)
			}
			if tmp := patternParens.FindString("(" + s + ")"); tmp != "("+s+")" {
				// If parentheses are in the string and required in the search
				// pattern, the match must contain them
				t.Errorf(
					"\"(%s)\" was matched as %q when requiring parentheses",
					s, tmp,
				)
			}
		}
	}
}

func TestParseOutput(t *testing.T) {
	field, _ := Define(256)
	for _, varName := range []string{"a", "α", "\\beta"} {
		field.SetVarName(varName)
		for rep := 0; rep < 25; rep++ {
			a := field.RandElement()

			if b, err := field.ElementFromString(a.String()); err != nil {
				t.Errorf("Parsing formatted output of %v returns error %q", a, err)
			} else if !a.Equal(b) {
				t.Errorf("Formatted output of %v is parsed as %v", a, b)
			}
		}
	}
}

func hardcodedTableTest(
	f *Field,
	elems []ff.Element,
	sumTable, prodTable [][]ff.Element,
	invList []ff.Element,
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
		for i := 1; i < len(elems); i++ {
			if t1 := elems[i].Inv(); !t1.Equal(invList[i-1]) {
				t.Errorf(
					"GF(%d) failed: inv(%v) = %v (Expected %v)",
					f.Card(), elems[i], t1, invList[i-1],
				)
			}
		}
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
	invList := []ff.Element{elems[1], elems[3], elems[2]}

	hardcodedTableTest(
		field,
		elems,
		sumTable,
		prodTable,
		invList,
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
	invList := []ff.Element{elems[1], elems[5], elems[6], elems[7], elems[2], elems[3], elems[4]}

	hardcodedTableTest(
		field,
		elems,
		sumTable,
		prodTable,
		invList,
		t,
	)
}
