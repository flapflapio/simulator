package simulationcontroller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/automata/dfa"
	"github.com/obonobo/mux"
)

var defaultService = newMockmockSimulatorService(0)

type testCaseDoSimulation struct {
	name     string
	status   int
	method   string
	tape     string
	machine  string
	response string
	service  *mockSimulatorService
}

var testCasesDoSimulation = []testCaseDoSimulation{
	{
		name:    "valid",
		service: defaultService,
		method:  "POST",
		tape:    "aaba",
		machine: dfa.ODDA,
		status:  http.StatusOK,
		response: `{
			"Accepted": true,
			"Path": ["q0", "q1", "q2", "q3", "q4"],
			"RemainingInput": ""
		}`,
	},
	{
		name:     "invalid-doesn't-load",
		service:  defaultService,
		method:   "POST",
		tape:     "aaba",
		machine:  dfa.ODDA,
		status:   http.StatusUnprocessableEntity,
		response: INVALID_MACHINE_MSG,
	},
	{
		name:     "invalid-no-tape-provided",
		service:  defaultService,
		method:   "POST",
		tape:     "",
		machine:  dfa.ODDA,
		status:   http.StatusBadRequest,
		response: PLEASE_PROVIDE_A_TAPE_MSG,
	},
	{
		name:     "invalid-simulator-start-fails",
		service:  newMockmockSimulatorService(FAIL_ON_START),
		method:   "POST",
		tape:     "aaba",
		machine:  dfa.ODDA,
		status:   http.StatusInternalServerError,
		response: FAILED_TO_CREATE_A_NEW_SIMULATION,
	},
	{
		name:     "invalid-simulator-result-fails",
		service:  newMockmockSimulatorService(FAIL_ON_RESULT),
		method:   "POST",
		tape:     "aaba",
		machine:  dfa.ODDA,
		status:   http.StatusInternalServerError,
		response: FAILED_TO_OBTAIN_RESULTS_OF_SIMULATION,
	},
	{
		name:     "invalid-json-serialization-fails",
		service:  newMockmockSimulatorService(RESULT_UNSERIALIZABLE),
		method:   "POST",
		tape:     "aaba",
		machine:  dfa.ODDA,
		status:   http.StatusInternalServerError,
		response: FAILED_TO_CREATE_A_RESPONSE,
	},
}

func TestDoSimulation(t *testing.T) {
	router := mux.NewRouter()

	test := func(tc testCaseDoSimulation) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()

			controller := New(tc.service)
			controller.Attach(router)

			recorder := httptest.NewRecorder()

			tt := ""
			if tc.tape != "" {
				tt = fmt.Sprintf("?tape=%v", tc.tape)
			}

			req, err := http.NewRequest(tc.method,
				fmt.Sprintf("simulate%v", tt),
				bytes.NewBufferString(tc.machine))

			if err != nil {
				t.Fatalf("Expected no error while building request, but got: %v", err)
			}

			router.ServeHTTP(recorder, req)
			assertStatusCode(t, tc.status, recorder)
			assertResponse(t, tc.response, recorder.Body.String())
			assertMockService(t, tc.service, 1, 1, 1)
		}
	}

	for _, tc := range testCasesDoSimulation {
		t.Run(tc.name, test(tc))
	}
}

func assertStatusCode(
	t *testing.T,
	expected int,
	recorder *httptest.ResponseRecorder,
) {
	if expected != recorder.Result().StatusCode {
		t.Fatalf("Expected status code %v but got %v",
			expected, recorder.Result().Status)
	}
}

func assertResponse(t *testing.T, expected, actual string) {
	var exp map[string]interface{}
	err := json.Unmarshal([]byte(expected), &exp)
	if err != nil {
		t.Fatalf("Expected no error while unmarshaling 'expected' but got %v", err)
	}

	var act map[string]interface{}
	err = json.Unmarshal([]byte(actual), &act)
	if err != nil {
		t.Fatalf("Expected no error while unmarshaling 'actual' but got %v", err)
	}

	if !reflect.DeepEqual(exp, act) {
		t.Fatalf("Expected response to be '%v' but got '%v'", exp, act)
	}
}

// Asserts that the specified methods have been called
func assertMockService(t *testing.T, m *mockSimulatorService, Start, Get, End int) {
	msg := "Expected SimulatorService.%v " +
		"to be called exactly %v times, " +
		"but was called %v times"

	if Start != m.methodsCalled.Start {
		t.Errorf(msg, "Start", Start, m.methodsCalled.Start)
	}
	if Get != m.methodsCalled.Get {
		t.Errorf(msg, "Get", Get, m.methodsCalled.Get)
	}
	if End != m.methodsCalled.End {
		t.Errorf(msg, "End", End, m.methodsCalled.End)
	}
}

const (
	FAIL_ON_START         = 1 << iota
	FAIL_ON_RESULT        = 1 << iota
	RESULT_UNSERIALIZABLE = 1 << iota
)

type mockSimulatorService struct {
	nextId int
	sims   map[int]simulation.Simulation

	methodsCalled struct {
		Start int
		Get   int
		End   int
	}

	failOnStart          bool
	failOnResult         bool
	resultUnserializable bool
}

func newMockmockSimulatorService(flags int) *mockSimulatorService {
	return &mockSimulatorService{
		sims:                 map[int]simulation.Simulation{},
		failOnStart:          flags&FAIL_ON_START == FAIL_ON_START,
		failOnResult:         flags&FAIL_ON_RESULT == FAIL_ON_RESULT,
		resultUnserializable: flags&RESULT_UNSERIALIZABLE == RESULT_UNSERIALIZABLE,
	}
}

func (s *mockSimulatorService) Start(
	machine simulation.Machine,
	input string,
) (id int, err error) {
	if s.failOnStart {
		return 0, errors.New("mock failure")
	}
	s.methodsCalled.Start++
	i := s.nextId
	s.nextId++
	mockMachine := &simulation.PhonyMachine{FailOnResult: s.failOnResult}
	s.sims[i] = mockMachine.Simulate(input)
	return i, nil
}

func (s *mockSimulatorService) Get(simulationId int) simulation.Simulation {
	s.methodsCalled.Get++
	return s.sims[simulationId]
}

func (s *mockSimulatorService) End(simulationId int) error {
	s.methodsCalled.End++
	sim := s.sims[simulationId]
	if sim == nil {
		return fmt.Errorf("simulation with id '%v' does not exist", simulationId)
	}
	delete(s.sims, simulationId)
	return nil
}
