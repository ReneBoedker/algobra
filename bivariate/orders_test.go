package bivariate

import (
	"testing"
)

func TestOrders(t *testing.T) {
	degrees := [][2][2]uint{
		{{1, 1}, {1, 1}},
		{{2, 1}, {1, 1}},
		{{2, 1}, {1, 2}},
		{{0, 4}, {1, 0}},
		{{4, 3}, {3, 4}},
		{{4, 0}, {0, 3}},
	}
	orders := []Order{
		Lex(true),
		Lex(false),
		DegLex(true),
		DegLex(false),
		DegRevLex(true),
		DegRevLex(false),
		WDegLex(3, 4, true),
		WDegLex(3, 4, false),
		WDegRevLex(3, 4, true),
		WDegRevLex(3, 4, false),
	}
	ordStr := []string{
		"Lex (X>Y)",
		"Lex (Y>X)",
		"DegLex (X>Y)",
		"DegLex (Y>X)",
		"DegRevLex (X>Y)",
		"DegRevLex (Y>X)",
		"WDegLex(3,4) (X>Y)",
		"WDegLex(3,4) (Y>X)",
		"WDegRevLex(3,4) (X>Y)",
		"WDegRevLex(3,4) (Y>X)",
	}
	expectedOrd := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, -1, 1, -1, 1, -1, -1, -1, -1, -1},
		{-1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, -1, 1, -1, 1, -1, -1, -1, -1, -1},
		{1, -1, 1, 1, 1, 1, 1, -1, 1, -1},
	}

	for i, d := range degrees {
		for j, o := range orders {
			if res := o(d[0], d[1]); res != expectedOrd[i][j] {
				t.Errorf(
					"Wrong ordering of %v and %v using %s (got %d, expected %d)",
					d[0], d[1], ordStr[j], res, expectedOrd[i][j],
				)
			}
		}
	}
}
