package taskdomain

type CronTaskPersister interface {
	LoadCrons() ([]CronTask, error)
	SaveCrons(crons []CronTask) error
}
