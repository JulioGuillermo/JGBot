package telegramchannel

import "JGBot/conf"

type TelegramConf struct {
	Token string
}

func GetTelegramConf() *TelegramConf {
	config := &TelegramConf{}
	changed := false

	config.Token = conf.Conf.Channels.Telegram.Config["token"]
	if config.Token == "" {
		config.Token = ""
		changed = true
	}

	if changed {
		conf.Conf.Channels.Telegram.Config["token"] = config.Token
		conf.Conf.Save()
	}
	return config
}
