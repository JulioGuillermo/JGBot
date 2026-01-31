package whatsappchannel

import "JGBot/config"

type WhatsappConf struct {
	DBPath   string
	LogLevel string
}

func GetWhatsappConf() *WhatsappConf {
	conf := &WhatsappConf{}
	change := false

	conf.DBPath = config.Conf.Channels.Whatsapp.Config["DBPath"]
	if conf.DBPath == "" {
		conf.DBPath = "db/whatsapp.db"
		change = true
	}

	conf.LogLevel = config.Conf.Channels.Whatsapp.LogLevel
	if conf.LogLevel == "" {
		conf.LogLevel = "info"
		change = true
	}

	if change {
		config.Conf.Channels.Whatsapp.Config["DBPath"] = conf.DBPath
		config.Conf.Channels.Whatsapp.LogLevel = conf.LogLevel
		config.Conf.Save()
	}

	return conf
}
