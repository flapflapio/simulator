package types

import (
	"github.com/flapflapio/simulator/core/simulation/machine"
)

type SimulationFactory func(machine *machine.Machine, input string) (Simulation, error)
