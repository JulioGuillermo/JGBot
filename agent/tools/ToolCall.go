package tools

import "encoding/json"

type ToolCall struct {
	ID    string
	Tool  string
	Input string
	Log   string
}

type ToolResult struct {
	ToolCall ToolCall

	Output string
	Error  string
}

func ToolCallFromJson(data string) *ToolCall {
	if data == "" {
		return nil
	}
	var toolCall ToolCall
	err := json.Unmarshal([]byte(data), &toolCall)
	if err != nil {
		return nil
	}
	return &toolCall
}

func ToolResultFromJson(data string) *ToolResult {
	if data == "" {
		return nil
	}
	var toolResult ToolResult
	err := json.Unmarshal([]byte(data), &toolResult)
	if err != nil {
		return nil
	}
	return &toolResult
}

func (t *ToolCall) ToJson() string {
	data, _ := json.Marshal(t)
	return string(data)
}

func (t *ToolResult) ToJson() string {
	data, _ := json.Marshal(t)
	return string(data)
}
