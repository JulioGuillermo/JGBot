package taskdomain

import (
	sessiondomain "JGBot/session/domain"
)

type Task struct {
	TaskOriginInfo
	TaskInfo
}

// ToActivationContext converts a Task to sessiondomain.ActivationContext
func (t *Task) ToActivationContext(schedule string) *sessiondomain.ActivationContext {
	return &sessiondomain.ActivationContext{
		Origin:      t.TaskOriginInfo.Origin,
		Channel:     t.TaskOriginInfo.Channel,
		ChatID:      t.TaskOriginInfo.ChatID,
		ChatName:    t.TaskOriginInfo.ChatName,
		SenderID:    t.TaskOriginInfo.SenderID,
		MessageID:   t.TaskOriginInfo.MessageID,
		Name:        t.TaskInfo.Name,
		Schedule:    schedule,
		Description: t.TaskInfo.Description,
		Message:     t.TaskInfo.Message,
	}
}
