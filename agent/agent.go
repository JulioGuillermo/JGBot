package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/input"
	"JGBot/agent/provider"
	"JGBot/agent/tools"
	"JGBot/log"
	"JGBot/session/sessionconf/sc"
	"JGBot/session/sessiondb"
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
	agentTools "github.com/tmc/langchaingo/tools"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
)

type Agent struct {
	ctx      context.Context
	provider llms.Model
}

func NewAgent(customTools ...agentTools.Tool) (*Agent, error) {
	agent := &Agent{}
	agent.ctx = context.Background()

	provider, err := provider.GetProvider(agent.ctx)
	if err != nil {
		return nil, err
	}
	agent.provider = provider

	return agent, nil
}

func (a *Agent) Respond(sessionConf *sc.SessionConf, history []*sessiondb.SessionMessage, message *sessiondb.SessionMessage, onResponse func(text, role, extra string) error, onReact func(msg uint, reaction string) error) {
	log.Info("Agent responding...")

	handler := handler.NewAgentHandler()
	handler.OnToolCall = func(toolCall tools.ToolCall) {
		onResponse("", "assistant", toolCall.ToJson())
	}
	handler.OnToolResult = func(toolResult tools.ToolResult) {
		onResponse("", "tool", toolResult.ToJson())
	}

	executor := a.getAgent(
		handler,
		tools.ReactionTool{
			CallbacksHandler: handler,
			OnReact:          onReact,
		},
		agentTools.Calculator{
			CallbacksHandler: handler,
		},
	)

	bytes, _ := json.Marshal(history)
	result, err := chains.Predict(
		a.ctx,
		executor,
		map[string]any{
			"input":       message.String(),
			"ChatHistory": string(bytes),
		},
	)
	fmt.Println()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
	onResponse(result, "assistant", "")
}

func (a *Agent) getAgent(handler callbacks.Handler, tools ...agentTools.Tool) *agents.Executor {
	agent := agents.NewOpenAIFunctionsAgent(
		a.provider,
		tools,
		agents.NewOpenAIOption().WithExtraMessages([]prompts.MessageFormatter{
			input.NewHistoryInput(),
		}),
	)
	return agents.NewExecutor(
		agent,
		agents.WithCallbacksHandler(handler),
	)
}
