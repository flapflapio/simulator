package machine

import (
	"encoding/json"
	"fmt"
)

type Transition struct {
	Start  *State `json:"Start"`
	End    *State `json:"End"`
	Symbol string `json:"Symbol"`
}

func (s Transition) String() string {
	return fmt.Sprintf("Transition%v", fmt.Sprintf("%v", s.JsonMap())[3:])
}

func (s Transition) Json() string {
	m := s.JsonMap()
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}

func (s Transition) JsonMap() map[string]interface{} {
	return map[string]interface{}{
		"Start":  s.Start.Id,
		"End":    s.End.Id,
		"Symbol": s.Symbol,
	}
}

func (s Transition) Copy() *Transition {
	return &Transition{
		Start:  s.Start,
		End:    s.End,
		Symbol: s.Symbol,
	}
}
