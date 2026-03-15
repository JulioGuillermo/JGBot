package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"sync"
	"testing"
	"time"
)

// mockCronMng is a mock implementation of taskdomain.CronMng for testing
type mockCronMng struct {
	entries map[int]func()
	nextID  int
	mu      sync.Mutex
}

func newMockCronMng() *mockCronMng {
	return &mockCronMng{
		entries: make(map[int]func()),
		nextID:  1,
	}
}

func (m *mockCronMng) AddFunc(spec string, cmd func()) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := m.nextID
	m.nextID++
	m.entries[id] = cmd
	return id, nil
}

func (m *mockCronMng) Remove(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, id)
}

func (m *mockCronMng) HasEntry(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.entries[id]
	return exists
}

func (m *mockCronMng) EntryCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.entries)
}

func (m *mockCronMng) TriggerEntry(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if fn, exists := m.entries[id]; exists {
		fn()
		return true
	}
	return false
}

func TestCronTaskTask(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin:    "telegram",
			Channel:   "general",
			ChatID:    12345,
			ChatName:  "Test Chat",
			SenderID:  67890,
			MessageID: 111,
		},
		TaskInfo: taskdomain.TaskInfo{
			Name:        "test_cron",
			Description: "A test cron task",
			Message:     "Test message",
		},
	}

	cronTask := &CronTask{
		CronTask: task,
		Schedule: taskdomain.CronArgs{
			Minute:  "0",
			Hour:    "12",
			Day:     "1",
			Month:   "1",
			Weekday: "1",
		},
	}

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

func TestCronTaskGetSchedule(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := taskdomain.CronArgs{
		Minute:  "30",
		Hour:    "14",
		Day:     "15",
		Month:   "6",
		Weekday: "3",
	}

	cronTask := &CronTask{
		CronTask: task,
		Schedule: schedule,
	}

	expectedSchedule := schedule.String()
	gotSchedule := cronTask.GetSchedule()

	if gotSchedule != expectedSchedule {
		t.Errorf("GetSchedule() = %v, want %v", gotSchedule, expectedSchedule)
	}
}

func TestCronTaskActivate(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := taskdomain.CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := &CronTask{
		CronTask: task,
		Schedule: schedule,
	}

	cronMng := newMockCronMng()
	activated := false
	var activatedTask *taskdomain.Task
	var activatedSchedule string

	handler := func(t *taskdomain.Task, s string) {
		activated = true
		activatedTask = t
		activatedSchedule = s
	}

	cronTask.Activate(cronMng, handler)

	if cronTask.ID <= 0 {
		t.Error("CronTask should have a valid ID after Activate()")
	}

	if cronMng.EntryCount() != 1 {
		t.Errorf("CronMng should have 1 entry, got %d", cronMng.EntryCount())
	}

	// Manually trigger the function to test handler
	cronMng.TriggerEntry(cronTask.ID)

	if !activated {
		t.Error("Handler should have been called")
	}

	if activatedTask != task {
		t.Error("Handler should receive the same task instance")
	}

	expectedSchedule := schedule.String()
	if activatedSchedule != expectedSchedule {
		t.Errorf("Handler received schedule = %v, want %v", activatedSchedule, expectedSchedule)
	}
}

func TestCronTaskClose(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := taskdomain.CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := &CronTask{
		CronTask: task,
		Schedule: schedule,
	}

	cronMng := newMockCronMng()

	cronTask.Activate(cronMng, nil)

	if cronMng.EntryCount() != 1 {
		t.Errorf("CronMng should have 1 entry before Close(), got %d", cronMng.EntryCount())
	}

	cronTask.Close(cronMng)

	if cronMng.EntryCount() != 0 {
		t.Errorf("CronMng should have 0 entries after Close(), got %d", cronMng.EntryCount())
	}
}

func TestCronTaskCloseWithoutActivate(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := taskdomain.CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := &CronTask{
		CronTask: task,
		Schedule: schedule,
		ID:       5, // Set ID manually
	}

	cronMng := newMockCronMng()

	// Close without activating should not panic
	cronTask.Close(cronMng)

	// Should still try to remove the ID
	if cronMng.EntryCount() != 0 {
		t.Errorf("CronMng should have 0 entries, got %d", cronMng.EntryCount())
	}
}

func TestCronTaskActivateMultipleTimes(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := taskdomain.CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := &CronTask{
		CronTask: task,
		Schedule: schedule,
	}

	cronMng := newMockCronMng()

	// Activate multiple times
	cronTask.Activate(cronMng, nil)
	firstID := cronTask.ID

	cronTask.Activate(cronMng, nil)
	secondID := cronTask.ID

	if firstID == secondID {
		t.Error("Each Activate() should assign a new ID")
	}

	if cronMng.EntryCount() != 2 {
		t.Errorf("CronMng should have 2 entries, got %d", cronMng.EntryCount())
	}
}

func TestCronTaskWithDifferentSchedules(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	tests := []struct {
		name     string
		schedule taskdomain.CronArgs
	}{
		{
			name: "daily at midnight",
			schedule: taskdomain.CronArgs{
				Minute:  "0",
				Hour:    "0",
				Day:     "every 1",
				Month:   "every 1",
				Weekday: "every 1",
			},
		},
		{
			name: "every 5 minutes",
			schedule: taskdomain.CronArgs{
				Minute:  "every 5",
				Hour:    "every 1",
				Day:     "every 1",
				Month:   "every 1",
				Weekday: "every 1",
			},
		},
		{
			name: "weekly on Monday",
			schedule: taskdomain.CronArgs{
				Minute:  "0",
				Hour:    "9",
				Day:     "every 1",
				Month:   "every 1",
				Weekday: "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cronTask := &CronTask{
				CronTask: task,
				Schedule: tt.schedule,
			}

			cronMng := newMockCronMng()
			cronTask.Activate(cronMng, nil)

			if cronTask.ID <= 0 {
				t.Error("CronTask should have a valid ID after Activate()")
			}

			cronTask.Close(cronMng)

			if cronMng.EntryCount() != 0 {
				t.Errorf("CronMng should have 0 entries after Close(), got %d", cronMng.EntryCount())
			}
		})
	}
}

func TestCronTaskConcurrentAccess(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := taskdomain.CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask := &CronTask{
		CronTask: task,
		Schedule: schedule,
	}

	cronMng := newMockCronMng()

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Test concurrent Activate calls
	for range 5 {
		wg.Go(func() {
			cronTask.Activate(cronMng, nil)
		})
	}

	// Wait a bit for activations
	time.Sleep(10 * time.Millisecond)

	// Test concurrent Close calls
	for range 5 {
		wg.Go(func() {
			cronTask.Close(cronMng)
		})
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent access error: %v", err)
		}
	}
}
