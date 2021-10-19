package main

import (
	"fmt"
	"log"

	"github.com/flapflapio/simulator/core/app"
	"github.com/flapflapio/simulator/core/controllers/schemacontroller"
	"github.com/flapflapio/simulator/core/controllers/simulationcontroller"
	"github.com/flapflapio/simulator/core/services/simulatorservice"
	"github.com/flapflapio/simulator/core/simulation/machine"
	"github.com/flapflapio/simulator/core/types"
)

func main() {

	var (
		config           = configure()
		server           = app.New(config)
		simulatorService = createSimulatorService()

		// Add any new middlewares to this slice
		middleware = []app.Middleware{
			app.LoggerAndRecovery,
			app.TrimTrailingSlash(true),
		}

		// Add any new controllers to this slice
		controllers = []types.Controller{
			schemacontroller.New(),

			simulationcontroller.
				New(simulatorService).
				WithPrefix("/simulate"),
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

func createSimulatorService() types.Simulator {
	return simulatorservice.New(
		func(machine *machine.Machine, input string) (types.Simulation, error) {
			return &PhonySimulation{input: input}, nil
		})
}
