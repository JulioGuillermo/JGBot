package toolargs

import (
	"encoding/json"
	"errors"
	"strings"
)

func ExtractArg(args string) (string, bool, error) {
	trimmed := strings.TrimSpace(args)
	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		return args, false, nil
	}

	if !json.Valid([]byte(args)) {
		return "", false, errors.New("invalid json")
	}

	var argsMap map[string]any

	err := json.Unmarshal([]byte(args), &argsMap)
	if err != nil {
		return args, false, nil
	}

	content, ok := argsMap["__arg1"]
	if !ok {
		return args, false, nil
	}

	contentStr, ok := content.(string)
	if ok {
		return contentStr, true, nil
	}

	contentBytes, err := json.Marshal(content)
	if err != nil {
		return "", false, err
	}

	return string(contentBytes), true, nil
}

func ExtractContent(args string) (string, error) {
	extracting := true
	content := args
	var err error
	for extracting {
		content, extracting, err = ExtractArg(content)
		if err != nil {
			return "", err
		}
	}
	return content, err
}
