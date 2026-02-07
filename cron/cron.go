package cron

import (
	"errors"
	"slices"
	"strings"

	"github.com/robfig/cron/v3"
)

type CronCtl struct {
	C *cron.Cron

	Tasks map[string]CronTask
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
		Tasks: make(map[string]CronTask),
	}
}

func (c *CronCtl) AddJob(name, description string, args CronArgs, job func()) error {
	if _, ok := c.Tasks[name]; ok {
		return errors.New("task with name " + name + " already exists")
	}

	id, err := c.C.AddFunc(args.CronString(), job)
	if err != nil {
		return err
	}

	task := CronTask{
		Name:        name,
		Description: description,
		Schedule:    args,
		ID:          id,
	}
	c.Tasks[name] = task

	return nil
}

func (c *CronCtl) RemoveJob(name string) error {
	task, ok := c.Tasks[name]
	if !ok {
		return errors.New("task with name " + name + " not found")
	}

	c.C.Remove(task.ID)
	delete(c.Tasks, name)
	return nil
}

func (c *CronCtl) ListJobs() []CronTask {
	var tasks []CronTask
	for _, task := range c.Tasks {
		tasks = append(tasks, task)
	}
	slices.SortFunc(tasks, func(a, b CronTask) int {
		return strings.Compare(a.Name, b.Name)
	})
	return tasks
}
