package machine

type Transition struct {
	Start  *State
	End    *State
	Symbol string
}
