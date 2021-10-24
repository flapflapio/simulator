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

type GraphParams struct {
	Start       string
	States      []State
	Transitions []TransitionParams
}

type TransitionParams struct {
	Start  string
	End    string
	Symbol string
}

func From(params GraphParams) *Graph {
	err := validateParams(params)
	if err != nil {
		return nil
	}
	g := Graph{
		States: params.States,
	}
	if g.States == nil {
		g.States = []State{}
	}
	g.Start = g.FindState(params.Start)
	for _, t := range params.Transitions {
		g.Transitions = append(g.Transitions, Transition{
			Start:  g.FindState(t.Start),
			End:    g.FindState(t.End),
			Symbol: t.Symbol,
		})
	}
	return &g
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

func (g *Graph) FindState(id string) *State {
	for i, s := range g.States {
		if s.Id == id {
			return &g.States[i]
		}
	}
	return nil
}

func (g *Graph) WithStates(states ...State) *Graph {
	g.States = append(g.States, states...)
	return g
}

func (g *Graph) WithState(state State) *Graph {
	g.States = append(g.States, state)
	return g
}

func (g *Graph) WithTransition(transition Transition) *Graph {
	g.Transitions = append(g.Transitions, transition)
	return g
}

func (g *Graph) WithTransitions(transitions ...Transition) *Graph {
	g.Transitions = append(g.Transitions, transitions...)
	return g
}

func (g *Graph) String() string {
	return fmt.Sprintf(
		"Graph[Start:%v States:%v Transitions:%v]",
		g.Start.Id,
		g.States,
		g.Transitions)
}

func (g *Graph) Json() string {
	mm := g.JsonMap()
	data, err := json.Marshal(mm)
	if err != nil {
		return ""
	}
	return string(data)
}

func (g *Graph) JsonMap() map[string]interface{} {
	res := map[string]interface{}{
		"Start":       g.Start.Id,
		"States":      []map[string]interface{}{},
		"Transitions": []map[string]interface{}{},
	}
	for _, s := range g.States {
		res["States"] = append(res["States"].([]map[string]interface{}), s.JsonMap())
	}
	for _, t := range g.Transitions {
		res["Transitions"] = append(
			res["Transitions"].([]map[string]interface{}), t.JsonMap())
	}
	return res
}

// Deep copy a Machine
func (g *Graph) Copy() *Graph {
	gg := Graph{
		States:      []State{},
		Transitions: []Transition{},
	}

	copyTransition := func(t Transition) Transition {
		tt := Transition{Symbol: t.Symbol}
		for i, s := range gg.States {
			if s.Id == t.Start.Id {
				tt.Start = &gg.States[i]
			}
			if s.Id == t.End.Id {
				tt.End = &gg.States[i]
			}
		}
		return tt
	}

	// Add states
	var oldStartingStateId *string
	for _, s := range g.States {
		if s.Id == g.Start.Id {
			oldStartingStateId = &g.Start.Id
		}
		gg.States = append(gg.States, s.Copy())
	}

	// Add starting state
	if oldStartingStateId != nil {
		for i, s := range gg.States {
			if s.Id == *oldStartingStateId {
				gg.Start = &gg.States[i]
			}
		}
	}

	// Add transitions
	for _, t := range g.Transitions {
		gg.Transitions = append(gg.Transitions, copyTransition(t))
	}

	return &gg
}

func validateParams(params GraphParams) error {
	d, err := json.Marshal(params)
	if err != nil {
		return err
	}
	_, err = Load(d)
	if err != nil {
		return err
	}
	return nil
}
