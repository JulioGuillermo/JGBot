package sessionconf

import (
	"JGBot/conf"
	"JGBot/session/sessionconf/sc"
)

type SessionCtl struct {
	Config   *SessionConfs
	Unconfig *SessionConfs
}

func NewSessionCtl() (*SessionCtl, error) {
	ctl := &SessionCtl{}
	var err error

	ctl.Config, err = NewSessionConfs(conf.SessionFile)
	if err != nil {
		return nil, err
	}

	ctl.Unconfig, err = NewSessionConfs(conf.UnconfigSessionFile)
	if err != nil {
		return nil, err
	}

	ctl.Config.OnChange = ctl.onConfigChange
	ctl.Config.Watch()

	return ctl, nil
}

func (s *SessionCtl) GetConfigOrigin(origin string) *sc.SessionConf {
	return s.Config.GetOrigin(origin)
}

func (s *SessionCtl) AddUnconfig(name string, id string, origin string) {
	if s.Unconfig.HasOrigin(origin) {
		return
	}
	s.Unconfig.Add(sc.NewSessionConf(name, id, origin))
}

func (s *SessionCtl) onConfigChange() {
	for _, conf := range s.Config.Sessions {
		s.Unconfig.RemoveOrigin(conf.Origin)
	}
}

func (s *SessionCtl) Close() {
	s.Config.Close()
	s.Unconfig.Close()
}
