package toolargs

import "encoding/json"

type ToolArg struct {
	Arg string `json:"__arg1"`
}

func NewToolArg(arg string) *ToolArg {
	return &ToolArg{Arg: arg}
}

func ToolArgFromJSON(s string) *ToolArg {
	var t ToolArg
	err := json.Unmarshal([]byte(s), &t)
	if err != nil {
		return nil
	}
	return &t
}

func ToolArgFromAny(m any) *ToolArg {
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return NewToolArg(string(b))
}

func (t *ToolArg) String() string {
	return t.Arg
}

func (t *ToolArg) ToJSON() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(bytes)
}
