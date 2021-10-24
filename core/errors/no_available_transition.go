package errors

const NO_POSSIBLE_TRANSITION = "no possible transition"

type NoPossibleTransition struct {
	Msg string
}

func (e NoPossibleTransition) Error() string {
	return e.Msg
}

func NoTrans() NoPossibleTransition {
	return NoPossibleTransition{Msg: NO_POSSIBLE_TRANSITION}
}
