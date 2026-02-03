package subagenttool

import (
	"JGBot/agent/handler"
	"JGBot/agent/prompt"
	"JGBot/agent/subagent"
	"JGBot/agent/tools"
	"JGBot/ctxs"
	"JGBot/log"
	"context"

	"github.com/tmc/langchaingo/llms"
	langchainTools "github.com/tmc/langchaingo/tools"
)

type SubAgentArgs struct {
	Name string `json:"name" description:"The name of the subagent to execute."`
	Task string `json:"task" description:"The task to pass to the subagent. It must describe the task in a way that the subagent can understand it."`
}

type SubAgentInitializerConf struct {
	Ctx      context.Context
	Handler  *handler.AgentHandler
	Tools    []langchainTools.Tool
	Provider llms.Model
}

func (c *SubAgentInitializerConf) Name() string {
	return "subagent"
}

func (c *SubAgentInitializerConf) ToolInitializer(rCtx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[SubAgentArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Allows you to execute a subagent with a task. The subagent will execute the task and return the result.",
		ToolFunc: func(ctx context.Context, args SubAgentArgs) (string, error) {
			log.Info("SubAgent responding...", "name", args.Name)

			sysPrompt := prompt.GetSubAgentPrompt(rCtx.SessionConf)

			agent := &subagent.SubAgent{
				Name:         "SubAgent",
				Ctx:          c.Ctx,
				Handler:      c.Handler,
				Provider:     c.Provider,
				MaxIters:     max(rCtx.SessionConf.AgentMaxIters, 3),
				SystemPrompt: sysPrompt,
			}

			agent.AddTools(c.Tools...)
			agent.Init()

			result, err := agent.RunSimple(args.Task)
			if err != nil {
				return "", err
			}

			log.Info("SUB AGENT RESPONDED", "subagent", args.Name, "result", result)
			return result, nil
		},
	}
}
