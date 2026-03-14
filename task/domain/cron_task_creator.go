package taskdomain

type CronTaskCreator interface {
	CreateCronTask(
		task *Task,
		schedule CronArgs,
	) (CronTask, error)
}
