package taskapplication

import (
	"errors"
	"testing"

	taskdomain "JGBot/task/domain"
)

// Mock implementations

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

type mockTimerTask struct {
	task     *taskdomain.Task
	schedule string
	closed   bool
	active   bool
}

func (m *mockTimerTask) Task() *taskdomain.Task { return m.task }
func (m *mockTimerTask) GetSchedule() string    { return m.schedule }
func (m *mockTimerTask) SetSchedule()           { m.schedule = "set" }
func (m *mockTimerTask) Activate(_ taskdomain.TaskActivationHandler) {
	m.active = true
}
func (m *mockTimerTask) Close() { m.closed = true }

type mockCronCreator struct {
	err error
}

func (m *mockCronCreator) CreateCronTask(task *taskdomain.Task, args taskdomain.CronArgs) (taskdomain.CronTask, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &mockCronTask{task: task, schedule: args.CronString()}, nil
}

type mockTimerCreator struct {
	err error
}

func (m *mockTimerCreator) CreateTimerTask(task *taskdomain.Task, timerType taskdomain.TimerType, timerTime taskdomain.TimerTime) (taskdomain.TimerTask, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &mockTimerTask{task: task}, nil
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

type mockTimerPersister struct {
	timers []taskdomain.TimerTask
	err    error
}

func (m *mockTimerPersister) LoadTimers() ([]taskdomain.TimerTask, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.timers, nil
}

func (m *mockTimerPersister) SaveTimers(timers []taskdomain.TimerTask) error {
	if m.err != nil {
		return m.err
	}
	m.timers = timers
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
		tt := tt
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
		tt := tt
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
		tt := tt
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := svc.ListJob(tt.origin)
			if len(got) != tt.wantCount {
				t.Errorf("ListJob() len = %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}

// ==================== TimerService Tests ====================

func TestTimerService_AddTimeout(t *testing.T) {
	t.Parallel()

	existingTimer := &mockTimerTask{task: newTestTask("origin1", "timer1")}

	tests := []struct {
		name       string
		timers     []taskdomain.TimerTask
		task       *taskdomain.Task
		timerTime  taskdomain.TimerTime
		wantErr    bool
		wantErrStr string
	}{
		{
			name:      "add new timeout",
			timers:    []taskdomain.TimerTask{},
			task:      newTestTask("origin1", "timer1"),
			timerTime: taskdomain.TimerTime{},
			wantErr:   false,
		},
		{
			name:       "duplicate timer",
			timers:     []taskdomain.TimerTask{existingTimer},
			task:       newTestTask("origin1", "timer1"),
			timerTime:  taskdomain.TimerTime{},
			wantErr:    true,
			wantErrStr: "timer with name `timer1` already exists",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewTimerService(
				&mockTimerCreator{},
				&mockTimerPersister{},
				nil,
			)
			svc.Timers = tt.timers

			err := svc.AddTimeout(tt.task, tt.timerTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.wantErrStr {
				t.Errorf("AddTimeout() error = %v, want %v", err.Error(), tt.wantErrStr)
			}
		})
	}
}

func TestTimerService_AddAlarm(t *testing.T) {
	t.Parallel()

	svc := NewTimerService(
		&mockTimerCreator{},
		&mockTimerPersister{},
		nil,
	)

	err := svc.AddAlarm(newTestTask("origin1", "alarm1"), taskdomain.TimerTime{})
	if err != nil {
		t.Errorf("AddAlarm() error = %v, wantErr false", err)
	}
}

func TestTimerService_GetTimer(t *testing.T) {
	t.Parallel()

	task1 := newTestTask("origin1", "timer1")
	task2 := newTestTask("origin1", "timer2")

	svc := NewTimerService(
		&mockTimerCreator{},
		&mockTimerPersister{},
		nil,
	)
	svc.Timers = []taskdomain.TimerTask{
		&mockTimerTask{task: task1},
		&mockTimerTask{task: task2},
	}

	tests := []struct {
		name      string
		origin    string
		timerName string
		wantFound bool
	}{
		{
			name:      "existing timer",
			origin:    "origin1",
			timerName: "timer1",
			wantFound: true,
		},
		{
			name:      "non-existing timer",
			origin:    "origin1",
			timerName: "timer3",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := svc.GetTimer(tt.origin, tt.timerName)
			if (got != nil) != tt.wantFound {
				t.Errorf("GetTimer() got = %v, wantFound %v", got, tt.wantFound)
			}
		})
	}
}

func TestTimerService_RemoveTimer(t *testing.T) {
	t.Parallel()

	task1 := newTestTask("origin1", "timer1")

	svc := NewTimerService(
		&mockTimerCreator{},
		&mockTimerPersister{},
		nil,
	)
	svc.Timers = []taskdomain.TimerTask{
		&mockTimerTask{task: task1},
	}

	tests := []struct {
		name       string
		origin     string
		timerName  string
		wantErr    bool
		wantErrStr string
	}{
		{
			name:      "remove existing timer",
			origin:    "origin1",
			timerName: "timer1",
			wantErr:   false,
		},
		{
			name:       "remove non-existing timer",
			origin:     "origin1",
			timerName:  "timer2",
			wantErr:    true,
			wantErrStr: "timer with name `timer2` not found",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := svc.RemoveTimer(tt.origin, tt.timerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveTimer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.wantErrStr {
				t.Errorf("RemoveTimer() error = %v, want %v", err.Error(), tt.wantErrStr)
			}
		})
	}
}

func TestTimerService_ListTimers(t *testing.T) {
	t.Parallel()

	svc := NewTimerService(
		&mockTimerCreator{},
		&mockTimerPersister{},
		nil,
	)
	svc.Timers = []taskdomain.TimerTask{
		&mockTimerTask{task: newTestTask("origin1", "timer1")},
		&mockTimerTask{task: newTestTask("origin1", "timer2")},
		&mockTimerTask{task: newTestTask("origin2", "timer3")},
	}

	tests := []struct {
		name      string
		origin    string
		wantCount int
	}{
		{
			name:      "list timers for origin1",
			origin:    "origin1",
			wantCount: 2,
		},
		{
			name:      "list timers for non-existing origin",
			origin:    "origin3",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := svc.ListTimers(tt.origin)
			if len(got) != tt.wantCount {
				t.Errorf("ListTimers() len = %d, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestTimerService_LoadTimers(t *testing.T) {
	t.Parallel()

	t.Run("load timers with error", func(t *testing.T) {
		t.Parallel()

		svc := NewTimerService(
			&mockTimerCreator{},
			&mockTimerPersister{err: errors.New("load error")},
			nil,
		)

		err := svc.LoadTimers()
		if err == nil {
			t.Error("LoadTimers() should return error")
		}
	})

	t.Run("load timers success", func(t *testing.T) {
		t.Parallel()

		timers := []taskdomain.TimerTask{
			&mockTimerTask{task: newTestTask("origin1", "timer1")},
		}
		persister := &mockTimerPersister{timers: timers}

		svc := NewTimerService(
			&mockTimerCreator{},
			persister,
			nil,
		)

		err := svc.LoadTimers()
		if err != nil {
			t.Errorf("LoadTimers() error = %v, want nil", err)
		}
		if len(svc.Timers) != 1 {
			t.Errorf("LoadTimers() timers len = %d, want 1", len(svc.Timers))
		}
	})
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

func TestTimerService_SaveTimers(t *testing.T) {
	t.Parallel()

	t.Run("save timers with error", func(t *testing.T) {
		t.Parallel()

		svc := NewTimerService(
			&mockTimerCreator{},
			&mockTimerPersister{err: errors.New("save error")},
			nil,
		)
		svc.Timers = []taskdomain.TimerTask{&mockTimerTask{task: newTestTask("origin1", "timer1")}}

		err := svc.SaveTimers()
		if err == nil {
			t.Error("SaveTimers() should return error")
		}
	})

	t.Run("save timers success", func(t *testing.T) {
		t.Parallel()

		persister := &mockTimerPersister{}

		svc := NewTimerService(
			&mockTimerCreator{},
			persister,
			nil,
		)
		svc.Timers = []taskdomain.TimerTask{&mockTimerTask{task: newTestTask("origin1", "timer1")}}

		err := svc.SaveTimers()
		if err != nil {
			t.Errorf("SaveTimers() error = %v, want nil", err)
		}
		if len(persister.timers) != 1 {
			t.Errorf("SaveTimers() persisted timers len = %d, want 1", len(persister.timers))
		}
	})
}
