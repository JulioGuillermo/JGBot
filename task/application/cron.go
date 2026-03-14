package taskapplication

import (
	taskdomain "JGBot/task/domain"
	"errors"
	"slices"
	"strings"
)

type CronService struct {
	TaskCreator  taskdomain.CronTaskCreator
	Persister    taskdomain.CronTaskPersister
	OnActivation taskdomain.TaskActivationHandler
	CronMng      taskdomain.CronMng
	Crons        []taskdomain.CronTask
}

func (s *CronService) LoadCrons() error {
	crons, err := s.Persister.LoadCrons()
	if err != nil {
		return err
	}
	s.Crons = crons
	for _, cron := range crons {
		cron.Activate(s.CronMng, s.onActivation)
	}
	return nil
}

func (s *CronService) SaveCrons() error {
	return s.Persister.SaveCrons(s.Crons)
}

func (s *CronService) GetJob(origin, name string) taskdomain.CronTask {
	for _, cron := range s.Crons {
		if cron.Task().Origin == origin && cron.Task().Name == name {
			return cron
		}
	}
	return nil
}

func (s *CronService) RemoveJob(origin, name string) error {
	cron := s.GetJob(origin, name)
	if cron == nil {
		return errors.New("cron with name `" + name + "` not found")
	}
	defer s.SaveCrons()

	cron.Close(s.CronMng)

	s.Crons = slices.DeleteFunc(s.Crons, func(cron taskdomain.CronTask) bool {
		return cron.Task().Origin == origin && cron.Task().Name == name
	})
	return nil
}

func (s *CronService) onActivation(task *taskdomain.Task, schedule string) {
	if s.OnActivation != nil {
		s.OnActivation(task, schedule)
	}
}

func (s *CronService) AddJob(
	task *taskdomain.Task,
	args taskdomain.CronArgs,
) error {
	if s.GetJob(task.Origin, task.Name) != nil {
		return errors.New("cron with name `" + task.Name + "` already exists")
	}
	defer s.SaveCrons()

	cron, err := s.TaskCreator.CreateCronTask(task, args)
	if err != nil {
		return err
	}

	cron.Activate(s.CronMng, s.OnActivation)

	s.Crons = append(s.Crons, cron)
	return nil
}

func (s *CronService) ListJob(origin string) []*taskdomain.Task {
	crons := []*taskdomain.Task{}
	for _, cron := range s.Crons {
		task := cron.Task()
		if task.Origin == origin {
			crons = append(crons, task)
		}
	}
	slices.SortFunc(crons, func(a, b *taskdomain.Task) int {
		return strings.Compare(a.Name, b.Name)
	})
	return crons
}
