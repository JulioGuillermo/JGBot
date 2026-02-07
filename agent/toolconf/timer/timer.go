package timer

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
	"JGBot/timer"
	"context"
	"fmt"
	"strings"
)

type TimerInitializerConf struct {
	OnExecute func(ctx *ctxs.RespondCtx, args TimerArgs)
}

func (c *TimerInitializerConf) Name() string {
	return "timer"
}

func (c *TimerInitializerConf) listTimers(rCtx *ctxs.RespondCtx) string {
	jobs := timer.Timer.ListTimers(rCtx.Origin)
	if len(jobs) == 0 {
		return "No timers are active."
	}
	var sb strings.Builder
	sb.WriteString("Active timers:\n")
	for _, job := range jobs {
		fmt.Fprintf(&sb, "- %s: %s\n", job.Name, job.Description)
	}
	return sb.String()
}

func (c *TimerInitializerConf) readTimer(rCtx *ctxs.RespondCtx, name string) string {
	job := timer.Timer.GetTimer(rCtx.Origin, name)
	if job == nil {
		return fmt.Sprintf("Fail to read timer %s: not found", name)
	}
	return fmt.Sprintf("Name: %s\nDescription: %s\nSchedule: (%s)", job.Name, job.Description, job.Time.String())
}

func (c *TimerInitializerConf) addTimer(rCtx *ctxs.RespondCtx, args TimerArgs) string {
	var err error

	switch args.Type {
	case "timeout":
		err = timer.Timer.AddTimeout(rCtx.Origin, args.Name, args.Description, args.TimerTime.ToTime(), func() {
			c.OnExecute(rCtx, args)
		})
	case "alarm":
		err = timer.Timer.AddAlarm(rCtx.Origin, args.Name, args.Description, args.TimerTime.ToTime(), func() {
			c.OnExecute(rCtx, args)
		})
	default:
		return fmt.Sprintf("Fail to add timer %s: invalid type %s, must be 'timeout' or 'alarm'", args.Name, args.Type)
	}

	if err != nil {
		return fmt.Sprintf("Fail to add timer %s: %s", args.Name, err.Error())
	}

	return fmt.Sprintf("Timer %s added", args.Name)
}

func (c *TimerInitializerConf) removeTimer(rCtx *ctxs.RespondCtx, name string) string {
	err := timer.Timer.RemoveTimer(rCtx.Origin, name)
	if err != nil {
		return fmt.Sprintf("Fail to remove timer %s: %s", name, err.Error())
	}
	return fmt.Sprintf("Timer %s removed", name)
}

func (c *TimerInitializerConf) ToolInitializer(rCtx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[TimerArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Allows you to list, read, or execute timers.",
		ToolFunc: func(ctx context.Context, args TimerArgs) (string, error) {
			args.Action = strings.ToLower(args.Action)
			args.Type = strings.ToLower(args.Type)

			switch args.Action {
			case "list":
				return c.listTimers(rCtx), nil
			case "read":
				return c.readTimer(rCtx, args.Name), nil
			case "add":
				return c.addTimer(rCtx, args), nil
			case "remove":
				return c.removeTimer(rCtx, args.Name), nil
			}

			return "Invalid action, please use 'list', 'read', 'add' or 'remove'", nil
		},
	}
}
