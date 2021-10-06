package types

// A Simulation is used for inspecting the state of your machine throughout the
// processing of an input string
type Simulation interface {
	// Perform a transition
	Step()

	// Get the current status (state + other info) of a simulation
	Stat() Report

	// Get the final result of your simulation.
	// Returns a SimulationIncomplete error if the simulation is not done
	Result() (Result, error)

	// Check if a simulation is finished
	Done() bool
}
