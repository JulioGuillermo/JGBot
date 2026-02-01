package tools

import (
	"context"

	"github.com/tmc/langchaingo/callbacks"
)

type Tool interface {
	Name() string
	Description() string
	SetHandler(handler callbacks.Handler)
	Call(ctx context.Context, input string) (string, error)
}
