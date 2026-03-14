package taskports

import (
	taskapplication "JGBot/task/application"
	taskinfrastructure "JGBot/task/infrastructure"
)

var CronService *taskapplication.CronService

func InitCronService() {
	CronService = &taskapplication.CronService{
		TaskCreator:  &taskinfrastructure.CronTaskCreator{},
		Persister:    &taskinfrastructure.CronTaskPersistence{},
		CronMng:      taskinfrastructure.NewCronMng(),
		OnActivation: nil,
	}
}
