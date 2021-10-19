package machine

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

const defaultSchema = "machine.schema.json"
const errMsg = "invalid %v, must be of type " +
	"`string`, " +
	"`map[string]interface{}`, " +
	"`[]byte`, " +
	"or `io.Reader`"

var cachedSchema map[string]interface{} = nil

// Loads a machine from a document of type `map[string]interface{}`, `string`,
// `[]byte`, or `io.Reader`. Uses the default schema.
func Load(document interface{}) (*Machine, error) {
	switch d := document.(type) {
	case map[string]interface{}:
		return loadMap(d)
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

// Loads a machine from a document of type `map[string]interface{}`, `string`,
// `[]byte`, or `io.Reader`. Uses the schema provided in the `schema` parameter,
// which has the same restrictions on type as the `document` parameter
func LoadWithSchema(document interface{}, schema interface{}) (*Machine, error) {
	switch d := document.(type) {
	case map[string]interface{}:
		return loadMapWithSchema(d, schema)
	case string:
		return loadFileWithSchema(d, document)
	case []byte:
		return loadBufferWithSchema(d, schema)
	case io.Reader:
		return loadReaderWithSchema(d, schema)
	default:
		return nil, fmt.Errorf(errMsg, "document")
	}
}

func ValidateJsonWithMachineSchema(
	document map[string]interface{},
) (*gojsonschema.Result, error) {
	schema, err := GetSchema()
	if err != nil {
		return nil, err
	}
	return ValidateJson(schema, document)
}

func ValidateJson(
	schema map[string]interface{},
	document map[string]interface{},
) (*gojsonschema.Result, error) {
	validator := gojsonschema.NewGoLoader(schema)
	validatee := gojsonschema.NewGoLoader(document)
	return gojsonschema.Validate(validator, validatee)
}

func ReadFileIntoMap(file string) (map[string]interface{}, error) {
	var loaded map[string]interface{}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		return nil, err
	}
	return loaded, nil
}

func loadBuffer(buf []byte) (*Machine, error) {
	var m map[string]interface{}
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}
	return loadMap(m)
}

func loadBufferWithSchema(buf []byte, schema interface{}) (*Machine, error) {
	var m map[string]interface{}
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}
	return loadMapWithSchema(m, schema)
}

func loadData(document map[string]interface{}) (*Machine, error) {
	var m Machine
	if err := loadStates(&m, document); err != nil {
		return nil, err
	}
	if err := loadTransitions(&m, document); err != nil {
		return nil, err
	}
	if err := loadStartState(&m, document); err != nil {
		return nil, err
	}
	return &m, nil
}

func loadStates(mach *Machine, document map[string]interface{}) error {
	states, ok := document["States"]
	if !ok {
		return fmt.Errorf("'States' key not found within document")
	}
	unknownStates, ok := states.([]interface{})
	if !ok {
		return fmt.Errorf("'States' key in document is not valid: %v", states)
	}
	for _, s := range unknownStates {
		ss, err := loadState(mach, s)
		if err != nil {
			return fmt.Errorf("invalid States set: %v", err)
		}
		mach.States = append(mach.States, ss)
	}
	return nil
}

func loadState(mach *Machine, unknown interface{}) (State, error) {
	m, ok := unknown.(map[string]interface{})
	if !ok {
		return State{}, fmt.Errorf("error casting unknown State to 'map', invalid State")
	}
	id, ok := m["Id"].(string)
	if !ok {
		return State{},
			fmt.Errorf("error casting 'Id' field of" +
				" unknown state to 'string', invalid State")
	}
	ending, ok := m["Ending"].(bool)
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

func loadTransitions(mach *Machine, document map[string]interface{}) error {
	transitions, ok := document["Transitions"]
	if !ok {
		return fmt.Errorf("'Transitions' key not found within document")
	}
	unknownTransitions, ok := transitions.([]interface{})
	if !ok {
		return fmt.Errorf("'Transitions' key in document is not valid: %v", transitions)
	}
	for _, t := range unknownTransitions {
		tt, err := loadTransition(mach, t)
		if err != nil {
			return fmt.Errorf("invalid Transitions set: %v", err)
		}
		mach.Transitions = append(mach.Transitions, tt)
	}
	return nil
}

func loadTransition(mach *Machine, unknown interface{}) (Transition, error) {
	m, ok := unknown.(map[string]interface{})
	if !ok {
		return Transition{},
			fmt.Errorf("error casting unknown Transition to 'map', invalid Transition")
	}
	symbol, ok := m["Symbol"].(string)
	if !ok {
		return Transition{},
			fmt.Errorf("error casting 'Symbol' field of unknown " +
				"Transition to 'map', invalid Transition")
	}
	t := Transition{
		Symbol: symbol,
	}

	start, err := findStateById(mach, m["Start"].(string))
	if err != nil {
		return Transition{}, err
	}

	end, err := findStateById(mach, m["End"].(string))
	if err != nil {
		return Transition{}, err
	}

	t.Start = start
	t.End = end

	return t, nil
}

func loadStartState(mach *Machine, document map[string]interface{}) error {
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

func findStateById(mach *Machine, id string) (*State, error) {
	for i, s := range mach.States {
		if s.Id == id {
			return &mach.States[i], nil
		}
	}
	return nil, fmt.Errorf("state with id '%v' was not found in machine", id)
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

func GetSchema() (map[string]interface{}, error) {
	if cachedSchema == nil {
		s, err := ReadFileIntoMap(defaultSchema)
		if err != nil {
			return nil, err
		}
		cachedSchema = s
	}
	return cachedSchema, nil
}

func parseSchema(schema interface{}) (map[string]interface{}, error) {
	switch s := schema.(type) {
	case string:
		ss, err := ReadFileIntoMap(s)
		if err != nil {
			return nil, err
		}
		return ss, nil
	case map[string]interface{}:
		return s, nil
	case []byte:
		var m map[string]interface{}
		err := json.Unmarshal(s, &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	case io.Reader:
		data, err := io.ReadAll(s)
		if err != nil {
			return nil, err
		}
		return parseSchema(data)
	default:
		ss, err := GetSchema()
		if err != nil {
			return nil, err
		}
		return ss, nil
	}
}

func loadFileWithSchema(file string, schema interface{}) (*Machine, error) {
	document, err := ReadFileIntoMap(file)
	if err != nil {
		return nil, err
	}
	return loadMapWithSchema(document, schema)
}

func loadMapWithSchema(
	document map[string]interface{},
	schema interface{},
) (*Machine, error) {
	s, err := parseSchema(schema)
	if err != nil {
		return nil, err
	}
	result, err := ValidateJson(document, s)
	if err != nil {
		return nil, err
	}
	if !result.Valid() {
		return nil, unifyErrors(result.Errors())
	}
	return loadData(document)
}

func loadMap(document map[string]interface{}) (*Machine, error) {
	s, err := GetSchema()
	if err != nil {
		return nil, err
	}
	return loadMapWithSchema(document, s)
}

func loadFile(file string) (*Machine, error) {
	document, err := ReadFileIntoMap(file)
	if err != nil {
		return nil, err
	}
	return loadMap(document)
}

func loadReader(document io.Reader) (*Machine, error) {
	buf, err := io.ReadAll(document)
	if err != nil {
		return nil, err
	}
	return loadBuffer(buf)
}

func loadReaderWithSchema(
	document io.Reader,
	schema interface{},
) (*Machine, error) {
	buf, err := io.ReadAll(document)
	if err != nil {
		return nil, err
	}
	return loadBufferWithSchema(buf, schema)
}
