package ftools

import (
	"fmt"
	"regexp"
)

func ParseStyle(msg, startPattern, endPattern, replaceStart, replaceEnd string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(([^\w\d]|^)?)%s([^\s](?:.*?[^\s])?)%s(([^\w\d]|$)?)`, regexp.QuoteMeta(startPattern), regexp.QuoteMeta(endPattern)))
	result := re.ReplaceAllStringFunc(msg, func(match string) string {
		matchs := re.FindStringSubmatch(match)
		return matchs[2] + replaceStart + matchs[3] + replaceEnd + matchs[4]
	})
	// fmt.Println("Parsing", startPattern, endPattern, "with", replaceStart, replaceEnd, "in", msg, "->", result)
	return result
}
