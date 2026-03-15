package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"JGBot/tools"
)

type TimerTaskPersistence struct {
	FilePath string
}

func (p *TimerTaskPersistence) LoadTimers() ([]taskdomain.TimerTask, error) {
	filePath := p.FilePath
	if filePath == "" {
		filePath = taskdomain.TIMER_FILE
	}

	var Timers []*TimerTask
	err := tools.ReadJSONFile(filePath, &Timers)
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
	filePath := p.FilePath
	if filePath == "" {
		filePath = taskdomain.TIMER_FILE
	}

	return tools.WriteJSONFile(filePath, timers)
}
