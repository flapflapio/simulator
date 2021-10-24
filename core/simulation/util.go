// Utilities for manipulating simulations
package simulation

func RunToCompletion(sim Simulation) {
	for ; !sim.Done(); sim.Step() {
	}
}

func ResultOf(sim Simulation) *Result {
	RunToCompletion(sim)
	res, err := sim.Result()
	if err != nil {
		return nil
	}
	return &res
}
