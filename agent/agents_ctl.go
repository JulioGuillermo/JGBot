package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/prompt"
	"JGBot/agent/provider"
	"JGBot/agent/subagent"
	"JGBot/agent/subagent/subagenttool"
	"JGBot/agent/toolconf"
	"JGBot/agent/toolconf/tools_conf"
	"JGBot/agent/tools"
	"JGBot/ctxs"
	"JGBot/log"
	"context"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
	lcTools "github.com/tmc/langchaingo/tools"
)

type AgentsCtl struct {
	ctx       context.Context
	toolsConf map[string]tools_conf.ToolInitializerConf
}

func NewAgentsCtl() (*AgentsCtl, error) {
	agent := &AgentsCtl{}
	agent.ctx = context.Background()

	provider.InitProviders(agent.ctx)

	agent.toolsConf = toolconf.GetToolMap()

	return agent, nil
}

func (a *AgentsCtl) getProvider(providerName string) (llms.Model, error) {
	prov, ok := provider.Providers[providerName]
	if !ok {
		return nil, fmt.Errorf("Provider %s not found", providerName)
	}
	return prov, nil
}

func (a *AgentsCtl) Respond(ctx *ctxs.RespondCtx) error {
	log.Info("Agent responding...")

	sysPrompt := prompt.GetSystemPrompt(ctx.SessionConf)

	handler := handler.NewAgentHandler()
	handler.OnToolCall = func(toolCall tools.ToolCall) {
		ctx.OnResponse("", "assistant", toolCall.ToJson())
	}
	handler.OnToolResult = func(toolResult tools.ToolResult) {
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

	tools := a.GetTools(ctx, handler)
	agent.AddTools(tools...)

	if conf := ctx.SessionConf.GetToolConf("subagent"); conf != nil && conf.Enabled {
		subAgentToolConf := subagenttool.SubAgentInitializerConf{
			Ctx:      a.ctx,
			Handler:  handler,
			Tools:    tools,
			Provider: provider,
		}
		subAgentTool := subAgentToolConf.ToolInitializer(ctx)
		agent.AddTools(subAgentTool)
	}

	agent.Init()

	result, err := agent.Run(ctx.History, ctx.Message)
	if err != nil {
		return err
	}

	log.Info("AGENT RESPONDED", "result", result)
	return ctx.OnResponse(result, "assistant", "")
}

func (a *AgentsCtl) GetTools(ctx *ctxs.RespondCtx, handler callbacks.Handler) []lcTools.Tool {
	tools := make([]lcTools.Tool, 0)

	for _, toolConf := range ctx.SessionConf.Tools {
		if !toolConf.Enabled {
			continue
		}

		toolInitializer, ok := a.toolsConf[toolConf.Name]
		if !ok {
			continue
		}

		tool := toolInitializer.ToolInitializer(ctx)
		tool.SetHandler(handler)

		tools = append(tools, tool)
	}

	return tools
}
