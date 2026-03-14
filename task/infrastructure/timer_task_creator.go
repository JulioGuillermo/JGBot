package taskinfrastructure

import taskdomain "JGBot/task/domain"

type TimerTaskCreator struct {
}

func (t *TimerTaskCreator) CreateTimerTask(
	task *taskdomain.Task,
	timerType taskdomain.TimerType,
	timerTime taskdomain.TimerTime,
) (taskdomain.TimerTask, error) {
	return &TimerTask{
		TimerTask: task,
		Type:      timerType,
		Time:      timerTime,
	}, nil
}
