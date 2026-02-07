package agent

import (
	"JGBot/agent/toolconf/timer"
	"JGBot/ctxs"
	"JGBot/log"
	"JGBot/session/sessiondb"
	"fmt"

	"github.com/tmc/langchaingo/tools"
)

func (a *AgentsCtl) AddTimerTool(ctx *ctxs.RespondCtx) tools.Tool {
	timerToolConf := timer.TimerInitializerConf{
		OnExecute: func(ctx *ctxs.RespondCtx, args timer.TimerArgs) {
			history, err := ctx.GetHistory()
			if err != nil {
				log.Error("TIMER TOOL ERROR: Get history error", "err", err)
				return
			}

			var msg string
			switch args.Type {
			case "timeout":
				msg = fmt.Sprintf("TIMEOUT EXECUTION: %s\n\nSCHEDULE: %s\n\nDESCRIPTION: %s\n\nMESSAGE: %s", args.Name, args.TimerTime.String(), args.Description, args.Message)
			case "alarm":
				msg = fmt.Sprintf("ALARM EXECUTION: %s\n\nSCHEDULE: %s\n\nDESCRIPTION: %s\n\nMESSAGE: %s", args.Name, args.TimerTime.String(), args.Description, args.Message)
			default:
				msg = fmt.Sprintf("TIMER EXECUTION: %s\n\nSCHEDULE: %s\n\nDESCRIPTION: %s\n\nMESSAGE: %s", args.Name, args.TimerTime.String(), args.Description, args.Message)
			}

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
				log.Error("TIMER TOOL ERROR: Respond error", "err", err)
			}
		},
	}
	timerTool := timerToolConf.ToolInitializer(ctx)
	return timerTool
}
