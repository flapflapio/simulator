package schemacontroller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/flapflapio/simulator/core/simulation/automata/dfa"
	"github.com/flapflapio/simulator/core/simulation/machine"
	"github.com/obonobo/mux"
)

type assertion struct {
	method      string
	path        string
	status      int
	sendBody    *string
	receiveBody *string
}

var testCasesPostValidate = []struct {
	name    string
	status  int
	machine string
}{
	{
		name:    "valid-dfa",
		status:  http.StatusOK,
		machine: dfa.ODDA,
	},
	{
		name:    "empty-machine",
		status:  http.StatusUnprocessableEntity,
		machine: "",
	},
	{
		name:   "invalid-alphabet-not-a-string",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "DFA",
			"Alphabet": true,
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
		}`,
	},
	{
		name:   "invalid-type-not-a-string",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": 25,
			"Alphabet": "ab",
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
		}`,
	},
	{
		name:   "invalid-type-missing-field",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Alphabet": "ab",
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
		}`,
	},
	{
		name:   "invalid-type-not-a-machine",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "Not a machine",
			"Alphabet": "ab",
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
		}`,
	},
	{
		name:   "invalid-start-state-doesn't-match-pattern",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "DFA,
			"Alphabet": "ab",
			"Start": "doesn't match pattern",
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
		}`,
	},
	{
		name:   "invalid-start-state-not-a-string",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "DFA",
			"Alphabet": "ab",
			"Start": ["q0"],
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
		}`,
	},
	{
		name:   "invalid-states-not-an-array",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "DFA",
			"Alphabet": "ab",
			"Start": "q0",
			"States": 100,
			"Transitions": [
			  { "Start": "q0", "End": "q1", "Symbol": "a" },
			  { "Start": "q0", "End": "q0", "Symbol": "b" },
			  { "Start": "q1", "End": "q1", "Symbol": "b" },
			  { "Start": "q1", "End": "q0", "Symbol": "a" }
			]
		}`,
	},
	{
		name:   "invalid-states-duplicates",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "DFA",
			"Alphabet": "ab",
			"Start": "q0",
			"States": [
				{ "Id": "q0", "Ending": false },
				{ "Id": "q0", "Ending": false },
				{ "Id": "q1", "Ending": true }
			  ],,
			"Transitions": [
			  { "Start": "q0", "End": "q1", "Symbol": "a" },
			  { "Start": "q0", "End": "q0", "Symbol": "b" },
			  { "Start": "q1", "End": "q1", "Symbol": "b" },
			  { "Start": "q1", "End": "q0", "Symbol": "a" }
			]
		}`,
	},
	{
		name:   "invalid-transitions-missing-transitions",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "DFA",
			"Alphabet": "ab",
			"Start": "q0",
			"States": [
				{ "Id": "q0", "Ending": false },
				{ "Id": "q1", "Ending": true }
			  ],,
			"Transitions": [
			  { "Start": "q0", "End": "q1", "Symbol": "a" },
			  { "Start": "q1", "End": "q1", "Symbol": "b" },
			  { "Start": "q1", "End": "q0", "Symbol": "a" }
			]
		}`,
	},
	{
		name:   "invalid-transitions-duplicates",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "DFA",
			"Alphabet": "ab",
			"Start": "q0",
			"States": [
				{ "Id": "q0", "Ending": false },
				{ "Id": "q1", "Ending": true }
			  ],,
			"Transitions": [
			  { "Start": "q0", "End": "q1", "Symbol": "a" },
			  { "Start": "q0", "End": "q1", "Symbol": "a" },
			  { "Start": "q0", "End": "q0", "Symbol": "b" },
			  { "Start": "q1", "End": "q1", "Symbol": "b" },
			  { "Start": "q1", "End": "q0", "Symbol": "a" }
			]
		}`,
	},
	{
		name:   "invalid-transitions-wrong-type",
		status: http.StatusUnprocessableEntity,
		machine: `
		{
			"Type": "DFA",
			"Alphabet": "ab",
			"Start": "q0",
			"States": [
			  { "Id": "q0", "Ending": false },
			  { "Id": "q1", "Ending": true }
			],
			"Transitions": "asdasd"
		}`,
	},
}

func TestGetMachineSchema(t *testing.T) {
	prefix := "/some/path"
	router := mux.NewRouter()
	controller := WithPrefix(prefix)
	controller.Attach(router)

	body := string(machine.SCHEMA)
	assertEndpoint(t, router, assertion{
		method:      "GET",
		path:        prefix + "/machine.schema.json",
		status:      http.StatusOK,
		receiveBody: &body,
	})
}

func TestPostValidate(t *testing.T) {
	prefix := "/some/path"
	router := mux.NewRouter()
	controller := WithPrefix(prefix)
	controller.Attach(router)

	test := func(machine string, status int) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			assertEndpoint(t, router, assertion{
				method:   "POST",
				path:     prefix + "/validate",
				status:   status,
				sendBody: &machine,
			})
		}
	}

	for _, tc := range testCasesPostValidate {
		t.Run(tc.name, test(tc.machine, tc.status))
	}
}

func assertEndpoint(
	t *testing.T,
	r *mux.Router,
	expectation assertion,
) {
	receive := expectation.receiveBody
	send := expectation.sendBody
	if send == nil {
		send = new(string)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(
		expectation.method,
		expectation.path,
		bytes.NewReader([]byte(*send)))

	if err != nil {
		t.Fatalf("Expected request to be created properly, but got error: %v", err)
	}

	r.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != expectation.status {
		t.Fatalf(
			"Expected response status %v, but got %v",
			expectation.status,
			recorder.Result().StatusCode)
	}

	// Test recieve body, if it is intended to be asserted
	if receive != nil {
		var expected map[string]interface{}
		err := json.Unmarshal([]byte(*receive), &expected)
		if err != nil {
			t.Fatalf(
				"Expected expectation body to unmarshal correctly, got err: %v",
				err)
		}
		var actual map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &actual)
		if err != nil {
			t.Fatalf("Expected actual body to unmarshal correctly, got err: %v", err)
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected body to be '%v' but got '%v'", expected, actual)
		}
	}
}
