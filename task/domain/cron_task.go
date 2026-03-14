package taskdomain

type CronTask interface {
	Activate(CronMng, TaskActivationHandler)
	Close(CronMng)
}
