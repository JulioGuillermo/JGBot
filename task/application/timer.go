package taskapplication

import (
	taskdomain "JGBot/task/domain"
	"errors"
	"slices"
	"strings"
)

type TimerService struct {
	TaskCreator  taskdomain.TimerTaskCreator
	Persister    taskdomain.TimerPersistence
	OnActivation taskdomain.TaskActivationHandler
	Timers       []taskdomain.TimerTask
}

func NewTimerService(taskCreator taskdomain.TimerTaskCreator, persister taskdomain.TimerPersistence, onActivation taskdomain.TaskActivationHandler) *TimerService {
	return &TimerService{
		TaskCreator:  taskCreator,
		Persister:    persister,
		OnActivation: onActivation,
	}
}

func (s *TimerService) LoadTimers() error {
	timers, err := s.Persister.LoadTimers()
	if err != nil {
		return err
	}
	s.Timers = timers
	for _, timer := range timers {
		timer.Activate(s.onActivation)
	}
	return nil
}

func (s *TimerService) SaveTimers() error {
	return s.Persister.SaveTimers(s.Timers)
}

func (s *TimerService) GetTimer(origin, name string) taskdomain.TimerTask {
	for _, timer := range s.Timers {
		if timer.Task().Origin == origin && timer.Task().Name == name {
			return timer
		}
	}
	return nil
}

func (s *TimerService) RemoveTimer(origin, name string) error {
	timer := s.GetTimer(origin, name)
	if timer == nil {
		return errors.New("timer with name " + name + " not found")
	}
	defer s.SaveTimers()

	timer.Close()

	s.Timers = slices.DeleteFunc(s.Timers, func(timer taskdomain.TimerTask) bool {
		return timer.Task().Name == name && timer.Task().Origin == origin
	})
	return nil
}

func (s *TimerService) onActivation(task *taskdomain.Task, schedule string) {
	if s.OnActivation != nil {
		s.OnActivation(task, schedule)
	}
	s.RemoveTimer(task.Origin, task.Name)
}

func (s *TimerService) addTimer(
	task *taskdomain.Task,
	timerType taskdomain.TimerType,
	timerTime taskdomain.TimerTime,
) error {
	if s.GetTimer(task.Origin, task.Name) != nil {
		return errors.New("timer with name `" + task.Name + "` already exists")
	}
	defer s.SaveTimers()

	timerTask, err := s.TaskCreator.CreateTimerTask(
		task,
		timerType,
		timerTime,
	)

	if err != nil {
		return err
	}

	timerTask.SetSchedule()
	timerTask.Activate(s.onActivation)

	s.Timers = append(s.Timers, timerTask)
	return nil
}

func (s *TimerService) AddTimeout(
	task *taskdomain.Task,
	timerTime taskdomain.TimerTime,
) error {
	return s.addTimer(
		task,
		taskdomain.TIMEOUT,
		timerTime,
	)
}

func (s *TimerService) AddAlarm(
	task *taskdomain.Task,
	timerTime taskdomain.TimerTime,
) error {
	return s.addTimer(
		task,
		taskdomain.ALARM,
		timerTime,
	)
}

func (s *TimerService) ListTimers(origin string) []*taskdomain.Task {
	var timers []*taskdomain.Task
	for _, timer := range s.Timers {
		task := timer.Task()
		if task.Origin == origin {
			timers = append(timers, task)
		}
	}
	slices.SortFunc(timers, func(a, b *taskdomain.Task) int {
		return strings.Compare(a.Name, b.Name)
	})
	return timers
}
