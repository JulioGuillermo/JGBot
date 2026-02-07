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
	Timers map[string]TimerTask
}

func NewTimerCtl() *TimerCtl {
	return &TimerCtl{
		Timers: make(map[string]TimerTask),
	}
}

func (t *TimerCtl) addTimer(name string, description string, timeout bool, timerTime TimerTime, job func()) error {
	if _, ok := t.Timers[name]; ok {
		return errors.New("timer with name " + name + " already exists")
	}

	var timer *time.Timer
	if timeout {
		// is time out... so timerTime is the duration
		timer = time.AfterFunc(timerTime.ToDuration(), func() {
			job()
			t.RemoveTimer(name)
		})
	} else {
		// is an alarm... so timerTime is the time
		timer = time.AfterFunc(time.Until(timerTime.ToTime()), func() {
			job()
			t.RemoveTimer(name)
		})
	}

	t.Timers[name] = TimerTask{
		Name:        name,
		Description: description,
		Time:        timerTime,
		Timer:       timer,
	}

	return nil
}

func (t *TimerCtl) AddTimeout(name string, description string, timerTime TimerTime, job func()) error {
	return t.addTimer(name, description, true, timerTime, job)
}

func (t *TimerCtl) AddAlarm(name string, description string, timerTime TimerTime, job func()) error {
	return t.addTimer(name, description, false, timerTime, job)
}

func (t *TimerCtl) RemoveTimer(name string) error {
	if _, ok := t.Timers[name]; !ok {
		return errors.New("timer with name " + name + " not found")
	}

	t.Timers[name].Timer.Stop()
	delete(t.Timers, name)
	return nil
}

func (t *TimerCtl) ListTimers() []TimerTask {
	var timers []TimerTask
	for _, timer := range t.Timers {
		timers = append(timers, timer)
	}
	slices.SortFunc(timers, func(a, b TimerTask) int {
		return strings.Compare(a.Name, b.Name)
	})
	return timers
}
