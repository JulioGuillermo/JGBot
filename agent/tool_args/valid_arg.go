package toolargs

import "fmt"

func GetToolArgContent(args string) (string, error) {
	return ExtractContent(args)
}

func GetMsgValidArg(msg string) string {
	content, err := GetToolArgContent(msg)
	if err != nil {
		return NewToolArgError(fmt.Sprintf("Could not parse arg as a valid JSON Arguments: %s\n%s", err.Error(), msg)).ToJSON()
	}

	return NewToolArg(content).ToJSON()
}

func GetMsgToolCallArg(msg string) string {
	content, err := GetToolArgContent(msg)
	if err != nil {
		return NewToolArgError(fmt.Sprintf("Could not parse arg as a valid JSON Arguments: %s\n%s", err.Error(), msg)).ToJSON()
	}

	return NewToolArg(content).ToJSON()
}
