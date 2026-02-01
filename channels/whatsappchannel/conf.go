package whatsappchannel

import "JGBot/conf"

type WhatsappConf struct {
	DBPath   string
	LogLevel string
}

func GetWhatsappConf() *WhatsappConf {
	config := &WhatsappConf{}
	change := false

	config.DBPath = conf.Conf.Channels.Whatsapp.Config["DBPath"]
	if config.DBPath == "" {
		config.DBPath = "db/whatsapp.db"
		change = true
	}

	config.LogLevel = conf.Conf.Channels.Whatsapp.LogLevel
	if config.LogLevel == "" {
		config.LogLevel = "info"
		change = true
	}

	if change {
		conf.Conf.Channels.Whatsapp.Config["DBPath"] = config.DBPath
		conf.Conf.Channels.Whatsapp.LogLevel = config.LogLevel
		conf.Conf.Save()
	}

	return config
}
