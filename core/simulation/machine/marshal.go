package machine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

const errMsg = "invalid %v, must be of type " +
	"`string`, " +
	"`map[string]interface{}`, " +
	"`[]byte`, " +
	"or `io.Reader`"

var defaultSchema = []byte(SCHEMA)
var cachedSchema = struct {
	*sync.Mutex
	value map[string]interface{}
}{&sync.Mutex{}, nil}

func Load(document interface{}) (*Graph, error) {
	return LoadWithSchema(document, nil)
}

func LoadWithSchema(document interface{}, schema interface{}) (*Graph, error) {
	documentMap, err := LoadMap(document)
	if err != nil {
		return nil, err
	}

	schemaMap, err := schemaOrDefault(schema)
	if err != nil {
		return nil, err
	}

	validationResult, err := ValidateJson(documentMap, schemaMap)
	if err != nil {
		return nil, err
	}

	if !validationResult.Valid() {
		return nil, unifyErrors(validationResult.Errors())
	}

	return createGraph(documentMap)
}

func LoadMap(document interface{}) (map[string]interface{}, error) {
	switch d := document.(type) {
	case map[string]interface{}:
		return d, nil
	case string:
		return loadFile(d)
	case []byte:
		return loadBuffer(d)
	case io.Reader:
		return loadReader(d)
	default:
		return nil, fmt.Errorf(errMsg, "document")
	}
}

func GetSchema() map[string]interface{} {
	cachedSchema.Lock()
	defer cachedSchema.Unlock()
	if cachedSchema.value == nil {
		err := json.Unmarshal(defaultSchema, &cachedSchema.value)
		if err != nil {
			cachedSchema.value = nil
			return map[string]interface{}{}
		}
	}
	return cachedSchema.value
}

func ValidateJson(
	schema map[string]interface{},
	document map[string]interface{},
) (*gojsonschema.Result, error) {
	validator := gojsonschema.NewGoLoader(schema)
	validatee := gojsonschema.NewGoLoader(document)
	return gojsonschema.Validate(validator, validatee)
}

func loadFile(path string) (map[string]interface{}, error) {
	chugged, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadBuffer(chugged)
}

func loadReader(reader io.Reader) (map[string]interface{}, error) {
	chugged, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return loadBuffer(chugged)
}

func loadBuffer(buf []byte) (map[string]interface{}, error) {
	if len(buf) == 0 {
		return nil, errors.New("cannot load machine from an empty buffer")
	}
	var m map[string]interface{}
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func unifyErrors(errors []gojsonschema.ResultError) error {
	err := "("
	for i, e := range errors {
		if err += fmt.Sprintf("%v", e); i != len(errors)-1 {
			err += "... "
		} else {
			err += ")"
		}
	}
	return fmt.Errorf("JSON schema is not valid: %v", err)
}

func createGraph(document map[string]interface{}) (*Graph, error) {
	var g Graph
	if err := addStates(&g, document); err != nil {
		return nil, err
	}
	if err := addTransitions(&g, document); err != nil {
		return nil, err
	}
	if err := addStartState(&g, document); err != nil {
		return nil, err
	}
	return &g, nil
}

func addStates(mach *Graph, document map[string]interface{}) error {
	states, ok := document["States"]
	if !ok {
		return fmt.Errorf("'States' key not found within document")
	}
	unknownStates, ok := states.([]interface{})
	if !ok {
		return fmt.Errorf("'States' key in document is not valid: %v", states)
	}
	for _, s := range unknownStates {
		ss, err := addState(mach, s)
		if err != nil {
			return fmt.Errorf("invalid States set: %v", err)
		}
		mach.States = append(mach.States, ss)
	}
	return nil
}

func addState(mach *Graph, unknown interface{}) (State, error) {
	s, ok := unknown.(map[string]interface{})
	if !ok {
		return State{},
			fmt.Errorf("error casting unknown State to 'map', invalid State")
	}
	id, ok := s["Id"].(string)
	if !ok {
		return State{},
			fmt.Errorf("error casting 'Id' field of" +
				" unknown state to 'string', invalid State")
	}
	ending, ok := s["Ending"].(bool)
	if !ok {
		return State{},
			fmt.Errorf("error casting 'Ending' field of unknown " +
				"state to 'bool', invalid State")
	}
	return State{
		Id:     id,
		Ending: ending,
	}, nil
}

func addTransitions(mach *Graph, document map[string]interface{}) error {
	transitions, ok := document["Transitions"]
	if !ok {
		return fmt.Errorf("'Transitions' key not found within document")
	}
	unknownTransitions, ok := transitions.([]interface{})
	if !ok {
		return fmt.Errorf("'Transitions' key in document is not valid: %v", transitions)
	}
	for _, t := range unknownTransitions {
		tt, err := addTransition(mach, t)
		if err != nil {
			return fmt.Errorf("invalid Transitions set: %v", err)
		}
		mach.Transitions = append(mach.Transitions, tt)
	}
	return nil
}

func addTransition(mach *Graph, unknown interface{}) (Transition, error) {
	t, ok := unknown.(map[string]interface{})
	if !ok {
		return Transition{},
			fmt.Errorf("error casting unknown Transition to 'map', invalid Transition")
	}
	symbol, ok := t["Symbol"].(string)
	if !ok {
		return Transition{},
			fmt.Errorf("error casting 'Symbol' field of unknown " +
				"Transition to 'map', invalid Transition")
	}
	tt := Transition{
		Symbol: symbol,
	}

	start, err := findStateById(mach, t["Start"].(string))
	if err != nil {
		return Transition{}, err
	}

	end, err := findStateById(mach, t["End"].(string))
	if err != nil {
		return Transition{}, err
	}

	tt.Start = start
	tt.End = end

	return tt, nil
}

func addStartState(mach *Graph, document map[string]interface{}) error {
	state, ok := document["Start"]
	if !ok {
		return fmt.Errorf("'Start' key not found within document")
	}

	stateId, ok := state.(string)
	if !ok {
		return fmt.Errorf("'Start' key in document is not valid: %v", state)
	}

	foundState, err := findStateById(mach, stateId)
	if err != nil {
		return fmt.Errorf("state with id '%v' was not found in state set", stateId)
	}

	mach.Start = foundState
	return nil
}

func findStateById(mach *Graph, id string) (*State, error) {
	for i, s := range mach.States {
		if s.Id == id {
			return &mach.States[i], nil
		}
	}
	return nil, fmt.Errorf("state with id '%v' was not found in machine", id)
}

func schemaOrDefault(schema interface{}) (map[string]interface{}, error) {
	if schema != nil {
		return LoadMap(schema)
	} else {
		return GetSchema(), nil
	}
}
