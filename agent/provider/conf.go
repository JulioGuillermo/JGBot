package provider

import "JGBot/conf"

func GetConfig() []conf.Provider {
	return conf.Conf.Providers
}
