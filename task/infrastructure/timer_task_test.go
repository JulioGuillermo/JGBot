package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"sync"
	"testing"
	"time"
)

func TestTimerTaskTask(t *testing.T) {
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
			Name:        "test_timer",
			Description: "A test timer task",
			Message:     "Test message",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.ALARM,
		Time: taskdomain.TimerTime{
			Hour:   10,
			Minute: 30,
			Second: 0,
		},
	}

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

func TestTimerTaskGetSchedule(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.ALARM,
		Time: taskdomain.TimerTime{
			Hour:   10,
			Minute: 30,
			Second: 0,
		},
	}

	// Before SetSchedule, Schedule should be zero time
	expectedSchedule := time.Time{}.String()
	gotSchedule := timerTask.GetSchedule()

	if gotSchedule != expectedSchedule {
		t.Errorf("GetSchedule() before SetSchedule() = %v, want %v", gotSchedule, expectedSchedule)
	}

	// After SetSchedule
	timerTask.SetSchedule()
	expectedSchedule = timerTask.Schedule.String()
	gotSchedule = timerTask.GetSchedule()

	if gotSchedule != expectedSchedule {
		t.Errorf("GetSchedule() after SetSchedule() = %v, want %v", gotSchedule, expectedSchedule)
	}
}

func TestTimerTaskSetScheduleAlarm(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.ALARM,
		Time: taskdomain.TimerTime{
			Day:    15,
			Month:  6,
			Year:   2024,
			Hour:   14,
			Minute: 30,
			Second: 0,
		},
	}

	timerTask.SetSchedule()

	// For ALARM, Schedule should be the time from TimerTime
	expectedTime := timerTask.Time.ToTime()
	if !timerTask.Schedule.Equal(expectedTime) {
		t.Errorf("SetSchedule() for ALARM: Schedule = %v, want %v", timerTask.Schedule, expectedTime)
	}
}

func TestTimerTaskSetScheduleTimeout(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.TIMEOUT,
		Time: taskdomain.TimerTime{
			Hour:   1,
			Minute: 30,
			Second: 0,
		},
	}

	beforeSet := time.Now()
	timerTask.SetSchedule()

	// For TIMEOUT, Schedule should be now + duration
	expectedTime := beforeSet.Add(timerTask.Time.ToDuration())

	// Allow some tolerance for execution time
	tolerance := 2 * time.Second
	if timerTask.Schedule.Before(expectedTime.Add(-tolerance)) || timerTask.Schedule.After(expectedTime.Add(tolerance)) {
		t.Errorf("SetSchedule() for TIMEOUT: Schedule = %v, want approximately %v", timerTask.Schedule, expectedTime)
	}
}

func TestTimerTaskActivate(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.TIMEOUT,
		Time: taskdomain.TimerTime{
			Second: 1, // 1 second timeout
		},
	}

	timerTask.SetSchedule()

	activated := false
	var activatedTask *taskdomain.Task
	var activatedSchedule string

	handler := func(t *taskdomain.Task, s string) {
		activated = true
		activatedTask = t
		activatedSchedule = s
	}

	timerTask.Activate(handler)

	if timerTask.Timer == nil {
		t.Error("Timer should not be nil after Activate()")
	}

	// Wait for the timer to fire
	time.Sleep(2 * time.Second)

	if !activated {
		t.Error("Handler should have been called after timeout")
	}

	if activatedTask != task {
		t.Error("Handler should receive the same task instance")
	}

	if activatedSchedule == "" {
		t.Error("Handler should receive a non-empty schedule string")
	}
}

func TestTimerTaskClose(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.TIMEOUT,
		Time: taskdomain.TimerTime{
			Second: 10, // 10 second timeout
		},
	}

	timerTask.SetSchedule()
	timerTask.Activate(func(*taskdomain.Task, string) {})

	if timerTask.Timer == nil {
		t.Fatal("Timer should not be nil before Close()")
	}

	timerTask.Close()

	if timerTask.Timer != nil {
		t.Error("Timer should be nil after Close()")
	}
}

func TestTimerTaskCloseWithoutActivate(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.TIMEOUT,
		Time: taskdomain.TimerTime{
			Second: 10,
		},
	}

	// Close without activating should not panic
	timerTask.Close()

	if timerTask.Timer != nil {
		t.Error("Timer should be nil")
	}
}

func TestTimerTaskCloseAfterFire(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.TIMEOUT,
		Time: taskdomain.TimerTime{
			Second: 1,
		},
	}

	timerTask.SetSchedule()
	timerTask.Activate(func(*taskdomain.Task, string) {})

	// Wait for the timer to fire
	time.Sleep(2 * time.Second)

	// Close after fire should not panic
	timerTask.Close()
}

func TestTimerTaskWithDifferentTypes(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	tests := []struct {
		name string
		typ  taskdomain.TimerType
		time taskdomain.TimerTime
	}{
		{
			name: "ALARM type",
			typ:  taskdomain.ALARM,
			time: taskdomain.TimerTime{
				Hour:   10,
				Minute: 30,
				Second: 0,
			},
		},
		{
			name: "TIMEOUT type",
			typ:  taskdomain.TIMEOUT,
			time: taskdomain.TimerTime{
				Second: 5,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			timerTask := &TimerTask{
				TimerTask: task,
				Type:      tt.typ,
				Time:      tt.time,
			}

			timerTask.SetSchedule()

			if timerTask.Schedule.IsZero() {
				t.Error("Schedule should not be zero after SetSchedule()")
			}
		})
	}
}

func TestTimerTaskConcurrentAccess(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.TIMEOUT,
		Time: taskdomain.TimerTime{
			Second: 5,
		},
	}

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	timerTask.SetSchedule()

	// Test concurrent Activate and Close
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			timerTask.Activate(func(*taskdomain.Task, string) {})
		}()
	}

	time.Sleep(10 * time.Millisecond)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			timerTask.Close()
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent access error: %v", err)
		}
	}
}

func TestTimerTaskGetScheduleAfterSetSchedule(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.ALARM,
		Time: taskdomain.TimerTime{
			Day:    15,
			Month:  6,
			Year:   2024,
			Hour:   14,
			Minute: 30,
			Second: 0,
		},
	}

	beforeSchedule := timerTask.GetSchedule()
	timerTask.SetSchedule()
	afterSchedule := timerTask.GetSchedule()

	if beforeSchedule == afterSchedule {
		t.Error("Schedule should change after SetSchedule()")
	}
}

func TestTimerTaskTimerNilAfterClose(t *testing.T) {
	t.Parallel()

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.TIMEOUT,
		Time: taskdomain.TimerTime{
			Second: 10,
		},
	}

	timerTask.SetSchedule()
	timerTask.Activate(func(*taskdomain.Task, string) {})

	if timerTask.Timer == nil {
		t.Fatal("Timer should not be nil before Close()")
	}

	timerTask.Close()

	if timerTask.Timer != nil {
		t.Error("Timer should be nil after Close()")
	}
}
