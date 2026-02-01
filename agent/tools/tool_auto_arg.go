package tools

import (
	"JGBot/agent/tools/args"
	"JGBot/agent/tools/templ"
	"JGBot/log"
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
)

type ToolAutoArgs[Args any] struct {
	Handler         callbacks.Handler
	ToolName        string
	ToolDescription string
	ToolFunc        func(ctx context.Context, args Args) (string, error)
}

func NewToolAutoArgs[Args any](toolName, toolDescription string, toolFunc func(ctx context.Context, args Args) (string, error)) *ToolAutoArgs[Args] {
	return &ToolAutoArgs[Args]{
		ToolName:        toolName,
		ToolDescription: toolDescription,
		ToolFunc:        toolFunc,
	}
}

func (t *ToolAutoArgs[Args]) Name() string {
	return t.ToolName
}

func (t *ToolAutoArgs[Args]) Description() string {
	var argsVar Args
	argsDescription := args.GetArgsMetaDataString(argsVar)
	return templ.GetNativeToolDescription(t.ToolName, t.ToolDescription, argsDescription)
}

func (t *ToolAutoArgs[Args]) SetHandler(handler callbacks.Handler) {
	t.Handler = handler
}

func (t *ToolAutoArgs[Args]) Call(ctx context.Context, input string) (string, error) {
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

func (t *ToolAutoArgs[Args]) fail(ctx context.Context, msg string, err error) (string, error) {
	log.Error(msg, "error", err)
	if t.Handler != nil {
		t.Handler.HandleToolError(ctx, err)
	}
	return fmt.Sprintf(msg, err.Error()), nil
}

func (t *ToolAutoArgs[Args]) success(ctx context.Context, output string) (string, error) {
	if t.Handler != nil {
		t.Handler.HandleToolEnd(ctx, output)
	}
	return output, nil
}

func (t *ToolAutoArgs[Args]) getArgs(input string) (Args, error) {
	input = GetInputString(input)

	var args Args
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return args, err
	}
	return args, nil
}
