package agent

import (
	"JGBot/agent/provider"
	"JGBot/agent/tools"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/memory"
	agentTools "github.com/tmc/langchaingo/tools"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
)

type Agent struct {
	ctx      context.Context
	provider llms.Model
	tools    []agentTools.Tool
}

func NewAgent(customTools ...agentTools.Tool) (*Agent, error) {
	agent := &Agent{}
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

func (a *Agent) Run() {
	memory.NewChatMessageHistory()
	agent := agents.NewConversationalAgent(a.provider, a.tools)
	executor := agents.NewExecutor(agent)
	result, err := chains.Run(a.ctx, executor, "What is your name?")
	fmt.Println(err)
	fmt.Println(result)
}
