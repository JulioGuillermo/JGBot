package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/input"
	"JGBot/agent/provider"
	"JGBot/agent/tools"
	"JGBot/session/sessiondb"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
	agentTools "github.com/tmc/langchaingo/tools"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
)

type Agent struct {
	logger   *slog.Logger
	ctx      context.Context
	provider llms.Model
	tools    []agentTools.Tool
}

func NewAgent(logger *slog.Logger, customTools ...agentTools.Tool) (*Agent, error) {
	agent := &Agent{
		logger: logger,
	}
	agent.ctx = context.Background()

	provider, err := provider.GetProvider(agent.ctx)
	if err != nil {
		return nil, err
	}
	agent.provider = provider

	agent.tools, err = tools.GetTools(customTools)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (a *Agent) Respond(history []*sessiondb.SessionMessage, message *sessiondb.SessionMessage, onResponse func(text, role, extra string) error, onReact func(msg uint, reaction string) error) {
	a.logger.Info("Agent responding...")

	handler := handler.NewAgentHandler()
	handler.OnToolCall = func(toolCall tools.ToolCall) {
		onResponse("", "assistant", toolCall.ToJson())
	}
	handler.OnToolResult = func(toolResult tools.ToolResult) {
		if toolResult.Error != "" {
			onResponse("", "assistant", toolResult.ToJson())
		} else {
			onResponse("", "assistant", toolResult.ToJson())
		}
	}

	agent := agents.NewOpenAIFunctionsAgent(
		a.provider,
		append(
			a.tools,
			tools.NewReactionTool(onReact),
		),
		agents.NewOpenAIOption().WithExtraMessages([]prompts.MessageFormatter{
			input.NewHistoryInput(),
		}),
	)
	executor := agents.NewExecutor(
		agent,
		agents.WithCallbacksHandler(handler),
	)
	bytes, _ := json.Marshal(history)

	result, err := chains.Predict(
		a.ctx,
		executor,
		map[string]any{
			"input":       message.String(),
			"ChatHistory": string(bytes),
		},
		// message.String(),
	)
	fmt.Println()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
	onResponse(result, "assistant", "")
}
