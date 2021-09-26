package main

import (
	"fmt"
	"log"

	"github.com/flapflapio/simulator/core/app"
	"github.com/flapflapio/simulator/core/controllers/simulationcontroller"
	"github.com/flapflapio/simulator/core/services/simulatorservice"
	"github.com/flapflapio/simulator/core/types"
)

func main() {

	var (
		server           = app.New()
		config           = configure()
		simulatorService = createSimulatorService()

		// Add any new controllers to this slice
		controllers = []types.Controller{
			simulationcontroller.
				New(simulatorService).
				WithPrefix("/simulate"),
		}
	)

	server.AttachControllers(controllers)
	log.Fatal(server.Run(config))
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
	return simulatorservice.New(func(machine types.Machine, input string) (types.Simulation, error) {
		return &PhonySimulation{input: input}, nil
	})
}

// A "phony" simulation that accepts any input
type PhonySimulation struct {
	path  []string
	input string
	i     int
}

func (ps *PhonySimulation) Step() {
	ps.path = append(ps.path, fmt.Sprintf("q%v", ps.i))
	ps.input = ps.input[1:]
	ps.i++
}

func (ps *PhonySimulation) Stat() types.Report {
	return types.Report{}
}

func (ps *PhonySimulation) Result() (types.Result, error) {
	return types.Result{
		Accepted: true,
		Path:     ps.path,
	}, nil
}

func (ps *PhonySimulation) Done() bool {
	return len(ps.input) == 0
}

func (ps *PhonySimulation) Kill() error {
	panic("not implemented")
}
