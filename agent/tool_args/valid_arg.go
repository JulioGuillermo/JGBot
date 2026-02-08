package toolargs

import "fmt"

func GetToolArgContent(args string) (string, error) {
	return ExtractContent(args)
}

func GetMsgValidArg(msg string) string {
	content, err := GetToolArgContent(msg)
	if err != nil {
		return ToolArgErrFromStr(
			msg,
			fmt.Sprintf("Could not parse arg as a valid JSON Arguments. Please ensure the arg is a valid JSON and matches the tool arguments schema. Error: %s", err.Error()),
		).ToJSON()
	}

	return NewToolArg(content).ToJSON()
}

func GetMsgToolCallArg(msg string) string {
	content, err := GetToolArgContent(msg)
	if err != nil {
		return ToolArgErrFromStr(
			msg,
			fmt.Sprintf("Could not parse arg as a valid JSON Arguments. Please ensure the arg is a valid JSON and matches the tool arguments schema. Error: %s", err.Error()),
		).ToJSON()
	}

	if !LooksJson(content) && !IsJson(content) {
		return ToolArgErrFromStr(
			content,
			"Could not parse arg as a valid JSON Arguments. Please ensure the arg is a valid JSON and matches the tool arguments schema.",
		).ToJSON()
		// return NewFromStrArg(content).ToJSON()
	}

	return NewToolArg(content).ToJSON()
}
