package whatsappchannel

import "JGBot/conf"

type WhatsappConf struct {
	DBPath string
}

func GetWhatsappConf() *WhatsappConf {
	config := &WhatsappConf{}
	change := false

	config.DBPath = conf.Conf.Channels.Whatsapp.Config["DBPath"]
	if config.DBPath == "" {
		config.DBPath = "db/whatsapp.db"
		change = true
	}

	if change {
		conf.Conf.Channels.Whatsapp.Config["DBPath"] = config.DBPath
		conf.Conf.Save()
	}

	return config
}
