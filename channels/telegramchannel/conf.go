package telegramchannel

import "JGBot/config"

type TelegramConf struct {
	Token    string
	LogLevel string
}

func GetTelegramConf() *TelegramConf {
	conf := &TelegramConf{}
	changed := false

	conf.Token = config.Conf.Channels.Telegram.Config["token"]
	if conf.Token == "" {
		conf.Token = ""
		changed = true
	}

	conf.LogLevel = config.Conf.Channels.Telegram.LogLevel
	if conf.LogLevel == "" {
		conf.LogLevel = "info"
		changed = true
	}

	if changed {
		config.Conf.Channels.Telegram.Config["token"] = conf.Token
		config.Conf.Channels.Telegram.LogLevel = conf.LogLevel
		config.Conf.Save()
	}
	return conf
}
