package provider

import "JGBot/config"

func GetConfig() config.Provider {
	return config.Conf.Provider
}
