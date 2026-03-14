package taskdomain

type TimerTask interface {
	Task() *Task
	SetSchedule()
	Activate(TaskActivationHandler)
	Close()
}
