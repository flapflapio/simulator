package main

import (
	"fmt"
	"log"

	"github.com/flapflapio/simulator/core/app"
	"github.com/flapflapio/simulator/core/controllers"
	"github.com/flapflapio/simulator/core/controllers/schemacontroller"
	"github.com/flapflapio/simulator/core/controllers/simulationcontroller"
	"github.com/flapflapio/simulator/core/services/simulatorservice"
)

func main() {

	var (
		config           = configure()
		server           = app.New(config)
		simulatorService = simulatorservice.New()

		// Add any new middlewares to this slice - middleware is added in
		// reverse order (i.e. middleware at the top of this slice is applied
		// first)
		middleware = []app.Middleware{
			app.LoggerAndRecovery,
			app.TrimTrailingSlash(true),
		}

		// Add any new controllers to this slice
		controllers = []controllers.Controller{
			schemacontroller.New(),
			simulationcontroller.New(simulatorService),
		}
	)

	fmt.Println(config)
	server.Attach(controllers, middleware)
	log.Fatal(server.Run())
}

func configure() app.Config {
	config, err := app.GetConfig()
	if err != nil {
		log.Println("An error occured while configuring the app")
		log.Fatalf("%v\n", err)
	}
	return config
}
