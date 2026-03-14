package cron

import (
	"JGBot/agent/tools"
	"JGBot/ctxs"
	taskdomain "JGBot/task/domain"
	taskports "JGBot/task/ports"
	"context"
	"fmt"
	"strings"
)

type CronInitializerConf struct{}

func (c *CronInitializerConf) Name() string {
	return "cron"
}

func (c *CronInitializerConf) listCronJobs(ctx *ctxs.RespondCtx, args CronArgs) string {
	origin := ctx.GetOrigin(args.Session)

	jobs := taskports.CronService.ListJob(origin)
	if len(jobs) == 0 {
		return "Active cron jobs:\nCron jobs list is empty."
	}
	var sb strings.Builder
	sb.WriteString("Active cron jobs:\n")
	for _, job := range jobs {
		fmt.Fprintf(&sb, "- %s: %s\n", job.Name, job.Description)
	}
	return sb.String()
}

func (c *CronInitializerConf) readCronJob(ctx *ctxs.RespondCtx, args CronArgs) string {
	origin := ctx.GetOrigin(args.Session)

	job := taskports.CronService.GetJob(origin, args.Name)
	if job == nil {
		return fmt.Sprintf("Fail to read cron job %s: not found", args.Name)
	}
	task := job.Task()
	return fmt.Sprintf("Name: %s\nDescription: %s\nSchedule: (%s)", task.Name, task.Description, job.GetSchedule())
}

func (c *CronInitializerConf) addCronJob(rCtx *ctxs.RespondCtx, args CronArgs) string {
	ctx, err := rCtx.GetSessionCtx(args.Session)
	if err != nil {
		return err.Error()
	}
	if ctx.SessionCtl == nil {
		return "Error: Session controller not initialized."
	}

	err = taskports.CronService.AddJob(
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
		taskdomain.CronArgs(args.Schedule),
	)
	if err != nil {
		return fmt.Sprintf("Fail to add cron job %s: %s", args.Name, err.Error())
	}
	return fmt.Sprintf("Cron job %s added for session %s", args.Name, ctx.Origin)
}

func (c *CronInitializerConf) removeCronJob(ctx *ctxs.RespondCtx, args CronArgs) string {
	origin := ctx.GetOrigin(args.Session)

	err := taskports.CronService.RemoveJob(origin, args.Name)
	if err != nil {
		return fmt.Sprintf("Fail to remove cron job %s: %s", args.Name, err.Error())
	}
	return fmt.Sprintf("Cron job %s removed", args.Name)
}

func (c *CronInitializerConf) ToolInitializer(rCtx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[CronArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Allows you to list, read, or execute cron jobs.",
		IsAdmin:         rCtx.IsAdmin,
		ToolFunc: func(ctx context.Context, args CronArgs) (string, error) {
			switch args.Action {
			case "list":
				return c.listCronJobs(rCtx, args), nil
			case "read":
				return c.readCronJob(rCtx, args), nil
			case "add":
				return c.addCronJob(rCtx, args), nil
			case "remove":
				return c.removeCronJob(rCtx, args), nil
			}

			return "Invalid action, please use 'list', 'read', 'add' or 'remove'", nil
		},
	}
}
