package formatter

func FormatMD2TelegramHTML(msg string) string {
	msg = FormatMD2Telegram(msg)
	msg = FormatMD2HTML(msg)
	return msg
}
