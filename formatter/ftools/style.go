package ftools

import (
	"regexp"
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
