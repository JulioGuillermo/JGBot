package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"time"
)

type TimerTask struct {
	TimerTask *taskdomain.Task

	Type     taskdomain.TimerType
	Time     taskdomain.TimerTime
	Schedule time.Time
	Timer    *time.Timer `json:"-"`
}

func (t *TimerTask) Task() *taskdomain.Task {
	return t.TimerTask
}

func (t *TimerTask) SetSchedule() {
	switch t.Type {
	case taskdomain.ALARM:
		t.Schedule = t.Time.ToTime()
	case taskdomain.TIMEOUT:
		t.Schedule = time.Now().Add(t.Time.ToDuration())
	}
}

func (t *TimerTask) Activate(fun taskdomain.TaskActivationHandler) {
	t.Timer = time.AfterFunc(time.Until(t.Schedule), func() {
		fun(t.TimerTask, t.Schedule.String())
		t.Timer = nil
	})
}

func (t *TimerTask) Close() {
	if t.Timer == nil {
		return
	}

	t.Timer.Stop()
	t.Timer = nil
}
