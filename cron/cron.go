package cron

import (
	"errors"
	"slices"
	"strings"

	"github.com/robfig/cron/v3"
)

type CronCtl struct {
	C *cron.Cron

	Tasks []CronTask
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
		Tasks: make([]CronTask, 0),
	}
}

func (c *CronCtl) GetJob(origin, name string) *CronTask {
	for _, task := range c.Tasks {
		if task.Name == name && task.Origin == origin {
			return &task
		}
	}
	return nil
}

func (c *CronCtl) AddJob(origin, name, description string, args CronArgs, job func()) error {
	if c.GetJob(origin, name) != nil {
		return errors.New("task with name " + name + " already exists")
	}

	id, err := c.C.AddFunc(args.CronString(), job)
	if err != nil {
		return err
	}

	task := CronTask{
		Origin:      origin,
		Name:        name,
		Description: description,
		Schedule:    args,
		ID:          id,
	}
	c.Tasks = append(c.Tasks, task)

	return nil
}

func (c *CronCtl) RemoveJob(origin, name string) error {
	task := c.GetJob(origin, name)
	if task == nil {
		return errors.New("task with name " + name + " not found")
	}

	c.C.Remove(task.ID)
	c.Tasks = slices.DeleteFunc(c.Tasks, func(t CronTask) bool {
		return t.Name == name && t.Origin == origin
	})
	return nil
}

func (c *CronCtl) ListJobs(origin string) []CronTask {
	var tasks []CronTask
	for _, task := range c.Tasks {
		if task.Origin == origin {
			tasks = append(tasks, task)
		}
	}
	slices.SortFunc(tasks, func(a, b CronTask) int {
		return strings.Compare(a.Name, b.Name)
	})
	return tasks
}
