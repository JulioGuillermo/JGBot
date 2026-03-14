package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"JGBot/tools"
)

type TimerTaskPersistence struct {
}

func (p *TimerTaskPersistence) LoadTimers() ([]taskdomain.TimerTask, error) {
	var Timers []*TimerTask
	err := tools.ReadJSONFile(taskdomain.TIMER_FILE, &Timers)
	if err != nil {
		return nil, err
	}

	var timerTasks []taskdomain.TimerTask = make([]taskdomain.TimerTask, len(Timers))
	for i, timer := range Timers {
		timerTasks[i] = timer
	}

	return timerTasks, nil
}

func (p *TimerTaskPersistence) SaveTimers(timers []taskdomain.TimerTask) error {
	return tools.WriteJSONFile(taskdomain.TIMER_FILE, timers)
}
