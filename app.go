package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/flapflapio/simulator/core/app"
	"github.com/flapflapio/simulator/core/controllers"
	"github.com/flapflapio/simulator/core/controllers/schemacontroller"
	"github.com/flapflapio/simulator/core/controllers/simulationcontroller"
	"github.com/flapflapio/simulator/core/services/simulatorservice"
)

var (
	cfg = configure()
	srv = app.New(cfg)
	sim = simulatorservice.New()

	// Add any new middlewares to this slice - mids is added in
	// reverse order (i.e. mids at the top of this slice is applied
	// first)
	mids = []app.Middleware{
		// app.Timeout(60 * time.Second),
		// app.LoggerAndRecovery,
		// app.TrimTrailingSlash,
		// app.CORS(*cfg.CORS...),
	}

	// Add any new cntrls to this slice
	cntrls = []controllers.Controller{
		schemacontroller.New(),
		simulationcontroller.New(sim),
	}
)

func main() {
	if exit := healthcheckMode(); exit > -1 {
		os.Exit(exit)
	}
	setupServer()
	srv.Run()
	panic(<-srv.Wait())
}

func setupServer() {
	log.Println(srv.Config)
	srv.Attach(cntrls, mids)
}

func configure() app.Config {
	config, err := app.GetConfig()
	if err != nil {
		log.Println("An error occured while configuring the app")
		log.Fatalf("%v\n", err)
	}
	return config
}

var (
	healthcheck = flag.Bool(
		"health",
		false,
		"Runs a healthcheck on the server and then exits")
)

// Healthcheck stuff
func healthcheckMode() int {
	route := fmt.Sprintf("http://localhost:%v/healthcheck", srv.Port())
	flag.Parse()
	if !*healthcheck {
		return -1
	}
	r, err := http.Get(route)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}
	if r.StatusCode != http.StatusOK {
		fmt.Printf("Get \"%v\": %v\n", route, r.Status)
		return 1
	}
	if bod, err := io.ReadAll(r.Body); err == nil {
		fmt.Print(string(bod))
	}
	return 0
}
