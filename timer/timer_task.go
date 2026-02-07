package timer

import "time"

type TimerTask struct {
	Name        string
	Description string
	Time        TimerTime
	Timer       *time.Timer
}
