package simulation

import (
	"errors"
	"fmt"
)

// A "phony" simulation that accepts any input
type PhonySimulation struct {
	Path          []string
	Input         string
	I             int
	MethodsCalled struct {
		Step   int
		Stat   int
		Result int
		Done   int
		Kill   int
	}
	FailOnResult bool
}

func NewPhonySimulation(input string, failOnResult bool) *PhonySimulation {
	return &PhonySimulation{
		Path:         make([]string, len(input)),
		Input:        input,
		I:            0,
		FailOnResult: failOnResult,
	}
}

func (ps *PhonySimulation) Step() {
	ps.Path[ps.I] = fmt.Sprintf("q%v", ps.I)
	ps.Input = ps.Input[1:]
	ps.I++
}

func (ps *PhonySimulation) Stat() Report {
	return Report{}
}

func (ps *PhonySimulation) Result() (Result, error) {
	if ps.FailOnResult {
		return Result{}, errors.New("mock failure")
	}
	return Result{
		Accepted: true,
		Path:     ps.Path,
	}, nil
}

func (ps *PhonySimulation) Done() bool {
	return len(ps.Input) == 0
}

func (ps *PhonySimulation) Kill() error {
	return nil
}

type PhonyMachine struct {
	FailOnResult bool
}

func (m *PhonyMachine) Simulate(input string) Simulation {
	return NewPhonySimulation(input, m.FailOnResult)
}

func (m *PhonyMachine) Json() string {
	return ""
}

func (m *PhonyMachine) JsonMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (m *PhonyMachine) String() string {
	return "PhonyMachine[...]"
}
