package cron

import "github.com/robfig/cron/v3"

type CronTask struct {
	Origin    string
	Channel   string
	ChatID    uint
	ChatName  string
	SenderID  uint
	MessageID uint

	Name        string
	Description string
	Message     string

	Schedule CronArgs

	ID cron.EntryID `json:"-"`
}

func (c *CronTask) activate(cron *cron.Cron, onActivate func(*CronTask)) (err error) {
	c.ID, err = cron.AddFunc(c.Schedule.CronString(), func() {
		onActivate(c)
	})
	return err
}

func (c *CronTask) close(cron *cron.Cron) {
	cron.Remove(c.ID)
}
