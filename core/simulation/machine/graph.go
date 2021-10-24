package machine

import (
	"encoding/json"
	"fmt"
)

type Graph struct {
	Start       *State       `json:"Start"`
	States      []State      `json:"States"`
	Transitions []Transition `json:"Transitions"`
}

func NewGraph(start *State, states []State, transitions []Transition) *Graph {
	return &Graph{
		Start:       start,
		States:      states,
		Transitions: transitions,
	}
}

func NewBlankGraph() *Graph {
	return &Graph{
		States:      []State{},
		Transitions: []Transition{},
	}
}

func (m1 *Graph) WithStates(states ...State) *Graph {
	m1.States = append(m1.States, states...)
	return m1
}

func (m *Graph) WithState(state State) *Graph {
	m.States = append(m.States, state)
	return m
}

func (m *Graph) WithTransition(transition Transition) *Graph {
	m.Transitions = append(m.Transitions, transition)
	return m
}

func (m *Graph) WithTransitions(transitions ...Transition) *Graph {
	m.Transitions = append(m.Transitions, transitions...)
	return m
}

func (m *Graph) String() string {
	return fmt.Sprintf(
		"Graph[Start:%v States:%v Transitions:%v]",
		m.Start.Id,
		m.States,
		m.Transitions)
}

func (m *Graph) Json() string {
	mm := m.JsonMap()
	data, err := json.Marshal(mm)
	if err != nil {
		return ""
	}
	return string(data)
}

func (m *Graph) JsonMap() map[string]interface{} {
	res := map[string]interface{}{
		"Start":       m.Start.Id,
		"States":      []map[string]interface{}{},
		"Transitions": []map[string]interface{}{},
	}
	for _, s := range m.States {
		res["States"] = append(res["States"].([]map[string]interface{}), s.JsonMap())
	}
	for _, t := range m.Transitions {
		res["Transitions"] = append(
			res["Transitions"].([]map[string]interface{}), t.JsonMap())
	}
	return res
}

// Deep copy a Machine
func (m *Graph) Copy() *Graph {
	mm := Graph{
		States:      []State{},
		Transitions: []Transition{},
	}

	copyTransition := func(t Transition) Transition {
		tt := Transition{Symbol: t.Symbol}
		for i, s := range mm.States {
			if s.Id == t.Start.Id {
				tt.Start = &mm.States[i]
			}
			if s.Id == t.End.Id {
				tt.End = &mm.States[i]
			}
		}
		return tt
	}

	// Add states
	var oldStartingStateId *string
	for _, s := range m.States {
		if s.Id == m.Start.Id {
			oldStartingStateId = &m.Start.Id
		}
		mm.States = append(mm.States, s.Copy())
	}

	// Add starting state
	if oldStartingStateId != nil {
		for i, s := range mm.States {
			if s.Id == *oldStartingStateId {
				mm.Start = &mm.States[i]
			}
		}
	}

	// Add transitions
	for _, t := range m.Transitions {
		mm.Transitions = append(mm.Transitions, copyTransition(t))
	}

	return &mm
}
