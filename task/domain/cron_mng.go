package taskdomain

type CronMng interface {
	AddFunc(string, func()) (int, error)
	Remove(int)
}
