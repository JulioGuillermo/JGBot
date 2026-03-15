package taskinfrastructure

import (
	taskdomain "JGBot/task/domain"
	"path/filepath"
	"testing"
)

func TestCronTaskPersistence(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "cron_test.json")

	persistence := &CronTaskPersistence{
		FilePath: filePath,
	}

	// 1. Test LoadCrons when file doesn't exist (should return empty list and create file)
	crons, err := persistence.LoadCrons()
	if err != nil {
		t.Fatalf("LoadCrons() error = %v", err)
	}
	if len(crons) != 0 {
		t.Errorf("LoadCrons() should return empty list initially, got %d", len(crons))
	}

	// 2. Test SaveCrons
	task := &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: "test",
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: "test_cron",
		},
	}

	schedule := taskdomain.CronArgs{
		Minute: "0",
		Hour:   "12",
	}

	cronTask := &CronTask{
		CronTask: task,
		Schedule: schedule,
	}

	err = persistence.SaveCrons([]taskdomain.CronTask{cronTask})
	if err != nil {
		t.Fatalf("SaveCrons() error = %v", err)
	}

	// 3. Test LoadCrons again (should return saved task)
	crons, err = persistence.LoadCrons()
	if err != nil {
		t.Fatalf("LoadCrons() error = %v", err)
	}

	if len(crons) != 1 {
		t.Fatalf("LoadCrons() should return 1 task, got %d", len(crons))
	}

	loadedTask := crons[0]
	if loadedTask.Task().Name != "test_cron" {
		t.Errorf("Loaded task name = %v, want test_cron", loadedTask.Task().Name)
	}

	if loadedTask.GetSchedule() != schedule.String() {
		t.Errorf("Loaded task schedule = %v, want %v", loadedTask.GetSchedule(), schedule.String())
	}
}
