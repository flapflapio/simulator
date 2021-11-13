package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrNoTransition(t *testing.T) {
	err := thrower(ErrNoTransition)
	if !errors.Is(err, ErrNoTransition) {
		t.Fail()
	}
}

func TestErrSimulationIncomplete(t *testing.T) {
	err := thrower(ErrSimulationIncomplete)
	if !errors.Is(err, ErrSimulationIncomplete) {
		t.Fail()
	}
}

func thrower(err error) error {
	return fmt.Errorf("err: %w", err)
}
