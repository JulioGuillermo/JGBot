package toolargs

import "encoding/json"

type ToolArg struct {
	Arg any `json:"__arg1"`
}

type ToolArgError struct {
	Error string `json:"error"`
}

func NewToolArg(arg string) *ToolArg {
	var a any
	err := json.Unmarshal([]byte(arg), &a)
	if err != nil {
		a = arg
	}
	return &ToolArg{Arg: a}
}

func NewToolArgError(errMsg string) *ToolArgError {
	return &ToolArgError{Error: errMsg}
}

func (t *ToolArg) ToJSON() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (t *ToolArgError) ToJSON() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(bytes)
}
