package machine

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleMachine = `
{
  "Start": "q0",
  "States": [
    { "Id": "q0", "Ending": false },
    { "Id": "q1", "Ending": true }
  ],
  "Transitions": [
    {
      "Start": "q0",
      "End": "q1",
      "Symbol": "a"
    }
  ]
}
`

func TestUnmarshalStates(t *testing.T) {
	assert.True(t, true) // TODO: implement method
}

func TestUnmarshalTransitions(t *testing.T) {
	assert.True(t, true) // TODO: implement method
}

func TestUnmarshalMachines(t *testing.T) {
	assert.True(t, true) // TODO: implement method
}

func TestMarshalingStates(t *testing.T) {
	assert.True(t, true) // TODO: implement method
}

func TestMarshalingTransitions(t *testing.T) {
	assert.True(t, true) // TODO: implement method
}

func TestMarshalingMachines(t *testing.T) {
}

func states(howMany int, endingMapper func(index int) bool) []State {
	var states []State
	mapper := endingMapper
	if mapper == nil {
		mapper = func(index int) bool { return index%2 == 0 }
	}
	for i := 0; i < howMany; i++ {
		states = append(states, State{fmt.Sprintf("q%v", i), mapper(i)})
	}
	return states
}

func transitions(howMany int, symbolMapper func(index int) string) []Transition {
	var transitions []Transition
	mapper := symbolMapper
	if mapper == nil {
		mapper = func(index int) string {
			if index%2 == 0 {
				return "a"
			} else {
				return "b"
			}
		}
	}
	for i := 0; i < howMany; i++ {
		transitions = append(transitions, Transition{Symbol: mapper(i)})
	}
	return transitions
}

func machinesWithDifferentStartingStates(howMany int) []*Machine {
	var machines []*Machine
	states := states(5, nil)
	transitions := transitions(4, nil)
	connectInStraightLine(&states, &transitions)
	for i := range states {
		machines = append(machines, &Machine{
			Start:       &states[i],
			States:      states,
			Transitions: transitions,
		})
	}
	return machines
}

func connectInStraightLine(states *[]State, transitions *[]Transition) {
	if len(*states) < 2 ||
		len(*transitions) < 1 ||
		len(*states) < len(*transitions)+1 {
		return
	}

	var prev *State
	var t *Transition
	for i := range *states {
		if prev != nil && t != nil {
			t = &(*transitions)[i-1]
			t.Start = prev
			t.End = &(*states)[i]
		}
		prev = &(*states)[i]
	}
}

func MktempFile(t *testing.T, testname, data string) *os.File {
	f, err := os.CreateTemp("./", fmt.Sprintf("%v-*", testname))
	assert.NoError(t, err, "Error occured creating temporary file")
	f.Write([]byte(data))
	return f
}
