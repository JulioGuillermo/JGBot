package taskapplication

import (
	"errors"
	"testing"

	taskdomain "JGBot/task/domain"
)

// Mock implementations for TimerService tests

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

type mockTimerCreator struct {
	err error
}

func (m *mockTimerCreator) CreateTimerTask(task *taskdomain.Task, timerType taskdomain.TimerType, timerTime taskdomain.TimerTime) (taskdomain.TimerTask, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &mockTimerTask{task: task}, nil
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

// Helper to create test tasks
func newTimerTestTask(origin, name string) *taskdomain.Task {
	return &taskdomain.Task{
		TaskOriginInfo: taskdomain.TaskOriginInfo{
			Origin: origin,
		},
		TaskInfo: taskdomain.TaskInfo{
			Name: name,
		},
	}
}

// ==================== TimerService Tests ====================

func TestTimerService_AddTimeout(t *testing.T) {
	t.Parallel()

	existingTimer := &mockTimerTask{task: newTimerTestTask("origin1", "timer1")}

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
			task:      newTimerTestTask("origin1", "timer1"),
			timerTime: taskdomain.TimerTime{},
			wantErr:   false,
		},
		{
			name:       "duplicate timer",
			timers:     []taskdomain.TimerTask{existingTimer},
			task:       newTimerTestTask("origin1", "timer1"),
			timerTime:  taskdomain.TimerTime{},
			wantErr:    true,
			wantErrStr: "timer with name `timer1` already exists",
		},
	}

	for _, tt := range tests {
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

	err := svc.AddAlarm(newTimerTestTask("origin1", "alarm1"), taskdomain.TimerTime{})
	if err != nil {
		t.Errorf("AddAlarm() error = %v, wantErr false", err)
	}
}

func TestTimerService_GetTimer(t *testing.T) {
	t.Parallel()

	task1 := newTimerTestTask("origin1", "timer1")
	task2 := newTimerTestTask("origin1", "timer2")

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

	task1 := newTimerTestTask("origin1", "timer1")

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
		&mockTimerTask{task: newTimerTestTask("origin1", "timer1")},
		&mockTimerTask{task: newTimerTestTask("origin1", "timer2")},
		&mockTimerTask{task: newTimerTestTask("origin2", "timer3")},
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
			&mockTimerTask{task: newTimerTestTask("origin1", "timer1")},
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

func TestTimerService_SaveTimers(t *testing.T) {
	t.Parallel()

	t.Run("save timers with error", func(t *testing.T) {
		t.Parallel()

		svc := NewTimerService(
			&mockTimerCreator{},
			&mockTimerPersister{err: errors.New("save error")},
			nil,
		)
		svc.Timers = []taskdomain.TimerTask{&mockTimerTask{task: newTimerTestTask("origin1", "timer1")}}

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
		svc.Timers = []taskdomain.TimerTask{&mockTimerTask{task: newTimerTestTask("origin1", "timer1")}}

		err := svc.SaveTimers()
		if err != nil {
			t.Errorf("SaveTimers() error = %v, want nil", err)
		}
		if len(persister.timers) != 1 {
			t.Errorf("SaveTimers() persisted timers len = %d, want 1", len(persister.timers))
		}
	})
}
