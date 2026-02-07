package formatter

import (
	"JGBot/formatter/ftools"
)

func FormatMD2WhatsApp(msg string) string {
	msg, codeBlocks := ftools.MapCodeBlocks(msg)
	msg, links := ftools.MapLinks(msg)

	msg = ftools.FormatHeaders(msg)
	msg = ftools.FormatTable(msg)
	msg = ftools.FormatList(msg)

	msg = ftools.ProtectBold(msg)

	msg = ftools.FormatItalic(msg, "_")
	msg = ftools.FormatStrike(msg, "~")

	msg = ftools.RestoreBold(msg, "*")

	msg = ftools.RestoreLinks(msg, links, false)
	msg = ftools.RestoreCodeBlocks(msg, codeBlocks, nil)

	return msg
}
