package machine

import "strings"

// Machine types
const (
	// Deterministic Finite Automaton
	DFA = "DFA"

	// Non-Deterministic Finite Automaton
	NFA = "NFA"

	// Pushdown Automaton
	PDA = "PDA"

	// Turing Machine
	TM = "TM"
)

func ParseMachineType(machineType string) string {
	switch strings.ToLower(machineType) {
	case "d", "dfa", "deterministic finite automaton":
		return DFA
	case "n", "nfa", "non-deterministic finite automaton":
		return NFA
	case "p", "pd", "pda", "pushdown automaton":
		return PDA
	case "t", "tm", "turingmachine", "turing machine":
		return TM
	}
	return ""
}
