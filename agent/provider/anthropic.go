package provider

import (
	"JGBot/conf"

	"github.com/tmc/langchaingo/llms/anthropic"
)

func GetAnthropic(conf conf.Provider) (*anthropic.LLM, error) {
	opts := []anthropic.Option{}

	if conf.ApiKey != nil {
		opts = append(opts, anthropic.WithToken(*conf.ApiKey))
	}
	if conf.Model != nil {
		opts = append(opts, anthropic.WithModel(*conf.Model))
	}

	return anthropic.New(opts...)
}
