package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestInput struct {
	str      string
	expected Result
}

var testInputs = []struct {
	name    string
	machine Machine
	inputs  []TestInput
}{
	{
		name:    "phony-machine",
		machine: &PhonyMachine{},
		inputs: func() []TestInput {
			res := make([]TestInput, 20)
			for i := 0; i < cap(res); i++ {
				res[i].expected.Path = make([]string, i)
				res[i].expected.Accepted = true
				for j := 0; j < i; j++ {
					res[i].str += "a"
					res[i].expected.Path[j] = fmt.Sprintf("q%v", j)
				}
			}
			return res
		}(),
	},
}

func TestRunToCompletion(t *testing.T) {
	test := func(machine Machine, input TestInput) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			sim := machine.Simulate(input.str)
			RunToCompletion(sim)
			res, err := sim.Result()
			assert.NoError(t, err)
			assert.Equal(t, input.expected.Path, res.Path)
		}
	}

	for _, tc := range testInputs {
		for _, in := range tc.inputs {
			t.Run(
				fmt.Sprintf("RunToCompletion[%v,%v]", tc.name, in.str),
				test(tc.machine, in))
		}
	}
}

func TestResultOf(t *testing.T) {
	test := func(machine Machine, input TestInput) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			sim := machine.Simulate(input.str)
			res := ResultOf(sim)
			assert.NotNil(t, res)
			assert.Equal(t, input.expected, *res)
		}
	}

	for _, tc := range testInputs {
		for _, in := range tc.inputs {
			t.Run(
				fmt.Sprintf("ResultOf[%v,%v]", tc.name, in.str),
				test(tc.machine, in))
		}
	}
}
