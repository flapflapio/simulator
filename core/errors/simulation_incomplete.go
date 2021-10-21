package errors

const SIMULATION_INCOMPLETE = "simulation incomplete"

type SimulationNotDone struct{ Msg string }

func NotDone() SimulationNotDone {
	return SimulationNotDone{Msg: SIMULATION_INCOMPLETE}
}

func (s SimulationNotDone) Error() string {
	return s.Msg
}
