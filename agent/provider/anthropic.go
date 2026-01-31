package provider

import (
	"JGBot/config"

	"github.com/tmc/langchaingo/llms/anthropic"
)

func GetAnthropic(conf config.Provider) (*anthropic.LLM, error) {
	opts := []anthropic.Option{}

	if conf.ApiKey != nil {
		opts = append(opts, anthropic.WithToken(*conf.ApiKey))
	}
	if conf.Model != nil {
		opts = append(opts, anthropic.WithModel(*conf.Model))
	}

	return anthropic.New(opts...)
}
