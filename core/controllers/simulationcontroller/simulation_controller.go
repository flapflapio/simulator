package simulationcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/flapflapio/simulator/core/controllers/utils"
	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/automata"
	"github.com/gorilla/mux"
)

const (
	INVALID_MACHINE_MSG = `` +
		`{"Err":"The machine that was sent is not ` +
		`valid or otherwise could not be processed"}`

	PLEASE_PROVIDE_A_TAPE_MSG = `` +
		`{"Err":"Please provide an input with query param 'tape'"}`

	FAILED_TO_CREATE_A_NEW_SIMULATION = `` +
		`{"Err":"Failed to create a new simulation"}`

	FAILED_TO_OBTAIN_RESULTS_OF_SIMULATION = `` +
		`{"Err":"Failed to obtain results of simulation"}`

	FAILED_TO_CREATE_A_RESPONSE = `` +
		`{"Err":"Failed to create a response"}`
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
	r := utils.CreateSubrouter(router, c.prefix)
	r.Methods("POST").Path("/simulate").HandlerFunc(c.DoSimulation)
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
		rw.Write([]byte(INVALID_MACHINE_MSG))
		return
	}

	tape := r.URL.Query().Get("tape")
	if tape == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(PLEASE_PROVIDE_A_TAPE_MSG))
		return
	}

	var sim simulation.Simulation

	// Create a new simulation
	id, err := c.simulator.Start(m, tape)
	if check(err, rw, FAILED_TO_CREATE_A_NEW_SIMULATION) {
		return
	}

	// Run the simulation from start to finish
	for sim = c.simulator.Get(id); !sim.Done(); sim.Step() {
	}

	// Grab the result of the simulation
	res, err := sim.Result()
	if check(err, rw, FAILED_TO_OBTAIN_RESULTS_OF_SIMULATION) {
		return
	}

	// Serialize result
	data, err := json.Marshal(res)
	if check(err, rw, FAILED_TO_CREATE_A_RESPONSE) {
		return
	}

	// Write result to response body
	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(append(data, '\n'))
	c.simulator.End(id)
}

func check(err error, rw http.ResponseWriter, msg string) bool {
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(msg))
		return true
	}
	return false
}
