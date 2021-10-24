package dfa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/machine"
	"github.com/stretchr/testify/assert"
)

type TestCaseOddA struct {
	str      string
	accepted bool
}

// MACHINES
var machines = struct {
	oddA string
}{
	oddA: `
	{
		"Start": "q0",
		"States": [
		  { "Id": "q0", "Ending": false },
		  { "Id": "q1", "Ending": true }
		],
		"Transitions": [
		  { "Start": "q0", "End": "q1", "Symbol": "a" },
		  { "Start": "q0", "End": "q0", "Symbol": "b" },
		  { "Start": "q1", "End": "q1", "Symbol": "b" },
		  { "Start": "q1", "End": "q0", "Symbol": "a" }
		]
	  }
	`,
}

// Some inputs for testing the odd accepting machine
func dataProviderOddA() []TestCaseOddA {
	cases := []TestCaseOddA{
		// ACCEPTED
		{"abaa", true},
		{"aaba", true},
		{"ababa", true},
		{"babababb", true},
		{"bbbbaaa", true},
		{"bbbbaabbaaabbaababa", true},
		{strings.Repeat("ab", 1000) + "a", true},

		// REJECTED
		{"abaaba", false},
		{"aabaa", false},
		{"abaabwposa", false},
		{"babaxubabba", false},
		{"bbbbaaaa||||||!?#?!", false},
		{"babbbaabbaaabbaababa", false},
		{strings.Repeat("ab", 10000), false},
	}
	cases = append(cases, aaa("odd")("accepted")(100)...)
	cases = append(cases, aaa("even")("rejected")(100)...)
	return cases
}

// Tests a DFA that accepts strings containing an odd number of a's
func TestMachineOddA(t *testing.T) {
	// Setup
	m := createMachine(t, machines.oddA)

	// Testing function
	test := func(str string, accepted bool) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			res := simulation.ResultOf(m.Simulate(str))
			assert.Equal(t, accepted, res.Accepted)
		}
	}

	// Run tests
	for _, tc := range dataProviderOddA() {
		t.Run(
			fmt.Sprintf("TestCase[str:'%v',accepted:%v]", tc.str, tc.accepted),
			test(tc.str, tc.accepted))
	}
}

func createMachine(t *testing.T, fromString string) *DFA {
	m, err := machine.Load([]byte(fromString))
	assert.NoError(t, err, "machine should build okay")
	return &DFA{
		Graph:    m,
		Alphabet: "",
	}
}

// A function for generating repeating strings of a's e.g. "aaaaaaaaa".
// Usage: aaa("even")("accepted")(25), aaa("odd")("rejected")(25)
func aaa(evenOrOdd string) func(accepted string) func(times int) []TestCaseOddA {
	return func(accepted string) func(times int) []TestCaseOddA {
		return func(times int) []TestCaseOddA {
			var cases []TestCaseOddA
			add := ""
			if evenOrOdd == "odd" {
				add = "a"
			}
			for i := 0; i < times; i++ {
				cases = append(cases, TestCaseOddA{
					str:      add + strings.Repeat("aa", i),
					accepted: accepted == "accepted",
				})
			}
			return cases
		}
	}
}
