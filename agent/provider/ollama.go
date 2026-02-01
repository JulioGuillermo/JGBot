package provider

import (
	"JGBot/conf"

	"github.com/tmc/langchaingo/llms/ollama"
)

func GetOllama(conf conf.Provider) (*ollama.LLM, error) {
	opts := []ollama.Option{}

	if conf.BaseUrl != nil {
		opts = append(opts, ollama.WithServerURL(*conf.BaseUrl))
	}
	if conf.Model != nil {
		opts = append(opts, ollama.WithModel(*conf.Model))
	}

	return ollama.New(opts...)
}
