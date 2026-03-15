package taskdomain

import (
	"testing"
)

// MockCronMng is a mock implementation of CronMng for testing
type MockCronMng struct {
	entries map[int]func()
	nextID  int
}

func NewMockCronMng() *MockCronMng {
	return &MockCronMng{
		entries: make(map[int]func()),
		nextID:  1,
	}
}

func (m *MockCronMng) AddFunc(spec string, cmd func()) (int, error) {
	id := m.nextID
	m.nextID++
	m.entries[id] = cmd
	return id, nil
}

func (m *MockCronMng) Remove(id int) {
	delete(m.entries, id)
}

func (m *MockCronMng) HasEntry(id int) bool {
	_, exists := m.entries[id]
	return exists
}

func (m *MockCronMng) EntryCount() int {
	return len(m.entries)
}

// MockCronTask is a concrete implementation of CronTask for testing
type MockCronTask struct {
	task     *Task
	schedule CronArgs
	id       int
	closed   bool
	active   bool
}

func NewMockCronTask(task *Task, schedule CronArgs) *MockCronTask {
	return &MockCronTask{
		task:     task,
		schedule: schedule,
		closed:   false,
		active:   false,
	}
}

func (m *MockCronTask) Task() *Task {
	return m.task
}

func (m *MockCronTask) GetSchedule() string {
	return m.schedule.String()
}

func (m *MockCronTask) Activate(cron CronMng, handler TaskActivationHandler) {
	id, _ := cron.AddFunc(m.schedule.CronString(), func() {
		handler(m.task, m.schedule.String())
	})
	m.id = id
	m.active = true
}

func (m *MockCronTask) Close(cron CronMng) {
	cron.Remove(m.id)
	m.closed = true
}

func (m *MockCronTask) GetID() int {
	return m.id
}

func (m *MockCronTask) IsClosed() bool {
	return m.closed
}

func (m *MockCronTask) IsActive() bool {
	return m.active
}

func TestCronTaskInterface(t *testing.T) {
	t.Parallel()

	// This test verifies that MockCronTask implements CronTask interface
	var _ CronTask = (*MockCronTask)(nil)
}

func TestCronTaskActivate(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := NewMockCronTask(task, schedule)
	cronMng := NewMockCronMng()
	activated := false

	handler := func(t *Task, s string) {
		activated = true
	}

	cronTask.Activate(cronMng, handler)

	if !cronTask.IsActive() {
		t.Error("CronTask should be active after Activate()")
	}

	if cronTask.GetID() <= 0 {
		t.Error("CronTask should have a valid ID after Activate()")
	}

	if cronMng.EntryCount() != 1 {
		t.Errorf("CronMng should have 1 entry, got %d", cronMng.EntryCount())
	}

	// Manually trigger the function to test handler
	if entry, exists := cronMng.entries[cronTask.GetID()]; exists {
		entry()
	}

	if !activated {
		t.Error("Handler should have been called")
	}
}

func TestCronTaskClose(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := NewMockCronTask(task, schedule)
	cronMng := NewMockCronMng()

	cronTask.Activate(cronMng, nil)

	if cronMng.EntryCount() != 1 {
		t.Errorf("CronMng should have 1 entry before Close(), got %d", cronMng.EntryCount())
	}

	cronTask.Close(cronMng)

	if !cronTask.IsClosed() {
		t.Error("CronTask should be closed after Close()")
	}

	if cronMng.EntryCount() != 0 {
		t.Errorf("CronMng should have 0 entries after Close(), got %d", cronMng.EntryCount())
	}
}

func TestCronTaskGetSchedule(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := CronArgs{
		Minute:  "30",
		Hour:    "14",
		Day:     "15",
		Month:   "6",
		Weekday: "3",
	}

	cronTask := NewMockCronTask(task, schedule)

	expectedSchedule := schedule.String()
	gotSchedule := cronTask.GetSchedule()

	if gotSchedule != expectedSchedule {
		t.Errorf("GetSchedule() = %v, want %v", gotSchedule, expectedSchedule)
	}
}

func TestCronTaskTask(t *testing.T) {
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
			Name:        "test_cron",
			Description: "A test cron task",
			Message:     "Test message",
		},
	}

	schedule := CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := NewMockCronTask(task, schedule)

	gotTask := cronTask.Task()

	if gotTask != task {
		t.Error("Task() should return the same task instance")
	}

	if gotTask.Origin != "telegram" {
		t.Errorf("Task().Origin = %v, want telegram", gotTask.Origin)
	}

	if gotTask.Name != "test_cron" {
		t.Errorf("Task().Name = %v, want test_cron", gotTask.Name)
	}
}

func TestCronTaskMultipleActivate(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := NewMockCronTask(task, schedule)
	cronMng := NewMockCronMng()

	// Activate multiple times
	cronTask.Activate(cronMng, nil)
	firstID := cronTask.GetID()

	cronTask.Activate(cronMng, nil)
	secondID := cronTask.GetID()

	if firstID == secondID {
		t.Error("Each Activate() should assign a new ID")
	}

	if cronMng.EntryCount() != 2 {
		t.Errorf("CronMng should have 2 entries, got %d", cronMng.EntryCount())
	}
}

func TestCronTaskCloseWithoutActivate(t *testing.T) {
	t.Parallel()

	task := &Task{
		TaskOriginInfo: TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := NewMockCronTask(task, schedule)
	cronMng := NewMockCronMng()

	// Close without activating should not panic
	cronTask.Close(cronMng)

	if !cronTask.IsClosed() {
		t.Error("CronTask should be marked as closed")
	}
}
