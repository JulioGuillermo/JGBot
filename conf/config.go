package conf

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database string
	LogLevel string

	Channels Channels

	Provider Provider
}

func newConfig() *Config {
	return &Config{
		Database: "",
		LogLevel: "Info",

		Channels: *newChannels(),

		Provider: Provider{},
	}
}

func (conf *Config) Save() error {
	jsonFile, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(byteValue)
	if err != nil {
		return err
	}

	return nil
}
