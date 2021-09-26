package errors

type SimulationNotDone struct{ Message string }

func (si SimulationNotDone) Error() string {
	return si.Message
}
