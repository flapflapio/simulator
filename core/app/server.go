package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flapflapio/simulator/core/controllers"
	"github.com/obonobo/mux"
)

const (
	HEALTHCHECK_MESSAGE = "All good in the hood"
)

var (
	ErrServerNotStarted = fmt.Errorf("server not started")
)

type Server struct {
	Name       string
	Router     *mux.Router
	Config     Config
	Middleware []Middleware
	ec         chan error
	srv        *http.Server
}

func New(config Config) *Server {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)

	router.
		Methods("GET").
		Name("Healthcheck").
		Path("/healthcheck").
		HandlerFunc(healthcheck)

	return NewNoHandlers(router, config)
}

func NewNoHandlers(router *mux.Router, config Config) *Server {
	s := &Server{
		Router: router,
		Config: config,
		ec:     make(chan error, 1),
	}

	if config.Name != nil {
		s.Name = *config.Name
	}

	return s
}

// Runs the server asychronously
func (s *Server) Run() {
	log.Printf("Starting %v", s.Name)
	log.Printf("Listening on %v\n", s.Config.Port)

	s.srv = &http.Server{
		Addr:           fmt.Sprintf(":%v", s.Config.Port),
		Handler:        s.applyMiddleware(),
		ReadTimeout:    time.Second * time.Duration(s.Config.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(s.Config.WriteTimeout),
		MaxHeaderBytes: s.Config.MaxHeaderBytes,
	}

	log.Printf(
		"Server config: ReadTimeout=%v, WriteTimeout=%v, MaxHeaderBytes=%v",
		s.srv.ReadTimeout,
		s.srv.WriteTimeout,
		s.srv.MaxHeaderBytes)

	go func() { s.ec <- s.srv.ListenAndServe() }()
}

// Runs the server asynchronously, with a cancel function to shutdown the server
func (s *Server) RunWithCancel() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	s.Run()
	go func() {
		select {
		case <-ctx.Done():
			s.Stop()
			s.ec <- ctx.Err()
		case err := <-s.ec:
			s.ec <- err // Emit the error again
		}
	}()
	return cancel
}

func (s *Server) Wait() <-chan error {
	return s.ec
}

func (s *Server) Stop() error {
	if s.srv == nil {
		return ErrServerNotStarted
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

// `stuff` is of type `controllers.Controller` or `Middleware` or a slice of
// either
func (s *Server) Attach(stuff ...interface{}) {
	for _, x := range stuff {
		switch xx := x.(type) {
		case controllers.Controller:
			s.AttachController(xx)
		case []controllers.Controller:
			s.AttachControllers(xx...)
		case Middleware:
			s.AttachMiddleware(xx)
		case []Middleware:
			s.AttachMiddlewares(xx...)
		}
	}
}

func (s *Server) AttachController(controller controllers.Controller) {
	controller.Attach(s.Router)
}

func (s *Server) AttachControllers(controllers ...controllers.Controller) {
	for _, c := range controllers {
		s.AttachController(c)
	}
}

func (s *Server) AttachMiddleware(middleware Middleware) {
	s.Middleware = append(s.Middleware, middleware)
}

func (s *Server) AttachMiddlewares(middlewares ...Middleware) {
	for _, m := range middlewares {
		s.AttachMiddleware(m)
	}
}

// Returns the port that this server has been configured to run on
func (s *Server) Port() int {
	return s.Config.Port
}

// Creates an `http.Handler` by wrapping the servers Router in all registered
// middlwares
func (s *Server) applyMiddleware() http.Handler {
	var mx http.Handler = s.Router
	for i := len(s.Middleware) - 1; i >= 0; i-- {
		mx = s.Middleware[i](mx)
	}
	return mx
}

func healthcheck(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(HEALTHCHECK_MESSAGE + "\n"))
}

func notFound(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusNotFound)
	rw.Write(notFoundMessage())
}

func methodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusMethodNotAllowed)
	rw.Write(methodNotAllowedMessage(r.Method))
}

func methodNotAllowedMessage(method string) []byte {
	return ErrFormat(
		fmt.Sprintf("Method '%v' is not allowed on this route", method))
}

func notFoundMessage() []byte {
	return ErrFormat("The route you have requested could not be found")
}

func ErrFormat(message string) []byte {
	return MsgFormat("Err", message)
}

func MsgFormat(key, value string) []byte {
	return []byte(fmt.Sprintf(`{"%v":"%v"}`, key, value))
}
