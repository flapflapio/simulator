package dfa

import (
	"github.com/flapflapio/simulator/core/errors"
	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/machine"
)

type DFASimulation struct {
	machine      *DFA
	currentState *machine.State
	input        string
	path         []string
	rejected     bool
}

// Perform a transition
func (dfa *DFASimulation) Step() {
	if dfa.Done() {
		return
	}
	dfa.logState()
	dfa.takeNextTransition()
	if dfa.Done() {
		dfa.logState()
	}
}

// Get the current status (state + other info) of a simulation
func (dfa *DFASimulation) Stat() simulation.Report {
	return simulation.Report{
		Result: simulation.Result{
			Accepted:       dfa.isAccepted(),
			Path:           dfa.path,
			RemainingInput: dfa.input,
		},
	}
}

// Get the final result of your simulation.
// Returns a SimulationIncomplete error if the simulation is not done
func (dfa *DFASimulation) Result() (simulation.Result, error) {
	if !dfa.Done() {
		return simulation.Result{}, errors.ErrSimulationIncomplete
	}
	return simulation.Result{
		Accepted:       dfa.isAccepted(),
		Path:           dfa.path,
		RemainingInput: dfa.input,
	}, nil
}

// Check if a simulation is finished
func (dfa *DFASimulation) Done() bool {
	return dfa.rejected || len(dfa.input) == 0
}

func (dfa *DFASimulation) takeNextTransition() {
	if dfa.rejected {
		return
	}
	next, err := dfa.nextTransition()
	if err != nil {
		dfa.rejected = true
		return
	}
	dfa.takeTransition(next)
}

func (dfa *DFASimulation) takeTransition(t machine.Transition) {
	dfa.currentState = t.End
	dfa.input = dfa.input[1:]
}

func (dfa *DFASimulation) nextTransition() (machine.Transition, error) {
	for _, t := range dfa.machine.Transitions {
		if dfa.shouldTakeTransition(t) {
			return t, nil
		}
	}
	return machine.Transition{}, errors.ErrNoTransition
}

func (dfa *DFASimulation) shouldTakeTransition(t machine.Transition) bool {
	return !dfa.rejected &&
		len(dfa.input) > 0 &&
		len(t.Symbol) > 0 &&
		dfa.currentState == t.Start &&
		dfa.input[0] == t.Symbol[0]
}

func (dfa *DFASimulation) isAccepted() bool {
	return !dfa.rejected &&
		len(dfa.input) == 0 &&
		dfa.currentState.Ending
}

// Appends the current state of the DFA onto the path
func (dfa *DFASimulation) logState() {
	dfa.path = append(dfa.path, dfa.currentState.Id)
}
