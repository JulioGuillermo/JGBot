package cron

type CronArgs struct {
	Action      string       `json:"action" description:"The action to execute. 'list' to list all the available cron jobs, 'read' to read a cron's description, 'add' to add a new cron job or 'remove' to remove a cron job. (Note: The action is required)"`
	Name        string       `json:"name" description:"The name of the cron job to read, add or remove (Note: Not required for 'list'). This is used to identify the cron job."`
	Description string       `json:"description" description:"The description of the cron job (Note: Required just for 'add'). This is for you to know what the cron job does."`
	Message     string       `json:"message" description:"The message to send (Note: Required just for 'add'). This is the message you will receive when the cron job is executed."`
	Schedule    CronSchedule `json:"schedule" description:"The schedule of the cron job (Note: Required just for 'add'). This is the schedule of the cron job."`
}
