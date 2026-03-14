package taskdomain

type TimerPersistence interface {
	LoadTimers() ([]TimerTask, error)
	SaveTimers(timers []TimerTask) error
}
