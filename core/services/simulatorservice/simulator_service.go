package simulatorservice

import (
	"fmt"

	"github.com/flapflapio/simulator/core/simulation"
)

type SimulatorService struct {
	sims   map[int]simulation.Simulation
	nextId int
}

func New() *SimulatorService {
	return &SimulatorService{
		sims: map[int]simulation.Simulation{},
	}
}

// Begins a new simulation
func (ss *SimulatorService) Start(machine simulation.Machine, input string) (id int, err error) {
	i := ss.nextId
	ss.nextId++
	ss.sims[i] = machine.Simulate(input)
	return i, nil
}

// Get a simulation by id
func (ss *SimulatorService) Get(simulationId int) simulation.Simulation {
	return ss.sims[simulationId]
}

// Ends a simulation
func (ss *SimulatorService) End(simulationId int) error {
	sim := ss.sims[simulationId]
	if sim == nil {
		return fmt.Errorf("simulation with id '%v' does not exist", simulationId)
	}
	delete(ss.sims, simulationId)
	return nil
}
