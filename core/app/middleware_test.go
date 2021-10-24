package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/obonobo/mux"
)

func TestTrim(t *testing.T) {
	type testcase struct{ name, input, output string }

	somePath := "/some/path/"
	tests := []testcase{
		{
			name:   "no-change",
			input:  somePath,
			output: somePath,
		},
		{
			name:   "no-trailing-slash",
			input:  "/some/path",
			output: "/some/path",
		},
		{
			name:   "no-leading-slash",
			input:  "some/path/",
			output: "some/path/",
		},
	}

	// Add some test cases for extra slashes
	for i := 1; i <= 20; i++ {
		tests = append(tests, testcase{
			name:   fmt.Sprintf("extra-slashes_count=%v", i),
			input:  somePath + strings.Repeat("/", i),
			output: somePath,
		})
	}
	tests = append(tests, testcase{
		name:   fmt.Sprintf("extra-slashes_count=%v", 99999),
		input:  somePath + strings.Repeat("/", 999999),
		output: somePath,
	})

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if out := Trim(tc.input); tc.output != out {
				t.Errorf("For input '%v', expected output '%v' but got '%v'",
					tc.input, tc.output, out)
			}
		})
	}
}

func TestMiddlewareLoggingAndRecovery(t *testing.T) {
	router := mux.NewRouter()
	recorder := httptest.NewRecorder()
	req := createEmptyRequest(t, "GET", "/")
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		panic("Something went wrong inside a handler")
	})
	LoggerAndRecovery(router).ServeHTTP(recorder, req)
	assertRequestStatus(t, recorder, http.StatusInternalServerError)
}

func TestMiddlewareTrimTrailingSlashes(t *testing.T) {
	type testcase struct{ name, path, requestPath string }
	tests := []testcase{
		{
			name:        "plain-route",
			path:        "/",
			requestPath: "/",
		},
		{
			name:        "plain-route-2",
			path:        "/some/path/",
			requestPath: "/some/path/",
		},
		{
			name:        "no-trailing-slashes",
			path:        "/some/path",
			requestPath: "/some/path",
		},
	}

	for i := 1; i <= 20; i++ {
		tests = append(tests, testcase{
			name:        fmt.Sprintf("extra-slashes_count=%v", i),
			path:        "/some/path/",
			requestPath: "/some/path/" + strings.Repeat("/", i),
		})
		tests = append(tests, testcase{
			name:        fmt.Sprintf("with-query-string_count=%v", i),
			path:        "/some/path/",
			requestPath: "/some/path/" + strings.Repeat("/", i) + "?key=value&key2=value2",
		})
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			router := mux.NewRouter()
			router.HandleFunc(tc.path, func(rw http.ResponseWriter, r *http.Request) {})
			recorder := httptest.NewRecorder()
			req := createEmptyRequest(t, "GET", tc.requestPath)
			TrimTrailingSlash(router).ServeHTTP(recorder, req)
			assertRequestStatus(t, recorder, http.StatusOK)
		})
	}
}

func assertRequestStatus(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	expected int,
) {
	if recorder.Result().StatusCode != expected {
		t.Errorf("Expected response status code '%v' but got '%v'",
			expected, recorder.Result().StatusCode)
	}
}

func createEmptyRequest(t *testing.T, method, path string) *http.Request {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Logf("Got an error while creating request: %v", err)
		t.Fatal("There should be no error while createing the request")
	}
	return req
}
