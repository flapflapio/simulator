package machine

type Machine struct {
	Start       *State       `json:"Start"`
	States      []*State     `json:"States"`
	Transitions []Transition `json:"Transitions"`
}

func NewMachine(start *State, state []*State, transitions []Transition) *Machine {
	return &Machine{
		Start:       start,
		States:      state,
		Transitions: transitions,
	}
}

func NewBlankMachine(state *State) *Machine {
	return &Machine{Start: state}
}

func (m1 *Machine) WithStates(states ...*State) *Machine {
	m1.States = append(m1.States, states...)
	return m1
}

func (m *Machine) WithState(state *State) *Machine {
	m.States = append(m.States, state)
	return m
}

func (m *Machine) WithTransition(transition Transition) *Machine {
	m.Transitions = append(m.Transitions, transition)
	return m
}

func (m *Machine) WithTransitions(transitions ...Transition) *Machine {
	m.Transitions = append(m.Transitions, transitions...)
	return m
}
