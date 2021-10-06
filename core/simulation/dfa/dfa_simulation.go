package dfa

import (
	"github.com/flapflapio/simulator/core/simulation/machine"
	"github.com/flapflapio/simulator/core/types"
)

type DFA struct {
	machine      *machine.Machine
	currentState *machine.State
	input        string
}

func New(machine *machine.Machine, input string) *DFA {
	return &DFA{
		machine:      machine,
		input:        input,
		currentState: nil,
	}
}

// Perform a transition
func (dfa *DFA) Step() {
	panic("not implemented") // TODO: Implement
}

// Get the current status (state + other info) of a simulation
func (dfa *DFA) Stat() types.Report {
	panic("not implemented") // TODO: Implement
}

// Get the final result of your simulation.
// Returns a SimulationIncomplete error if the simulation is not done
func (dfa *DFA) Result() (types.Result, error) {
	panic("not implemented") // TODO: Implement
}

// Check if a simulation is finished
func (dfa *DFA) Done() bool {
	panic("not implemented") // TODO: Implement
}
