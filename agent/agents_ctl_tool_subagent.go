package agent

import (
	"JGBot/agent/handler"
	"JGBot/agent/subagent"
	"JGBot/agent/subagent/subagenttool"
	"JGBot/ctxs"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

func (a *AgentsCtl) AddSubAgentTool(agent *subagent.SubAgent, handler *handler.AgentHandler, tools []tools.Tool, provider llms.Model, ctx *ctxs.RespondCtx) {
	subAgentToolConf := subagenttool.SubAgentInitializerConf{
		Ctx:      a.ctx,
		Handler:  handler,
		Tools:    tools,
		Provider: provider,
	}
	subAgentTool := subAgentToolConf.ToolInitializer(ctx)
	agent.AddTools(subAgentTool)
}
