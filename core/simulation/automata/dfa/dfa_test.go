package dfa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/flapflapio/simulator/core/simulation"
	"github.com/stretchr/testify/assert"
)

const (
	machineShouldBuildOkay = "machine should build okay"
)

type MachineOddA struct {
	machine string
	inputs  func() []TestCaseOddA
}

type TestCaseOddA struct {
	str      string
	accepted bool
}

type MachineInvalidAlphabet struct {
	machine string
	inputs  func() []TestCaseInvalidAlphabet
}

type TestCaseInvalidAlphabet struct {
	alphabet string
	valid    bool
}

// MACHINES
var machines = struct {
	OddA            MachineOddA
	InvalidAlphabet MachineInvalidAlphabet
}{
	OddA: MachineOddA{
		machine: ODDA,
		inputs: func() []TestCaseOddA {
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
		},
	},

	InvalidAlphabet: MachineInvalidAlphabet{
		machine: `
		{
			"Type": "DFA",
			"Alphabet": "%v",
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
		inputs: func() []TestCaseInvalidAlphabet {
			return []TestCaseInvalidAlphabet{
				{alphabet: "ab", valid: true},
				{alphabet: "ba", valid: true},
				{alphabet: "", valid: false},
				{alphabet: "a", valid: false},
				{alphabet: "abc", valid: false},
				{alphabet: "def", valid: false},
			}
		},
	},
}

// Tests a DFA that accepts strings containing an odd number of a's
func TestMachineOddA(t *testing.T) {
	// Setup
	m := createMachine(t, machines.OddA.machine)

	// Testing function
	test := func(tc TestCaseOddA) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			res := simulation.ResultOf(m.Simulate(tc.str))
			assert.Equal(t, tc.accepted, res.Accepted)
		}
	}

	// Run tests
	for _, tc := range machines.OddA.inputs() {
		t.Run(
			fmt.Sprintf("TestCase[str:%v,accepted:%v]", tc.str, tc.accepted),
			test(tc))
	}
}

func TestMachineInvalidAlphabet(t *testing.T) {
	// Testing function
	test := func(tc TestCaseInvalidAlphabet) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			m := fmt.Sprintf(machines.InvalidAlphabet.machine, tc.alphabet)
			_, err := Load([]byte(m))
			if tc.valid {
				assert.NoError(t, err, machineShouldBuildOkay)
			} else {
				assert.Error(t, err, "machine should be invalid, unmarshaling should fail")
				assert.Contains(t, err.Error(), "DFA is invalid")
			}
		}
	}

	// Run tests
	for _, tc := range machines.InvalidAlphabet.inputs() {
		t.Run(
			fmt.Sprintf("TestCase[alphabet:%v,valid:%v]", tc.alphabet, tc.valid),
			test(tc))
	}
}

func createMachine(t *testing.T, fromString string) simulation.Machine {
	m, err := Load([]byte(fromString))
	assert.NoError(t, err, machineShouldBuildOkay)
	return m
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
