package taskdomain

import (
	"fmt"
	"time"
)

type TimerTime struct {
	Day   int
	Month int
	Year  int

	Hour   int
	Minute int
	Second int
}

func (t TimerTime) String() string {
	if t.Day <= 0 && t.Month <= 0 && t.Year <= 0 {
		return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
	}
	return fmt.Sprintf("%02d:%02d:%02d %02d/%02d/%04d", t.Hour, t.Minute, t.Second, t.Month, t.Day, t.Year)
}

func (t TimerTime) ToDuration() time.Duration {
	return time.Duration(
		t.Year*Year+
			t.Month*Month+
			t.Day*Day+
			t.Hour*Hour+
			t.Minute*Minute+
			t.Second,
	) * time.Second
}

func (t TimerTime) ToTime() time.Time {
	now := time.Now()

	if t.Day <= 0 {
		t.Day = now.Day()
	}
	if t.Month <= 0 {
		t.Month = int(now.Month())
	}
	if t.Year <= 0 {
		t.Year = now.Year()
	}

	return time.Date(t.Year, time.Month(t.Month), t.Day, t.Hour, t.Minute, t.Second, 0, now.Location())
}
