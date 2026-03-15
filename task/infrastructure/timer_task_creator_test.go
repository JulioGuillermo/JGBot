package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"testing"
)

func TestTimerTaskCreatorCreateTimerTask(t *testing.T) {
	t.Parallel()

	creator := &TimerTaskCreator{}

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

	timerTime := taskdomain.TimerTime{
		Day:    15,
		Month:  6,
		Year:   2024,
		Hour:   14,
		Minute: 30,
		Second: 0,
	}

	timerTask, err := creator.CreateTimerTask(task, taskdomain.ALARM, timerTime)
	if err != nil {
		t.Fatalf("CreateTimerTask() error = %v", err)
	}

	if timerTask == nil {
		t.Fatal("CreateTimerTask() returned nil")
	}

	// Verify the returned task
	returnedTask := timerTask.Task()
	if returnedTask != task {
		t.Error("Returned task should be the same instance as input")
	}
}

func TestTimerTaskCreatorCreateTimerTaskWithDifferentTypes(t *testing.T) {
	t.Parallel()

	creator := &TimerTaskCreator{}

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTime := taskdomain.TimerTime{
		Hour:   10,
		Minute: 30,
		Second: 0,
	}

	tests := []struct {
		name string
		typ  taskdomain.TimerType
	}{
		{
			name: "ALARM type",
			typ:  taskdomain.ALARM,
		},
		{
			name: "TIMEOUT type",
			typ:  taskdomain.TIMEOUT,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			timerTask, err := creator.CreateTimerTask(task, tt.typ, timerTime)
			if err != nil {
				t.Errorf("CreateTimerTask() error = %v", err)
				return
			}

			if timerTask == nil {
				t.Error("CreateTimerTask() returned nil")
				return
			}

			returnedTask := timerTask.Task()
			if returnedTask != task {
				t.Error("Returned task should be the same instance as input")
			}
		})
	}
}

func TestTimerTaskCreatorCreateTimerTaskWithEmptyTask(t *testing.T) {
	t.Parallel()

	creator := &TimerTaskCreator{}

	task := &taskdomain.Task{}

	timerTime := taskdomain.TimerTime{
		Hour:   10,
		Minute: 30,
		Second: 0,
	}

	timerTask, err := creator.CreateTimerTask(task, taskdomain.ALARM, timerTime)
	if err != nil {
		t.Fatalf("CreateTimerTask() error = %v", err)
	}

	if timerTask == nil {
		t.Fatal("CreateTimerTask() returned nil")
	}

	returnedTask := timerTask.Task()
	if returnedTask != task {
		t.Error("Returned task should be the same instance as input")
	}
}

func TestTimerTaskCreatorCreateTimerTaskWithZeroTimerTime(t *testing.T) {
	t.Parallel()

	creator := &TimerTaskCreator{}

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTime := taskdomain.TimerTime{}

	timerTask, err := creator.CreateTimerTask(task, taskdomain.ALARM, timerTime)
	if err != nil {
		t.Fatalf("CreateTimerTask() error = %v", err)
	}

	if timerTask == nil {
		t.Fatal("CreateTimerTask() returned nil")
	}
}

func TestTimerTaskCreatorMultipleCreations(t *testing.T) {
	t.Parallel()

	creator := &TimerTaskCreator{}

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTime := taskdomain.TimerTime{
		Hour:   10,
		Minute: 30,
		Second: 0,
	}

	// Create multiple timer tasks
	timerTasks := make([]taskdomain.TimerTask, 5)
	for i := 0; i < 5; i++ {
		timerTask, err := creator.CreateTimerTask(task, taskdomain.ALARM, timerTime)
		if err != nil {
			t.Fatalf("CreateTimerTask() error = %v", err)
		}
		timerTasks[i] = timerTask
	}

	// Verify all tasks are created
	for i, timerTask := range timerTasks {
		if timerTask == nil {
			t.Errorf("TimerTask[%d] is nil", i)
		}
	}

	// Verify each task has the correct type
	for i, timerTask := range timerTasks {
		returnedTask := timerTask.Task()
		if returnedTask != task {
			t.Errorf("TimerTask[%d] Task() should return the same instance", i)
		}
	}
}

func TestTimerTaskCreatorInterfaceCompliance(t *testing.T) {
	t.Parallel()

	// Verify that TimerTaskCreator implements taskdomain.TimerTaskCreator interface
	var _ taskdomain.TimerTaskCreator = &TimerTaskCreator{}
}

func TestTimerTaskCreatorReturnsTimerTaskType(t *testing.T) {
	t.Parallel()

	creator := &TimerTaskCreator{}

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTime := taskdomain.TimerTime{
		Hour:   10,
		Minute: 30,
		Second: 0,
	}

	timerTask, err := creator.CreateTimerTask(task, taskdomain.ALARM, timerTime)
	if err != nil {
		t.Fatalf("CreateTimerTask() error = %v", err)
	}

	// Verify the returned type is *TimerTask
	if _, ok := timerTask.(*TimerTask); !ok {
		t.Error("CreateTimerTask() should return *TimerTask")
	}
}

func TestTimerTaskCreatorWithComplexTimerTime(t *testing.T) {
	t.Parallel()

	creator := &TimerTaskCreator{}

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	timerTime := taskdomain.TimerTime{
		Year:   2024,
		Month:  12,
		Day:    31,
		Hour:   23,
		Minute: 59,
		Second: 59,
	}

	timerTask, err := creator.CreateTimerTask(task, taskdomain.ALARM, timerTime)
	if err != nil {
		t.Fatalf("CreateTimerTask() error = %v", err)
	}

	if timerTask == nil {
		t.Fatal("CreateTimerTask() returned nil")
	}

	returnedTask := timerTask.Task()
	if returnedTask != task {
		t.Error("Returned task should be the same instance as input")
	}
}

func TestTimerTaskCreatorWithDuration(t *testing.T) {
	t.Parallel()

	creator := &TimerTaskCreator{}

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_timer",
		},
	}

	// Create a timer with a complex duration
	timerTime := taskdomain.TimerTime{
		Year:   1,
		Month:  2,
		Day:    3,
		Hour:   4,
		Minute: 5,
		Second: 6,
	}

	timerTask, err := creator.CreateTimerTask(task, taskdomain.TIMEOUT, timerTime)
	if err != nil {
		t.Fatalf("CreateTimerTask() error = %v", err)
	}

	if timerTask == nil {
		t.Fatal("CreateTimerTask() returned nil")
	}

	// Verify the timer can be used
	timerTask.SetSchedule()

	if timerTask.GetSchedule() == "" {
		t.Error("Schedule should not be empty after SetSchedule()")
	}
}
