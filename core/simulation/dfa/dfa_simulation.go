package dfa

import (
	"github.com/flapflapio/simulator/core/simulation/machine"
	"github.com/flapflapio/simulator/core/types"
)

type DFA struct {
	machine      *machine.Machine
	currentState *machine.State
	input        string
	path         []string
}

func New(machine *machine.Machine, input string) *DFA {
	return &DFA{
		machine:      machine,
		input:        input,
		currentState: machine.Start,
	}
}

// Perform a transition
func (dfa *DFA) Step() {
	var transition machine.Transition
	for _, t := range dfa.machine.Transitions {
		if dfa.currentState == t.Start {
			transition = t
			break
		}
	}
	dfa.path = append(dfa.path, dfa.currentState.Id)
	dfa.currentState = transition.End
	dfa.input = dfa.input[1:]
}

// Get the current status (state + other info) of a simulation
func (dfa *DFA) Stat() types.Report {
	return types.Report{
		Accepted: len(dfa.input) == 0 && dfa.currentState.Ending,
		Path:     dfa.path,
	}
}

// Get the final result of your simulation.
// Returns a SimulationIncomplete error if the simulation is not done
func (dfa *DFA) Result() (types.Result, error) {
	return types.Result{
		Accepted: len(dfa.input) == 0 && dfa.currentState.Ending,
		Path:     dfa.path,
	}, nil
}

// Check if a simulation is finished
func (dfa *DFA) Done() bool {
	return len(dfa.input) == 0
}
