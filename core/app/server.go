package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flapflapio/simulator/core/types"
	"github.com/flapflapio/simulator/core/util"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const (
	HEALTHCHECK_MESSAGE = "All good in the hood"
)

type Server struct {
	Name   string
	Router *mux.Router
}

func New(name string) *Server {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowed)

	router.
		Methods("GET").
		Name("Healthcheck").
		Path("/healthcheck").
		HandlerFunc(healthcheck)

	return &Server{Name: name, Router: router}
}

func (s *Server) Run(config Config) error {
	log.Printf("Starting %v", s.Name)
	log.Printf("Listening on %v\n", config.Port)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%v", config.Port),
		Handler:        newLoggerAndRecoveryMiddlewareWrapper(s.Router),
		ReadTimeout:    time.Second * time.Duration(config.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(config.WriteTimeout),
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	log.Printf(
		"Server config: ReadTimeout = %v, WriteTimeout = %v, MaxHeaderBytes = %v",
		server.ReadTimeout,
		server.WriteTimeout,
		server.MaxHeaderBytes)

	return server.ListenAndServe()
}

func (s *Server) Attach(controller types.Controller) *Server {
	controller.Attach(s.Router)
	return s
}

func (s *Server) AttachControllers(controllers []types.Controller) *Server {
	for _, controller := range controllers {
		s.Attach(controller)
	}
	return s
}

func newLoggerAndRecoveryMiddlewareWrapper(h http.Handler) http.Handler {
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(h)
	return n
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
