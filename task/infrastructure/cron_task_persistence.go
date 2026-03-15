package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"JGBot/tools"
)

type CronTaskPersistence struct {
	FilePath string
}

func (p *CronTaskPersistence) LoadCrons() ([]taskdomain.CronTask, error) {
	filePath := p.FilePath
	if filePath == "" {
		filePath = taskdomain.CRON_FILE
	}

	var crons []*CronTask
	err := tools.ReadJSONFile(filePath, &crons)
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
	filePath := p.FilePath
	if filePath == "" {
		filePath = taskdomain.CRON_FILE
	}
	return tools.WriteJSONFile(filePath, crons)
}
