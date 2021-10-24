package simulation

import "github.com/flapflapio/simulator/core/simulation/machine"

type Machine interface {
	machine.Marshalable
	Simulate(input string) Simulation
}
