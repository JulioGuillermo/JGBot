package timer

import (
	"fmt"
	"regexp"
	"strconv"
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

func ParseTimerTime(time string) (TimerTime, error) {
	hour, minute, second := parseTimes(time)

	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		hourInt = 0
	}

	minuteInt, err := strconv.Atoi(minute)
	if err != nil {
		minuteInt = 0
	}

	secondInt, err := strconv.Atoi(second)
	if err != nil {
		secondInt = 0
	}

	if secondInt > 59 {
		minuteInt += secondInt / 60
		secondInt = secondInt % 60
	}

	if minuteInt > 59 {
		hourInt += minuteInt / 60
		minuteInt = minuteInt % 60
	}

	return TimerTime{
		Hour:   hourInt,
		Minute: minuteInt,
		Second: secondInt,
	}, nil
}

func parseTimes(time string) (hour, minute, second string) {
	reFull := regexp.MustCompile(`^(\d+):(\d+):(\d+)$`)
	reHourMin := regexp.MustCompile(`^(\d+):(\d+)$`)
	reMin := regexp.MustCompile(`^(\d+)$`)

	matches := reFull.FindStringSubmatch(time)
	if len(matches) == 4 {
		hour = matches[1]
		minute = matches[2]
		second = matches[3]
		return
	}

	matches = reHourMin.FindStringSubmatch(time)
	if len(matches) == 3 {
		hour = matches[1]
		minute = matches[2]
		second = "0"
		return
	}

	matches = reMin.FindStringSubmatch(time)
	if len(matches) == 2 {
		hour = "0"
		minute = matches[1]
		second = "0"
		return
	}

	return "0", "0", "0"
}

func (t TimerTime) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second)
}

func (t TimerTime) ToDuration() time.Duration {
	return time.Duration(t.Hour*60*60+t.Minute*60+t.Second) * time.Second
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
