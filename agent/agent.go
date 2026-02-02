package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/input"
	"JGBot/session/sessiondb"
	"context"
	"encoding/json"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
)

type Agent struct {
	Ctx          context.Context
	Name         string
	Handler      *handler.AgentHandler
	Provider     llms.Model
	MaxIters     int
	SystemPrompt string
	tools        []tools.Tool
	agent        agents.Agent
	executor     *agents.Executor
}

func (a *Agent) AddTools(tool ...tools.Tool) {
	if a.tools == nil {
		a.tools = tool
		return
	}
	a.tools = append(a.tools, tool...)
}

func (a *Agent) Init() {
	a.initAgent()
	a.initExecutor()
}

func (a *Agent) initAgent() {
	a.agent = agents.NewOpenAIFunctionsAgent(
		a.Provider,
		a.tools,
		agents.NewOpenAIOption().
			WithSystemMessage(a.SystemPrompt),
		agents.NewOpenAIOption().
			WithExtraMessages([]prompts.MessageFormatter{
				input.NewHistoryInput(),
			}),
	)
}

func (a *Agent) initExecutor() {
	a.executor = agents.NewExecutor(
		a.agent,
		agents.WithCallbacksHandler(a.Handler),
		agents.WithMaxIterations(a.MaxIters),
	)
}

func (a *Agent) Run(history []*sessiondb.SessionMessage, message *sessiondb.SessionMessage) (string, error) {
	bytes, _ := json.Marshal(history)
	return chains.Predict(
		a.Ctx,
		a.executor,
		map[string]any{
			"input":       message.String(),
			"ChatHistory": string(bytes),
		},
	)
}
