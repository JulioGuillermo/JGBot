package conf

import (
	"encoding/json"
	"io"
	"os"
)

var Conf *Config

func InitConfig() error {
	var err error

	err = os.MkdirAll(ConfigDir, 0755)
	if err != nil {
		return err
	}

	Conf, err = loadConfig()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if Conf != nil {
		return nil
	}

	Conf = newConfig()
	return Conf.Save()
}

func loadConfig() (*Config, error) {
	jsonFile, err := os.Open(ConfigFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		return nil, err
	}
	conf.Save()
	return &conf, nil
}
