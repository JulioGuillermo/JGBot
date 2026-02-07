package agent

import (
	"JGBot/agent/toolconf/cron"
	"JGBot/ctxs"
	"JGBot/log"
	"JGBot/session/sessiondb"
	"fmt"

	"github.com/tmc/langchaingo/tools"
)

func (a *AgentsCtl) AddCronTool(ctx *ctxs.RespondCtx) tools.Tool {
	cronToolConf := cron.CronInitializerConf{
		OnExecute: func(ctx *ctxs.RespondCtx, args cron.CronArgs) {
			history, err := ctx.GetHistory()
			if err != nil {
				log.Error("CRON TOOL ERROR: Get history error", "err", err)
				return
			}

			msg := fmt.Sprintf("CRON EXECUTION: %s\n\nSCHEDULE: %s\n\nDESCRIPTION: %s\n\nMESSAGE: %s", args.Name, args.Schedule.String(), args.Description, args.Message)

			newCtx := ctx.Copy()
			newCtx.History = history
			newCtx.Message = &sessiondb.SessionMessage{
				Channel:    ctx.Message.Channel,
				ChatID:     ctx.Message.ChatID,
				ChatName:   ctx.Message.ChatName,
				SenderID:   ctx.Message.SenderID,
				SenderName: "tool",
				Role:       "tool",
				MessageID:  ctx.Message.MessageID,
				Message:    msg,
				Extra:      "",
			}

			err = a.Respond(newCtx)
			if err != nil {
				log.Error("CRON TOOL ERROR: Respond error", "err", err)
			}
		},
	}
	cronTool := cronToolConf.ToolInitializer(ctx)
	return cronTool
}
