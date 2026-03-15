package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"testing"
)

func TestCronTaskCreatorCreateCronTask(t *testing.T) {
	t.Parallel()

	creator := &CronTaskCreator{}

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

	schedule := taskdomain.CronArgs{
		Minute:  "30",
		Hour:    "14",
		Day:     "15",
		Month:   "6",
		Weekday: "3",
	}

	cronTask, err := creator.CreateCronTask(task, schedule)
	if err != nil {
		t.Fatalf("CreateCronTask() error = %v", err)
	}

	if cronTask == nil {
		t.Fatal("CreateCronTask() returned nil")
	}

	// Verify the returned task
	returnedTask := cronTask.Task()
	if returnedTask != task {
		t.Error("Returned task should be the same instance as input")
	}

	// Verify the schedule
	returnedSchedule := cronTask.GetSchedule()
	expectedSchedule := schedule.String()
	if returnedSchedule != expectedSchedule {
		t.Errorf("GetSchedule() = %v, want %v", returnedSchedule, expectedSchedule)
	}
}

func TestCronTaskCreatorCreateCronTaskWithDifferentSchedules(t *testing.T) {
	t.Parallel()

	creator := &CronTaskCreator{}

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
		{
			name: "specific date",
			schedule: taskdomain.CronArgs{
				Minute:  "30",
				Hour:    "14",
				Day:     "15",
				Month:   "6",
				Weekday: "*",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cronTask, err := creator.CreateCronTask(task, tt.schedule)
			if err != nil {
				t.Errorf("CreateCronTask() error = %v", err)
				return
			}

			if cronTask == nil {
				t.Error("CreateCronTask() returned nil")
				return
			}

			returnedSchedule := cronTask.GetSchedule()
			expectedSchedule := tt.schedule.String()
			if returnedSchedule != expectedSchedule {
				t.Errorf("GetSchedule() = %v, want %v", returnedSchedule, expectedSchedule)
			}
		})
	}
}

func TestCronTaskCreatorCreateCronTaskWithEmptyTask(t *testing.T) {
	t.Parallel()

	creator := &CronTaskCreator{}

	task := &taskdomain.Task{}

	schedule := taskdomain.CronArgs{
		Minute:  "0",
		Hour:    "12",
		Day:     "1",
		Month:   "1",
		Weekday: "1",
	}

	cronTask, err := creator.CreateCronTask(task, schedule)
	if err != nil {
		t.Fatalf("CreateCronTask() error = %v", err)
	}

	if cronTask == nil {
		t.Fatal("CreateCronTask() returned nil")
	}

	returnedTask := cronTask.Task()
	if returnedTask != task {
		t.Error("Returned task should be the same instance as input")
	}
}

func TestCronTaskCreatorCreateCronTaskWithEmptySchedule(t *testing.T) {
	t.Parallel()

	creator := &CronTaskCreator{}

	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := taskdomain.CronArgs{}

	cronTask, err := creator.CreateCronTask(task, schedule)
	if err != nil {
		t.Fatalf("CreateCronTask() error = %v", err)
	}

	if cronTask == nil {
		t.Fatal("CreateCronTask() returned nil")
	}

	returnedSchedule := cronTask.GetSchedule()
	expectedSchedule := schedule.String()
	if returnedSchedule != expectedSchedule {
		t.Errorf("GetSchedule() = %v, want %v", returnedSchedule, expectedSchedule)
	}
}

func TestCronTaskCreatorMultipleCreations(t *testing.T) {
	t.Parallel()

	creator := &CronTaskCreator{}

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

	// Create multiple cron tasks
	cronTasks := make([]taskdomain.CronTask, 5)
	for i := range 5 {
		cronTask, err := creator.CreateCronTask(task, schedule)
		if err != nil {
			t.Fatalf("CreateCronTask() error = %v", err)
		}
		cronTasks[i] = cronTask
	}

	// Verify all tasks are created
	for i, cronTask := range cronTasks {
		if cronTask == nil {
			t.Errorf("CronTask[%d] is nil", i)
		}
	}

	// Verify each task has the correct schedule
	for i, cronTask := range cronTasks {
		returnedSchedule := cronTask.GetSchedule()
		expectedSchedule := schedule.String()
		if returnedSchedule != expectedSchedule {
			t.Errorf("CronTask[%d] GetSchedule() = %v, want %v", i, returnedSchedule, expectedSchedule)
		}
	}
}

func TestCronTaskCreatorInterfaceCompliance(t *testing.T) {
	t.Parallel()

	// Verify that CronTaskCreator implements taskdomain.CronTaskCreator interface
	var _ taskdomain.CronTaskCreator = &CronTaskCreator{}
}

func TestCronTaskCreatorReturnsCronTaskType(t *testing.T) {
	t.Parallel()

	creator := &CronTaskCreator{}

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

	cronTask, err := creator.CreateCronTask(task, schedule)
	if err != nil {
		t.Fatalf("CreateCronTask() error = %v", err)
	}

	// Verify the returned type is *CronTask
	if _, ok := cronTask.(*CronTask); !ok {
		t.Error("CreateCronTask() should return *CronTask")
	}
}
