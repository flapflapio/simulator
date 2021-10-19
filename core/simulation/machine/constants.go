package machine

import "strings"

// Machine types
const (
	// Deterministic Finite Automaton
	DFA = "DFA"

	// Non-Deterministic Finite Automaton
	NFA = "NFA"

	// Turing Machine
	TM = "TM"

	// Pushdown Automaton
	PD = "PD"
)

func ParseMachineType(machineType string) string {
	switch strings.ToLower(machineType) {
	case "d", "dfa", "deterministic finite automaton":
		return DFA
	case "n", "nfa", "non-deterministic finite automaton":
		return NFA
	case "p", "pd", "pushdown automaton":
		return PD
	case "t", "tm", "turingmachine", "turing machine":
		return TM
	}
	return ""
}
