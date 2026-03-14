package taskdomain

const (
	TIMER_FILE = "config/timers.json"
	CRON_FILE  = "config/cron.json"
)

const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Month  = 30 * Day
	Year   = 12 * Month
)
