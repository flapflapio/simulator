package simulatorservice

import "github.com/flapflapio/simulator/core/types"

type SimulatorService struct {
	sims    map[int]types.Simulation
	factory types.SimulationFactory
	nextId  int
}

func New(simulationFactory types.SimulationFactory) *SimulatorService {
	return &SimulatorService{
		sims:    map[int]types.Simulation{},
		factory: simulationFactory,
	}
}

// Begins a new simulation
func (ss *SimulatorService) Start(machine types.Machine, input string) (id int, err error) {
	i := ss.nextId
	ss.nextId++

	sim, err := ss.factory(machine, input)
	if err != nil {
		return -1, err
	}

	ss.sims[i] = sim
	return i, nil
}

// Get a simulation by id
func (ss *SimulatorService) Get(simulationId int) types.Simulation {
	return ss.sims[simulationId]
}

// Ends a simulation
func (ss *SimulatorService) End(simulationId int) error {
	return ss.sims[simulationId].Kill()
}
