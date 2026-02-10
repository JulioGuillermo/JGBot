package timer

import "time"

type TimerType string

const (
	ALARM   TimerType = "ALARM"
	TIMEOUT TimerType = "TIMEOUT"
)

type TimerTask struct {
	Origin    string
	Channel   string
	ChatID    uint
	ChatName  string
	SenderID  uint
	MessageID uint

	Name        string
	Description string
	Message     string

	Type     TimerType
	Time     TimerTime
	Schedule time.Time
	Timer    *time.Timer `json:"-"`
}

func (tt *TimerTask) setSchedule() {
	if tt.Type == ALARM {
		tt.Schedule = tt.Time.ToTime()
		return
	}

	tt.Schedule = time.Now().Add(tt.Time.ToDuration())
}

func (tt *TimerTask) activate(f func(*TimerTask)) {
	tt.Timer = time.AfterFunc(time.Until(tt.Schedule), func() {
		f(tt)
		tt.Timer = nil
	})
}

func (tt *TimerTask) close() {
	if tt.Timer == nil {
		return
	}
	tt.Timer.Stop()
	tt.Timer = nil
}
