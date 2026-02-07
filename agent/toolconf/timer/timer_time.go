package timer

import (
	"JGBot/timer"
	"fmt"
)

type TimerTime struct {
	Hour   int `json:"hour" description:"The hour to execute the timer if is an alarm, or the hours duration if is a timeout"`
	Minute int `json:"minute" description:"The minute to execute the timer if is an alarm, or the minutes duration if is a timeout"`
	Second int `json:"second" description:"The second to execute the timer if is an alarm, or the seconds duration if is a timeout"`
}

func (t *TimerTime) ToTime() timer.TimerTime {
	return timer.TimerTime{
		Hour:   t.Hour,
		Minute: t.Minute,
		Second: t.Second,
	}
}

func (t *TimerTime) String() string {
	return fmt.Sprintf("%d:%d:%d", t.Hour, t.Minute, t.Second)
}
