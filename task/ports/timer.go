package taskports

import (
	taskapplication "JGBot/task/application"
	taskinfrastructure "JGBot/task/infrastructure"
)

var TimerService *taskapplication.TimerService

func InitTimerService() {
	TimerService = taskapplication.NewTimerService(
		&taskinfrastructure.TimerTaskCreator{},
		&taskinfrastructure.TimerTaskPersistence{},
		nil,
	)
}
