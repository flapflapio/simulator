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
	FailOnResult         bool
	ResultUnserializable bool
}

func NewPhonySimulation(input string, failOnResult, resultUnserializable bool) *PhonySimulation {
	return &PhonySimulation{
		Path:                 make([]string, len(input)),
		Input:                input,
		I:                    0,
		FailOnResult:         failOnResult,
		ResultUnserializable: resultUnserializable,
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
	res := Result{
		Accepted: true,
		Path:     ps.Path,
	}
	if ps.ResultUnserializable {
		res = Result{
			Path: []string{""},
		}
	}

	if ps.FailOnResult {
		return Result{}, errors.New("mock failure")
	}
	return res, nil
}

func (ps *PhonySimulation) Done() bool {
	return len(ps.Input) == 0
}

func (ps *PhonySimulation) Kill() error {
	return nil
}

type PhonyMachine struct {
	Wrapped              Machine
	FailOnResult         bool
	ResultUnserializable bool
	MethodsCalled        struct {
		Simulate int
		Json     int
		JsonMap  int
		String   int
	}
}

func NewPhonyMachine() *PhonyMachine {
	return &PhonyMachine{}
}

func NewPhonyMachineWrapper(wrapping Machine) *PhonyMachine {
	return &PhonyMachine{Wrapped: wrapping}
}

func (m *PhonyMachine) Simulate(input string) Simulation {
	m.MethodsCalled.Simulate++
	if m.Wrapped != nil {
		return m.Wrapped.Simulate(input)
	}
	return NewPhonySimulation(input, m.FailOnResult, m.ResultUnserializable)
}

func (m *PhonyMachine) Json() string {
	m.MethodsCalled.Json++
	return ""
}

func (m *PhonyMachine) JsonMap() map[string]interface{} {
	m.MethodsCalled.JsonMap++
	return map[string]interface{}{}
}

func (m *PhonyMachine) String() string {
	m.MethodsCalled.String++
	return "PhonyMachine[...]"
}
