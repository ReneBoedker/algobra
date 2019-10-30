package bivariate

import (
	"testing"

	"github.com/ReneBoedker/algobra/errors"
	"github.com/ReneBoedker/algobra/finitefield"
	"github.com/ReneBoedker/algobra/finitefield/ff"
)

func defineField(char uint, t *testing.T) ff.Field {
	field, err := finitefield.Define(char)
	if err != nil {
		t.Fatalf("Failed to define finite field of %d elements", char)
	}
	return field
}

func assertError(t *testing.T, err error, k errors.Kind, desc string, args ...interface{}) {
	if err == nil {
		t.Errorf(desc+" returned no error", args)
	} else if !errors.Is(k, err) {
		t.Errorf(desc+" returned an error but not of the correct type", args...)
	}
}

func fieldLoop(do func(field ff.Field), minCard ...uint) {
	for _, card := range [...]uint{2, 3, 4, 5, 9, 16, 25, 49, 64, 125} {
		if len(minCard) > 0 && card < minCard[0] {
			continue
		}
		field, err := finitefield.Define(card)
		if err != nil {
			// Error is in tests, so panic is OK
			panic(err)
		}

		do(field)
	}
	return
}

func TestReduce(t *testing.T) {
	field := defineField(3, t)
	r := DefRing(field, WDegLex(3, 4, false))
	mod := r.PolynomialFromUnsigned(map[[2]uint]uint{{9, 0}: 1, {1, 0}: 2})
	id, err := r.NewIdeal(mod)
	if err != nil {
		t.Errorf("Failed to construct ideal. Error message %q", err.Error())
	}
	qr, err := r.Quotient(id)
	if err != nil {
		t.Errorf("Failed to construct quotient ring")
	}
	f := qr.PolynomialFromUnsigned(map[[2]uint]uint{
		{12, 3}: 1,
	})

	if f.Ld() != [2]uint{4, 3} {
		t.Errorf("Reduce failed: Got %v", f.Ld())
	}
}

func TestGroebner1(t *testing.T) {
	field := defineField(7, t)
	r := DefRing(field, Lex(true))
	id, _ := r.NewIdeal(
		r.PolynomialFromSigned(map[[2]uint]int{
			{1, 2}: 1,
			{0, 3}: -1,
		}),
		r.PolynomialFromSigned(map[[2]uint]int{
			{0, 3}: 1,
			{0, 2}: -1,
		}),
	)
	expectedGens := []*Polynomial{
		r.PolynomialFromSigned(map[[2]uint]int{
			{1, 2}: 1,
			{0, 2}: -1,
		}),
		r.PolynomialFromSigned(map[[2]uint]int{
			{0, 3}: 1,
			{0, 2}: -1,
		}),
	}
	id = id.GroebnerBasis()
	id.ReduceBasis()
	if len(id.generators) != 2 {
		t.Fatalf("Gröbner basis has wrong number of elements")
	}
	if (!id.generators[0].Equal(expectedGens[0]) && !id.generators[0].Equal(expectedGens[1])) || (!id.generators[1].Equal(expectedGens[0]) && !id.generators[1].Equal(expectedGens[1])) {
		t.Errorf("Got Gröbner basis %v", id.generators)
	}
}

// Examples 5.7.1 and 5.8.4 in [Lauritzen]
func TestGroebner2(t *testing.T) {
	field := defineField(9, t)
	r := DefRing(field, Lex(true))

	generators := []*Polynomial{
		r.Polynomial(map[[2]uint]ff.Element{
			{2, 0}: field.One(),
			{0, 1}: field.One(),
		}),
		r.Polynomial(map[[2]uint]ff.Element{
			{2, 1}: field.One(),
			{0, 0}: field.One(),
		}),
	}

	id, _ := r.NewIdeal(generators...)

	if id.IsGroebner() {
		t.Errorf("IsGroebner returned true, but false was expected")
	}

	id = id.GroebnerBasis()

	// Add the final generator
	generators = append(
		generators,
		r.Polynomial(map[[2]uint]ff.Element{
			{0, 2}: field.One(),
			{0, 0}: field.One().SetNeg(),
		}),
	)

	if len(id.Generators()) != 3 {
		t.Errorf("Gröbner basis has %d generators, but expected 3", len(id.Generators()))
	}
outer:
	for _, g := range generators {
		for _, h := range id.generators {
			if g.Equal(h) {
				continue outer
			}
		}
		t.Errorf("Generator %v was not in the Gröbner basis", g)
	}

	if id.IsMinimal() {
		t.Errorf("IsMinimal returned true, but false was expected")
	}
	if id.IsReduced() {
		t.Errorf("IsReduced returned true, but false was expected")
	}

	// Now compute the reduced basis
	id.ReduceBasis()
	// We expect the second generator to be removed
	generators = append(generators[:1], generators[2])
	if len(id.Generators()) != 2 {
		t.Errorf(
			"Reduced Gröbner basis has %d generators, but expected 3",
			len(id.Generators()),
		)
	}
outerReduced:
	for _, g := range generators {
		for _, h := range id.generators {
			if g.Equal(h) {
				continue outerReduced
			}
		}
		t.Errorf("Generator %v was not in the reduced Gröbner basis", g)
	}

	if !id.IsMinimal() {
		t.Errorf("IsMinimal returned false, but true was expected")
	}
	if !id.IsReduced() {
		t.Errorf("IsReduced returned false, but true was expected")
	}
}

func TestQuotientErrors(t *testing.T) {
	field1 := defineField(49, t)
	field2 := defineField(25, t)

	ring1 := DefRing(field1, Lex(true))
	ring2 := DefRing(field2, Lex(false))

	id, err := ring1.NewIdeal(ring1.PolynomialFromSigned(map[[2]uint]int{
		{2, 0}: 1,
		{0, 1}: 2,
	}))
	if err != nil {
		t.Fatalf("Could not create initial ideal")
	}

	qr, err := ring1.Quotient(id)
	if err != nil {
		t.Fatalf("Could not create initial quotient ring")
	}

	if _, err := qr.Quotient(id); true {
		assertError(t, err, errors.InputValue, "Constructing quotient of ring with non-nil ideal")
	}

	if _, err := ring2.Quotient(id); true {
		assertError(t, err, errors.InputIncompatible, "Constructing quotient with ideal from different ring")
	}
}
