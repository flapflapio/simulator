package simulationcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/flapflapio/simulator/core/types"
	"github.com/flapflapio/simulator/core/util"
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
			Queries("tape", "{tape}").
			HandlerFunc(controller.DoSimulation)
	}
}

func (controller *SimulationController) WithPrefix(prefix string) types.Controller {
	return &SimulationController{
		prefix:    prefix,
		simulator: controller.simulator,
	}
}

// TODO: This method is a basically a placeholder for now. It doesn't really
// TODO: simulate anything.
func (controller *SimulationController) DoSimulation(rw http.ResponseWriter, r *http.Request) {
	var sim types.Simulation

	// Create a new simulation on a non-existant machine
	id, err := controller.simulator.Start(struct{}{}, r.URL.Query().Get("tape"))
	if ok := check(err, rw); !ok {
		return
	}

	// Run the simulation from start to finish
	for sim = controller.simulator.Get(id); !sim.Done(); sim.Step() {
	}

	// Grab the result of the simulation
	res, err := sim.Result()
	if ok := check(err, rw); !ok {
		return
	}

	writeSimResult(rw, res)
}

func check(err error, rw http.ResponseWriter) bool {
	if err != nil {
		util.MustWriteJSON(http.StatusInternalServerError, rw, map[string]interface{}{
			"error": err.Error(),
		})
		return false
	}
	return true
}

// Writes the result of a simulation to your http response body
func writeSimResult(rw http.ResponseWriter, res types.Result) {
	data, err := json.Marshal(res)
	if ok := check(err, rw); !ok {
		return
	}
	_, err = rw.Write(data)
	if ok := check(err, rw); !ok {
		return
	}
	rw.WriteHeader(http.StatusOK)
}
