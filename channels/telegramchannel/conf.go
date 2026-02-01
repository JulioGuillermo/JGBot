package telegramchannel

import "JGBot/conf"

type TelegramConf struct {
	Token    string
	LogLevel string
}

func GetTelegramConf() *TelegramConf {
	config := &TelegramConf{}
	changed := false

	config.Token = conf.Conf.Channels.Telegram.Config["token"]
	if config.Token == "" {
		config.Token = ""
		changed = true
	}

	config.LogLevel = conf.Conf.Channels.Telegram.LogLevel
	if config.LogLevel == "" {
		config.LogLevel = "info"
		changed = true
	}

	if changed {
		conf.Conf.Channels.Telegram.Config["token"] = config.Token
		conf.Conf.Channels.Telegram.LogLevel = config.LogLevel
		conf.Conf.Save()
	}
	return config
}
