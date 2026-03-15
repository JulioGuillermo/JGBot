package taskdomain

import (
	"testing"
)

func TestTaskInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		taskInfo    TaskInfo
		wantName    string
		wantDesc    string
		wantMessage string
	}{
		{
			name: "empty task info",
			taskInfo: TaskInfo{
				Name:        "",
				Description: "",
				Message:     "",
			},
			wantName:    "",
			wantDesc:    "",
			wantMessage: "",
		},
		{
			name: "full task info",
			taskInfo: TaskInfo{
				Name:        "test_task",
				Description: "A test task",
				Message:     "Hello, World!",
			},
			wantName:    "test_task",
			wantDesc:    "A test task",
			wantMessage: "Hello, World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.taskInfo.Name != tt.wantName {
				t.Errorf("TaskInfo.Name = %v, want %v", tt.taskInfo.Name, tt.wantName)
			}
			if tt.taskInfo.Description != tt.wantDesc {
				t.Errorf("TaskInfo.Description = %v, want %v", tt.taskInfo.Description, tt.wantDesc)
			}
			if tt.taskInfo.Message != tt.wantMessage {
				t.Errorf("TaskInfo.Message = %v, want %v", tt.taskInfo.Message, tt.wantMessage)
			}
		})
	}
}

func TestTaskOriginInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		originInfo TaskOriginInfo
		wantOrigin string
		wantChan   string
		wantChatID uint
	}{
		{
			name: "empty origin info",
			originInfo: TaskOriginInfo{
				Origin:    "",
				Channel:   "",
				ChatID:    0,
				ChatName:  "",
				SenderID:  0,
				MessageID: 0,
			},
			wantOrigin: "",
			wantChan:   "",
			wantChatID: 0,
		},
		{
			name: "full origin info",
			originInfo: TaskOriginInfo{
				Origin:    "telegram",
				Channel:   "general",
				ChatID:    12345,
				ChatName:  "Test Chat",
				SenderID:  67890,
				MessageID: 111,
			},
			wantOrigin: "telegram",
			wantChan:   "general",
			wantChatID: 12345,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.originInfo.Origin != tt.wantOrigin {
				t.Errorf("TaskOriginInfo.Origin = %v, want %v", tt.originInfo.Origin, tt.wantOrigin)
			}
			if tt.originInfo.Channel != tt.wantChan {
				t.Errorf("TaskOriginInfo.Channel = %v, want %v", tt.originInfo.Channel, tt.wantChan)
			}
			if tt.originInfo.ChatID != tt.wantChatID {
				t.Errorf("TaskOriginInfo.ChatID = %v, want %v", tt.originInfo.ChatID, tt.wantChatID)
			}
		})
	}
}

func TestTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		task        Task
		wantOrigin  string
		wantName    string
		wantDesc    string
		wantMessage string
	}{
		{
			name: "empty task",
			task: Task{
				TaskOriginInfo: TaskOriginInfo{},
				TaskInfo:       TaskInfo{},
			},
			wantOrigin:  "",
			wantName:    "",
			wantDesc:    "",
			wantMessage: "",
		},
		{
			name: "full task",
			task: Task{
				TaskOriginInfo: TaskOriginInfo{
					Origin:    "telegram",
					Channel:   "general",
					ChatID:    12345,
					ChatName:  "Test Chat",
					SenderID:  67890,
					MessageID: 111,
				},
				TaskInfo: TaskInfo{
					Name:        "test_task",
					Description: "A test task",
					Message:     "Hello, World!",
				},
			},
			wantOrigin:  "telegram",
			wantName:    "test_task",
			wantDesc:    "A test task",
			wantMessage: "Hello, World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.task.Origin != tt.wantOrigin {
				t.Errorf("Task.Origin = %v, want %v", tt.task.Origin, tt.wantOrigin)
			}
			if tt.task.Name != tt.wantName {
				t.Errorf("Task.Name = %v, want %v", tt.task.Name, tt.wantName)
			}
			if tt.task.Description != tt.wantDesc {
				t.Errorf("Task.Description = %v, want %v", tt.task.Description, tt.wantDesc)
			}
			if tt.task.Message != tt.wantMessage {
				t.Errorf("Task.Message = %v, want %v", tt.task.Message, tt.wantMessage)
			}
		})
	}
}

func TestTaskEquality(t *testing.T) {
	t.Parallel()

	task1 := Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin:    "telegram",
			Channel:   "general",
			ChatID:    12345,
			ChatName:  "Test Chat",
			SenderID:  67890,
			MessageID: 111,
		},
		TaskInfo: TaskInfo{
			Name:        "test_task",
			Description: "A test task",
			Message:     "Hello, World!",
		},
	}

	task2 := Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin:    "telegram",
			Channel:   "general",
			ChatID:    12345,
			ChatName:  "Test Chat",
			SenderID:  67890,
			MessageID: 111,
		},
		TaskInfo: TaskInfo{
			Name:        "test_task",
			Description: "A test task",
			Message:     "Hello, World!",
		},
	}

	task3 := Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin:    "discord",
			Channel:   "general",
			ChatID:    12345,
			ChatName:  "Test Chat",
			SenderID:  67890,
			MessageID: 111,
		},
		TaskInfo: TaskInfo{
			Name:        "test_task",
			Description: "A test task",
			Message:     "Hello, World!",
		},
	}

	if task1.Origin != task2.Origin {
		t.Error("task1 and task2 should have same origin")
	}
	if task1.Name != task2.Name {
		t.Error("task1 and task2 should have same name")
	}
	if task1.Origin == task3.Origin {
		t.Error("task1 and task3 should have different origins")
	}
}
