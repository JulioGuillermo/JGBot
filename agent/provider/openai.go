package provider

import (
	"JGBot/conf"

	"github.com/tmc/langchaingo/llms/openai"
)

func GetOpenAI(conf conf.Provider) (*openai.LLM, error) {
	opts := []openai.Option{}

	if conf.BaseUrl != nil {
		opts = append(opts, openai.WithBaseURL(*conf.BaseUrl))
	}
	if conf.ApiKey != nil {
		opts = append(opts, openai.WithToken(*conf.ApiKey))
	}
	if conf.Model != nil {
		opts = append(opts, openai.WithModel(*conf.Model))
	}

	return openai.New(opts...)
}
