package taskdomain

type CronTask interface {
	Task() *Task
	GetSchedule() string
	Activate(CronMng, TaskActivationHandler)
	Close(CronMng)
}
