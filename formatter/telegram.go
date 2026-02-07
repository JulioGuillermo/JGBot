package formatter

import (
	"JGBot/formatter/ftools"
)

func FormatMD2Telegram(msg string) string {
	msg, codeBlocks := ftools.MapCodeBlocks(msg)
	msg, links := ftools.MapLinks(msg)

	msg = ftools.FormatHeaders(msg)
	msg = ftools.FormatTable(msg)
	msg = ftools.FormatList(msg)

	// By default telegram support bold, italic, strikethrough, etc.
	// msg = ftools.ProtectBold(msg)

	// msg = ftools.FormatItalic(msg, "_")
	// msg = ftools.FormatStrike(msg, "~")

	// msg = ftools.RestoreBold(msg, "*")

	msg = ftools.RestoreLinks(msg, links, true)
	msg = ftools.RestoreCodeBlocks(msg, codeBlocks, nil)

	return msg
}
