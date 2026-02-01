package tools

import (
	"JGBot/log"
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
)

type Tool[Args any] struct {
	Handler     callbacks.Handler
	ToolFunc    func(ctx context.Context, args Args) (string, error)
	name        string
	description string
}

func (t *Tool[Args]) Name() string {
	return t.name
}

func (t *Tool[Args]) Description() string {
	return t.description
}

func (t *Tool[Args]) Call(ctx context.Context, input string) (string, error) {
	if t.Handler != nil {
		t.Handler.HandleToolStart(ctx, input)
	}

	if t.ToolFunc == nil {
		return "Tool execution success", nil
	}

	args, err := t.getArgs(input)
	if err != nil {
		return t.fail(ctx, "ARGS ERROR: Fail to parse args from input JSON: %s", err)
	}

	result, err := t.ToolFunc(ctx, args)
	if err != nil {
		return t.fail(ctx, "TOOL ERROR: Fail to execute tool: %s", err)
	}

	return t.success(ctx, result)
}

func (t *Tool[Args]) fail(ctx context.Context, msg string, err error) (string, error) {
	log.Error(msg, "error", err)
	if t.Handler != nil {
		t.Handler.HandleToolError(ctx, err)
	}
	return fmt.Sprintf(msg, err.Error()), nil
}

func (t *Tool[Args]) success(ctx context.Context, output string) (string, error) {
	if t.Handler != nil {
		t.Handler.HandleToolEnd(ctx, output)
	}
	return output, nil
}

func (t *Tool[Args]) getArgs(input string) (Args, error) {
	input = GetInputString(input)

	var args Args
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return args, err
	}
	return args, nil
}
