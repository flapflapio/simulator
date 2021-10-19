package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flapflapio/simulator/core/types"
	"github.com/flapflapio/simulator/core/util"
	"github.com/gorilla/mux"
)

const (
	HEALTHCHECK_MESSAGE = "All good in the hood"
)

type Middleware func(http.Handler) http.Handler

type Server struct {
	Name       string
	Router     *mux.Router
	Config     Config
	Middleware []Middleware
}

func New(config Config) *Server {
	prefix := *config.Prefix
	if prefix == "" {
		prefix = "/"
	}

	router := mux.NewRouter().PathPrefix(prefix).Subrouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)

	router.
		Methods("GET").
		Name("Healthcheck").
		Path("/healthcheck").
		HandlerFunc(healthcheck)

	return &Server{Name: *config.Name, Router: router, Config: config}
}

func (s *Server) Run() error {
	log.Printf("Starting %v", s.Name)
	log.Printf("Listening on %v\n", s.Config.Port)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%v", s.Config.Port),
		Handler:        s.applyMiddleware(),
		ReadTimeout:    time.Second * time.Duration(s.Config.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(s.Config.WriteTimeout),
		MaxHeaderBytes: s.Config.MaxHeaderBytes,
	}

	log.Printf(
		"Server config: Prefix=%v, ReadTimeout=%v, WriteTimeout=%v, MaxHeaderBytes=%v",
		*s.Config.Prefix,
		server.ReadTimeout,
		server.WriteTimeout,
		server.MaxHeaderBytes)

	return server.ListenAndServe()
}

// `stuff` is of type `types.Controller` or `Middleware` or a slice of either
func (s *Server) Attach(stuff ...interface{}) {
	for _, x := range stuff {
		switch xx := x.(type) {
		case types.Controller:
			s.AttachController(xx)
		case []types.Controller:
			s.AttachControllers(xx...)
		case Middleware:
			s.AttachMiddleware(xx)
		case []Middleware:
			s.AttachMiddlewares(xx...)
		}
	}
}

func (s *Server) AttachController(controller types.Controller) {
	controller.Attach(s.Router)
}

func (s *Server) AttachControllers(controllers ...types.Controller) {
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
	util.MustWriteOkJSON(rw, map[string]interface{}{
		"message": HEALTHCHECK_MESSAGE,
	})
}

func notFound(rw http.ResponseWriter, r *http.Request) {
	util.MustWriteJSON(http.StatusNotFound, rw, map[string]interface{}{
		"message": "The route you have requested could not be found",
	})
}

func methodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "" {
		method = "GET"
	}
	util.MustWriteJSON(http.StatusMethodNotAllowed, rw, map[string]interface{}{
		"message": fmt.Sprintf("Method '%v' is not allowed on this route", method),
	})
}
