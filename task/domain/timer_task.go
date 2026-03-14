package taskdomain

type TimerTask interface {
	Task() *Task
	GetSchedule() string
	SetSchedule()
	Activate(TaskActivationHandler)
	Close()
}
