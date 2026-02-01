package conf

import (
	"JGBot/tools"
	"os"
)

var Conf *Config

func InitConfig() error {
	var err error

	Conf, err = loadConfig()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if Conf == nil {
		Conf = newConfig()
	}

	return Conf.Save()
}

func loadConfig() (*Config, error) {
	var conf Config
	err := tools.ReadJSONFile(ConfigFile, &conf)
	return &conf, err
}
