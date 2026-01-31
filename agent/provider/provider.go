package provider

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

func GetProvider(ctx context.Context) (llm llms.Model, err error) {
	conf := GetConfig()

	switch conf.Type {
	case "openai":
		return GetOpenAI(conf)
	case "google":
		return GetGoogle(ctx, conf)
	case "anthropic":
		return GetAnthropic(conf)
	case "ollama":
		return GetOllama(conf)
	case "mistral":
		return GetMistral(conf)
	}

	return nil, fmt.Errorf("Invalid provider type %s", conf.Type)
}
