package machine

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var basicGraph = (func() Graph {
	g := Graph{
		States: []State{
			{Id: "q0", Ending: false},
			{Id: "q1", Ending: true},
		},
	}
	g.Start = &g.States[0]
	g.Transitions = []Transition{
		{Start: &g.States[0], End: &g.States[1], Symbol: "a"},
		{Start: &g.States[1], End: &g.States[0], Symbol: "a"},
		{Start: &g.States[0], End: &g.States[0], Symbol: "b"},
		{Start: &g.States[1], End: &g.States[1], Symbol: "b"},
	}
	return g
})()

var testCasesFromSuccess = []struct {
	name   string
	params GraphParams
	graph  Graph
}{
	{
		name: "easy-success-case",
		params: GraphParams{
			Start: "q0",
			States: []State{
				{Id: "q0", Ending: false},
				{Id: "q1", Ending: true},
			},
			Transitions: []TransitionParams{
				{Start: "q0", End: "q1", Symbol: "a"},
				{Start: "q1", End: "q0", Symbol: "a"},
				{Start: "q0", End: "q0", Symbol: "b"},
				{Start: "q1", End: "q1", Symbol: "b"},
			},
		},
		graph: basicGraph,
	},
}

func TestFromSuccess(t *testing.T) {
	test := func(params GraphParams, graph Graph) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			g := From(params)
			assert.NotNil(t, g)
			assert.Equal(t, graph.Start, g.Start)
			assert.Equal(t, len(graph.States), len(g.States))
			assert.Equal(t, len(graph.Transitions), len(g.Transitions))
			for i, s := range graph.States {
				assert.Equal(t, s, g.States[i])
			}
			for i, tt := range graph.Transitions {
				assert.Equal(t, tt.Start.Id, g.Transitions[i].Start.Id)
				assert.Equal(t, tt.End.Id, g.Transitions[i].End.Id)
				assert.Equal(t, tt.Symbol, g.Transitions[i].Symbol)
			}
		}
	}

	for _, tc := range testCasesFromSuccess {
		t.Run(tc.name, test(tc.params, tc.graph))
	}
}

var testCasesFromFailure = []struct {
	name   string
	params GraphParams
}{
	{
		name:   "empty-params",
		params: GraphParams{},
	},
	{
		name: "bad-state-names",
		params: GraphParams{
			Start: "q0",
			States: []State{
				{Id: "q0", Ending: false},
				{Id: "asd", Ending: true},
			},
			Transitions: []TransitionParams{
				{Start: "q0", End: "q1", Symbol: "a"},
				{Start: "q1", End: "q0", Symbol: "a"},
				{Start: "q0", End: "q0", Symbol: "b"},
				{Start: "q1", End: "q1", Symbol: "b"},
			},
		},
	},
}

func TestFromFailure(t *testing.T) {
	test := func(params GraphParams) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			g := From(params)
			assert.Nil(t, g, "machine should not be built")
		}
	}

	for _, tc := range testCasesFromFailure {
		t.Run(tc.name, test(tc.params))
	}
}

var testCasesLoadMachine = []struct {
	name     string
	data     []byte
	expected Graph
}{
	{
		name: "easy-success-case",
		data: []byte(`
		{
			"$schema": "https://raw.githubusercontent.com/flapflapio/simulator/main/core/simulation/machine/machine.schema.json",
			"Type": "DFA",
			"Alphabet": "ab",
			"Start": "q0",
			"States": [
			  { "Id": "q0", "Ending": false },
			  { "Id": "q1", "Ending": true }
			],
			"Transitions": [
			  { "Start": "q0", "End": "q1", "Symbol": "a" },
			  { "Start": "q1", "End": "q0", "Symbol": "a" },
			  { "Start": "q0", "End": "q0", "Symbol": "b" },
			  { "Start": "q1", "End": "q1", "Symbol": "b" }
			]
		  }
		`),
		expected: basicGraph,
	},
}

// Tests the `machine.Load` function
func TestLoadMachine(t *testing.T) {
	testLoadingViaMap := func(
		data []byte,
		expected Graph,
	) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			var g map[string]interface{}
			err := json.Unmarshal(data, &g)
			assert.NoError(t, err)
			gg, err := Load(g)
			assert.NoError(t, err)
			assertMarshalablesEqual(t, &expected, gg)
		}
	}

	testLoadingViaTempFile := func(
		name string,
		data []byte,
		expected Graph,
	) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			f := MktempFile(t, name, data)
			assert.NoError(t, f.Close())
			defer os.Remove(f.Name())
			g, err := Load(f.Name())
			assert.NoError(t, err)
			assertMarshalablesEqual(t, &expected, g)
		}
	}

	testLoadingViaBytes := func(
		name string,
		data []byte,
		expected Graph,
	) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			g, err := Load(data)
			assert.NoError(t, err)
			assertMarshalablesEqual(t, &expected, g)
		}
	}

	testLoadingViaReader := func(
		name string,
		data []byte,
		expected Graph,
	) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			f := MktempFile(t, name, data)
			defer os.Remove(f.Name())
			defer f.Close()
			g, err := Load(f)
			assert.NoError(t, err)
			assertMarshalablesEqual(t, &expected, g)
		}
	}

	for _, tc := range testCasesLoadMachine {
		var (
			mapTestName    = fmt.Sprintf("%v_load-from-map", tc.name)
			fileTestName   = fmt.Sprintf("%v_load-from-file", tc.name)
			bytesTestName  = fmt.Sprintf("%v_load-from-bytes", tc.name)
			readerTestName = fmt.Sprintf("%v_load-from-reader", tc.name)
		)

		t.Run(mapTestName, testLoadingViaMap(tc.data, tc.expected))
		t.Run(fileTestName, testLoadingViaTempFile(fileTestName, tc.data, tc.expected))
		t.Run(bytesTestName, testLoadingViaBytes(bytesTestName, tc.data, tc.expected))
		t.Run(readerTestName, testLoadingViaReader(readerTestName, tc.data, tc.expected))
	}
}

func TestLoadMachineEmptyFile(t *testing.T) {
	const (
		loadingFromErrMsg = "loading from an empty %v should throw an error"
		emptyBufferErrMsg = "empty buffer"
	)

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "load-from-empty-reader",
			testFunc: func(t *testing.T) {
				t.Parallel()
				f := MktempFile(t, "load-from-empty-reader", []byte{})
				defer os.Remove(f.Name())
				defer f.Close()
				_, err := Load(f)
				assert.Error(t, err, fmt.Sprintf(loadingFromErrMsg, "reader"))
				assert.Contains(t, err.Error(), "empty buffer")
			},
		},
		{
			name: "load-from-empty-file",
			testFunc: func(t *testing.T) {
				t.Parallel()
				f := MktempFile(t, "load-from-empty-file", []byte{})
				defer os.Remove(f.Name())
				f.Close()
				_, err := Load(f.Name())
				assert.Error(t, err, fmt.Sprintf(loadingFromErrMsg, "file"))
				assert.Contains(t, err.Error(), emptyBufferErrMsg)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}
}

func MktempFile(t *testing.T, testname string, data []byte) *os.File {
	f, err := os.CreateTemp("./", fmt.Sprintf("%v-*", testname))
	assert.NoError(t, err, "Error occured creating temporary file")
	assert.NotNil(t, f)
	_, err = f.Write(data)
	assert.NoError(t, err)
	assert.NoError(t, f.Sync())
	_, err = f.Seek(0, 0)
	assert.NoError(t, err)
	return f
}

func assertMarshalablesEqual(t *testing.T, expected, actual Marshalable) {
	assert.Equal(t,
		expected.Json(), actual.Json(), "json value of graphs should be equal")
}
