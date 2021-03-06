package simulatorservice

import (
	"fmt"
	"sync"

	"github.com/flapflapio/simulator/core/simulation"
)

type SimulatorService struct {
	sims   map[int]simulation.Simulation
	nextId int
	lock   sync.Mutex
}

func New() *SimulatorService {
	return &SimulatorService{
		sims: map[int]simulation.Simulation{},
	}
}

// Begins a new simulation
func (ss *SimulatorService) Start(
	machine simulation.Machine,
	input string,
) (id int, err error) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	i := ss.nextId
	ss.nextId++
	ss.sims[i] = machine.Simulate(input)
	return i, nil
}

// Get a simulation by id
func (ss *SimulatorService) Get(simulationId int) simulation.Simulation {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	sim := ss.sims[simulationId]
	return sim
}

// Ends a simulation
func (ss *SimulatorService) End(simulationId int) error {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	if _, ok := ss.sims[simulationId]; !ok {
		return fmt.Errorf("simulation with id '%v' does not exist", simulationId)
	}
	delete(ss.sims, simulationId)
	return nil
}
