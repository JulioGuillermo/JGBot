package timer

type TimerArgs struct {
	Action      string    `json:"action" description:"The action to execute. 'list' to list all the available timers, 'read' to read a timer's description, 'add' to add a new timer or 'remove' to remove a timer. (Note: The action is required)"`
	Name        string    `json:"name" description:"The name of the timer to read, add or remove (Note: Not required for 'list'). This is used to identify the timer."`
	Type        string    `json:"type" description:"The type of the timer (Note: Required just for 'add'). This is the type of the timer (timeout or alarm). If is timeout, the timer will be executed after the time has passed. If is alarm, the timer will be executed at the time."`
	Description string    `json:"description" description:"The description of the timer (Note: Required just for 'add'). This is for you to know what the timer does."`
	Message     string    `json:"message" description:"The message to send (Note: Required just for 'add'). This is the message you will receive when the timer is executed."`
	TimerTime   TimerTime `json:"time" description:"The timer time of the timer (Note: Required just for 'add'). If is timeout, the timer will be executed after the time has passed. If is alarm, the timer will be executed at the time."`
	Session     string    `json:"session" description:"The origin of the session to act on (Note: Admin only). Omit if you want to act on the current session." admin:"true"`
}
