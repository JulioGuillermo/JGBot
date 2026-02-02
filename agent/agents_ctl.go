package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/provider"
	"JGBot/agent/toolconf"
	"JGBot/agent/tools"
	"JGBot/log"
	"JGBot/session/sessionconf/sc"
	"JGBot/session/sessiondb"
	"JGBot/skill"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type AgentsCtl struct {
	ctx       context.Context
	providers map[string]llms.Model
	toolsConf map[string]toolconf.ToolInitializerConf
	skills    []*skill.Skill
}

func NewAgentsCtl(skills []*skill.Skill) (*AgentsCtl, error) {
	agent := &AgentsCtl{
		skills: skills,
	}
	agent.ctx = context.Background()
	agent.providers = provider.GetProviders(agent.ctx)
	agent.toolsConf = toolconf.GetToolMap()

	return agent, nil
}

func (a *AgentsCtl) getProvider(provider string) (llms.Model, error) {
	prov, ok := a.providers[provider]
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
		MaxIters: max(sessionConf.AgentMaxIters, 3),
	}

	for _, toolConf := range sessionConf.Tools {
		if !toolConf.Enabled {
			continue
		}

		toolInitializer, ok := a.toolsConf[toolConf.Name]
		if !ok {
			log.Warn("Tool not found", "tool", toolConf.Name)
			continue
		}

		tool := toolInitializer.ToolInitializer(sessionConf, history, message, onResponse, onReact)
		tool.SetHandler(handler)

		agent.AddTools(tool)
	}
	agent.Init()

	result, err := agent.Run(history, message)
	if err != nil {
		return err
	}

	log.Info("AGENT RESPONDED", "result", result)
	return onResponse(result, "assistant", "")
}
