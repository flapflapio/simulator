package types

type SimulationFactory func(machine Machine, input string) (Simulation, error)
