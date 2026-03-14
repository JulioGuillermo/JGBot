package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"JGBot/tools"
)

type CronTaskPersistence struct{}

func (p *CronTaskPersistence) LoadCrons() ([]taskdomain.CronTask, error) {
	var crons []*CronTask
	err := tools.ReadJSONFile(taskdomain.CRON_FILE, &crons)
	if err != nil {
		return nil, err
	}

	var cronTasks []taskdomain.CronTask = make([]taskdomain.CronTask, len(crons))
	for i, crons := range crons {
		cronTasks[i] = crons
	}

	return cronTasks, nil
}

func (p *CronTaskPersistence) SaveCrons(crons []taskdomain.CronTask) error {
	return tools.WriteJSONFile(taskdomain.CRON_FILE, crons)
}
