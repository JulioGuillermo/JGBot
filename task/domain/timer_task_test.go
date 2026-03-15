package taskdomain

import (
	"testing"
)

// MockTimerTask is a concrete implementation of TimerTask for testing
type MockTimerTask struct {
	task      *Task
	timerType TimerType
	timerTime TimerTime
	schedule  string
	closed    bool
	active    bool
}

func NewMockTimerTask(task *Task, timerType TimerType, timerTime TimerTime) *MockTimerTask {
	return &MockTimerTask{
		task:      task,
		timerType: timerType,
		timerTime: timerTime,
		closed:    false,
		active:    false,
	}
}

func (m *MockTimerTask) Task() *Task {
	return m.task
}

func (m *MockTimerTask) GetSchedule() string {
	return m.schedule
}

func (m *MockTimerTask) SetSchedule() {
	m.schedule = "set"
}

func (m *MockTimerTask) Activate(handler TaskActivationHandler) {
	m.active = true
}

func (m *MockTimerTask) Close() {
	m.closed = true
}

func (m *MockTimerTask) IsClosed() bool {
	return m.closed
}

func (m *MockTimerTask) IsActive() bool {
	return m.active
}

func TestTimerTaskInterface(t *testing.T) {
	t.Parallel()

	// This test verifies that MockTimerTask implements TimerTask interface
	var _ TimerTask = (*MockTimerTask)(nil)
}

func TestTimerTaskActivate(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := NewMockTimerTask(task, ALARM, TimerTime{})

	handler := func(t *Task, s string) {
		// Handler called
	}

	timerTask.Activate(handler)

	if !timerTask.IsActive() {
		t.Error("TimerTask should be active after Activate()")
	}
}

func TestTimerTaskClose(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := NewMockTimerTask(task, ALARM, TimerTime{})

	timerTask.Close()

	if !timerTask.IsClosed() {
		t.Error("TimerTask should be closed after Close()")
	}
}

func TestTimerTaskGetSchedule(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := NewMockTimerTask(task, ALARM, TimerTime{})
	timerTask.schedule = "10:30:00"

	expectedSchedule := "10:30:00"
	gotSchedule := timerTask.GetSchedule()

	if gotSchedule != expectedSchedule {
		t.Errorf("GetSchedule() = %v, want %v", gotSchedule, expectedSchedule)
	}
}

func TestTimerTaskTask(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin:    "telegram",
			Channel:   "general",
			ChatID:    12345,
			ChatName:  "Test Chat",
			SenderID:  67890,
			MessageID: 111,
		},
		TaskInfo: TaskInfo{
			Name:        "test_timer",
			Description: "A test timer task",
			Message:     "Test message",
		},
	}

	timerTask := NewMockTimerTask(task, TIMEOUT, TimerTime{})

	gotTask := timerTask.Task()

	if gotTask != task {
		t.Error("Task() should return the same task instance")
	}

	if gotTask.Origin != "telegram" {
		t.Errorf("Task().Origin = %v, want telegram", gotTask.Origin)
	}

	if gotTask.Name != "test_timer" {
		t.Errorf("Task().Name = %v, want test_timer", gotTask.Name)
	}
}

func TestTimerTaskSetSchedule(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := NewMockTimerTask(task, ALARM, TimerTime{})

	if timerTask.GetSchedule() != "" {
		t.Error("Schedule should be empty before SetSchedule()")
	}

	timerTask.SetSchedule()

	if timerTask.GetSchedule() != "set" {
		t.Errorf("GetSchedule() after SetSchedule() = %v, want 'set'", timerTask.GetSchedule())
	}
}

func TestTimerTaskCloseWithoutActivate(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := NewMockTimerTask(task, ALARM, TimerTime{})

	// Close without activating should not panic
	timerTask.Close()

	if !timerTask.IsClosed() {
		t.Error("TimerTask should be marked as closed")
	}
}

func TestTimerTaskWithDifferentTypes(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_timer",
		},
	}

	// Test with ALARM type
	alarmTask := NewMockTimerTask(task, ALARM, TimerTime{})
	alarmTask.Activate(nil)
	if !alarmTask.IsActive() {
		t.Error("ALARM timer should be active after Activate()")
	}

	// Test with TIMEOUT type
	timeoutTask := NewMockTimerTask(task, TIMEOUT, TimerTime{})
	timeoutTask.Activate(nil)
	if !timeoutTask.IsActive() {
		t.Error("TIMEOUT timer should be active after Activate()")
	}
}
