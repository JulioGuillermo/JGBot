package conf

import (
	"JGBot/tools"
)

type Config struct {
	Database string
	LogLevel string

	Channels Channels

	Providers []Provider
}

func newConfig() *Config {
	return &Config{
		Database: "",
		LogLevel: "Info",

		Channels: *newChannels(),

		Providers: []Provider{},
	}
}

func (conf *Config) Save() error {
	return tools.WriteJSONFile(ConfigFile, conf)
}
