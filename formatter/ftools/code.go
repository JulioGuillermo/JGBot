package ftools

import (
	"fmt"
	"regexp"
)

const CodePlaceholder = "CODEBLOCK_PH_"

// MapCodeBlocks finds code blocks and replaces them with placeholders.
func MapCodeBlocks(msg string) (string, []string) {
	var codeBlocks []string
	reCode := regexp.MustCompile("(?s)```.*?```|`[^`]*`")
	msg = reCode.ReplaceAllStringFunc(msg, func(match string) string {
		codeBlocks = append(codeBlocks, match)
		return fmt.Sprintf("%s%d", CodePlaceholder, len(codeBlocks)-1)
	})
	return msg, codeBlocks
}

// RestoreCodeBlocks replaces code placeholders with their original content.
func RestoreCodeBlocks(msg string, codeBlocks []string, formatFunc func(string) string) string {
	reRestoreCode := regexp.MustCompile(CodePlaceholder + `(\d+)`)
	return reRestoreCode.ReplaceAllStringFunc(msg, func(match string) string {
		var idx int
		n, err := fmt.Sscanf(match, CodePlaceholder+"%d", &idx)
		if err != nil || n != 1 {
			return match
		}
		if idx >= 0 && idx < len(codeBlocks) {
			content := codeBlocks[idx]
			if formatFunc != nil {
				return formatFunc(content)
			}
			return content
		}
		return match
	})
}
