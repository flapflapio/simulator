package automata

import (
	"errors"

	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/automata/dfa"
	"github.com/flapflapio/simulator/core/simulation/machine"
)

// Loads a machine from the given `document`
func Load(document interface{}) (simulation.Machine, error) {
	return LoadWithSchema(document, nil)
}

// Loads a machine from the given `document` using a specific schema
func LoadWithSchema(
	document interface{},
	schema interface{},
) (simulation.Machine, error) {
	documentMap, err := machine.LoadMap(document)
	if err != nil {
		return nil, err
	}
	t, err := extractType(documentMap)
	if err != nil {
		return nil, err
	}
	return createMachineOfType(t, documentMap, schema)
}

func Dump(mach simulation.Machine) string {
	return mach.Json()
}

func createMachineOfType(
	machineType string,
	document map[string]interface{},
	schema interface{},
) (simulation.Machine, error) {
	switch machineType {
	case machine.DFA:
		return dfa.LoadWithSchema(document, schema)
	case machine.NFA:
	case machine.PDA:
	case machine.TM:
	default:
		return dfa.LoadWithSchema(document, schema)
	}
	return nil, errors.New(
		"machine was not able to be created, unrecognized machine type")
}

func extractType(document map[string]interface{}) (string, error) {
	unknown, ok := document["Type"]
	if !ok {
		return "", errors.New("invalid document, field 'Type' is required")
	}
	t, ok := unknown.(string)
	if !ok {
		return "", errors.New("invalid document, field 'Type' should be a string")
	}
	return machine.ParseMachineType(t), nil
}
