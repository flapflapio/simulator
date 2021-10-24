package simulationcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/flapflapio/simulator/core/controllers"
	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/automata"
	"github.com/gorilla/mux"
)

type SimulationController struct {
	prefix    string
	simulator simulation.Simulator
}

func New(simulator simulation.Simulator) *SimulationController {
	return &SimulationController{
		prefix:    "/",
		simulator: simulator,
	}
}

// Attaches this controller to the given router
func (c *SimulationController) Attach(router *mux.Router) {
	r := controllers.CreateSubrouter(router, c.prefix)
	r.Methods("POST").Path("").HandlerFunc(c.DoSimulation)
}

func (c *SimulationController) WithPrefix(prefix string) *SimulationController {
	return &SimulationController{
		prefix:    prefix,
		simulator: c.simulator,
	}
}

func (c *SimulationController) DoSimulation(rw http.ResponseWriter, r *http.Request) {
	m, err := automata.Load(r.Body)

	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte("The machine that was sent is not " +
			"valid or could not be processed\n"))
		return
	}

	tape := r.URL.Query().Get("tape")
	if tape == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Please provide an input with query param 'tape'\n"))
		return
	}

	var sim simulation.Simulation

	// Create a new simulation on a non-existant machine
	id, err := c.simulator.Start(m, tape)
	check(err)

	// Run the simulation from start to finish
	for sim = c.simulator.Get(id); !sim.Done(); sim.Step() {
	}

	// Grab the result of the simulation
	res, err := sim.Result()
	check(err)

	// Serialize result
	data, err := json.Marshal(res)
	check(err)

	// Write result to response body
	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(append(data, '\n'))
	c.simulator.End(id)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
