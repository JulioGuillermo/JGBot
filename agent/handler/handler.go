package handler

import (
	"JGBot/agent/tools"
	"JGBot/log"
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type AgentHandler struct {
	ToolCall     tools.ToolCall
	OnToolCall   func(toolCall tools.ToolCall)
	OnToolResult func(toolResult tools.ToolResult)
}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{}
}

func (h *AgentHandler) HandleText(ctx context.Context, text string) {
}

func (h *AgentHandler) HandleLLMStart(ctx context.Context, prompts []string) {
}

func (h *AgentHandler) HandleLLMGenerateContentStart(ctx context.Context, ms []llms.MessageContent) {
}

func (h *AgentHandler) HandleLLMGenerateContentEnd(ctx context.Context, res *llms.ContentResponse) {
}

func (h *AgentHandler) HandleLLMError(ctx context.Context, err error) {
}

func (h *AgentHandler) HandleChainStart(ctx context.Context, inputs map[string]any) {
}

func (h *AgentHandler) HandleChainEnd(ctx context.Context, outputs map[string]any) {
}

func (h *AgentHandler) HandleChainError(ctx context.Context, err error) {
}

func (h *AgentHandler) HandleToolStart(ctx context.Context, input string) {
}

func (h *AgentHandler) HandleToolEnd(ctx context.Context, output string) {
	log.Info("TOOL RESULT", "output", output)
	if h.OnToolResult != nil {
		toolResult := tools.ToolResult{
			ToolCall: h.ToolCall,
			Output:   output,
			Error:    "",
		}
		h.OnToolResult(toolResult)
	}
}

func (h *AgentHandler) HandleToolError(ctx context.Context, error error) {
	log.Error("TOOL ERROR", "error", error)
	if h.OnToolResult != nil {
		toolResult := tools.ToolResult{
			ToolCall: h.ToolCall,
			Output:   "",
			Error:    error.Error(),
		}
		h.OnToolResult(toolResult)
	}
}

func (h *AgentHandler) HandleAgentAction(ctx context.Context, action schema.AgentAction) {
	h.ToolCall = tools.ToolCall{
		ID:    action.ToolID,
		Tool:  action.Tool,
		Input: action.ToolInput,
		// Log:   action.Log,
	}
	log.Info("TOOL CALL", "tool", action.Tool, "id", action.ToolID, "input", action.ToolInput, "log", action.Log)
	if h.OnToolCall != nil {
		h.OnToolCall(h.ToolCall)
	}
}

func (h *AgentHandler) HandleAgentFinish(ctx context.Context, finish schema.AgentFinish) {
}

func (h *AgentHandler) HandleRetrieverStart(ctx context.Context, query string) {
}

func (h *AgentHandler) HandleRetrieverEnd(ctx context.Context, query string, documents []schema.Document) {
}

func (h *AgentHandler) HandleStreamingFunc(ctx context.Context, chunk []byte) {
	// fmt.Print(string(chunk))
}
