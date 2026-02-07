package timer

import "time"

type TimerTask struct {
	Origin      string
	Name        string
	Description string
	Time        TimerTime
	Timer       *time.Timer
}
