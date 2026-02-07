package provider

import (
	toolargs "JGBot/agent/tool_args"
	"context"

	"github.com/tmc/langchaingo/llms"
)

type ProviderWrapper struct {
	model llms.Model
}

func (p *ProviderWrapper) GenerateContent(ctx context.Context, messages []llms.MessageContent, opts ...llms.CallOption) (*llms.ContentResponse, error) {
	// return p.model.GenerateContent(ctx, messages, opts...)
	resp, err := p.model.GenerateContent(ctx, messages, opts...)
	if err != nil {
		return nil, err
	}
	for i := range resp.Choices {
		if resp.Choices[i].FuncCall != nil {
			resp.Choices[i].FuncCall.Arguments = toolargs.GetMsgToolCallArg(resp.Choices[i].FuncCall.Arguments)
		}
		for j := range resp.Choices[i].ToolCalls {
			if resp.Choices[i].ToolCalls[j].FunctionCall != nil {
				resp.Choices[i].ToolCalls[j].FunctionCall.Arguments = toolargs.GetMsgToolCallArg(resp.Choices[i].ToolCalls[j].FunctionCall.Arguments)
			}
		}
	}
	return resp, nil
}

func (p *ProviderWrapper) Call(ctx context.Context, prompt string, opts ...llms.CallOption) (string, error) {
	return p.model.Call(ctx, prompt, opts...)
}
