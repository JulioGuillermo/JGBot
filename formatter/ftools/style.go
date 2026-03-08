package ftools

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	BoldPlaceholder   = "«BOLD»"
	ItalicPlaceholder = "«ITALIC»"
)

func ProtectBold(msg string) string {
	msg = ParseStyle(msg, `**`, `**`, BoldPlaceholder, BoldPlaceholder)
	msg = ParseStyle(msg, `__`, `__`, BoldPlaceholder, BoldPlaceholder)
	return msg
}

func RestoreBold(msg string, replacementWrapper string) string {
	msg = ParseStyle(msg, BoldPlaceholder, BoldPlaceholder, replacementWrapper, replacementWrapper)
	return msg
}

func FormatItalic(msg string, replacementWrapper string) string {
	msg = ParseStyle(msg, `*`, `*`, replacementWrapper, replacementWrapper)
	msg = ParseStyle(msg, `_`, `_`, replacementWrapper, replacementWrapper)
	return msg
}

func FormatStrike(msg string, replacementWrapper string) string {
	msg = ParseStyle(msg, `~~`, `~~`, replacementWrapper, replacementWrapper)
	return msg
}

func FormatStyleTags(msg string, startPattern, endPattern string, replaceStart, replaceEnd string) string {
	// We want to match startPattern + content + endPattern.
	// For underscore emphasis, avoid matching inside words/identifiers like `ak_arch`.
	if startPattern == "_" {
		re := regexp.MustCompile(fmt.Sprintf(
			`(^|[^\w\d])%s([^\s](?:.*?[^\s])?)%s([^\w\d]|$)`,
			regexp.QuoteMeta(startPattern),
			regexp.QuoteMeta(endPattern),
		))
		return re.ReplaceAllStringFunc(msg, func(match string) string {
			matchs := re.FindStringSubmatch(match)
			if len(matchs) < 4 {
				return match
			}
			return matchs[1] + replaceStart + matchs[2] + replaceEnd + matchs[3]
		})
	}

	pattern := regexp.QuoteMeta(startPattern) + `([^\s](?:.*?[^\s])?)` + regexp.QuoteMeta(endPattern)
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllStringFunc(msg, func(match string) string {
		content := match[len(startPattern) : len(match)-len(endPattern)]
		return replaceStart + content + replaceEnd
	})
}

func EscapeHTML(msg string) string {
	msg = strings.ReplaceAll(msg, "&", "&amp;")
	msg = strings.ReplaceAll(msg, "<", "&lt;")
	msg = strings.ReplaceAll(msg, ">", "&gt;")
	return msg
}
