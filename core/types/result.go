package types

// A report of the current state of a simulation
type Report struct {
	Accepted bool     `json:"Accepted"`
	Path     []string `json:"Path"`
}

// A report of the end result of a simulation
type Result struct {
	Accepted bool     `json:"Accepted"`
	Path     []string `json:"Path"`
}
