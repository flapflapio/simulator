package app

import (
	"encoding/json"
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

type Server struct {
	Name       string
	Router     *mux.Router
	Config     Config
	Middleware []Middleware
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
		"Server config: ReadTimeout=%v, WriteTimeout=%v, MaxHeaderBytes=%v",
		server.ReadTimeout,
		server.WriteTimeout,
		server.MaxHeaderBytes)

	return server.ListenAndServe()
}

// `stuff` is of type `controllers.Controller` or `Middleware` or a slice of either
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
	data, err := json.Marshal(map[string]string{
		"message": HEALTHCHECK_MESSAGE,
	})
	if err != nil {
		panic(err)
	}
	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func notFound(rw http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(map[string]string{
		"message": "The route you have requested could not be found",
	})
	if err != nil {
		panic(err)
	}

	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusNotFound)
	rw.Write(data)
}

func methodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "" {
		method = "GET"
	}

	data, err := json.Marshal(map[string]string{
		"message": fmt.Sprintf("Method '%v' is not allowed on this route", method),
	})
	if err != nil {
		panic(err)
	}

	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusMethodNotAllowed)
	rw.Write(data)
}
