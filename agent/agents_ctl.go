package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/provider"
	"JGBot/agent/tools"
	"JGBot/log"
	"JGBot/session/sessionconf/sc"
	"JGBot/session/sessiondb"
	"context"
	"fmt"

	agentTools "github.com/tmc/langchaingo/tools"

	"github.com/tmc/langchaingo/llms"
)

type AgentsCtl struct {
	ctx      context.Context
	provider map[string]llms.Model
}

func NewAgentsCtl(customTools ...agentTools.Tool) (*AgentsCtl, error) {
	agent := &AgentsCtl{}
	agent.ctx = context.Background()
	agent.provider = provider.GetProviders(agent.ctx)

	return agent, nil
}

func (a *AgentsCtl) getProvider(provider string) (llms.Model, error) {
	prov, ok := a.provider[provider]
	if !ok {
		return nil, fmt.Errorf("Provider %s not found", provider)
	}
	return prov, nil
}

func (a *AgentsCtl) Respond(sessionConf *sc.SessionConf, history []*sessiondb.SessionMessage, message *sessiondb.SessionMessage, onResponse func(text, role, extra string) error, onReact func(msg uint, reaction string) error) error {
	log.Info("Agent responding...")

	handler := handler.NewAgentHandler()
	handler.OnToolCall = func(toolCall tools.ToolCall) {
		onResponse("", "assistant", toolCall.ToJson())
	}
	handler.OnToolResult = func(toolResult tools.ToolResult) {
		onResponse("", "tool", toolResult.ToJson())
	}

	provider, err := a.getProvider(sessionConf.Provider)
	if err != nil {
		return err
	}

	agent := &Agent{
		Name:     "Main Agent",
		Ctx:      a.ctx,
		Handler:  handler,
		Provider: provider,
	}
	agent.AddTools(
		tools.ReactionTool{
			CallbacksHandler: handler,
			OnReact:          onReact,
		},
		agentTools.Calculator{
			CallbacksHandler: handler,
		},
	)
	agent.Init()

	result, err := agent.Run(history, message)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return onResponse(result, "assistant", "")
}
