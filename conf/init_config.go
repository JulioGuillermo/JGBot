package conf

import (
	"JGBot/tools"
)

var Conf *Config

func InitConfig() error {
	var err error
	Conf, err = loadConfig()
	return err
}

func loadConfig() (*Config, error) {
	conf := newConfig()
	err := tools.ReadJSONFile(ConfigFile, conf)
	return conf, err
}
