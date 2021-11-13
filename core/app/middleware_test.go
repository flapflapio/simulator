package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flapflapio/simulator/internal/simtest"
	"github.com/obonobo/mux"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name, route, origin string
		origins             []string
	}{
		{
			name:    "wildcard",
			route:   "/",
			origin:  "http://localhost:3000",
			origins: []string{"*"},
		},
		{
			name:    "localhost",
			route:   "/some/path",
			origin:  "http://localhost:3000",
			origins: []string{"http://localhost:3000"},
		},
		{
			name:   "web app",
			route:  "/some/path",
			origin: "https://machinist.flapflap.io",
			origins: []string{
				"http://localhost:3000",
				"https://localhost:3000",
				"https://machinist.flapflap.io",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc(tc.route, func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusOK)
			})

			req := simtest.MustCreateRequest(t, "GET", tc.route, nil)
			req.Header.Add("Origin", tc.origin)

			handler := CORS(tc.origins...)(router)
			handler.ServeHTTP(recorder, req)

			// Assertions
			assertRequestStatus(t, recorder, http.StatusOK)
			if len(tc.origins) > 1 {
				assertResponseHeader(t, recorder, "Access-Control-Allow-Origin", tc.origin)
				assertResponseHeader(t, recorder, "Vary", "Origin")
			} else if len(tc.origins) == 1 {
				assertResponseHeader(t, recorder, "Access-Control-Allow-Origin", tc.origins[0])
			}
		})
	}
}

func TestTrim(t *testing.T) {
	somePath := "/some/path/"

	type testcase struct{ name, input, output string }
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
	req := simtest.MustCreateRequest(t, "GET", "/", nil)
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
			req := simtest.MustCreateRequest(t, "GET", tc.requestPath, nil)
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

func assertResponseHeader(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	key string,
	value string,
) {
	v, ok := recorder.Result().Header[key]
	if !ok {
		t.Fatalf("Expected response header map to contain key '%v' "+
			"with value '%v', but no such header was found. Response headers: %v",
			key, value, recorder.Result().Header)
	} else if len(v) == 0 {
		t.Fatalf("Expected response header map to contain key '%v' "+
			"with value '%v', but got an empty map for this key. Response headers: %v",
			key, value, recorder.Result().Header)
	} else if v[0] != value {
		t.Fatalf("Expected response header map to contain key '%v' "+
			"with value '%v', but got value '%v'",
			key, value, v[0])
	}
}
