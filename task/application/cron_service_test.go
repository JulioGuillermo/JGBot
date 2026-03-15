package taskapplication

import (
	"errors"
	"testing"

	taskdomain "JGBot/task/domain"
)

// Mock implementations for CronService tests

type mockCronTask struct {
	task     *taskdomain.Task
	schedule string
	closed   bool
	active   bool
}

func (m *mockCronTask) Task() *taskdomain.Task { return m.task }
func (m *mockCronTask) GetSchedule() string    { return m.schedule }
func (m *mockCronTask) Activate(_ taskdomain.CronMng, _ taskdomain.TaskActivationHandler) {
	m.active = true
}
func (m *mockCronTask) Close(_ taskdomain.CronMng) { m.closed = true }

type mockCronCreator struct {
	err error
}

func (m *mockCronCreator) CreateCronTask(task *taskdomain.Task, args taskdomain.CronArgs) (taskdomain.CronTask, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &mockCronTask{task: task, schedule: args.CronString()}, nil
}

type mockCronPersister struct {
	crons []taskdomain.CronTask
	err   error
}

func (m *mockCronPersister) LoadCrons() ([]taskdomain.CronTask, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.crons, nil
}

func (m *mockCronPersister) SaveCrons(crons []taskdomain.CronTask) error {
	if m.err != nil {
		return m.err
	}
	m.crons = crons
	return nil
}

type mockCronMng struct {
	entries map[int]func()
}

func (m *mockCronMng) AddFunc(s string, f func()) (int, error) {
	id := len(m.entries) + 1
	m.entries[id] = f
	return id, nil
}

func (m *mockCronMng) Remove(id int) {
	delete(m.entries, id)
}

// Helper to create test tasks
func newTestTask(origin, name string) *taskdomain.Task {
	return &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: origin,
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: name,
		},
	}
}

// ==================== CronService Tests ====================

func TestCronService_AddJob(t *testing.T) {
	t.Parallel()

	existingCron := &mockCronTask{task: newTestTask("origin1", "job1")}

	tests := []struct {
		name       string
		crons      []taskdomain.CronTask
		task       *taskdomain.Task
		args       taskdomain.CronArgs
		wantErr    bool
		wantErrStr string
	}{
		{
			name:    "add new cron job",
			crons:   []taskdomain.CronTask{},
			task:    newTestTask("origin1", "job1"),
			args:    taskdomain.CronArgs{Minute: "1", Hour: "1", Day: "1", Month: "1", Weekday: "1"},
			wantErr: false,
		},
		{
			name:       "duplicate cron job",
			crons:      []taskdomain.CronTask{existingCron},
			task:       newTestTask("origin1", "job1"),
			args:       taskdomain.CronArgs{Minute: "1", Hour: "1", Day: "1", Month: "1", Weekday: "1"},
			wantErr:    true,
			wantErrStr: "cron with name `job1` already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := &CronService{
				TaskCreator: &mockCronCreator{},
				Persister:   &mockCronPersister{},
				CronMng:     &mockCronMng{entries: make(map[int]func())},
				Crons:       tt.crons,
			}

			err := svc.AddJob(tt.task, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.wantErrStr {
				t.Errorf("AddJob() error = %v, want %v", err.Error(), tt.wantErrStr)
			}
		})
	}
}

func TestCronService_GetJob(t *testing.T) {
	t.Parallel()

	task1 := newTestTask("origin1", "job1")
	task2 := newTestTask("origin1", "job2")

	svc := &CronService{
		Crons: []taskdomain.CronTask{
			&mockCronTask{task: task1},
			&mockCronTask{task: task2},
		},
	}

	tests := []struct {
		name      string
		origin    string
		jobName   string
		wantFound bool
	}{
		{
			name:      "existing job",
			origin:    "origin1",
			jobName:   "job1",
			wantFound: true,
		},
		{
			name:      "non-existing job",
			origin:    "origin1",
			jobName:   "job3",
			wantFound: false,
		},
		{
			name:      "wrong origin",
			origin:    "origin2",
			jobName:   "job1",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := svc.GetJob(tt.origin, tt.jobName)
			if (got != nil) != tt.wantFound {
				t.Errorf("GetJob() got = %v, wantFound %v", got, tt.wantFound)
			}
		})
	}
}

func TestCronService_RemoveJob(t *testing.T) {
	t.Parallel()

	task1 := newTestTask("origin1", "job1")
	task2 := newTestTask("origin1", "job2")

	svc := &CronService{
		TaskCreator: &mockCronCreator{},
		Persister:   &mockCronPersister{},
		CronMng:     &mockCronMng{entries: make(map[int]func())},
		Crons: []taskdomain.CronTask{
			&mockCronTask{task: task1},
			&mockCronTask{task: task2},
		},
	}

	tests := []struct {
		name       string
		origin     string
		jobName    string
		wantErr    bool
		wantErrStr string
	}{
		{
			name:    "remove existing job",
			origin:  "origin1",
			jobName: "job1",
			wantErr: false,
		},
		{
			name:       "remove non-existing job",
			origin:     "origin1",
			jobName:    "job3",
			wantErr:    true,
			wantErrStr: "cron with name `job3` not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := svc.RemoveJob(tt.origin, tt.jobName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.wantErrStr {
				t.Errorf("RemoveJob() error = %v, want %v", err.Error(), tt.wantErrStr)
			}
		})
	}
}

func TestCronService_ListJob(t *testing.T) {
	t.Parallel()

	svc := &CronService{
		Crons: []taskdomain.CronTask{
			&mockCronTask{task: newTestTask("origin1", "job1")},
			&mockCronTask{task: newTestTask("origin1", "job2")},
			&mockCronTask{task: newTestTask("origin2", "job3")},
		},
	}

	tests := []struct {
		name      string
		origin    string
		wantCount int
	}{
		{
			name:      "list jobs for origin1",
			origin:    "origin1",
			wantCount: 2,
		},
		{
			name:      "list jobs for origin2",
			origin:    "origin2",
			wantCount: 1,
		},
		{
			name:      "list jobs for non-existing origin",
			origin:    "origin3",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := svc.ListJob(tt.origin)
			if len(got) != tt.wantCount {
				t.Errorf("ListJob() len = %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestCronService_LoadCrons(t *testing.T) {
	t.Parallel()

	t.Run("load crons with error", func(t *testing.T) {
		t.Parallel()

		svc := &CronService{
			TaskCreator: &mockCronCreator{},
			Persister:   &mockCronPersister{err: errors.New("load error")},
			CronMng:     &mockCronMng{entries: make(map[int]func())},
		}

		err := svc.LoadCrons()
		if err == nil {
			t.Error("LoadCrons() should return error")
		}
	})

	t.Run("load crons success", func(t *testing.T) {
		t.Parallel()

		crons := []taskdomain.CronTask{
			&mockCronTask{task: newTestTask("origin1", "job1")},
		}
		persister := &mockCronPersister{crons: crons}

		svc := &CronService{
			TaskCreator: &mockCronCreator{},
			Persister:   persister,
			CronMng:     &mockCronMng{entries: make(map[int]func())},
		}

		err := svc.LoadCrons()
		if err != nil {
			t.Errorf("LoadCrons() error = %v, want nil", err)
		}
		if len(svc.Crons) != 1 {
			t.Errorf("LoadCrons() crons len = %d, want 1", len(svc.Crons))
		}
	})
}

func TestCronService_SaveCrons(t *testing.T) {
	t.Parallel()

	t.Run("save crons with error", func(t *testing.T) {
		t.Parallel()

		svc := &CronService{
			TaskCreator: &mockCronCreator{},
			Persister:   &mockCronPersister{err: errors.New("save error")},
			CronMng:     &mockCronMng{entries: make(map[int]func())},
			Crons:       []taskdomain.CronTask{&mockCronTask{task: newTestTask("origin1", "job1")}},
		}

		err := svc.SaveCrons()
		if err == nil {
			t.Error("SaveCrons() should return error")
		}
	})

	t.Run("save crons success", func(t *testing.T) {
		t.Parallel()

		persister := &mockCronPersister{}

		svc := &CronService{
			TaskCreator: &mockCronCreator{},
			Persister:   persister,
			CronMng:     &mockCronMng{entries: make(map[int]func())},
			Crons:       []taskdomain.CronTask{&mockCronTask{task: newTestTask("origin1", "job1")}},
		}

		err := svc.SaveCrons()
		if err != nil {
			t.Errorf("SaveCrons() error = %v, want nil", err)
		}
		if len(persister.crons) != 1 {
			t.Errorf("SaveCrons() persisted crons len = %d, want 1", len(persister.crons))
		}
	})
}
