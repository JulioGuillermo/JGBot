package timer

import (
	"JGBot/ctxs"
	"JGBot/log"
	"JGBot/tools"
	"errors"
	"slices"
	"strings"
)

const (
	TimerSave = "config/timers.json"
)

var Timer *TimerCtl

func InitTimerCtl() {
	Timer = NewTimerCtl()
}

type TimerCtl struct {
	OnActivation func(
		origin string,
		channel string,
		chatID uint,
		chatName string,
		senderID uint,
		messageID uint,

		name,
		schedule,
		description,
		message string,
	)
	Timers []TimerTask
}

func NewTimerCtl() *TimerCtl {
	return &TimerCtl{
		Timers: make([]TimerTask, 0),
	}
}

func (t *TimerCtl) Save() {
	tools.WriteJSONFile(TimerSave, t.Timers)
}

func (t *TimerCtl) Load() {
	err := tools.ReadJSONFile(TimerSave, &t.Timers)
	if err != nil {
		log.Error("Fail to load timers", "error", err)
		return
	}
	for _, timer := range t.Timers {
		log.Info("Timer loaded", "name", timer.Name, "type", timer.Type, "timer", timer.Time.String(), "schedule", timer.Schedule.String())
		timer.activate(t.onActivate)
	}
}

func (t *TimerCtl) onActivate(tt *TimerTask) {
	if t.OnActivation != nil {
		t.OnActivation(
			// Chat info
			tt.Origin,
			tt.Channel,
			tt.ChatID,
			tt.ChatName,
			tt.SenderID,
			tt.MessageID,

			// Timer info
			tt.Name,
			tt.Time.String(),
			tt.Description,
			tt.Message,
		)
	}
	t.RemoveTimer(tt.Origin, tt.Name)
}

func (t *TimerCtl) GetTimer(origin string, name string) *TimerTask {
	for _, timer := range t.Timers {
		if timer.Name == name && timer.Origin == origin {
			return &timer
		}
	}
	return nil
}

func (t *TimerCtl) RemoveTimer(origin, name string) error {
	timer := t.GetTimer(origin, name)
	if timer == nil {
		return errors.New("timer with name " + name + " not found")
	}
	defer t.Save()

	timer.close()
	t.Timers = slices.DeleteFunc(t.Timers, func(timer TimerTask) bool {
		return timer.Name == name && timer.Origin == origin
	})
	return nil
}

func (t *TimerCtl) addTimer(
	// ctx
	ctx *ctxs.RespondCtx,

	// task info
	name string,
	description string,
	message string,

	// timer info
	timerType TimerType,
	timerTime TimerTime,
) error {
	if t.GetTimer(ctx.Origin, name) != nil {
		return errors.New("timer with name " + name + " already exists")
	}
	defer t.Save()

	timer := TimerTask{
		Origin:    ctx.Origin,
		Channel:   ctx.Channel,
		ChatID:    ctx.ChatID,
		ChatName:  ctx.ChatName,
		SenderID:  ctx.Message.SenderID,
		MessageID: ctx.Message.ID,

		Name:        name,
		Description: description,
		Message:     message,

		Time: timerTime,
		Type: timerType,
	}
	timer.setSchedule()
	timer.activate(t.onActivate)

	t.Timers = append(t.Timers, timer)
	return nil
}

func (t *TimerCtl) AddTimeout(
	ctx *ctxs.RespondCtx,
	name string,
	description string,
	message string,
	timerTime TimerTime,
) error {
	return t.addTimer(
		ctx,
		name,
		description,
		message,
		TIMEOUT,
		timerTime,
	)
}

func (t *TimerCtl) AddAlarm(
	ctx *ctxs.RespondCtx,
	name string,
	description string,
	message string,
	timerTime TimerTime,
) error {
	return t.addTimer(
		ctx,
		name,
		description,
		message,
		ALARM,
		timerTime,
	)
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
