package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoPossibleTransition(t *testing.T) {
	e := NoTrans()
	assert.Equal(t, NO_POSSIBLE_TRANSITION, e.Error())
	e = NoPossibleTransition{}
	assert.Equal(t, "", e.Error())
}

func TestSimulationNotDone(t *testing.T) {
	e := NotDone()
	assert.Equal(t, SIMULATION_INCOMPLETE, e.Error())
	e = SimulationNotDone{}
	assert.Equal(t, "", e.Error())
}
