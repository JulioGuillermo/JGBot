package cron

import (
	"JGBot/ctxs"
	"JGBot/log"
	"JGBot/tools"
	"errors"
	"slices"
	"strings"

	"github.com/robfig/cron/v3"
)

const (
	CRON_FILE = "config/cron.json"
)

type CronCtl struct {
	C *cron.Cron

	OnActivation func(
		origin string,
		channel string,
		chatID uint,
		chatName string,
		senderID uint,
		messageID uint,

		name,
		schedule,
		description,
		message string,
	)

	Tasks []*CronTask
}

var Cron *CronCtl

func InitCronCtl() {
	Cron = NewCronCtl()
}

func NewCronCtl() *CronCtl {
	c := cron.New()
	go c.Start()

	return &CronCtl{
		C:     c,
		Tasks: make([]*CronTask, 0),
	}
}

func (c *CronCtl) Save() {
	tools.WriteJSONFile(CRON_FILE, c.Tasks)
}

func (c *CronCtl) Load() {
	err := tools.ReadJSONFile(CRON_FILE, &c.Tasks)
	if err != nil {
		log.Error("Fail to load cron", "error", err)
		return
	}
	for _, task := range c.Tasks {
		log.Info("Cron loaded", "name", task.Name, "schedule", task.Schedule.String())
		task.activate(c.C, c.onActivate)
	}
}

func (c *CronCtl) onActivate(task *CronTask) {
	if c.OnActivation != nil {
		c.OnActivation(
			// Chat info
			task.Origin,
			task.Channel,
			task.ChatID,
			task.ChatName,
			task.SenderID,
			task.MessageID,

			// Timer info
			task.Name,
			task.Schedule.String(),
			task.Description,
			task.Message,
		)
	}
}

func (c *CronCtl) GetJob(origin, name string) *CronTask {
	for _, task := range c.Tasks {
		if task.Name == name && task.Origin == origin {
			return task
		}
	}
	return nil
}

func (c *CronCtl) RemoveJob(origin, name string) error {
	task := c.GetJob(origin, name)
	if task == nil {
		return errors.New("task with name " + name + " not found")
	}
	defer c.Save()

	task.close(c.C)
	c.Tasks = slices.DeleteFunc(c.Tasks, func(t *CronTask) bool {
		return t.Name == name && t.Origin == origin
	})
	return nil
}

func (c *CronCtl) AddJob(
	// ctx
	ctx *ctxs.RespondCtx,

	// task info
	name string,
	description string,
	message string,

	// cron info
	args CronArgs,
) error {
	if c.GetJob(ctx.Origin, name) != nil {
		return errors.New("task with name " + name + " already exists")
	}
	defer c.Save()

	task := &CronTask{
		Origin:    ctx.Origin,
		Channel:   ctx.Channel,
		ChatID:    ctx.ChatID,
		ChatName:  ctx.ChatName,
		SenderID:  ctx.Message.SenderID,
		MessageID: ctx.Message.ID,

		Name:        name,
		Description: description,
		Message:     message,

		Schedule: args,
	}
	task.activate(c.C, c.onActivate)
	c.Tasks = append(c.Tasks, task)

	return nil
}

func (c *CronCtl) ListJobs(origin string) []*CronTask {
	var tasks []*CronTask
	for _, task := range c.Tasks {
		if task.Origin == origin {
			tasks = append(tasks, task)
		}
	}
	slices.SortFunc(tasks, func(a, b *CronTask) int {
		return strings.Compare(a.Name, b.Name)
	})
	return tasks
}
