package tools

import "regexp"

func GetInputString(input string) string {
	re := regexp.MustCompile(`^\s*\{\s*"__arg1"\s*:\s*(.*?)\s*\}\s*$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return input
}
