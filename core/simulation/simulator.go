package simulation

import "github.com/flapflapio/simulator/core/simulation/machine"

// A Simulator is used for managing your simulations
type Simulator interface {
	// Begins a new simulation
	Start(machine *machine.Machine, input string) (id int, err error)

	// Get a simulation by id
	Get(simulationId int) Simulation

	// Ends a simulation
	End(simulationId int) error
}
