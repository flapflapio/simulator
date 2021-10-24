package simulation

import (
	"fmt"
)

// A "phony" simulation that accepts any input
type PhonySimulation struct {
	path  []string
	input string
	i     int
}

func NewPhonySimulation(input string) *PhonySimulation {
	return &PhonySimulation{
		path:  make([]string, len(input)),
		input: input,
		i:     0,
	}
}

func (ps *PhonySimulation) Step() {
	ps.path[ps.i] = fmt.Sprintf("q%v", ps.i)
	ps.input = ps.input[1:]
	ps.i++
}

func (ps *PhonySimulation) Stat() Report {
	return Report{}
}

func (ps *PhonySimulation) Result() (Result, error) {
	return Result{
		Accepted: true,
		Path:     ps.path,
	}, nil
}

func (ps *PhonySimulation) Done() bool {
	return len(ps.input) == 0
}

func (ps *PhonySimulation) Kill() error {
	return nil
}

type PhonyMachine struct{}

func (m *PhonyMachine) Simulate(input string) Simulation {
	return NewPhonySimulation(input)
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
