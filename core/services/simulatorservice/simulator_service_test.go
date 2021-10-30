package simulatorservice

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/flapflapio/simulator/core/simulation"
)

type testCase struct {
	name    string
	machine string
	inputs  []string
}

var testCases = []testCase{
	{
		name: "valid",
		inputs: []string{
			"aaba",
		},
	},
}

func TestStart(t *testing.T) {
	test := func(tc testCase) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			mach := simulation.NewPhonyMachine()
			service := New()
			ids := startInputs(t, service, mach, &tc)

			assertSimulations(t, service, ids, tc.inputs)
			assertIdsInOrder(t, ids)
			assertMethodWasCalled(t, "Simulate", len(ids), mach.MethodsCalled.Simulate)
			deleteAllSimulations(t, service, ids)
		}
	}
	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestEndNilSimulation(t *testing.T) {
	service := New()
	service.End(-1)
}

func deleteAllSimulations(t *testing.T, service *SimulatorService, ids []int) {
	for _, id := range ids {
		service.End(id)
		got := service.Get(id)
		if got != nil {
			t.Fatalf("Expected simulation with id '%v' to have ended, "+
				"but it is still present in the SimulatorService", id)
		}
	}
}

func assertSimulations(
	t *testing.T,
	service *SimulatorService,
	ids []int,
	inputs []string,
) {
	for i, id := range ids {
		expectedResult := expectedResult(inputs[i])
		res := simulation.ResultOf(service.Get(id))
		if !reflect.DeepEqual(expectedResult, *res) {
			t.Fatalf(
				"Expected simulation result to be %v but got %v",
				expectedResult, res)
		}
	}
}

func expectedResult(input string) simulation.Result {
	expectedResult := simulation.Result{
		Accepted:       true,
		Path:           make([]string, len(input)),
		RemainingInput: "",
	}
	for i := range input {
		expectedResult.Path[i] = fmt.Sprintf("q%v", i)
	}
	return expectedResult
}

func startInputs(
	t *testing.T,
	service *SimulatorService,
	mach simulation.Machine,
	tc *testCase,
) []int {
	ids := make([]int, len(tc.inputs))
	for i, in := range tc.inputs {
		id, err := service.Start(mach, in)
		if err != nil {
			t.Fatalf("Expected SimulatorService.Start not"+
				" to throw an error, but got: %v", err)
		}
		ids[i] = id
	}
	return ids
}

func assertIdsInOrder(t *testing.T, ids []int) {
	for i, id := range ids {
		if i != id {
			t.Fatalf("Expected simulation id to be %v, but got %v", i, id)
		}
	}
}

func assertMethodWasCalled(t *testing.T, methodName string, expected, actual int) {
	if expected != actual {
		t.Fatalf(""+
			"Expected Machine.%v to be called %v "+
			"times but was only called %v times",
			methodName, expected, actual)
	}
}
