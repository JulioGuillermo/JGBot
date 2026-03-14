package taskinfrastructure

import (
	"github.com/robfig/cron/v3"
)

type CronMng struct {
	*cron.Cron
}

func NewCronMng() *CronMng {
	cron := cron.New()
	go cron.Start()

	return &CronMng{
		Cron: cron,
	}
}

func (m *CronMng) AddFunc(spec string, cmd func()) (int, error) {
	id, err := m.Cron.AddFunc(spec, cmd)
	return int(id), err
}

func (m *CronMng) Remove(id int) {
	m.Cron.Remove(cron.EntryID(id))
}
