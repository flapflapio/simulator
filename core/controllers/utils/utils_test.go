package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/obonobo/mux"
	"github.com/stretchr/testify/assert"
)

var testCasesCompact = []struct {
	name     string
	data     string
	expected string
}{
	{
		name: "flat-json",
		data: `
		{
			"field1": "value1",
			"field2": 20,
			"field3": true
		}
		`,
		expected: `{"field1":"value1","field2":20,"field3":true}`,
	},
	{
		name: "nested-json",
		data: `
		{
			"field1": "value1",
			"field2": 20,
			"field3": {
				"nest1": ["yo", "bro"]
			}
		}
		`,
		expected: `{"field1":"value1","field2":20,"field3":{"nest1":["yo","bro"]}}`,
	},
}

func TestCompact(t *testing.T) {
	test := func(data, expected string) func(*testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, expected, string(Compact(data)))
		}
	}

	for _, tc := range testCasesCompact {
		t.Run(tc.name, test(tc.data, tc.expected))
	}
}

func TestCreateSubrouter(t *testing.T) {
	subrouter := CreateSubrouter(mux.NewRouter(), "/some/path")
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/some/path/route", nil)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	subrouter.Path("/route").HandlerFunc(func(
		rw http.ResponseWriter,
		r *http.Request,
	) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("good"))
	})

	subrouter.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf(
			"Expected route to be matched (200 OK), but got status %v",
			recorder.Result().Status)
	}

	if b := recorder.Body.String(); b != "good" {
		t.Fatalf(
			`Expected route to be matched (body="good"), but got body "%v"`, b)
	}
}
