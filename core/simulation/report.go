package simulation

import "fmt"

// A report of the current state of a simulation
// TODO: make a more detailed report that consists of more than just the result
type Report struct {
	Result
}

func (r Report) String() string {
	return fmt.Sprintf("Report[Result:%v]", r.Result)
}
