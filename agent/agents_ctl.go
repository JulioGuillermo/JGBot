package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/prompt"
	"JGBot/agent/provider"
	"JGBot/agent/subagent"
	"JGBot/agent/toolconf"
	"JGBot/agent/toolconf/tools_conf"
	"JGBot/agent/tools"
	channelsdomain "JGBot/channels/domain"
	"JGBot/ctxs"
	"JGBot/log"
	"JGBot/session/sessionconf"
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type AgentsCtl struct {
	ctx        context.Context
	toolsConf  map[string]tools_conf.ToolInitializerConf
	sessionCtl *sessionconf.SessionCtl
	channelCtl channelsdomain.ChannelController
}

func NewAgentsCtl() (*AgentsCtl, error) {
	agent := &AgentsCtl{}
	agent.ctx = context.Background()

	provider.InitProviders(agent.ctx)

	agent.toolsConf = toolconf.GetToolMap()

	return agent, nil
}

func (a *AgentsCtl) SetDependencies(sessionCtl *sessionconf.SessionCtl, channelCtl channelsdomain.ChannelController) {
	a.sessionCtl = sessionCtl
	a.channelCtl = channelCtl
}

func (a *AgentsCtl) getProvider(providerName string) (llms.Model, error) {
	prov, ok := provider.Providers[providerName]
	if !ok {
		return nil, fmt.Errorf("Provider %s not found", providerName)
	}
	return prov, nil
}

func (a *AgentsCtl) Respond(ctx *ctxs.RespondCtx) error {
	defer ctx.Status(channelsdomain.Normal)
	log.Info("Agent responding...")

	sysPrompt := prompt.GetSystemPrompt(ctx.SessionConf)

	handler := handler.NewAgentHandler()
	handler.OnToolCall = func(toolCall tools.ToolCall) {
		ctx.Status(channelsdomain.Writing)
		ctx.OnResponse("", "assistant", toolCall.ToJson())
	}
	handler.OnToolResult = func(toolResult tools.ToolResult) {
		ctx.Status(channelsdomain.Writing)
		ctx.OnResponse("", "tool", toolResult.ToJson())
	}

	provider, err := a.getProvider(ctx.SessionConf.Provider)
	if err != nil {
		return err
	}

	agent := &subagent.SubAgent{
		Name:         "Main Agent",
		Ctx:          a.ctx,
		Handler:      handler,
		Provider:     provider,
		MaxIters:     max(ctx.SessionConf.AgentMaxIters, 3),
		SystemPrompt: sysPrompt,
	}

	a.AddTools(agent, handler, provider, ctx)

	agent.Init()

	result, err := agent.Run(ctx.History, ctx.Message)
	if err != nil {
		return err
	}

	log.Info("AGENT RESPONDED", "result", result)
	return ctx.OnResponse(removeThink(result), "assistant", "")
}

func removeThink(text string) string {
	const Start = "<think>"
	const End = "</think>"
	if !strings.HasPrefix(text, Start) {
		return text
	}

	idx := strings.Index(text, End)
	if idx == -1 {
		return text
	}

	idx += len(End)
	if idx >= len(text) {
		return ""
	}

	text = text[idx:]
	return strings.TrimSpace(text)
}
