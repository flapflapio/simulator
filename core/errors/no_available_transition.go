package errors

const msg = "no possible transition"

type NoPossibleTransition struct {
	Msg string
}

func (e NoPossibleTransition) Error() string {
	return e.Msg
}

func NoTrans() NoPossibleTransition {
	return NoPossibleTransition{Msg: msg}
}
