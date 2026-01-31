package provider

import (
	"JGBot/config"

	"github.com/tmc/langchaingo/llms/mistral"
)

func GetMistral(conf config.Provider) (*mistral.Model, error) {
	opts := []mistral.Option{}

	if conf.ApiKey != nil {
		opts = append(opts, mistral.WithAPIKey(*conf.ApiKey))
	}
	if conf.Model != nil {
		opts = append(opts, mistral.WithModel(*conf.Model))
	}

	return mistral.New(opts...)
}
