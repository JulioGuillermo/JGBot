package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/subagent"
	"JGBot/agent/toolconf/admin"
	"JGBot/ctxs"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
	lcTools "github.com/tmc/langchaingo/tools"
)

func (a *AgentsCtl) AddTools(agent *subagent.SubAgent, handler *handler.AgentHandler, provider llms.Model, ctx *ctxs.RespondCtx) {
	tools := a.GetTools(ctx, handler)

	agent.AddTools(tools...)

	if conf := ctx.SessionConf.GetToolConf("subagent"); conf != nil && conf.Enabled {
		a.AddSubAgentTool(agent, handler, tools, provider, ctx)
	}
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

	if !ctx.IsAdmin {
		return tools
	}

	listSessions := admin.NewAdminListSessionsInitializerConf().ToolInitializer(ctx)
	listSessions.SetHandler(handler)
	tools = append(tools, listSessions)

	sendMessage := admin.NewAdminSendMessageInitializerConf().ToolInitializer(ctx)
	sendMessage.SetHandler(handler)
	tools = append(tools, sendMessage)

	return tools
}
