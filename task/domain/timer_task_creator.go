package taskdomain

type TimerTaskCreator interface {
	CreateTimerTask(
		task *Task,
		timerType TimerType,
		timerTime TimerTime,
	) (TimerTask, error)
}
