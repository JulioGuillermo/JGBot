package provider

import (
	"JGBot/config"
	"context"

	"github.com/tmc/langchaingo/llms/googleai"
)

func GetGoogle(ctx context.Context, conf config.Provider) (*googleai.GoogleAI, error) {
	opts := []googleai.Option{}

	if conf.ApiKey != nil {
		opts = append(opts, googleai.WithAPIKey(*conf.ApiKey))
	}
	if conf.Model != nil {
		opts = append(opts, googleai.WithDefaultModel(*conf.Model))
	}

	return googleai.New(ctx, opts...)
}
