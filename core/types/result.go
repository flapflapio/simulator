package types

// A report of the current state of a simulation
type Report struct{}

// A report of the end result of a simulation
type Result struct {
	Accepted bool     `json:"Accepted"`
	Path     []string `json:"Path"`
}
