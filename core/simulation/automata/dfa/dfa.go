package dfa

import (
	"encoding/json"
	"fmt"

	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/machine"
)

type DFA struct {
	*machine.Graph
	Alphabet string
}

func (d *DFA) Simulate(input string) simulation.Simulation {
	return &DFASimulation{
		machine:      d,
		currentState: d.Start,
		input:        input,
		path:         []string{},
		rejected:     false,
	}
}

func (d *DFA) String() string {
	return fmt.Sprintf(
		"DFA[Alphabet:%v Start:%v States:%v Transitions:%v]",
		d.Alphabet,
		d.Start.Id,
		d.States,
		d.Transitions)
}

func (d *DFA) Json() string {
	m := d.JsonMap()
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}

func (d *DFA) JsonMap() map[string]interface{} {
	g := d.Graph.JsonMap()
	g["Type"] = machine.DFA
	g["Alphabet"] = d.Alphabet
	return g
}
