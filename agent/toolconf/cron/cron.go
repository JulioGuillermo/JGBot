package cron

import (
	"JGBot/agent/tools"
	"JGBot/cron"
	"JGBot/ctxs"
	"context"
	"fmt"
	"strings"
)

type CronInitializerConf struct {
	OnExecute func(ctx *ctxs.RespondCtx, args CronArgs)
}

func (c *CronInitializerConf) Name() string {
	return "cron"
}

func (c *CronInitializerConf) listCronJobs() string {
	jobs := cron.Cron.ListJobs()
	if len(jobs) == 0 {
		return "No cron jobs are active."
	}
	var sb strings.Builder
	sb.WriteString("Active cron jobs:\n")
	for _, job := range jobs {
		fmt.Fprintf(&sb, "- %s: %s\n", job.Name, job.Description)
	}
	return sb.String()
}

func (c *CronInitializerConf) readCronJob(name string) string {
	job, ok := cron.Cron.Tasks[name]
	if !ok {
		return fmt.Sprintf("Fail to read cron job %s: not found", name)
	}
	return fmt.Sprintf("Name: %s\nDescription: %s\nSchedule: (%s)", job.Name, job.Description, job.Schedule.String())
}

func (c *CronInitializerConf) addCronJob(rCtx *ctxs.RespondCtx, args CronArgs) string {
	err := cron.Cron.AddJob(args.Name, args.Description, args.Schedule.ToCron(), func() {
		c.OnExecute(rCtx, args)
	})
	if err != nil {
		return fmt.Sprintf("Fail to add cron job %s: %s", args.Name, err.Error())
	}
	return fmt.Sprintf("Cron job %s added", args.Name)
}

func (c *CronInitializerConf) removeCronJob(name string) string {
	err := cron.Cron.RemoveJob(name)
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
				return c.listCronJobs(), nil
			case "read":
				return c.readCronJob(args.Name), nil
			case "add":
				return c.addCronJob(rCtx, args), nil
			case "remove":
				return c.removeCronJob(args.Name), nil
			}

			return "Invalid action, please use 'list', 'read', 'add' or 'remove'", nil
		},
	}
}
