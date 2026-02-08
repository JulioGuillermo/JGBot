package cron

import (
	"JGBot/agent/tools"
	"JGBot/cron"
	"JGBot/ctxs"
	"context"
	"fmt"
	"strings"
)

type CronInitializerConf struct{}

func (c *CronInitializerConf) Name() string {
	return "cron"
}

func (c *CronInitializerConf) listCronJobs(ctx *ctxs.RespondCtx) string {
	jobs := cron.Cron.ListJobs(ctx.Origin)
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

func (c *CronInitializerConf) readCronJob(ctx *ctxs.RespondCtx, name string) string {
	job := cron.Cron.GetJob(ctx.Origin, name)
	if job == nil {
		return fmt.Sprintf("Fail to read cron job %s: not found", name)
	}
	return fmt.Sprintf("Name: %s\nDescription: %s\nSchedule: (%s)", job.Name, job.Description, job.Schedule.String())
}

func (c *CronInitializerConf) addCronJob(rCtx *ctxs.RespondCtx, args CronArgs) string {
	err := cron.Cron.AddJob(rCtx, args.Name, args.Description, args.Message, args.Schedule.ToCron())
	if err != nil {
		return fmt.Sprintf("Fail to add cron job %s: %s", args.Name, err.Error())
	}
	return fmt.Sprintf("Cron job %s added", args.Name)
}

func (c *CronInitializerConf) removeCronJob(ctx *ctxs.RespondCtx, name string) string {
	err := cron.Cron.RemoveJob(ctx.Origin, name)
	if err != nil {
		return fmt.Sprintf("Fail to remove cron job %s: %s", name, err.Error())
	}
	return fmt.Sprintf("Cron job %s removed", name)
}

func (c *CronInitializerConf) ToolInitializer(rCtx *ctxs.RespondCtx) tools.Tool {
	return &tools.ToolAutoArgs[CronArgs]{
		ToolName:        c.Name(),
		ToolDescription: "Allows you to list, read, or execute cron jobs.",
		ToolFunc: func(ctx context.Context, args CronArgs) (string, error) {
			switch args.Action {
			case "list":
				return c.listCronJobs(rCtx), nil
			case "read":
				return c.readCronJob(rCtx, args.Name), nil
			case "add":
				return c.addCronJob(rCtx, args), nil
			case "remove":
				return c.removeCronJob(rCtx, args.Name), nil
			}

			return "Invalid action, please use 'list', 'read', 'add' or 'remove'", nil
		},
	}
}
