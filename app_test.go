package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/flapflapio/simulator/core/app"
	"github.com/flapflapio/simulator/core/simulation/automata/dfa"
	"github.com/flapflapio/simulator/internal/simtest"
	"github.com/obonobo/mux"
)

// Whether or not these tests should be run
var integration = flag.Bool(
	"skip-integration",
	false,
	"enables an integration test suite that spins up the entire app server")

var overrideCfg = app.Config{
	Name:           simtest.StringPointer("Test Server"),
	Port:           8181,
	ReadTimeout:    60,
	WriteTimeout:   60,
	MaxHeaderBytes: 4096,
}

func handleFlags(t *testing.T) {
	if *integration {
		t.Skip("Integration tests have not been enabled with the `-integration` flag")
	}
}

// Integration tests for the /simulate route
func TestSimulate(t *testing.T) {
	handleFlags(t)
	cancel := startServer(t)
	defer resetServer()
	defer cancel()

	type TestCase struct{ name, route, body string }
	tests := []TestCase{
		{
			name:  "POST success",
			route: "/simulate?tape=aaba",
			body:  `{"Accepted":true,"Path":["q0","q1","q0","q0","q1"],"RemainingInput":""}`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.Post(
				fmt.Sprintf("%v%v", simtest.AppUrl(overrideCfg.Port), tc.route),
				"application/json",
				bytes.NewBufferString(dfa.ODDA))
			assertError(t, err)
			assertBodyMatches(t, r, tc.body)
		})
	}
}

func TestHealthcheckCLISuccess(t *testing.T) {
	handleFlags(t)
	cancel := startServer(t)
	defer resetServer()
	defer cancel()
	exitCode, output, err := simtest.SubprocessWithEnv(
		map[string]string{"PORT": fmt.Sprintf("%v", overrideCfg.Port)},
		"go", "run", "app.go", "-health")
	t.Logf("OUTPUT: %v", output)
	assertError(t, err)
	assertStatus(t, 0, exitCode)
	assertOutput(t, app.HEALTHCHECK_MESSAGE, output)
}

func TestHealthcheckCLIServerNotRunning(t *testing.T) {
	handleFlags(t)
	exitCode, output, err := simtest.SubprocessWithEnv(
		map[string]string{"PORT": fmt.Sprintf("%v", overrideCfg.Port)},
		"go", "run", "app.go", "-health")
	t.Logf("OUTPUT: %v", output)
	assertError(t, err)
	assertStatus(t, 1, exitCode)
	assertOutput(t, regexp.MustCompile(fmt.Sprintf(""+
		`Error: Get "http:\/\/localhost:%v\/healthcheck": `+
		`dial tcp (\d+\.\d+\.\d+\.\d+|\[.*:.*:.*\]):%v: connect: connection refused`+
		"\nexit status 1", overrideCfg.Port, overrideCfg.Port)), output)
}

func TestHealthcheckCLIStatusCodeNot200(t *testing.T) {
	handleFlags(t)
	srv = app.NewNoHandlers(mux.NewRouter(), overrideCfg)
	srv.Router.
		Methods("GET").
		Path("/healthcheck").
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

	cancel :=  simtest.StartServer(t, srv)
	defer resetServer()
	defer cancel()

	exitCode, output, err := simtest.SubprocessWithEnv(
		map[string]string{"PORT": fmt.Sprintf("%v", overrideCfg.Port)},
		"go", "run", "app.go", "-health")

	t.Logf("OUTPUT: %v", output)
	assertError(t, err)
	assertStatus(t, 1, exitCode)
	assertOutput(t, fmt.Sprintf(""+
		`Get "http://localhost:%v/healthcheck": `+
		"500 Internal Server Error"+
		"\nexit status 1", overrideCfg.Port), output)
}

// Tests the CLI using in-process techniques (as opposed to creating a
// subprocess to run the tool)
func TestHealthcheckCLIInProcessSuccess(t *testing.T) {
	handleFlags(t)
	cancel := startServer(t)
	defer resetServer()
	defer cancel()

	old := *healthcheck
	defer func() { *healthcheck = old }()
	*healthcheck = true
	out, close := simtest.MockStdoutStderr(t)
	exitCode := healthcheckMode()
	close()

	assertStatus(t, 0, exitCode)
	assertOutput(t, app.HEALTHCHECK_MESSAGE, <-out)
}

func TestHealthcheckCLIInProcessServerNotRunning(t *testing.T) {
	handleFlags(t)

	old := *healthcheck
	defer func() { *healthcheck = old }()
	*healthcheck = true
	out, close := simtest.MockStdoutStderr(t)
	exitCode := healthcheckMode()
	close()

	assertStatus(t, 1, exitCode)
	assertOutput(t, regexp.MustCompile(fmt.Sprintf(""+
		`Error: Get "http:\/\/localhost:%v\/healthcheck": `+
		`dial tcp (\d+\.\d+\.\d+\.\d+|\[.*:.*:.*\]):%v: connect: connection refused`,
		overrideCfg.Port, overrideCfg.Port)),
		<-out)
}

func TestHealthcheckCLIInProcessStatusCodeNot200(t *testing.T) {
	handleFlags(t)
	srv = app.NewNoHandlers(mux.NewRouter(), overrideCfg)
	srv.Router.
		Methods("GET").
		Path("/healthcheck").
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

	cancel := simtest.StartServer(t, srv)

	defer resetServer()
	defer cancel()

	old := *healthcheck
	defer func() { *healthcheck = old }()
	*healthcheck = true
	out, close := simtest.MockStdoutStderr(t)
	exitCode := healthcheckMode()
	close()

	assertStatus(t, 1, exitCode)
	assertOutput(t, fmt.Sprintf(""+
		`Get "http://localhost:%v/healthcheck": `+
		"500 Internal Server Error", overrideCfg.Port),
		<-out)
}

func TestHealthcheckCLIInProcessNoHealthcheckFlag(t *testing.T) {
	handleFlags(t)
	old := *healthcheck
	defer func() { *healthcheck = old }()
	*healthcheck = false
	exitCode := healthcheckMode()
	assertStatus(t, -1, exitCode)
}

func TestMainServerAlreadyRunning(t *testing.T) {
	handleFlags(t)

	cancel := startServer(t)
	old := *healthcheck
	*healthcheck = false
	defer func() { *healthcheck = old }()

	errc := make(chan error, 1)
	go func() {
		defer func() {
			err := recover()
			if e, ok := err.(error); ok {
				errc <- e
			} else {
				errc <- nil
			}
			cancel()
		}()
		main()
	}()

	err := <-errc
	t.Log(err)
	if err == nil {
		t.Fatalf(""+
			"Expected to receive an error from running main "+
			"(port %v should be occupied in another goroutine), "+
			"but got no error: %v",
			overrideCfg.Port,
			err)
	}

	expected := fmt.Sprintf("listen tcp :%v: bind: address already in use", overrideCfg.Port)
	if actual := err.Error(); actual != expected {
		t.Fatalf("Expected error message '%v', but got '%v'", expected, actual)
	}
}

// Asserts that the exit status of the command matches
func assertStatus(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Fatalf("Expected command exit code %v but got %v", expected, actual)
	}
}

// Asserts a string match. The expected param can be a string or a regex pattern
func assertOutput(t *testing.T, expected interface{}, actual string) {
	trim := func(s string) string { return strings.Trim(s, " \n") }
	actual = trim(actual)
	switch exp := expected.(type) {
	case string:
		exp = trim(exp)
		if exp != actual {
			t.Fatalf("Expected command output '%v' but got '%v'", exp, actual)
		}
	case *regexp.Regexp:
		if !exp.MatchString(actual) {
			t.Fatalf(
				"Expected command output to match regex '%s' but got '%s'",
				exp.String(), actual)
		}
	}
}

// Fails your test if the error is not nil
func assertError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%v", err)
	}
}

// Asserts that the response body matches the expected string
func assertBodyMatches(t *testing.T, r *http.Response, expected string) {
	trim := func(s string) string { return strings.Trim(string(s), " \n") }
	bod, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if actual, expected := trim(string(bod)), trim(expected); actual != expected {
		t.Fatalf(
			"Response body does not match, expected '%v' but got '%v'",
			expected,
			actual)
	}
}

func startServer(t *testing.T) context.CancelFunc {
	resetServer()
	return simtest.StartServer(t, srv)
}

func resetServer() {
	srv = app.New(overrideCfg)
	setupServer()
}
