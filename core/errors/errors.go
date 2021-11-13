package errors

import "errors"

var ErrNoTransition = errors.New("no possible transition")
var ErrSimulationIncomplete = errors.New("simulation incomplete")
