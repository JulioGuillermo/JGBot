package timer

import (
	"JGBot/timer"
	"fmt"
)

type TimerTime struct {
	Day   int `json:"day" description:"This param is only for the alarm type, you can omit it if the alarm is for today. The day to execute the timer if is an alarm, or the days duration if is a timeout"`
	Month int `json:"month" description:"This param is only for the alarm type, you can omit it if the alarm is for this month. The month to execute the timer if is an alarm, or the months duration if is a timeout"`
	Year  int `json:"year" description:"This param is only for the alarm type, you can omit it if the alarm is for this year. The year to execute the timer if is an alarm, or the years duration if is a timeout"`

	Hour   int `json:"hour" description:"The hour to execute the timer if is an alarm, or the hours duration if is a timeout"`
	Minute int `json:"minute" description:"The minute to execute the timer if is an alarm, or the minutes duration if is a timeout"`
	Second int `json:"second" description:"The second to execute the timer if is an alarm, or the seconds duration if is a timeout"`
}

func (t *TimerTime) ToTime() timer.TimerTime {
	return timer.TimerTime{
		Day:    t.Day,
		Month:  t.Month,
		Year:   t.Year,
		Hour:   t.Hour,
		Minute: t.Minute,
		Second: t.Second,
	}
}

func (t *TimerTime) String() string {
	return fmt.Sprintf("%d:%d:%d", t.Hour, t.Minute, t.Second)
}
