package simulation

import "fmt"

// A report of the end result of a simulation
type Result struct {
	Accepted       bool     `json:"Accepted"`
	Path           []string `json:"Path"`
	RemainingInput string   `json:"RemainingInput"`
}

func (r Result) String() string {
	return fmt.Sprintf("Result[Accepted:%v Path:%v]", r.Accepted, r.Path)
}
