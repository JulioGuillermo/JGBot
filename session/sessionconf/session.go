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

func (s *SessionCtl) AddUnconfig(name, id, origin, channel string) *sc.SessionConf {
	session := s.getNewConfig(name, id, origin, channel)
	s.Unconfig.SetSession(session)
	return &session
}

func (s *SessionCtl) AddConfig(name, id, origin, channel string) *sc.SessionConf {
	session := s.getNewConfig(name, id, origin, channel)
	s.Config.SetSession(session)
	return &session
}

func (s *SessionCtl) getNewConfig(name, id, origin, channel string) sc.SessionConf {
	session := sc.NewSessionConf(name, id, origin)
	session.ApplyDefConf(channel)
	return session
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
