package timer

import (
	"errors"
	"slices"
	"strings"
	"time"
)

var Timer *TimerCtl

func InitTimerCtl() {
	Timer = NewTimerCtl()
}

type TimerCtl struct {
	Timers []TimerTask
}

func NewTimerCtl() *TimerCtl {
	return &TimerCtl{
		Timers: make([]TimerTask, 0),
	}
}

func (t *TimerCtl) GetTimer(origin string, name string) *TimerTask {
	for _, timer := range t.Timers {
		if timer.Name == name && timer.Origin == origin {
			return &timer
		}
	}
	return nil
}

func (t *TimerCtl) addTimer(origin string, name string, description string, timeout bool, timerTime TimerTime, job func()) error {
	if t.GetTimer(origin, name) != nil {
		return errors.New("timer with name " + name + " already exists")
	}

	var timer *time.Timer
	if timeout {
		// is time out... so timerTime is the duration
		timer = time.AfterFunc(timerTime.ToDuration(), func() {
			job()
			t.RemoveTimer(origin, name)
		})
	} else {
		// is an alarm... so timerTime is the time
		timer = time.AfterFunc(time.Until(timerTime.ToTime()), func() {
			job()
			t.RemoveTimer(origin, name)
		})
	}

	t.Timers = append(t.Timers, TimerTask{
		Origin:      origin,
		Name:        name,
		Description: description,
		Time:        timerTime,
		Timer:       timer,
	})

	return nil
}

func (t *TimerCtl) AddTimeout(origin string, name string, description string, timerTime TimerTime, job func()) error {
	return t.addTimer(origin, name, description, true, timerTime, job)
}

func (t *TimerCtl) AddAlarm(origin string, name string, description string, timerTime TimerTime, job func()) error {
	return t.addTimer(origin, name, description, false, timerTime, job)
}

func (t *TimerCtl) RemoveTimer(origin, name string) error {
	timer := t.GetTimer(origin, name)
	if timer == nil {
		return errors.New("timer with name " + name + " not found")
	}

	timer.Timer.Stop()
	t.Timers = slices.DeleteFunc(t.Timers, func(timer TimerTask) bool {
		return timer.Name == name && timer.Origin == origin
	})
	return nil
}

func (t *TimerCtl) ListTimers(origin string) []TimerTask {
	var timers []TimerTask
	for _, timer := range t.Timers {
		if timer.Origin == origin {
			timers = append(timers, timer)
		}
	}
	slices.SortFunc(timers, func(a, b TimerTask) int {
		return strings.Compare(a.Name, b.Name)
	})
	return timers
}
