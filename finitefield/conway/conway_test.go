package conway

import (
	"github.com/ReneBoedker/algobra/errors"
	"testing"
)

func TestConwayKnownValues(t *testing.T) {
	inputs := [][2]uint{
		{2, 23},
		{2, 409},
		{79, 37},
		{179, 8},
		{239, 47},
		{941, 9},
		{6089, 6},
		{77801, 4},
	}

	for _, i := range inputs {
		if _, err := Lookup(i[0], i[1]); err != nil {
			t.Errorf(
				"No polynomial was found for characteristic %d and extension degree %d.",
				i[0], i[1],
			)
		}
	}
}

func TestConwayUnknownValues(t *testing.T) {
	inputs := [][2]uint{
		{2, 405},
		{64301, 5},
		{77801, 3},
	}

	for _, i := range inputs {
		if _, err := Lookup(i[0], i[1]); err == nil {
			t.Errorf(
				"No error was returned for characteristic %d and extension degree %d.",
				i[0], i[1],
			)
		} else if !errors.Is(errors.InputValue, err) {
			t.Errorf(
				"Error was returned for characteristic %d and extension degree "+
					"%d, but it was of unexpected kind.",
				i[0], i[1],
			)
		}
	}
}

func TestConwayInternalErr(t *testing.T) {
	errList := `[2,1,[1,1,1]],
[2,2,[1,1,1.7]],
`
	inputs := [][2]uint{
		{2, 1},
		{2, 2},
	}

	for _, i := range inputs {
		if _, err := lookupInternal(i[0], i[1], errList); err == nil {
			t.Errorf(
				"No error was returned for characteristic %d and extension degree %d.",
				i[0], i[1],
			)
		} else if !errors.Is(errors.Internal, err) {
			t.Errorf(
				"Error was returned for characteristic %d and extension degree "+
					"%d, but it was of unexpected kind.",
				i[0], i[1],
			)
		}
	}
}
