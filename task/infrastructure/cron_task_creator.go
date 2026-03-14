package taskinfrastructure

import taskdomain "JGBot/task/domain"

type CronTaskCreator struct{}

func (c *CronTaskCreator) CreateCronTask(task *taskdomain.Task, args taskdomain.CronArgs) (taskdomain.CronTask, error) {
	cronTask := &CronTask{
		CronTask: task,
		Schedule: args,
	}

	return cronTask, nil
}
