package timer

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
	taskdomain "JGBot/task/domain"
	taskports "JGBot/task/ports"
	"context"
	"fmt"
	"strings"
)

type TimerInitializerConf struct{}

func (c *TimerInitializerConf) Name() string {
	return "timer"
}

func (c *TimerInitializerConf) listTimers(rCtx *ctxs.RespondCtx, args TimerArgs) string {
	origin := rCtx.GetOrigin(args.Session)
	if rCtx.IsAdmin && args.Session != "" {
		origin = args.Session
	}

	jobs := taskports.TimerService.ListTimers(origin)
	if len(jobs) == 0 {
		return "Active timers:\nTimer list is empty."
	}
	var sb strings.Builder
	sb.WriteString("Active timers:\n")
	for _, job := range jobs {
		fmt.Fprintf(&sb, "- %s: %s\n", job.Name, job.Description)
	}
	return sb.String()
}

func (c *TimerInitializerConf) readTimer(rCtx *ctxs.RespondCtx, args TimerArgs) string {
	origin := rCtx.GetOrigin(args.Session)

	job := taskports.TimerService.GetTimer(origin, args.Name)
	if job == nil {
		return fmt.Sprintf("Fail to read timer %s: not found", args.Name)
	}
	task := job.Task()
	return fmt.Sprintf("Name: %s\nDescription: %s\nSchedule: (%s)", task.Name, task.Description, job.GetSchedule())
}

func (c *TimerInitializerConf) addTimer(rCtx *ctxs.RespondCtx, args TimerArgs) string {
	ctx, err := rCtx.GetSessionCtx(args.Session)

	switch args.Type {
	case "timeout":
		err = taskports.TimerService.AddTimeout(
			&taskdomain.Task{
				TaskOriginInfo: taskdomain.TaskOriginInfo{
					Origin:    ctx.Origin,
					Channel:   ctx.Channel,
					ChatID:    ctx.ChatID,
					ChatName:  ctx.ChatName,
					SenderID:  ctx.Message.SenderID,
					MessageID: ctx.Message.MessageID,
				},
				TaskInfo: taskdomain.TaskInfo{
					Name:        args.Name,
					Description: args.Description,
					Message:     args.Message,
				},
			},
			taskdomain.TimerTime(args.TimerTime),
		)
	case "alarm":
		err = taskports.TimerService.AddAlarm(
			&taskdomain.Task{
				TaskOriginInfo: taskdomain.TaskOriginInfo{
					Origin:    ctx.Origin,
					Channel:   ctx.Channel,
					ChatID:    ctx.ChatID,
					ChatName:  ctx.ChatName,
					SenderID:  ctx.Message.SenderID,
					MessageID: ctx.Message.MessageID,
				},
				TaskInfo: taskdomain.TaskInfo{
					Name:        args.Name,
					Description: args.Description,
					Message:     args.Message,
				},
			},
			taskdomain.TimerTime(args.TimerTime),
		)
	default:
		return fmt.Sprintf("Fail to add timer %s: invalid type %s, must be 'timeout' or 'alarm'", args.Name, args.Type)
	}

	if err != nil {
		return fmt.Sprintf("Fail to add timer %s: %s", args.Name, err.Error())
	}

	return fmt.Sprintf("Timer %s added", args.Name)
}

func (c *TimerInitializerConf) removeTimer(rCtx *ctxs.RespondCtx, args TimerArgs) string {
	origin := rCtx.GetOrigin(args.Session)

	err := taskports.TimerService.RemoveTimer(origin, args.Name)
	if err != nil {
		return fmt.Sprintf("Fail to remove timer %s: %s", args.Name, err.Error())
	}
	return fmt.Sprintf("Timer %s removed", args.Name)
}

func (c *TimerInitializerConf) ToolInitializer(rCtx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[TimerArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Allows you to list, read, or execute timers.",
		IsAdmin:         rCtx.IsAdmin,
		ToolFunc: func(ctx context.Context, args TimerArgs) (string, error) {
			args.Action = strings.ToLower(args.Action)
			args.Type = strings.ToLower(args.Type)

			switch args.Action {
			case "list":
				return c.listTimers(rCtx, args), nil
			case "read":
				return c.readTimer(rCtx, args), nil
			case "add":
				return c.addTimer(rCtx, args), nil
			case "remove":
				return c.removeTimer(rCtx, args), nil
			}

			return "Invalid action, please use 'list', 'read', 'add' or 'remove'", nil
		},
	}
}
