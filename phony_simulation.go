package main

import (
	"fmt"

	"github.com/flapflapio/simulator/core/types"
)

// A "phony" simulation that accepts any input
type PhonySimulation struct {
	path  []string
	input string
	i     int
}

func (ps *PhonySimulation) Step() {
	ps.path = append(ps.path, fmt.Sprintf("q%v", ps.i))
	ps.input = ps.input[1:]
	ps.i++
}

func (ps *PhonySimulation) Stat() types.Report {
	return types.Report{}
}

func (ps *PhonySimulation) Result() (types.Result, error) {
	return types.Result{
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