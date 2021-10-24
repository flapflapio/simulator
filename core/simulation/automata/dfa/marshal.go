package dfa

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/flapflapio/simulator/core/simulation/machine"
)

func Load(document interface{}) (*DFA, error) {
	return LoadWithSchema(document, nil)
}

func LoadWithSchema(document interface{}, schema interface{}) (*DFA, error) {
	dfa := &DFA{}
	documentMap, err := machine.LoadMap(document)
	if err != nil {
		return nil, err
	}
	dfa.Graph, err = machine.LoadWithSchema(documentMap, schema)
	if err != nil {
		return nil, err
	}
	err = addAlphabet(dfa, documentMap)
	if err != nil {
		return nil, err
	}
	err = checkThatStatesHaveATransitionForEverySymbol(dfa)
	if err != nil {
		return nil, err
	}
	err = checkThatTransitionSymbolsMatchAlphabet(dfa)
	if err != nil {
		return nil, err
	}
	return dfa, nil
}

func checkThatTransitionSymbolsMatchAlphabet(dfa *DFA) error {
	for _, t := range dfa.Transitions {
		found := false
		for _, s := range dfa.Alphabet {
			if string(s) == t.Symbol {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf(
				"DFA is invalid, %v contains a symbol not present in the alphabet", t)
		}
	}
	return nil
}

func checkThatStatesHaveATransitionForEverySymbol(dfa *DFA) error {
	for _, state := range dfa.Graph.States {
		for _, symbol := range dfa.Alphabet {
			sym := string(symbol)
			found := false
			for _, transition := range dfa.Transitions {
				if transition.Start.Id == state.Id && transition.Symbol == sym {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf(
					"DFA is invalid, %v is missing a transition for symbol %v",
					state, sym)
			}
		}
	}
	return nil
}

func addAlphabet(dfa *DFA, document map[string]interface{}) error {
	unknown, ok := document["Alphabet"]
	if !ok {
		inferAlphabet(dfa)
		return nil
	}

	alphabet, ok := unknown.(string)
	if !ok {
		return errors.New("'Alphabet' field in json document is not valid " +
			"- it should be a string")
	}
	dfa.Alphabet = alphabet
	return nil
}

func inferAlphabet(dfa *DFA) {
	alphabet := ""
	appendSymbol := func(symbol string) {
		alreadyPresent := false
		for _, r := range alphabet {
			if strconv.QuoteRune(r) == symbol {
				alreadyPresent = true
				break
			}
		}
		if !alreadyPresent {
			alphabet += symbol
		}
	}
	for _, t := range dfa.Graph.Transitions {
		appendSymbol(string(t.Symbol[0]))
	}
}
