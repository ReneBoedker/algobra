package basic

import (
	"testing"
)

func TestCeilLog(t *testing.T) {
	testPairs := [][2]uint{
		{0, 0},
		{1, 0},
		{2, 1},
		{431, 9},
		{999999, 20},
		{1<<40 - 1, 40},
	}
	for _, pair := range testPairs {
		if tmp := CeilLog(pair[0]); tmp != pair[1] {
			t.Errorf("CeilLog(%d)=%d (Expected %d)", pair[0], tmp, pair[1])
		}
	}
}
