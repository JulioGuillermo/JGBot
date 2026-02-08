package toolargs

import (
	"encoding/json"
	"errors"
	"strings"
)

func LooksJson(args string) bool {
	trimmed := strings.TrimSpace(args)
	return strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")
}

func IsJson(args string) bool {
	return json.Valid([]byte(args))
}

func ExtractArg(args string) (string, bool, error) {
	if !LooksJson(args) {
		return args, false, nil
	}

	if !IsJson(args) {
		return "", false, errors.New("invalid json")
	}

	var argsMap map[string]any

	err := json.Unmarshal([]byte(args), &argsMap)
	if err != nil {
		return args, false, nil
	}

	content, ok := argsMap["string_arg"]
	if ok {
		return content.(string), true, nil
	}

	content, ok = argsMap["__arg1"]
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
