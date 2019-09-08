package errors

import (
	"fmt"
	"testing"
)

var kinds = []Kind{
	Input,
	InputValue,
	InputIncompatible,
	InputTooLarge,
	ArithmeticIncompat,
	Parsing,
	Conversion,
	Overflow,
	Internal,
}

func TestIs(t *testing.T) {
	for _, k := range kinds {
		const op = "Testing"
		err := New(
			op, k,
			"Testing errors",
		)
		errWrap := Wrap(
			op, Inherit,
			err,
		)
		if !Is(k, err) {
			t.Errorf("Following error does not have correct kind: %q", err.Error())
		}
		for _, i := range kinds {
			if i == k {
				continue
			}
			if Is(i, err) {
				t.Errorf("Following error had unexpected kind: %q",
					err.Error())
			}
		}
		if !Is(k, errWrap) {
			t.Errorf("Following wrapped error had unexpected kind: %q",
				errWrap.Error())
		}
	}
}

func TestBuiltInError(t *testing.T) {
	builtIn := fmt.Errorf("A built-in error")
	errWrap := Wrap("Testing", Inherit, builtIn)

	for _, k := range kinds {
		if Is(k, builtIn) {
			t.Errorf("Built-in error has kind: %q", builtIn.Error())
		}
		if Is(k, errWrap) {
			t.Errorf("Wrapped built-in error has kind: %q", errWrap.Error())
		}
	}
}
