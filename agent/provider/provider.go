package provider

import (
	"JGBot/conf"
	"JGBot/log"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

var Providers map[string]llms.Model

func InitProviders(ctx context.Context) {
	conf := GetConfig()

	Providers = map[string]llms.Model{}
	var prov llms.Model
	var err error
	for _, conf := range conf {
		prov, err = GetProvider(ctx, conf)
		if err != nil {
			log.Error("Fail to initialize provider", "provider", conf.Name, "error", err)
			continue
		}
		Providers[conf.Name] = &ProviderWrapper{model: prov}
	}
}

func GetProvider(ctx context.Context, conf conf.Provider) (llms.Model, error) {
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
