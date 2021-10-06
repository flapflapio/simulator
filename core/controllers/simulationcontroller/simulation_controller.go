package simulationcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/flapflapio/simulator/core/types"
	"github.com/gorilla/mux"
)

type SimulationController struct {
	prefix    string
	simulator types.Simulator
}

func New(simulator types.Simulator) *SimulationController {
	return &SimulationController{
		prefix:    "/",
		simulator: simulator,
	}
}

// Attaches this controller to the given router
func (controller *SimulationController) Attach(router *mux.Router) {
	r := router
	if controller.prefix != "" && controller.prefix != "/" {
		r = router.PathPrefix(controller.prefix).Subrouter()
	}

	for _, path := range []string{"", "/"} {
		r.Methods("GET").
			Path(path).
			HandlerFunc(controller.DoSimulation)
	}
}

func (controller *SimulationController) WithPrefix(prefix string) types.Controller {
	return &SimulationController{
		prefix:    prefix,
		simulator: controller.simulator,
	}
}

// TODO: This method is a basically a placeholder for now. It uses a phony
// TODO: simulator for now
func (controller *SimulationController) DoSimulation(rw http.ResponseWriter, r *http.Request) {
	var sim types.Simulation
	var tape = r.URL.Query().Get("tape")

	if tape == "" {
		rw.WriteHeader(400)
		rw.Write([]byte("Please provide an input with query param 'tape'\n"))
		return
	}

	// Create a new simulation on a non-existant machine
	id, err := controller.simulator.Start(nil, tape)
	check(err)

	// Run the simulation from start to finish
	for sim = controller.simulator.Get(id); !sim.Done(); sim.Step() {
	}

	// Grab the result of the simulation
	res, err := sim.Result()
	check(err)

	// Serialize result
	data, err := json.Marshal(res)
	check(err)

	// Write result to response body
	rw.WriteHeader(http.StatusOK)
	rw.Write(append(data, '\n'))
	controller.simulator.End(id)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
