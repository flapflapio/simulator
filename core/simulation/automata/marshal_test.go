package automata

import (
	"fmt"
	"testing"

	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/automata/dfa"
	"github.com/flapflapio/simulator/core/simulation/machine"
	"github.com/stretchr/testify/assert"
)

type TestCaseDFAMarshaling struct {
	success     bool
	marshaled   string
	unmarshaled simulation.Machine
}

var machineStrings = map[string]TestCaseDFAMarshaling{
	"OddA": TestCaseDFAMarshaling{
		success:   true,
		marshaled: dfa.ODDA,
		unmarshaled: dfa.From(dfa.DFAParams{
			Alphabet: "ab",
			GraphParams: machine.GraphParams{
				Start: "q0",
				States: []machine.State{
					{Id: "q0", Ending: false},
					{Id: "q1", Ending: true},
				},
				Transitions: []machine.TransitionParams{
					{Start: "q0", End: "q1", Symbol: "a"},
					{Start: "q0", End: "q0", Symbol: "b"},
					{Start: "q1", End: "q1", Symbol: "b"},
					{Start: "q1", End: "q0", Symbol: "a"},
				},
			},
		}),
	},
}

func TestMarshaling(t *testing.T) {
	// Test function
	test := func(tc TestCaseDFAMarshaling) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			m, err := Load([]byte(tc.marshaled))
			if tc.success {
				assert.NoError(t, err, "machine should be built properly")
			} else {
				assert.Error(t, err, "machine should not be build properly")
			}
			assertMachinesEqual(t, tc.unmarshaled, m)
		}
	}

	// Run the tests
	for name, machine := range machineStrings {
		t.Run(fmt.Sprintf("TestCase[name:%v]", name), test(machine))
	}
}

func assertMachinesEqual(t *testing.T, m1, m2 simulation.Machine) {
	assert.Equal(t, m1.Json(), m2.Json(), "json value of maps should be equal")
}
