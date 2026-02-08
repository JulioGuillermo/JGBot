package toolargs

import "encoding/json"

type ToolArg struct {
	Arg any `json:"__arg1"`
}

type PlainArg struct {
	Arg string `json:"string_arg"`
}

type ToolArgError struct {
	Arg   any    `json:"__arg1"`
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

func NewFromStrArg(arg string) *ToolArg {
	strArg := &PlainArg{
		Arg: arg,
	}
	return &ToolArg{
		Arg: strArg,
	}
}

func NewToolArgError(arg any, errMsg string) *ToolArgError {
	return &ToolArgError{
		Arg:   arg,
		Error: errMsg,
	}
}

func ToolArgErrFromStr(arg, errMsg string) *ToolArgError {
	strArg := &PlainArg{
		Arg: arg,
	}
	return &ToolArgError{
		Arg:   strArg,
		Error: errMsg,
	}
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
