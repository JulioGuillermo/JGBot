package conf

import (
	"JGBot/tools"
)

type Config struct {
	Database string
	LogLevel string

	Channels Channels

	Providers []Provider

	DefConf *DefConf
}

func newConfig() *Config {
	return &Config{
		Database: "",
		LogLevel: "Info",

		Channels: *newChannels(),

		Providers: []Provider{},
	}
}

func (conf *Config) GetChannelByName(name string) *Channel {
	switch name {
	case "Telegram":
		return &conf.Channels.Telegram
	case "WhatsApp":
		return &conf.Channels.Whatsapp
	default:
		return nil
	}
}

func (conf *Config) Save() error {
	return tools.WriteJSONFile(ConfigFile, conf)
}
