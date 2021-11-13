package app

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/flapflapio/simulator/core/controllers"
	"github.com/flapflapio/simulator/internal/simtest"
	"github.com/obonobo/mux"
)

var cfg = Config{
	Name:           simtest.StringPointer("Test Server"),
	Port:           8181,
	ReadTimeout:    60,
	WriteTimeout:   60,
	MaxHeaderBytes: 4096,
}

// Tests the healthcheck, method not allowed, and not found handlers in the
// default server
func TestBasicEndpoints(t *testing.T) {
	simtest.WaitForServerToStop(cfg.Port)
	srv := New(cfg)
	srv.Attach(LoggerAndRecovery)
	cancel := simtest.StartServer(t, srv)
	defer cancel()

	// Run a bunch of tests
	for _, tc := range []struct {
		name, route, body, method string
		status                    int
	}{
		{
			name:   "healthcheck",
			route:  "/healthcheck",
			method: "GET",
			status: http.StatusOK,
			body:   HEALTHCHECK_MESSAGE + "\n",
		},
		{
			name:   "method not allowed",
			route:  "/healthcheck",
			method: "POST",
			status: http.StatusMethodNotAllowed,
			body:   string(methodNotAllowedMessage("POST")),
		},
		{
			name:   "not found",
			route:  "/not/found",
			method: "GET",
			status: http.StatusNotFound,
			body:   string(notFoundMessage()),
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := simtest.DoRequest(t, simtest.MustCreateRequest(t,
				tc.method, simtest.AppUrl(cfg.Port)+tc.route, nil))

			defer r.Body.Close()

			bod, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			if actual := string(bod); tc.body != actual {
				t.Fatalf(
					"Expected response body to be '%v' but got '%v'",
					tc.body, actual)
			}
		})
	}
}

// Tries to attach controllers and middleware to the server
func TestAttach(t *testing.T) {

	mid := createMockMiddleware(3)
	mids := []*mockMiddleware{createMockMiddleware(1), createMockMiddleware(2)}
	midFuncs := make([]Middleware, len(mids))
	for i := range mids {
		midFuncs[i] = mids[i].middleware
	}

	ctl := &mockController{route: "/controller4"}
	ctls := []*mockController{
		{route: "/controller1"},
		{route: "/controller2"},
		{route: "/controller3"},
	}
	ctlsCast := make([]controllers.Controller, len(ctls))
	for i := range ctls {
		ctlsCast[i] = ctls[i]
	}

	srv := New(cfg)
	srv.Attach(midFuncs, ctlsCast, mid.middleware, ctl)

	cancel := simtest.StartServer(t, srv)
	defer cancel()

	ping := func(ctl *mockController) {
		r := simtest.DoRequest(t, simtest.MustCreateRequest(t,
			"GET", simtest.AppUrl(cfg.Port)+ctl.route, nil))

		if r.StatusCode != http.StatusOK {
			t.Fatalf("Controller %v: expected status code %v but got %v",
				ctl.route,
				http.StatusOK,
				r.StatusCode)
		}
	}

	ctls = append(ctls, ctl)
	for _, c := range ctls {
		ping(c)
	}

	for _, c := range ctls {
		if c.called != 1 {
			t.Fatalf("Controller route '%v' was not called exactly once", c.route)
		}
	}

	for _, m := range append(mids, mid) {
		if expected := len(ctls) + 1; m.called != expected {
			t.Fatalf("Middleware #%v: expected %v calls, but got %v",
				m.id, expected, m.called)
		}
	}
}

func TestStopServerNotStarted(t *testing.T) {
	srv := New(cfg)
	err := srv.Stop()
	if !errors.Is(err, ErrServerNotStarted) {
		t.Fatalf(""+
			"server.Stop() should throw an error if the server is not started, "+
			"but got err %v", err)
	}
}

func TestServerStartsWithError(t *testing.T) {
	srv := New(cfg)
	cancel := simtest.StartServer(t, srv)
	defer cancel()

	srv2 := New(cfg)
	cancel2 := srv2.RunWithCancel()
	defer cancel2()

	// Wait for the error to appear on the server-internal error channel
	errIsOnErrChannel := make(chan struct{})
	go func() {
		srv2.ec <- <-srv2.ec
		errIsOnErrChannel <- struct{}{}
	}()
	<-errIsOnErrChannel

	if err := <-srv2.Wait(); err == nil ||
		!regexp.MustCompile(`bind: address already in use`).Match([]byte(err.Error())) {
		t.Fatalf("Expected to receive a started up from running the server, but got %v", err)
	}
}

// A mock conforming to the controllers.Controller interface
type mockController struct {
	route  string
	called int
}

func (c *mockController) Attach(router *mux.Router) {
	router.Name(c.route).Methods("GET").Path(c.route).HandlerFunc(c.handler)
}

func (c *mockController) handler(rw http.ResponseWriter, r *http.Request) {
	c.called++
	rw.WriteHeader(http.StatusOK)
}

type mockMiddleware struct {
	middleware Middleware
	id, called int
}

// Creates an ad-hoc middleware mock
func createMockMiddleware(id int) *mockMiddleware {
	m := mockMiddleware{id: id}
	m.middleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			m.called++
			next.ServeHTTP(rw, r)
		})
	}
	return &m
}
