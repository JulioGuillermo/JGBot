package handler

import (
	"JGBot/agent/tools"
	"context"
	"fmt"

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
	// fmt.Println(text)
}

func (h *AgentHandler) HandleLLMStart(ctx context.Context, prompts []string) {
	// fmt.Println(prompts)
	fmt.Println()
}

func (h *AgentHandler) HandleLLMGenerateContentStart(ctx context.Context, ms []llms.MessageContent) {
	// fmt.Println(ms)
}

func (h *AgentHandler) HandleLLMGenerateContentEnd(ctx context.Context, res *llms.ContentResponse) {
	// fmt.Println(res.Choices)
}

func (h *AgentHandler) HandleLLMError(ctx context.Context, err error) {
	fmt.Println("\n\nLLM error:", err)
}

func (h *AgentHandler) HandleChainStart(ctx context.Context, inputs map[string]any) {
	// fmt.Println(inputs)
}

func (h *AgentHandler) HandleChainEnd(ctx context.Context, outputs map[string]any) {
	// fmt.Println(outputs)
	fmt.Println("### Chain end", outputs)
}

func (h *AgentHandler) HandleChainError(ctx context.Context, err error) {
	fmt.Println(err)
}

func (h *AgentHandler) HandleToolStart(ctx context.Context, input string) {
	// fmt.Println(input)
}

func (h *AgentHandler) HandleToolEnd(ctx context.Context, output string) {
	fmt.Printf("Tool Result: %s\n", output)
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
	fmt.Printf("ERROR: %s", error.Error())
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
		Log:   action.Log,
	}
	fmt.Printf("\nTool Call: %s[%s](Args: %s) Log: %s\n", action.Tool, action.ToolID, action.ToolInput, action.Log)
	if h.OnToolCall != nil {
		h.OnToolCall(h.ToolCall)
	}
}

func (h *AgentHandler) HandleAgentFinish(ctx context.Context, finish schema.AgentFinish) {
	fmt.Println("Agent end", finish)
}

func (h *AgentHandler) HandleRetrieverStart(ctx context.Context, query string) {
	// fmt.Println(query)
}

func (h *AgentHandler) HandleRetrieverEnd(ctx context.Context, query string, documents []schema.Document) {
	// fmt.Println(documents)
}

func (h *AgentHandler) HandleStreamingFunc(ctx context.Context, chunk []byte) {
	// fmt.Print(string(chunk))
}
