package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
)

type CronTask struct {
	CronTask *taskdomain.Task

	Schedule taskdomain.CronArgs
	ID       int `json:"-"`
}

func (t *CronTask) Task() *taskdomain.Task {
	return t.CronTask
}

func (t *CronTask) GetSchedule() string {
	return t.Schedule.String()
}

func (t *CronTask) Activate(cron taskdomain.CronMng, handler taskdomain.TaskActivationHandler) {
	t.ID, _ = cron.AddFunc(t.Schedule.CronString(), func() {
		handler(t.CronTask, t.Schedule.String())
	})
}

func (t *CronTask) Close(cron taskdomain.CronMng) {
	cron.Remove(t.ID)
}
