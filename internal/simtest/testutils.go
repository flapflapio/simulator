package simtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"testing"
	"time"
)

const (
	// An app url for use in testing
	appUrlFmt = "http://localhost:%v"

	// The full URL for the /healthcheck route of the app
	HealthcheckFmt = appUrlFmt + "/healthcheck"
)

// This library allows only one server to be running at a time. The mutex is
// locked when the StartServer method is called and unlocked when the cancel()
// function that is returned by StartServer is called. You must call the
// cancel() function whenever you use this library.
var thereCanBeOnlyOne = new(sync.Mutex)

func AppUrl(port int) string {
	return fmt.Sprintf(appUrlFmt, port)
}

// Polls the given route of the app until one of the following conditions: the
// server responds with a matching status code (-1 or no statusCodes for any
// status code), an error is emitted from errc, or 10 seconds have passed.
func WaitForServerToStart(route string, errc <-chan error, statusCodes ...int) error {
	maxWait := 10 * time.Second
	good := make(chan struct{}, 1)

	codes := make(map[int]struct{}, len(statusCodes)+1)
	if len(statusCodes) == 0 {
		codes[-1] = struct{}{}
	}
	for _, c := range statusCodes {
		codes[c] = struct{}{}
	}

	// Poll the healthcheck route until the server is responsive
	go func() {
		for {
			r, err := http.Get(route)
			if err != nil {
				continue
			}
			_, ok1 := codes[-1]
			_, ok2 := codes[r.StatusCode]
			if ok1 || ok2 {
				good <- struct{}{}
				break
			}
		}
	}()

	select {
	case e := <-errc:
		return e
	case <-good:
		return nil
	case <-time.After(maxWait):
		return fmt.Errorf(
			"timeout (%v) while waiting for server to start", maxWait)
	}
}

// Polls a TCP port until a connection error is received. Optional timeout can
// be specified.
//
// The only error that can be returned by this function is a timeout error if
// the port is still responsive after the provided timeout duration.
func WaitForServerToStop(port int, timeout ...time.Duration) error {
	maxWait := 10 * time.Second
	if len(timeout) > 0 {
		maxWait = timeout[0]
	}

	good := make(chan struct{}, 1)

	// Try to connect to the port to see if it is in use
	go func() {
		for {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%v", port), time.Second)
			if err != nil {
				// We are waiting for the port to be unresponsive, so we can
				// exit here if there is an error connecting to this port
				good <- struct{}{}
				break
			}
			if conn != nil {
				conn.Close()
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case <-good:
		return nil
	case <-time.After(maxWait):
		return fmt.Errorf(
			"timeout (%v) while waiting for server to start", maxWait)
	}
}

// Starts the server and waits for it to respond to healthchecks via
// WaitForServerToStart
func StartServer(
	t *testing.T,
	srv interface {
		RunWithCancel() context.CancelFunc
		Wait() <-chan error
		Port() int
	},
) context.CancelFunc {
	thereCanBeOnlyOne.Lock()
	port := srv.Port()
	cancel := srv.RunWithCancel()
	cancelFunc := func() {
		defer thereCanBeOnlyOne.Unlock()
		defer func() {
			err := recover()
			errr, ok := err.(error)

			// We are ignoring double mutex unlock
			if err != nil && (!ok || errr.Error() != "sync: unlock of unlocked mutex") {
				panic(err)
			}
		}()
		cancel()
		err := WaitForServerToStop(port)
		if err != nil {
			panic(err)
		}
	}

	err := WaitForServerToStart(fmt.Sprintf(HealthcheckFmt, port), srv.Wait())
	if err != nil {
		cancelFunc()
		t.Fatalf("Got an error during server startup: %v", err)
		return func() {}
	}
	return cancelFunc
}

// Consumes the stdout and stderr of the current process, dumping them as a
// single string on the output chan. Must call close() when you are done
// printing or the output chan will hang waiting for more output.
func MockStdoutStderr(t *testing.T) (output <-chan string, close func()) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	out := make(chan string, 1)
	ready, done := make(chan struct{}, 1), make(chan struct{}, 1)
	go func() {
		old := os.Stdout
		olde := os.Stderr
		defer func() {
			r.Close()
			os.Stdout = old
			os.Stderr = olde
			done <- struct{}{}
		}()

		os.Stdout = w
		os.Stderr = w
		ready <- struct{}{}

		var buf bytes.Buffer
		io.Copy(&buf, r)
		io.Copy(old, bytes.NewBuffer(buf.Bytes()))
		out <- buf.String()
	}()

	<-ready
	return out, func() {
		w.Close()
		<-done
	}
}

// Runs the specified command and returns the exit code, output (combined stdout
// and stderr), and an error if the command failed to complete (failed to return
// ANY exit code). A command returning a non-zero exit code is not considered an
// error by this function.
func Subprocess(cmd string, args ...string) (exitCode int, output string, err error) {
	return SubprocessWithEnv(nil, cmd, args...)
}

// Runs the specified command, with the given environment variables, and returns
// the exit code, output (combined stdout and stderr), and an error if the
// command failed to complete (failed to return ANY exit code). A command
// returning a non-zero exit code is not considered an error by this function.
func SubprocessWithEnv(
	env map[string]string, cmd string, args ...string,
) (exitCode int, output string, err error) {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	for k, v := range env {
		c.Env = append(c.Env, fmt.Sprintf("%v=%v", k, v))
	}

	out, err := c.CombinedOutput()
	if e := new(exec.ExitError); errors.As(err, &e) {
		status, ok := e.Sys().(syscall.WaitStatus)
		if !ok {
			return 0, string(out), fmt.Errorf("command failed to return an exit code")
		}
		return status.ExitStatus(), string(out), nil
	}

	if err != nil {
		return 0, string(out), fmt.Errorf(
			"command did not complete properly (no exit code was returned): %w", err)
	}

	return 0, string(out), nil
}

// Returns a pointer to a string variable with the provided content
func StringPointer(s string) *string {
	ss := s
	return &ss
}

// Creates a request using http.NewRequest, handles errors by failing the test
func MustCreateRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Expected request to be created successfully: %v", err)
	}
	return req
}

// Executes the request, fails your test if the request fails
func DoRequest(t *testing.T, r *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
