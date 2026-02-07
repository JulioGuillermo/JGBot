package cron

import "github.com/robfig/cron/v3"

type CronTask struct {
	Name        string
	Description string

	Schedule CronArgs

	ID cron.EntryID
}
