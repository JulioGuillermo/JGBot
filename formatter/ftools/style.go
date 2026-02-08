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
	reBoldStar := regexp.MustCompile(`\*\*(.*?)\*\*`)
	msg = reBoldStar.ReplaceAllStringFunc(msg, func(match string) string {
		content := match[2 : len(match)-2]
		return BoldPlaceholder + content + BoldPlaceholder
	})

	reBoldUnderscore := regexp.MustCompile(`__(.*?)__`)
	msg = reBoldUnderscore.ReplaceAllStringFunc(msg, func(match string) string {
		content := match[2 : len(match)-2]
		return BoldPlaceholder + content + BoldPlaceholder
	})
	return msg
}

func RestoreBold(msg string, replacementWrapper string) string {
	reRestoreBold := regexp.MustCompile(BoldPlaceholder + `(?s)(.*?)` + BoldPlaceholder)
	return reRestoreBold.ReplaceAllStringFunc(msg, func(match string) string {
		content := match[len(BoldPlaceholder) : len(match)-len(BoldPlaceholder)]
		return replacementWrapper + content + replacementWrapper
	})
}

func FormatItalic(msg string, replacementWrapper string) string {
	reItalicStar := regexp.MustCompile(`\*(.*?)\*`)
	msg = reItalicStar.ReplaceAllStringFunc(msg, func(match string) string {
		content := match[1 : len(match)-1]
		return replacementWrapper + content + replacementWrapper
	})
	reItalicUnderscore := regexp.MustCompile(`_(.*?)_`)
	msg = reItalicUnderscore.ReplaceAllStringFunc(msg, func(match string) string {
		content := match[1 : len(match)-1]
		return replacementWrapper + content + replacementWrapper
	})
	return msg
}

func FormatStrike(msg string, replacementWrapper string) string {
	reStrike := regexp.MustCompile(`~~(.*?)~~`)
	return reStrike.ReplaceAllStringFunc(msg, func(match string) string {
		content := match[2 : len(match)-2]
		return replacementWrapper + content + replacementWrapper
	})
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
