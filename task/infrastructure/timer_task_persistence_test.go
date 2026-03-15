package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"path/filepath"
	"testing"
	"time"
)

func TestTimerTaskPersistence(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "timer_test.json")

	persistence := &TimerTaskPersistence{
		FilePath: filePath,
	}

	// 1. Test LoadTimers when file doesn't exist (should return empty list and create file)
	timers, err := persistence.LoadTimers()
	if err != nil {
		t.Fatalf("LoadTimers() error = %v", err)
	}
	if len(timers) != 0 {
		t.Errorf("LoadTimers() should return empty list initially, got %d", len(timers))
	}

	// 2. Test SaveTimers
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
		Month:  6,
		Day:    15,
		Hour:   14,
		Minute: 30,
	}

	timerTask := &TimerTask{
		TimerTask: task,
		Type:      taskdomain.ALARM,
		Time:      timerTime,
		Schedule:  time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC),
	}

	err = persistence.SaveTimers([]taskdomain.TimerTask{timerTask})
	if err != nil {
		t.Fatalf("SaveTimers() error = %v", err)
	}

	// 3. Test LoadTimers again (should return saved task)
	timers, err = persistence.LoadTimers()
	if err != nil {
		t.Fatalf("LoadTimers() error = %v", err)
	}

	if len(timers) != 1 {
		t.Fatalf("LoadTimers() should return 1 task, got %d", len(timers))
	}

	loadedTask := timers[0]
	if loadedTask.Task().Name != "test_timer" {
		t.Errorf("Loaded task name = %v, want test_timer", loadedTask.Task().Name)
	}
}
