package machine

import (
	"encoding/json"
	"fmt"
)

type State struct {
	Id     string `json:"Id"`
	Ending bool   `json:"Ending"`
}

func (s State) String() string {
	return fmt.Sprintf("State%v", fmt.Sprintf("%v", s.JsonMap())[3:])
}

func (s State) Json() string {
	m := s.JsonMap()
	data, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(data)
}

func (s State) JsonMap() map[string]interface{} {
	return map[string]interface{}{
		"Id":     s.Id,
		"Ending": s.Ending,
	}
}

func (s State) Copy() State {
	return State{Id: s.Id, Ending: s.Ending}
}
