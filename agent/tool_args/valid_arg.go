package toolargs

import (
	"encoding/json"
	"strings"
)

func isInvalidJSON(content string) bool {
	if !strings.HasPrefix(content, "{") || !strings.HasSuffix(content, "}") {
		return false
	}

	var a any
	err := json.Unmarshal([]byte(content), &a)
	return err != nil
}

func GetValidArg(args string) string {
	content := FromArgFormat(args)

	if isInvalidJSON(content) {
		bytes, _ := json.Marshal(content)
		content = "Could not parse arg as a valid JSON Arguments: " + string(bytes)
		content = NewToolArg(content).ToJSON()
	}

	return NewToolArg(content).ToJSON()
}
