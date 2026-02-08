package ftools

import (
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
	// But if startPattern is "_", we want to avoid matching it if it's inside a word (like LINK_PH_0).
	// In Markdown, _italic_ requires it to be at the start of a word or preceded by space.

	pattern := regexp.QuoteMeta(startPattern) + `(.*?)` + regexp.QuoteMeta(endPattern)
	if startPattern == "_" {
		// More strict for underscore to avoid matching in the middle of words (placeholders)
		// It must be at the start of string or preceded by a non-word char,
		// and followed by a non-word char or end of string.
		// Since Go regexp doesn't support lookarounds, we'll use a simpler heuristic:
		// match only if it's not immediately preceded/followed by letters/numbers/underscores?
		// Actually, let's just use the fact that placeholders are usually Uppercase + PH + Number.

		re := regexp.MustCompile(pattern)
		return re.ReplaceAllStringFunc(msg, func(match string) string {
			// Check if this match is likely a placeholder (contains _PH_)
			if strings.Contains(match, "_PH_") {
				return match
			}
			content := match[len(startPattern) : len(match)-len(endPattern)]
			return replaceStart + content + replaceEnd
		})
	}

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
